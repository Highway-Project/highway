package rule

import (
	"fmt"
	"github.com/Highway-Project/highway/config"
	mw "github.com/Highway-Project/highway/internal/middlewares"
	"github.com/Highway-Project/highway/internal/service"
	"github.com/Highway-Project/highway/logging"
	"github.com/Highway-Project/highway/pkg/middlewares"
	"github.com/Highway-Project/highway/pkg/rules"
	"github.com/pkg/errors"
)

func NewRule(spec config.RuleSpec) (rules.Rule, error) {
	s, err := service.GetServiceByName(spec.ServiceName)
	if err != nil {
		msg := fmt.Sprintf("service %s does not exists", spec.ServiceName)
		logging.Logger.WithError(err).Error(msg)
		return rules.Rule{}, errors.Wrap(err, msg)
	}

	middlewareList := make([]middlewares.Middleware, 0)
	for _, refName := range spec.MiddlewareNames {
		middleware, err := mw.GetMiddlewareByName(refName)
		if err != nil {
			logging.Logger.Fatalf("middleware %s does not exist", refName)
		}
		middlewareList = append(middlewareList, middleware)
	}

	r, err := rules.NewRule(s, spec.Schema, spec.PathPrefix, spec.Hosts, spec.Methods, spec.Headers, spec.Queries, middlewareList)
	if err != nil {
		logging.Logger.WithError(err).Errorf("could not create rule for service %s", spec.ServiceName)
		return rules.Rule{}, errors.Wrapf(err, "could not create rule for service %s", spec.ServiceName)
	}

	return *r, nil
}

func NewRules(specs []config.RuleSpec) ([]rules.Rule, error) {
	rs := make([]rules.Rule, 0)
	for _, spec := range specs {
		r, err := NewRule(spec)
		if err != nil {
			logging.Logger.WithError(err).Error("could not create rule for service", spec.ServiceName)
			continue
		}
		rs = append(rs, r)
	}
	return rs, nil
}
