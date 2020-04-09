package rules

import (
	"github.com/fiust/highway/pkg/service"
	"net/http"
)

type Rule struct {
	service     service.Service
	schema      string
	pathPrefix  string
	host        string
	methods     []string
	headers     map[string]string
	queries     map[string]string
	middlewares []func(http.Handler) http.Handler
}
