package tpool

import (
	"fmt"
	"sync"
	"testing"
)

type benchcase struct {
	name       string
	threadPool ThreadPool
}

func BenchmarkCPUTasks(b *testing.B) {
	cases := []benchcase{
		{name: "FixedThreadPool-4Threads", threadPool: NewFixedThreadPool(4)},
		{name: "CachedThreadPool-UnlimitedThreads", threadPool: NewCachedThreadPool(WithCachedMaxIdleThreads(32))},
	}
	defer func() {
		for _, c := range cases {
			c.threadPool.Close()
		}
	}()

	for _, c := range cases {
		b.Run(fmt.Sprintf("%s", c.name), func(b *testing.B) {
			maxCPUTasks := 32
			for tasks := 1; tasks <= maxCPUTasks; tasks *= 2 {
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
