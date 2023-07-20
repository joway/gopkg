package tpool

import (
	"runtime"
	"sync"
	"sync/atomic"
)

func NewNativeThreadPool() ThreadPool {
	tp := new(nativeThreadPool)
	return tp
}

type nativeThreadPool struct {
	mu    sync.Mutex
	state int32 // 0: running, -1: closed
}

func (tp *nativeThreadPool) Size() (size int) {
	return runtime.GOMAXPROCS(0)
}

func (tp *nativeThreadPool) Submit(t task) {
	if atomic.LoadInt32(&tp.state) < 0 { // closed
		return
	}
	go func() {
		t()
	}()
}

func (tp *nativeThreadPool) Close() {
	atomic.CompareAndSwapInt32(&tp.state, 0, -1)
}
