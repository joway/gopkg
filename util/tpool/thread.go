package tpool

import "runtime"

type Thread struct {
	scheduler *scheduler
}

func newThread(scheduler *scheduler) *Thread {
	t := &Thread{
		scheduler: scheduler,
	}
	go t.run()
	return t
}

func (t *Thread) run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for {
		tsk := t.scheduler.Get()
		if tsk == nil {
			// closed
			return
		}
		tsk()
	}
}
