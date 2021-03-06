package provider

import "github.com/Highway-Project/highway/pkg/service"

type Message struct {
	Type         string
	ProviderName string
	Service      *service.Service
}

type ServiceProvider interface {
	Provide() ([]service.Service, error)
	Watch(messageChan chan<- Message) error
}
