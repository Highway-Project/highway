package config

type MiddlewareSpec struct {
	MiddlewareName   string            `mapstructure:"middlewareName"`
	RefName          string            `mapstructure:"refName"`
	MiddlewarePath   string            `mapstructure:"middlewarePath"`
	CustomMiddleware bool              `mapstructure:"customMiddleware"`
	Params           map[string]string `mapstructure:"params"`
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

	// TODO Validate PluginPath

	// TODO Check plugin existence

	if isValid {
		return nil
	}
	return validationError
}
