package loadbalancer

import (
	"errors"
	"github.com/Highway-Project/highway/pkg/service"
)

var loadBalancerConstructors map[string]func() (service.LoadBalancer, error)

func init() {
	loadBalancerConstructors = make(map[string]func() (service.LoadBalancer, error))
}

func RegisterLoadBalancer(name string, constructor func() (service.LoadBalancer, error)) error {
	if _, exists := loadBalancerConstructors[name]; exists {
		return errors.New("LoadBalancer with this name exists: " + name)
	}

	loadBalancerConstructors[name] = constructor
	return nil
}

func NewLoadBalancer(name string) (service.LoadBalancer, error) {
	constructor, exists := loadBalancerConstructors[name]
	if !exists {
		return nil, errors.New("LoadBalancer with this name does not exists: " + name)
	}
	lb, err := constructor()
	if err != nil {
		return nil, err
	}
	return lb, nil
}
