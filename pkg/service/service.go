package service

import "net/http"

type Service struct {
	Name string
	Backends []Backend
	LB LoadBalancer
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}





