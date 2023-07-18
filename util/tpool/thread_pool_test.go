package tpool

import (
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	runtime.GOMAXPROCS(2)
	threads := 4
	p := New(threads)
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
	t.Logf("Currnt thread count: %d", threadProfile.Count())
}

func TestCPUBond(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	runtime.GOMAXPROCS(2)
	threads := 4
	p := New(threads)
	defer p.Close()

	for round := threads; round <= 128; round *= 2 {
		var wg sync.WaitGroup
		begin := time.Now()
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
		cost := time.Now().Sub(begin)
		t.Logf("Round[%d]: cost %d ms", round, cost.Milliseconds())
	}
	t.Logf("Currnt thread count: %d", threadProfile.Count())
}

func TestHugeThreads(t *testing.T) {
	var threadProfile = pprof.Lookup("threadcreate")
	runtime.GOMAXPROCS(2)
	threads := 32
	p := New(threads)
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
	t.Logf("Currnt thread count: %d", threadProfile.Count())
}
