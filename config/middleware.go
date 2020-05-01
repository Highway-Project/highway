package config

import "github.com/asaskevich/govalidator"

type MiddlewareSpec struct {
	MiddlewareName string            `mapstructure:"name"`
	PluginPath     string            `mapstructure:"plugin-path"`
	Params         map[string]string `mapstructure:"params"`
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

	// Validate PluginPath
	if  ok, _ := govalidator.IsFilePath(m.PluginPath); m.PluginPath != "" && !ok {
		validationError.PluginPathError = true
		validationError.PluginPathErrorMessage = "plugin-path field for middleware should a valid file path"
		isValid = false
	}

	// TODO Check plugin existence

	if isValid {
		return nil
	}
	return validationError
}
