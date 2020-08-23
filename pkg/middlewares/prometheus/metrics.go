package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "highway",
		Name:      "http_requests_count",
		Help:      "Total requests count",
	}, []string{"rule", "service", "backend", "method", "status"})
)

func init() {
	_ = prometheus.Register(HTTPRequestsCount)
}
