package autopool

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type Pool struct {
	cache  []*container
	shards uint32
	pos    uint32
}

func New(newer func() interface{}) *Pool {
	procs := runtime.GOMAXPROCS(0)
	p := &Pool{
		cache:  make([]*container, procs),
		shards: uint32(procs),
	}
	for i := 0; i < procs; i++ {
		p.cache[i] = newContainer(newer)
	}
	return p
}

func (p *Pool) Get() interface{} {
	pos := atomic.AddUint32(&p.pos, 1) % p.shards
	return p.cache[pos].Get()
}

func newContainer(newer func() interface{}) *container {
	return &container{
		sp: sync.Pool{
			New: newer,
		},
	}
}

type container struct {
	sp sync.Pool
}

func (c *container) Get() interface{} {
	o := c.sp.Get()
	runtime.SetFinalizer(o, c.put)
	return o
}

func (c *container) put(o interface{}) {
	c.sp.Put(o)
}
