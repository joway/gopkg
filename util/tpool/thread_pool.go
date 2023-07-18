package tpool

func New(size int) *threadPool {
	s := newScheduler(size)
	threads := make([]*Thread, size)
	for i := 0; i < size; i++ {
		threads[i] = newThread(s)
	}

	pool := &threadPool{
		scheduler: s,
		threads:   threads,
	}
	return pool
}

type threadPool struct {
	scheduler *scheduler
	threads   []*Thread
}

func (tp *threadPool) Submit(task task) {
	tp.scheduler.Add(task)
}

func (tp *threadPool) Close() {
	tp.scheduler.Close()
}
