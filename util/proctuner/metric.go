package proctuner

import (
	"runtime/metrics"
)

const latencyMetricName = "/sched/latencies:seconds"

var (
	metricSamples = []metrics.Sample{{Name: latencyMetricName}}
)

func fetchSchedLatency() (p50, p90, p99, max float64) {
	metrics.Read(metricSamples)
	histogram := metricSamples[0].Value.Float64Histogram()

	var totalCount uint64
	var latestIdx int
	for idx, count := range histogram.Counts {
		if count > 0 {
			latestIdx = idx
		}
		totalCount += count
	}
	p50Count := totalCount / 2
	p90Count := uint64(float64(totalCount) * 0.90)
	p99Count := uint64(float64(totalCount) * 0.99)

	var cursor uint64
	for idx, count := range histogram.Counts {
		cursor += count
		if p99 == 0 && cursor >= p99Count {
			p99 = histogram.Buckets[idx]
		} else if p90 == 0 && cursor >= p90Count {
			p90 = histogram.Buckets[idx]
		} else if p50 == 0 && cursor >= p50Count {
			p50 = histogram.Buckets[idx]
		}
	}
	max = histogram.Buckets[latestIdx]
	return
}
