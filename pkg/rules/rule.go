package rules

import (
	"github.com/fiust/highway/pkg/service"
	"net/http"
)

type Rule struct {
	Service     service.Service
	Schema      string
	PathPrefix  string
	Host        string
	Methods     []string
	Headers     map[string]string
	Queries     map[string]string
	Middlewares []func(http.Handler) http.Handler
}
