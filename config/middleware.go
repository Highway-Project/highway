package config

import (
	"fmt"
	"os"
)

type MiddlewareSpec struct {
	MiddlewareName   string                 `mapstructure:"middlewareName"`
	RefName          string                 `mapstructure:"refName"`
	MiddlewarePath   string                 `mapstructure:"middlewarePath"`
	CustomMiddleware bool                   `mapstructure:"customMiddleware"`
	Params           map[string]interface{} `mapstructure:"params"`
}

func (m MiddlewareSpec) Validate() error {
	validationError := MiddlewareValidationError{}
	isValid := true

	// Validate Name
	if m.MiddlewareName == "" {
		validationError.NameError = true
		validationError.NameErrorMessage = "name field for middleware is required"
		isValid = false
	}

	// Validate plugin path
	if m.CustomMiddleware {
		_, err := os.Stat(m.MiddlewarePath)
		if err != nil {
			validationError.PluginPathError = true
			validationError.PluginPathErrorMessage = fmt.Sprintf("%s does not exists", m.MiddlewarePath)
			isValid = false
		}
	}

	if isValid {
		return nil
	}
	return validationError
}
