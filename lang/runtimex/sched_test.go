package runtimex

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

//go:noinline
func testStackFunction(n int) int {
	var stack [1024 * 4]int64
	for i := 0; i < len(stack); i++ {
		stack[i] = int64(n + i)
	}
	return int(stack[len(stack)/2])
}

func TestCPUBondSingleP(t *testing.T) {
	runtime.GOMAXPROCS(4)
	const (
		tasks    = 8
		interval = 1000000
		round    = 1000
	)
	begin := time.Now()
	var wg sync.WaitGroup
	for t := 0; t < tasks; t++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if n == 0 {
				MPreemptOff()
				defer MPreemptOn()
			}
			for i := 0; i < round*interval; i++ {
				if i%interval == 0 {
					testStackFunction(i)
				}
			}
		}(t)
	}
	wg.Wait()
	cost := time.Now().Sub(begin)
	t.Logf("Cost: %vms", cost.Milliseconds())
}

func TestIOBondSingleP(t *testing.T) {
	runtime.GOMAXPROCS(4)
	const (
		tasks    = 8
		interval = 1000000
		round    = 1000
	)
	begin := time.Now()
	var wg sync.WaitGroup
	for t := 0; t < tasks; t++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if n == 0 {
				MPreemptOff()
				defer MPreemptOn()
			}
			for i := 0; i < round*interval; i++ {
				if i%interval == 0 {
					testStackFunction(i)
				}
				if i%interval == interval/2 {
					time.Sleep(time.Millisecond)
				}
			}
		}(t)
	}
	wg.Wait()
	cost := time.Now().Sub(begin)
	t.Logf("Cost: %vms", cost.Milliseconds())
}
