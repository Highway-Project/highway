package service

import (
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/pkg/service"
)

func NewBackends(specs []config.BackendSpec) ([]service.Backend, error) {
	backends := make([]service.Backend, 0)
	for _, spec := range specs {
		backend := service.Backend{
			Name:   spec.BackendName,
			Addr:   spec.Address,
			Weight: spec.Weight,
			Status: service.Available,
		}
		backends = append(backends, backend)
	}
	return backends, nil
}
