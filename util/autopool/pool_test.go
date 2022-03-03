package autopool

import (
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testObj struct {
	Id int
}

func TestPool(t *testing.T) {
	runtime.GOMAXPROCS(1)
	var counter uint32
	p := New(func() interface{} {
		return &testObj{
			Id: int(atomic.AddUint32(&counter, 1)),
		}
	})

	refMap := make(map[int]int)
	for i := 1; i <= 100; i++ {
		o := p.Get().(*testObj)
		refMap[o.Id]++
	}
	assert.Equal(t, len(refMap), 100)
	runtime.GC()
	for i := 1; i <= 100; i++ {
		o := p.Get().(*testObj)
		refMap[o.Id]++
	}
	assert.True(t, len(refMap) < 150)
}
