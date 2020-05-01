package rules

import (
	"github.com/fiust/highway/pkg/middlewares"
	"github.com/fiust/highway/pkg/service"
	"net/http"
)

type Rule struct {
	Service     service.Service
	Schema      string
	PathPrefix  string
	Hosts       []string
	Methods     []string
	Headers     map[string][]string
	Queries     map[string]string
	Middlewares []middlewares.Middleware
}

func (rule Rule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rule.Service.ServeHTTP(w, r)
}
