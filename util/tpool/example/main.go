package main

import (
	"log"
	"runtime"
	"time"

	"github.com/bytedance/gopkg/util/tpool"
)

func cpuWork() int {
	var sum int
	for i := 0; ; i++ {
		sum++
	}
}

func main() {
	maxProcs := 2
	runtime.GOMAXPROCS(maxProcs)
	p := tpool.NewCachedThreadPool()
	defer p.Close()

	//for i := 0; i < 1; i++ {
	//	go cpuWork()
	//}
	for i := 0; i < maxProcs+1; i++ {
		log.Printf("create new thread: %d", i)
		p.Submit(func() {
			cpuWork()
		})
	}

	for i := 0; ; i++ {
		time.Sleep(time.Second)
		log.Printf("Main: %d", i)
	}
}
