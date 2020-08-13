package service

import (
	"errors"
	"fmt"
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/logging"
	"github.com/Highway-Project/highway/pkg/service"
)

var services map[string]*service.Service

func init() {
	services = make(map[string]*service.Service)
}

func NewService(spec config.ServiceSpec) (*service.Service, error) {
	if _, exists := services[spec.Name]; exists {
		msg := fmt.Sprintf("service %s already exists", spec.Name)
		logging.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	bs, err := NewBackends(spec.BackendsSpecs)
	if err != nil {
		msg := fmt.Sprintf("could not create backends for service %s", spec.Name)
		logging.Logger.WithError(err).Error(msg)
		return nil, errors.New(msg)
	}

	lb, err := NewLoadBalancer(spec.LoadBalancerName)
	if err != nil {
		msg := fmt.Sprintf("could not create loadbalancer %s for service %s", spec.LoadBalancerName, spec.Name)
		logging.Logger.WithError(err).Error(msg)
		return nil, errors.New(msg)
	}

	s := service.Service{
		Name:     spec.Name,
		Backends: bs,
		LB:       lb,
	}

	services[spec.Name] = &s
	return &s, nil
}

func GetServiceByName(name string) (*service.Service, error) {
	if s, exists := services[name]; exists {
		return s, nil
	}
	return nil, errors.New(fmt.Sprintf("service %s does not exists", name))
}
