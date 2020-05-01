package main

import (
	"github.com/creasty/defaults"
	"github.com/fiust/highway/config"
	_ "github.com/fiust/highway/config"
	"github.com/fiust/highway/internal/server"
	"github.com/fiust/highway/pkg/router/gorilla"
	"github.com/fiust/highway/pkg/rules"
	"github.com/fiust/highway/pkg/service"
	"net/http"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		panic("conf panic " + err.Error())
	}
	err = defaults.Set(conf)
	if err != nil {
		panic("default panic " + err.Error())
	}

	err = conf.Validate()
	if err != nil {
		panic("validate panic " + err.Error())
	}

	r := gorilla.New()

	r1 := rules.Rule{
		Service:     service.Service{},
		Schema:      "https",
		PathPrefix:  "/hi",
		Host:        []string{"localhost"},
		Methods:     []string{"GET", "POST"},
		Headers:     nil,
		Queries:     nil,
		Middlewares: nil,
	}
	s := server.Server{
		Router: r,
		Rules:  nil,
	}
	s.Router.AddRule(r1)
	http.ListenAndServe(":8080", r)
}
