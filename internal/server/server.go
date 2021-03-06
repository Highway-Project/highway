package server

import (
	"fmt"
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/internal/middlewares"
	"github.com/Highway-Project/highway/internal/router"
	"github.com/Highway-Project/highway/internal/rule"
	"github.com/Highway-Project/highway/internal/service"
	"github.com/Highway-Project/highway/logging"
	pkgRouter "github.com/Highway-Project/highway/pkg/router"
	"github.com/Highway-Project/highway/pkg/rules"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func NewServer(global config.GlobalConfig, routerSpec config.RouterSpec, servicesSpec []config.ServiceSpec, rulesSpec []config.RuleSpec, middlewaresSpec []config.MiddlewareSpec) (*http.Server, error) {
	r, err := router.NewRouter(routerSpec.Name, pkgRouter.RouterOptions{
		Options: routerSpec.RouterOpts,
	})
	if err != nil {
		logging.Logger.WithError(err).Fatal("could not create router")
	}

	metricRule := rules.Rule{
		Name:       "metrics",
		Schema:     "http",
		Methods:    []string{"GET"},
		PathPrefix: "/metrics",
	}
	err = metricRule.SetHandler(promhttp.Handler())
	if err != nil {
		logging.Logger.WithError(err).Fatal("could not set prometheus metrics handler")
	}
	err = r.AddRule(metricRule)
	if err != nil {
		logging.Logger.WithError(err).Fatal("could not set prometheus metrics rule")
	}

	for _, spec := range servicesSpec {
		_, err := service.NewService(spec)
		if err != nil {
			logging.Logger.WithError(err).Errorf("could not create service %s", spec.Name)
		}
	}

	err = middlewares.LoadMiddlewares(middlewaresSpec)
	if err != nil {
		logging.Logger.WithError(err).Fatal("could not load middlwares")
	}

	rules, err := rule.NewRules(rulesSpec)
	for _, ruleObj := range rules {
		err := r.AddRule(ruleObj)
		if err != nil {
			logging.Logger.WithError(err).Errorf("could not create rule for service %s", ruleObj.Service.Name)
		}
	}

	s := http.Server{
		Addr:              fmt.Sprintf(":%d", global.Port),
		Handler:           r,
		ReadTimeout:       global.ReadTimeout * time.Millisecond,
		ReadHeaderTimeout: global.ReadHeaderTimeout * time.Millisecond,
		WriteTimeout:      global.WriteTimeout * time.Millisecond,
		IdleTimeout:       global.IdleTimeout * time.Millisecond,
		MaxHeaderBytes:    global.MaxHeaderBytes,
	}

	return &s, nil
}
