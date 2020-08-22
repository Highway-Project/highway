package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"
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
