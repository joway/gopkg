// Copyright 2021 ByteDance Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gopool

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync/atomic"

	"github.com/bytedance/gopkg/util/logger"
)

type Pool interface {
	// Name returns the corresponding pool name.
	Name() string
	// SetCap sets the goroutine capacity of the pool.
	SetCap(cap int32)
	// Go executes f.
	Go(f func())
	// CtxGo executes f and accepts the context.
	CtxGo(ctx context.Context, f func())
	// SetPanicHandler sets the panic handler.
	SetPanicHandler(f func(context.Context, interface{}))
}

type task func()

type pool struct {
	// The name of the pool
	name string

	// Capacity of the pool, the maximum number of goroutines that are actually working
	cap int32
	// Configuration information
	config *Config
	// List of pending tasks
	taskList chan task

	// Record the number of running workers
	workerCount int32

	// This method will be called when the worker panic
	panicHandler func(context.Context, interface{})
}

// NewPool creates a new pool with the given name, cap and config.
func NewPool(name string, cap int32, config *Config) Pool {
	cpus := runtime.GOMAXPROCS(0)
	if cpus < 1 {
		cpus = 1
	}
	p := &pool{
		name:     name,
		cap:      cap,
		config:   config,
		taskList: make(chan task, cpus),
	}
	p.initDaemonWorkers(cpus)
	return p
}

func (p *pool) Name() string {
	return p.name
}

func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := p.spawnTask(ctx, f)

	// try to send task to other running worker
	select {
	case p.taskList <- t:
		return
	default:
		// no running workers waiting for tasked
	}

	// blocking when out of cap
	for p.WorkerCount() >= p.Cap() {
		runtime.Gosched()
	}
	// start a new worker when the attempt of reusing worker failed
	p.spawnWorker(t)
}

// SetPanicHandler the func here will be called after the panic has been recovered.
func (p *pool) SetPanicHandler(f func(context.Context, interface{})) {
	p.panicHandler = f
}

func (p *pool) Cap() int32 {
	return atomic.LoadInt32(&p.cap)
}

func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}

func (p *pool) spawnTask(ctx context.Context, f func()) task {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("GOPOOL: panic in pool: %s: %v: %s", p.name, r, debug.Stack())
				logger.CtxErrorf(ctx, msg)
				if p.panicHandler != nil {
					p.panicHandler(ctx, r)
				}
			}
		}()
		f()
	}
}

func (p *pool) initDaemonWorkers(size int) {
	for i := 0; i < size; i++ {
		go p.spawnDaemonWorker()
	}
}

func (p *pool) spawnDaemonWorker() {
	p.incWorkerCount()
	defer p.decWorkerCount()
	lifetime := 0
	for t := range p.taskList {
		lifetime++
		t()

		if lifetime >= DefaultWorkerLifetime {
			//start a new one to clear stack
			go p.spawnDaemonWorker()
			return
		}
	}
}

func (p *pool) spawnWorker(initialTask task) {
	p.incWorkerCount()
	go func() {
		if initialTask != nil {
			initialTask()
		}

		for {
			var t task
			select {
			case t = <-p.taskList:
				t()
			default:
				// if there's no task to do, exit
				p.decWorkerCount()
				return
			}
		}
	}()
}
