package tpool

import "sync"

type task func()

func newScheduler(threads int) *scheduler {
	s := &scheduler{
		threads:  threads,
		tasks:    make([]task, 0, 1024),
		notifier: make([]chan struct{}, 0, 1024),
	}
	return s
}

type scheduler struct {
	mu       sync.Mutex
	state    int    // 0: running, -1: closed
	threads  int    // threads number
	tasks    []task // LIFO: We want to make most tasks have the fastest latency
	notifier []chan struct{}
}

func (s *scheduler) Close() {
	s.mu.Lock()
	s.state = -1
	for i := 0; i < len(s.notifier); i++ {
		notify := s.notifier[i]
		notify <- struct{}{}
	}
	s.mu.Unlock()
}

func (s *scheduler) Add(t task) {
	var notify chan struct{}
	s.mu.Lock()
	if s.state < 0 { // closed
		return
	}

	waits := len(s.notifier)
	s.tasks = append(s.tasks, t)
	if waits > 0 {
		notify = s.notifier[waits-1]
		s.notifier = s.notifier[:waits-1]
	}
	s.mu.Unlock()
	if notify != nil {
		notify <- struct{}{}
	}
}

func (s *scheduler) Get() (t task) {
	var notify chan struct{}
GET:
	s.mu.Lock()
	if s.state < 0 { // closed
		return
	}

	size := len(s.tasks)
	if size > 0 {
		t = s.tasks[size-1]
		s.tasks = s.tasks[:size-1]
		s.mu.Unlock()
		if notify != nil {
			close(notify)
		}
		return t
	}
	if notify == nil {
		notify = make(chan struct{}, 1)
	}
	s.notifier = append(s.notifier, notify)
	s.mu.Unlock()

	<-notify // thread go to sleep
	// thread wakeup
	goto GET
}
