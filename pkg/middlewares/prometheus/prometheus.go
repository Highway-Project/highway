package prometheus

import (
	"github.com/Highway-Project/highway/pkg/middlewares"
	"net/http"
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
		method := r.Method
		HTTPRequestsCount.WithLabelValues(rule, service, backend, method, status).Inc()
	}
}

func New(params middlewares.MiddlewareParams) (middlewares.Middleware, error) {
	return PrometheusMiddleware{}, nil
}
