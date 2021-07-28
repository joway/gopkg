package gopool

import (
	"sync"
	"testing"
	"time"

	"github.com/loov/hrtime/hrtesting"
)

const benchmarkTimes = 10000

func doCopyStack(cursor, n int) int {
	if cursor < n {
		return doCopyStack(cursor+1, n)
	}
	return 0
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
	{"Empty", func() { doCopyStack(0, 0) }},
	{"StackCopyLight", func() { doCopyStack(0, 10) }},
	{"StackCopyHeavy", func() { doCopyStack(0, 100) }},
	{"PureCalcLight", func() { testCalcHeavyFunc(1024) }},
	{"PureCalcHeavy", func() { testCalcHeavyFunc(102400) }},
	{"LongRT-10ms", func() { time.Sleep(time.Millisecond * 10) }},
	{"LongRT-50ms", func() { time.Sleep(time.Millisecond * 50) }},
}

func BenchmarkGoPool(b *testing.B) {
	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			p := NewPool("bench", defaultPoolCap, NewConfig())
			var wg sync.WaitGroup
			bench := hrtesting.NewBenchmark(b)
			defer bench.Report()
			b.ReportAllocs()
			b.ResetTimer()
			for bench.Next() {
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
			bench := hrtesting.NewBenchmark(b)
			defer bench.Report()
			b.ReportAllocs()
			b.ResetTimer()
			for bench.Next() {
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
