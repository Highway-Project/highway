package main

import (
	"github.com/Highway-Project/highway/config"
	_ "github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/internal/server"
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/Highway-Project/highway/pkg/router"
	"github.com/Highway-Project/highway/pkg/router/gorilla"
	"github.com/Highway-Project/highway/pkg/rules"
	"github.com/Highway-Project/highway/pkg/service"
	"github.com/Highway-Project/highway/pkg/service/random"
	"github.com/creasty/defaults"
	"log"
)

func main() {
	initializeRouters()
	initializeLoadBalancer()

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

	r, err := createRouter(conf.RouterSpec)
	if err != nil {
		panic("router init " + err.Error())
	}

	ruleList, err := createRules(conf.RulesSpecs, conf)
	if err != nil {
		panic("ruleList init " + err.Error())
	}

	s, err := server.New(r, ruleList)
	log.Fatal(s.Run())

	//r := gorilla.New()
	//
	//r1 := ruleList.Rule{
	//	Service:     service.Service{},
	//	Schema:      "http",
	//	PathPrefix:  "/hi",
	//	Host:        []string{"localhost"},
	//	Methods:     []string{"GET", "POST"},
	//	Headers:     nil,
	//	Queries:     nil,
	//	Middlewares: nil,
	//}
	//s := server.Server{
	//	Router: r,
	//	Rules:  nil,
	//}
	//s.Router.AddRule(r1)
	//http.ListenAndServe(":8080", r)
}

func initializeRouters() {
	gorilla.Register()
}

func initializeLoadBalancer() {
	random.Register()
}

func createRouter(spec config.RouterSpec) (router.Router, error) {
	options := router.RouterOptions{Options: spec.RouterOpts}
	r, err := router.NewRouter(spec.Name, options)
	if err != nil {
		return nil, err
	}
	return r, err
}

func createRules(specs []config.RuleSpec, conf *config.Config) ([]rules.Rule, error) {
	result := make([]rules.Rule, len(specs))
	for i, ruleSpec := range specs {
		serviceSpec, err := conf.GetServiceSpec(ruleSpec.ServiceName)
		if err != nil {
			return nil, err
		}

		s, err := createService(serviceSpec)

		rule := rules.Rule{
			Service:     s,
			Schema:      ruleSpec.Schema,
			PathPrefix:  ruleSpec.PathPrefix,
			Hosts:       ruleSpec.Hosts,
			Methods:     ruleSpec.Methods,
			Headers:     ruleSpec.Headers,
			Queries:     ruleSpec.Queries,
			Middlewares: make([]middlewares.Middleware, 0),
		}

		result[i] = rule

	}
	return result, nil
}

func createService(spec *config.ServiceSpec) (service.Service, error) {
	backends := make([]service.Backend, len(spec.BackendsSpecs))

	for i, backendSpec := range spec.BackendsSpecs {
		backend, err := createBackend(backendSpec)
		if err != nil {
			return service.Service{}, err
		}
		backends[i] = backend
	}

	lb, err := service.NewLoadBalancer(spec.LoadBalancerName)
	if err != nil {
		return service.Service{}, err
	}

	s := service.Service{
		Name:     spec.ServiceName,
		Backends: backends,
		LB:       lb,
	}
	return s, nil
}

func createBackend(spec config.BackendSpec) (service.Backend, error) {
	backend := service.Backend{
		Name:   spec.BackendName,
		Addr:   spec.Address,
		Weight: spec.Weight,
	}
	return backend, nil
}
