package gorilla

import (
	"github.com/fiust/highway/pkg/router"
	"github.com/fiust/highway/pkg/rules"
	"github.com/gorilla/mux"
	"net/http"
)

func Register() {
	err := router.RegisterRouter("gorilla", New)
	if err != nil {
		panic("failed registering gorilla router")
	}
}

type MuxRouter struct {
	router *mux.Router
}

func (r *MuxRouter) AddRule(rule rules.Rule) error {
	r.router.Handle(rule.PathPrefix, &rule.Service).
		//Host(rule.Hosts[0]).
		Schemes(rule.Schema).
		Methods(rule.Methods...).Handler(&rule.Service)
	return nil
}

func (r *MuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func New(options router.RouterOptions) (router.Router, error) {
	var router router.Router
	router = &MuxRouter{
		router:mux.NewRouter(),
	}
	return router, nil
}