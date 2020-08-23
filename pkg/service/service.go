package service

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

type Service struct {
	Name     string
	Backends []Backend
	LB       LoadBalancer
}

func (s *Service) StartHealthCheck() {

}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := s.LB.Balance(s.getHealthyBackends())
	remote, err := url.Parse(uri.Addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`500 - Internal Server Error`))
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	director := proxy.Director
	var start time.Time
	var responseTime time.Duration
	proxy.Director = func(request *http.Request) {
		start = time.Now()
		director(request)
		ctx := r.Context()
		v := ctx.Value("metrics")
		var metrics map[string]string
		if v == nil {
			metrics = make(map[string]string)
		} else {
			metrics = v.(map[string]string)
		}
		metrics["service"] = s.Name
		metrics["backend"] = uri.Addr
		ctx = context.WithValue(ctx, "metrics", metrics)
		r = r.WithContext(ctx)
	}

	proxy.ModifyResponse = func(response *http.Response) error {
		responseTime = time.Since(start)
		ctx := r.Context()
		v := ctx.Value("metrics")
		var metrics map[string]string
		if v == nil {
			metrics = make(map[string]string)
		} else {
			metrics = v.(map[string]string)
		}
		metrics["status"] = strconv.Itoa(response.StatusCode)
		metrics["contentLength"] = strconv.FormatInt(response.ContentLength, 10)
		metrics["responseTime"] = strconv.FormatInt(int64(responseTime), 10)
		ctx = context.WithValue(ctx, "metrics", metrics)
		r = r.WithContext(ctx)

		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		responseTime = time.Since(start)
		ctx := r.Context()
		v := ctx.Value("metrics")
		var metrics map[string]string
		if v == nil {
			metrics = make(map[string]string)
		} else {
			metrics = v.(map[string]string)
		}
		metrics["status"] = "502"
		metrics["contentLength"] = "0"
		metrics["responseTime"] = strconv.FormatInt(int64(responseTime), 10)
		ctx = context.WithValue(ctx, "metrics", metrics)
		r = r.WithContext(ctx)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("502 - Bad Gateway"))
	}

	proxy.ServeHTTP(w, r)
}

func (s *Service) getHealthyBackends() []Backend {
	result := make([]Backend, 0)
	for _, backend := range s.Backends {
		if backend.Status == Available {
			result = append(result, backend)
		}
	}
	return result
}
