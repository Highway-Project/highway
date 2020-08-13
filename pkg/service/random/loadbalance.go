package random

import (
	"github.com/Highway-Project/highway/pkg/service"
	"math/rand"
)

type RandomLoadBalancer struct{}

func (r RandomLoadBalancer) Balance(backends []service.Backend) service.Backend {
	index := rand.Intn(len(backends))
	return backends[index]
}

func New() (service.LoadBalancer, error) {
	return RandomLoadBalancer{}, nil
}
