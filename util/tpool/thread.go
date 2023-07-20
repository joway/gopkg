package tpool

import (
	"runtime"
	"sync/atomic"
)

var idMaker int32

type Thread struct {
	id          int32
	distributor Distributor
	destructor  func()
}

func newThread(distributor Distributor, destructor func(), startup task) *Thread {
	id := atomic.AddInt32(&idMaker, 1)
	t := &Thread{
		id:          id,
		distributor: distributor,
		destructor:  destructor,
	}
	go t.run(startup)
	return t
}

func (t *Thread) ID() int {
	return int(t.id)
}

func (t *Thread) Distributor() Distributor {
	return t.distributor
}

func (t *Thread) run(startup task) {
	runtime.LockOSThread()
	defer func() {
		if t.destructor != nil {
			t.destructor()
		}
		runtime.UnlockOSThread()
	}()

	if startup != nil {
		startup()
	}
	var tsk task
	for {
		tsk = t.distributor.Get()
		if tsk == nil {
			// closed
			return
		}
		tsk()
	}
}
