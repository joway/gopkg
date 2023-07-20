package tpool

import (
	"sync"
)

var _ ThreadPool = (*cachedThreadPool)(nil)

type CachedThreadPoolOption func(tp *cachedThreadPool)

func WithCachedMaxIdleThreads(maxIdle int) CachedThreadPoolOption {
	return func(tp *cachedThreadPool) {
		tp.maxIdle = maxIdle
	}
}

func NewCachedThreadPool(opts ...CachedThreadPoolOption) ThreadPool {
	tp := &cachedThreadPool{
		state:   0,
		idle:    make([]Distributor, 0, 1024),
		threads: map[int]*Thread{},
		maxIdle: -1,
	}
	for _, opt := range opts {
		opt(tp)
	}
	return tp
}

type cachedThreadPool struct {
	mu      sync.Mutex
	state   int
	size    int
	idle    []Distributor
	threads map[int]*Thread

	maxIdle int // maxIdle < 0 means unlimited
}

func (p *cachedThreadPool) Size() (size int) {
	p.mu.Lock()
	size = p.size
	p.mu.Unlock()
	return size
}

func (p *cachedThreadPool) IdleSize() (size int) {
	p.mu.Lock()
	size = len(p.idle)
	p.mu.Unlock()
	return size
}

func (p *cachedThreadPool) Submit(tsk task) {
	p.mu.Lock()
	if p.state < 0 {
		p.mu.Unlock()
		return
	}

	// select distributor
	var dist Distributor
	wrapTask := func() {
		tsk()
		isIdle := false
		p.mu.Lock()
		if p.state >= 0 && (p.maxIdle < 0 || len(p.idle) < p.maxIdle) {
			isIdle = true
			p.idle = append(p.idle, dist)
		}
		p.mu.Unlock()
		if !isIdle {
			dist.Close()
		}
	}
	idleSize := len(p.idle)
	if idleSize == 0 {
		p.size++
		dist = newDistributor(1)
		thread := newThread(
			dist,
			func() {
				p.mu.Lock()
				p.size--
				p.mu.Unlock()
			},
			wrapTask,
		)
		p.threads[thread.ID()] = thread
		p.mu.Unlock()
		return
	}
	dist = p.idle[idleSize-1]
	p.idle = p.idle[:idleSize-1]
	p.mu.Unlock()
	dist.Add(wrapTask)
}

func (p *cachedThreadPool) Close() {
	p.mu.Lock()
	p.state = -1
	for _, dist := range p.idle {
		dist.Close()
	}
	p.mu.Unlock()
}
