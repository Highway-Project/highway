package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "highway",
		Name:      "http_requests_count",
		Help:      "Total requests count",
	}, []string{"rule", "service", "backend", "method", "status"})

	HTTPResponseSizeBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "highway",
		Name:        "http_response_size_bytes",
		Help:        "Total http responses size",
		ConstLabels: nil,
	}, []string{"rule", "service", "backend", "method", "status"})

	HTTPResponseTimeSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "highway",
		Name:      "http_response_time_seconds",
		Help:      "Total http response time",
		Buckets:   prometheus.DefBuckets,
	}, []string{"rule", "service", "backend", "method", "status"})
)

func init() {
	_ = prometheus.Register(HTTPRequestsCount)
	_ = prometheus.Register(HTTPResponseSizeBytes)
	_ = prometheus.Register(HTTPResponseTimeSeconds)
}
