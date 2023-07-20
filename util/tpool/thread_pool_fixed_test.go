package tpool

import (
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedThreadPool_Sleep(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	runtime.GOMAXPROCS(2)
	threads := 4
	p := NewFixedThreadPool(threads)
	defer p.Close()

	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		p.Submit(func() {
			time.Sleep(time.Millisecond * 10)
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, p.Size(), threads)
	t.Logf("Currnt threads count: %d", threadProfile.Count())
}

func TestFixedThreadPool_CPUBond(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	runtime.GOMAXPROCS(2)
	threads := 4
	p := NewFixedThreadPool(threads)
	defer p.Close()

	for tasks := threads; tasks <= 128; tasks *= 2 {
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
	assert.Equal(t, p.Size(), threads)
	t.Logf("Currnt threads count: %d", threadProfile.Count())
}

func TestFixedThreadPool_HugeThreads(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	runtime.GOMAXPROCS(2)
	threads := 32
	p := NewFixedThreadPool(threads)
	defer p.Close()

	var wg sync.WaitGroup
	round := threads * 16
	wg.Add(round)
	for i := 0; i < round; i++ {
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
	assert.Equal(t, p.Size(), threads)
	t.Logf("Currnt threads count: %d", threadProfile.Count())
}
