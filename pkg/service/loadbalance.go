package service

import (
	"errors"
)

var loadBalancerConstructors map[string]func() (LoadBalancer, error)

func init() {
	loadBalancerConstructors = make(map[string]func() (LoadBalancer, error))
}

type LoadBalancer interface {
	Balance([]Backend) Backend
}

func RegisterLoadBalancer(name string, constructor func() (LoadBalancer, error)) error {
	if _, exists := loadBalancerConstructors[name]; exists {
		return errors.New("LoadBalancer with this name exists: " + name)
	}

	loadBalancerConstructors[name] = constructor
	return nil
}

func NewLoadBalancer(name string) (LoadBalancer, error) {
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
