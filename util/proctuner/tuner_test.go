package proctuner

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestTuner(t *testing.T) {
	runtime.GOMAXPROCS(1)
	cpus := runtime.NumCPU()
	_ = Tuning(
		WithMaxProcs(cpus),
		WithMonitorFrequency(time.Millisecond*100),
		WithTuningFrequency(time.Millisecond),
	)

	var running = make([]int64, cpus)
	for i := 0; i < len(running); i++ {
		go func(id int) {
			sum := 0
			for x := 0; ; x++ {
				sum += x
				if x%1000000 == 0 {
					running[id]++
				}
			}
		}(i)
	}

	total := 0
	for x := 0; runtime.GOMAXPROCS(0) < cpus; x++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Println("main threads running", running[1:])
		total += x
	}
}
