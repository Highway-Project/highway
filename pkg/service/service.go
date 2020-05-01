package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Service struct {
	Name string
	Backends []Backend
	LB LoadBalancer
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := s.LB.Balance(s.Backends)
	remote, err := url.Parse(uri.Addr)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}





