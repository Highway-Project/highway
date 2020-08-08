package random

import (
	"github.com/Highway-Project/highway/pkg/service"
	"math/rand"
	"time"
)

func Register() {
	rand.Seed(time.Now().Unix())
	err := service.RegisterLoadBalancer("random", New)
	if err != nil {
		panic("could not register random load balancer")
	}
}

type RandomLoadBalancer struct{}

func (r RandomLoadBalancer) Balance(backends []service.Backend) service.Backend {
	index := rand.Intn(len(backends))
	return backends[index]
}

func New() (service.LoadBalancer, error) {
	return RandomLoadBalancer{}, nil
}
