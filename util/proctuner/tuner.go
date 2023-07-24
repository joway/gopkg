package proctuner

import (
	"log"
	"runtime"
	"time"
)

const (
	defaultLatencyAcceptable   float64 = 0.001 // 1ms
	defaultP50LatencyThreshold float64 = 0.01  // 10ms
	defaultP90LatencyThreshold float64 = 0.01  // 20ms
	defaultP99LatencyThreshold float64 = 0.1   // 100ms
	defaultMonitorFrequency            = time.Second * 10
	defaultTuningFrequency             = time.Minute
)

var originMaxProcs int
var MaxProcs = 32

type Option func(t *tuner)

func WithMaxProcs(maxProcs int) Option {
	return func(t *tuner) {
		t.maxProcs = maxProcs
	}
}

func WithP50LatencyThreshold(threshold float64) Option {
	return func(t *tuner) {
		t.p50Threshold = threshold
	}
}

func WithP90LatencyThreshold(threshold float64) Option {
	return func(t *tuner) {
		t.p90Threshold = threshold
	}
}

func WithP99LatencyThreshold(threshold float64) Option {
	return func(t *tuner) {
		t.p99Threshold = threshold
	}
}

func WithMaxLatencyThreshold(threshold float64) Option {
	return func(t *tuner) {
		t.maxThreshold = threshold
	}
}

func WithMonitorFrequency(duration time.Duration) Option {
	return func(t *tuner) {
		t.monitorFrequency = duration
	}
}

func WithTuningFrequency(duration time.Duration) Option {
	return func(t *tuner) {
		t.tuningFrequency = duration
	}
}

func Tuning(opts ...Option) error {
	t := new(tuner)
	t.acceptable = defaultLatencyAcceptable
	t.p50Threshold = defaultP50LatencyThreshold
	t.p90Threshold = defaultP90LatencyThreshold
	t.p99Threshold = defaultP99LatencyThreshold
	t.monitorFrequency = defaultMonitorFrequency
	t.tuningFrequency = defaultTuningFrequency
	for _, opt := range opts {
		opt(t)
	}
	maxProcs := t.maxProcs
	originMaxProcs = runtime.GOMAXPROCS(0)
	// default tuning
	if maxProcs <= 0 {
		maxProcs = originMaxProcs * 3
	}
	// reduce to MaxProcs
	if maxProcs > MaxProcs {
		maxProcs = MaxProcs
	}
	// no need to tuning
	if maxProcs <= originMaxProcs {
		return nil
	}
	t.minProcs = originMaxProcs
	t.maxProcs = maxProcs
	go t.tuning()
	return nil
}

type tuner struct {
	minProcs         int
	maxProcs         int
	acceptable       float64
	p50Threshold     float64
	p90Threshold     float64
	p99Threshold     float64
	maxThreshold     float64
	monitorFrequency time.Duration
	tuningFrequency  time.Duration
}

func (t *tuner) tuning() {
	log.Printf("ProcTuning: MinProcs=%d MaxProcs=%d", t.minProcs, t.maxProcs)

	var lastModify = time.Now()
	var p50, p90, p99, max float64
	var currentProcs = t.minProcs
	for {
		time.Sleep(t.monitorFrequency)
		p50, p90, p99, max = fetchSchedLatency()
		if p50 >= t.p50Threshold || p90 >= t.p90Threshold || p99 >= t.p99Threshold {
			if currentProcs == t.maxProcs {
				continue
			}
			now := time.Now()
			if now.Sub(lastModify) < t.tuningFrequency {
				continue
			}
			currentProcs += 1
			runtime.GOMAXPROCS(currentProcs)
			lastModify = now
			log.Printf("ProcTuning: GOMAXPROCS from %d to %d", currentProcs-1, currentProcs)
		} else if p99 <= t.acceptable {
			if currentProcs <= t.minProcs {
				continue
			}
			now := time.Now()
			if now.Sub(lastModify) >= t.tuningFrequency {
				continue
			}
			currentProcs -= 1
			runtime.GOMAXPROCS(currentProcs)
			lastModify = now
			log.Printf("ProcTuning: GOMAXPROCS from %d to %d", currentProcs+1, currentProcs)
		}
		log.Printf("ProcTuning: Scheduler Latency[p50=%.6f,p90=%.6f,p99=%.6f,max=%.6f]", p50, p90, p99, max)
	}
}
