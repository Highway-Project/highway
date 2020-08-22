package rules

import (
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/Highway-Project/highway/pkg/service"
	"net/http"
)

type Rule struct {
	Service     *service.Service
	Schema      string
	PathPrefix  string
	Hosts       []string
	Methods     []string
	Headers     map[string]string
	Queries     map[string]string
	Middlewares []middlewares.Middleware
	handler     http.HandlerFunc
}

func (rule *Rule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rule.handler.ServeHTTP(w, r)
}

func NewRule(srv *service.Service, schema string, pathPrefix string, hosts []string, methods []string,
	headers map[string]string, queries map[string]string, middlewareList []middlewares.Middleware) (*Rule, error) {
	for i, j := 0, len(middlewareList)-1; i < j; i, j = i+1, j-1 {
		middlewareList[i], middlewareList[j] = middlewareList[j], middlewareList[i]
	}

	rule := Rule{
		Service:     srv,
		Schema:      schema,
		PathPrefix:  pathPrefix,
		Hosts:       hosts,
		Methods:     methods,
		Headers:     headers,
		Queries:     queries,
		Middlewares: middlewareList,
	}

	handler := rule.Service.ServeHTTP
	for _, middleware := range middlewareList {
		handler = middleware.Process(handler)
	}
	rule.handler = handler

	return &rule, nil
}
