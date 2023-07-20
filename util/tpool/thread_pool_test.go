package tpool

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type benchcase struct {
	name       string
	threadPool ThreadPool
}

func BenchmarkCPUTasks(b *testing.B) {
	cases := []benchcase{
		{name: "NewNativeThreadPool", threadPool: NewFixedThreadPool(runtime.NumCPU())},
		//{name: "FixedThreadPool", threadPool: NewFixedThreadPool(runtime.NumCPU())},
		{name: "CachedThreadPool", threadPool: NewCachedThreadPool(WithCachedMaxIdleThreads(32))},
	}
	defer func() {
		for _, c := range cases {
			c.threadPool.Close()
		}
	}()

	for _, c := range cases {
		b.Run(fmt.Sprintf("%s", c.name), func(b *testing.B) {
			maxTasks := 32
			for tasks := 1; tasks <= maxTasks; tasks *= 2 {
				b.Run(fmt.Sprintf("Tasks[%d]", tasks), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						var wg sync.WaitGroup
						for t := 0; t < tasks; t++ {
							wg.Add(1)
							c.threadPool.Submit(func() {
								defer wg.Done()
								sum := 0
								for x := 0; x < 10000000; x++ {
									sum += x
								}
								_ = sum
							})
						}
						wg.Wait()
					}
				})
			}
		})
	}
}

func BenchmarkIOTasks(b *testing.B) {
	cases := []benchcase{
		{name: "NativeThreadPool", threadPool: NewFixedThreadPool(runtime.NumCPU())},
		//{name: "FixedThreadPool", threadPool: NewFixedThreadPool(runtime.NumCPU())},
		{name: "CachedThreadPool", threadPool: NewCachedThreadPool(WithCachedMaxIdleThreads(32))},
	}
	defer func() {
		for _, c := range cases {
			c.threadPool.Close()
		}
	}()

	for _, c := range cases {
		b.Run(fmt.Sprintf("%s", c.name), func(b *testing.B) {
			maxTasks := 32
			for tasks := 1; tasks <= maxTasks; tasks *= 2 {
				b.Run(fmt.Sprintf("Tasks[%d]", tasks), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						var wg sync.WaitGroup
						for t := 0; t < tasks; t++ {
							wg.Add(1)
							c.threadPool.Submit(func() {
								defer wg.Done()
								time.Sleep(time.Millisecond * 10)
							})
						}
						wg.Wait()
					}
				})
			}
		})
	}
}

func BenchmarkGoroutineVsThreads(b *testing.B) {
	runtime.GOMAXPROCS(2)
	ioWork := func() {
		time.Sleep(time.Millisecond * 10)
	}
	cpuWork := func() {
		var mem [1000]int
		for i := 0; i < 100000; i++ {
			mem[i%len(mem)] += i
		}
		_ = mem
	}
	type benchcase struct {
		kind   string
		runner func(workload func(), count int)
		count  int
	}
	var benchcases = []benchcase{
		{kind: "Goroutines", runner: goroutineRunner, count: 1},
		{kind: "Goroutines", runner: goroutineRunner, count: 4},
		{kind: "Goroutines", runner: goroutineRunner, count: 8},
		{kind: "Goroutines", runner: goroutineRunner, count: 16},
		{kind: "Goroutines", runner: goroutineRunner, count: 32},
		{kind: "Goroutines", runner: goroutineRunner, count: 100},

		{kind: "Threads", runner: threadRunner, count: 1},
		{kind: "Threads", runner: threadRunner, count: 4},
		{kind: "Threads", runner: threadRunner, count: 8},
		{kind: "Threads", runner: threadRunner, count: 16},
		{kind: "Threads", runner: threadRunner, count: 32},
		{kind: "Threads", runner: threadRunner, count: 100},
	}
	for _, c := range benchcases {
		b.Run(fmt.Sprintf("%s[%d]", c.kind, c.count), func(b *testing.B) {
			b.Run("IO Work", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					c.runner(ioWork, c.count)
				}
			})
			b.Run("CPU Work", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					c.runner(cpuWork, c.count)
				}
			})
		})
	}
}

func goroutineRunner(workload func(), count int) {
	var wg sync.WaitGroup
	for t := 0; t < count; t++ {
		wg.Add(1)
		go func() {
			workload()
			wg.Done()
		}()
	}
	wg.Wait()
}

var testThreadPool = NewCachedThreadPool(WithCachedMaxIdleThreads(32))

func threadRunner(workload func(), count int) {
	var wg sync.WaitGroup
	for t := 0; t < count; t++ {
		wg.Add(1)
		testThreadPool.Submit(func() {
			workload()
			wg.Done()
		})
	}
	wg.Wait()
}
