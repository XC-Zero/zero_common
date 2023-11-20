package metric

import "runtime/metrics"

func Metrics() []metrics.Description {
	all := metrics.All()
	return all
}
