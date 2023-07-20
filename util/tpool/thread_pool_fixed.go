package tpool

import (
	"sync"
	"sync/atomic"
)

func NewFixedThreadPool(threads int) ThreadPool {
	tp := new(fixedThreadPool)
	tp.distributor = newDistributor(1024)
	tp.threads = make([]*Thread, threads)
	for i := 0; i < threads; i++ {
		tp.threads[i] = newThread(tp.distributor, nil, nil)
	}
	return tp
}

type fixedThreadPool struct {
	mu          sync.Mutex
	state       int32 // 0: running, -1: closed
	distributor Distributor
	threads     []*Thread
}

func (tp *fixedThreadPool) Size() (size int) {
	return len(tp.threads)
}

func (tp *fixedThreadPool) Submit(t task) {
	if atomic.LoadInt32(&tp.state) < 0 { // closed
		return
	}
	tp.distributor.Add(t)
}

func (tp *fixedThreadPool) Close() {
	atomic.CompareAndSwapInt32(&tp.state, 0, -1)
	tp.distributor.Close()
}
