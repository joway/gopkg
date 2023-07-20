package tpool

type task func()

type ThreadPool interface {
	Size() int
	Submit(t task)
	Close()
}

type Distributor interface {
	Add(tsk task)
	Get() (tsk task)
	Close()
}

type distributor chan task

func newDistributor(size int) Distributor {
	return make(distributor, size)
}

func (d distributor) Add(tsk task) {
	select {
	case d <- tsk:
	default:
	}
}

func (d distributor) Get() task {
	return <-d
}

func (d distributor) Close() {
	close(d)
}
