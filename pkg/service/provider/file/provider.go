package file

import (
	"github.com/Highway-Project/highway/pkg/service"
	"github.com/Highway-Project/highway/pkg/service/provider"
)

type FileProvider struct {
	FilePath string
}

func (f FileProvider) Provide() ([]service.Service, error) {
	return nil, nil
}

func (f FileProvider) Watch(messageChan chan<- provider.Message) error {
	return nil
}

func New(filePath string) (provider.ServiceProvider, error) {
	return FileProvider{
		FilePath: filePath,
	}, nil
}
