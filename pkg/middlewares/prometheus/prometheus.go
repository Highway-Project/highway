package prometheus

import (
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

type PrometheusMiddleware struct {
	HTTPRequestsCount       *prometheus.CounterVec
	HTTPResponseSizeBytes   *prometheus.CounterVec
	HTTPResponseTimeSeconds *prometheus.HistogramVec
}

func (pm *PrometheusMiddleware) Process(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)

		v := r.Context().Value("metrics")
		var metrics map[string]string
		if v == nil {
			return
		}
		metrics, ok := v.(map[string]string)
		if !ok {
			return
		}

		rule, ok := metrics["rule"]
		if !ok {
			return
		}

		service, ok := metrics["service"]
		if !ok {
			return
		}

		backend, ok := metrics["backend"]
		if !ok {
			return
		}

		status, ok := metrics["status"]
		if !ok {
			return
		}

		length, ok := metrics["contentLength"]
		if !ok {
			return
		}
		contentLength, err := strconv.ParseFloat(length, 64)
		if err != nil {
			return
		}

		rt, ok := metrics["responseTime"]
		if !ok {
			return
		}
		responseTime, err := strconv.ParseFloat(rt, 64)
		if err != nil {
			return
		}
		responseTime = responseTime * 10e-9 // Convert nanosecond to second

		method := r.Method

		pm.HTTPRequestsCount.WithLabelValues(rule, service, backend, method, status).Inc()
		pm.HTTPResponseSizeBytes.WithLabelValues(rule, service, backend, method, status).Add(contentLength)
		pm.HTTPResponseTimeSeconds.WithLabelValues(rule, service, backend, method, status).Observe(responseTime)
	}
}

func (pm *PrometheusMiddleware) register() error {
	err := prometheus.Register(pm.HTTPRequestsCount)
	if err != nil {
		return err
	}

	err = prometheus.Register(pm.HTTPResponseSizeBytes)
	if err != nil {
		return err
	}

	err = prometheus.Register(pm.HTTPResponseTimeSeconds)
	if err != nil {
		return err
	}

	return nil
}

func New(params middlewares.MiddlewareParams) (middlewares.Middleware, error) {
	p := PrometheusMiddleware{}
	p.HTTPRequestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "highway",
		Name:      "http_requests_count",
		Help:      "Total requests count",
	}, []string{"rule", "service", "backend", "method", "status"})
	p.HTTPResponseSizeBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "highway",
		Name:        "http_response_size_bytes",
		Help:        "Total http responses size",
		ConstLabels: nil,
	}, []string{"rule", "service", "backend", "method", "status"})
	p.HTTPResponseTimeSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "highway",
		Name:      "http_response_time_seconds",
		Help:      "Total http response time",
		Buckets:   prometheus.DefBuckets,
	}, []string{"rule", "service", "backend", "method", "status"})

	err := p.register()
	if err != nil {
		return nil, err
	}

	return &p, nil
}
