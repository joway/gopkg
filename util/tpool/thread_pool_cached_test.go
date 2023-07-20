package tpool

import (
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCachedThreadPool_Sleep(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	runtime.GOMAXPROCS(2)
	round := 32
	p := NewCachedThreadPool()
	defer p.Close()

	var wg sync.WaitGroup
	wg.Add(round)
	for i := 0; i < round; i++ {
		p.Submit(func() {
			time.Sleep(time.Millisecond * 100)
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, round, p.Size())
	t.Logf("Currnt threads count: %d", threadProfile.Count())
}

func TestCachedThreadPool_CPUBond(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	//runtime.GOMAXPROCS(2)
	cpus := runtime.NumCPU()
	p := NewCachedThreadPool(WithCachedMaxIdleThreads(cpus))
	defer p.Close()

	maxTasks := 32
	for tasks := 1; tasks <= maxTasks; tasks *= 2 {
		var wg sync.WaitGroup
		begin := time.Now()
		wg.Add(tasks)
		for i := 0; i < tasks; i++ {
			p.Submit(func() {
				var sum int
				for x := 0; x <= 100000000; x++ {
					sum += x
				}
				_ = sum
				wg.Done()
			})
		}
		wg.Wait()
		cost := time.Now().Sub(begin)
		t.Logf("Tasks[%d]: cost %d ms", tasks, cost.Milliseconds())
	}
	assert.Equal(t, cpus, p.Size())
	t.Logf("Currnt threads count: %d", threadProfile.Count())
}
