package prometheus

import (
	"github.com/Highway-Project/highway/pkg/middlewares"
	"net/http"
	"strconv"
)

type PrometheusMiddleware struct {
}

func (pm PrometheusMiddleware) Process(handler http.HandlerFunc) http.HandlerFunc {
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

		HTTPRequestsCount.WithLabelValues(rule, service, backend, method, status).Inc()
		HTTPResponseSizeBytes.WithLabelValues(rule, service, backend, method, status).Add(contentLength)
		HTTPResponseTimeSeconds.WithLabelValues(rule, service, backend, method, status).Observe(responseTime)
	}
}

func New(params middlewares.MiddlewareParams) (middlewares.Middleware, error) {
	return PrometheusMiddleware{}, nil
}
