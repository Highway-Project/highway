package rules

import (
	"context"
	"errors"
	"fmt"
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/Highway-Project/highway/pkg/service"
	"net/http"
)

type Rule struct {
	Name        string
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
	ctx := r.Context()
	v := ctx.Value("metrics")
	var metrics map[string]string
	if v == nil {
		metrics = make(map[string]string)
	} else {
		metrics = v.(map[string]string)
	}
	metrics["rule"] = rule.Name
	ctx = context.WithValue(ctx, "metrics", metrics)
	r = r.WithContext(ctx)
	rule.handler.ServeHTTP(w, r)
}

func (rule *Rule) SetHandler(handler http.Handler) error {
	if handler == nil {
		msg := fmt.Sprintf("rule %s handler should not be nil", rule.Name)
		return errors.New(msg)
	}
	rule.handler = func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
	return nil
}

func NewRule(srv *service.Service, name string, schema string, pathPrefix string, hosts []string, methods []string,
	headers map[string]string, queries map[string]string, middlewareList []middlewares.Middleware) (*Rule, error) {
	for i, j := 0, len(middlewareList)-1; i < j; i, j = i+1, j-1 {
		middlewareList[i], middlewareList[j] = middlewareList[j], middlewareList[i]
	}

	rule := Rule{
		Name:        name,
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
