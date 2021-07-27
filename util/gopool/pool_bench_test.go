package gopool

import (
	"runtime"
	"sync"
	"testing"
)

const benchmarkTimes = 10000

func doCopyStack(a, b int) int {
	if b < 100 {
		return doCopyStack(0, b+1)
	}
	return 0
}

func testStackHeavyFunc() {
	doCopyStack(0, 0)
}

func testCalcHeavyFunc(n int) (sum int) {
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

type benchcase struct {
	name    string
	handler func()
}

var benchmarkCases = []benchcase{
	{"StackCopyHeavy", func() { testStackHeavyFunc() }},
	{"PureCalc", func() { testCalcHeavyFunc(1024) }},
}

func BenchmarkGoPool(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			config := NewConfig()
			config.ScaleThreshold = 1
			p := NewPool("benchmark", int32(runtime.GOMAXPROCS(0)), config)
			var wg sync.WaitGroup
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				wg.Add(benchmarkTimes)
				for j := 0; j < benchmarkTimes; j++ {
					p.Go(func() {
						bc.handler()
						wg.Done()
					})
				}
				wg.Wait()
			}
		})
	}
}

func BenchmarkGo(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			var wg sync.WaitGroup
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				wg.Add(benchmarkTimes)
				for j := 0; j < benchmarkTimes; j++ {
					go func() {
						bc.handler()
						wg.Done()
					}()
				}
				wg.Wait()
			}
		})
	}
}
