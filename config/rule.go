package config

import (
	"github.com/asaskevich/govalidator"
	"regexp"
)

type RuleSpec struct {
	ServiceName     string              `mapstructure:"service"`
	Schema          string              `default:"http" mapstructure:"schema"`
	PathPrefix      string              `default:"/" mapstructure:"pathPrefix"`
	Hosts           []string            `mapstructure:"hosts"`
	Methods         []string            `mapstructure:"methods"`
	Headers         map[string][]string `mapstructure:"headers"`
	Queries         map[string]string   `mapstructure:"queries"`
	MiddlewareNames []string            `mapstructure:"middlewares"`
}

func (r RuleSpec) Validate() error {
	validationError := RuleValidationError{}
	isValid := true

	// Validate Name
	if r.ServiceName == "" {
		validationError.NameError = true
		validationError.NameErrorMessage = "service field for rule is required"
		isValid = false
	}

	// Validate Schema
	if r.Schema == "" || (r.Schema != "http" && r.Schema != "https") {
		validationError.SchemaError = true
		validationError.SchemaErrorMessage = "schema field for rule should be http or https"
		isValid = false
	}

	// Validate PathPrefix
	if _, err := regexp.MatchString(govalidator.URLPath, r.PathPrefix); err != nil {
		validationError.PathPrefixError = true
		validationError.PathPrefixErrorMessage = err.Error()
		isValid = false
	}

	// Validate Hosts
	for _, host := range r.Hosts {
		if !govalidator.IsHost(host) {
			validationError.HostsError = true
			validationError.HostsErrorMessage = "host fields are invalid"
			isValid = false
			break
		}
	}

	if isValid {
		return nil
	}
	return validationError
}
