package service

type LoadBalancer interface {
	Balance([]Backend) Backend
}
