package runtimex

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"
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

func TestRuntimeStatus(t *testing.T) {
	runtime.GOMAXPROCS(4)
	var stop int32
	var wg sync.WaitGroup
	for t := 0; t < 8; t++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for i := 1; atomic.LoadInt32(&stop) == 0; i++ {
				testStackFunction(i)
				if i%100000 == 0 {
					g := GetG()
					m := g.M()
					p := m.P()
					log.Printf(
						"[G] gid=%d,gpreempt=%v | "+
							"[M] mid=%d,mpreemptoff=%s | "+
							"[P] pid=%d,pstatus=%d,psysmontick=%v,ppreempt=%v,qsize=%d",
						*g.Id(), *g.Preempt(),
						*m.Id(), *m.PreemptOff(),
						*p.Id(), *p.Status(), *p.Sysmontick(), *p.Preempt(), p.RunqSize(),
					)
				}
			}
		}(t)
	}
	time.Sleep(time.Second * 3)
	atomic.StoreInt32(&stop, 1)
	wg.Wait()
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
			g := GetG()
			m := g.M()
			if n == 0 {
				*m.PreemptOff() = "holding"
				defer func() {
					*m.PreemptOff() = ""
				}()
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
			g := GetG()
			m := g.M()
			p := m.P()
			if n == 0 {
				*m.PreemptOff() = "holding"
				defer func() {
					*m.PreemptOff() = ""
				}()
			}
			for i := 0; i < round*interval; i++ {
				if i%interval == 0 {
					testStackFunction(i)
					log.Printf(
						"[G] gid=%d,gpreempt=%v | "+
							"[M] mid=%d,mpreemptoff=%s | "+
							"[P] pid=%d,pstatus=%d,psysmontick=%v,ppreempt=%v,qsize=%d",
						*g.Id(), *g.Preempt(),
						*m.Id(), *m.PreemptOff(),
						*p.Id(), *p.Status(), *p.Sysmontick(), *p.Preempt(), p.RunqSize(),
					)
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
