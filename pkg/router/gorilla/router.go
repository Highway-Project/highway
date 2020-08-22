package gorilla

import (
	"github.com/Highway-Project/highway/pkg/router"
	"github.com/Highway-Project/highway/pkg/rules"
	"github.com/gorilla/mux"
	"net/http"
)

type GorillaRouter struct {
	router *mux.Router
}

func (r *GorillaRouter) AddRule(rule rules.Rule) error {
	route := r.router.Schemes(rule.Schema).PathPrefix(rule.PathPrefix)

	if rule.Hosts != nil {
		for _, host := range rule.Hosts {
			route.Host(host)
		}
	}

	if rule.Methods != nil {
		route.Methods(rule.Methods...)
	}

	if rule.Headers != nil {
		for k, v := range rule.Headers {
			route.Headers(k, v)
		}
	}

	if rule.Queries != nil {
		for k, v := range rule.Queries {
			route.Queries(k, v)
		}
	}

	route.Handler(&rule)

	return nil
}

func (r *GorillaRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func New(options router.RouterOptions) (router.Router, error) {
	return &GorillaRouter{
		router: mux.NewRouter(),
	}, nil
}
