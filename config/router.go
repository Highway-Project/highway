package config

type RouterSpec struct {
	Name       string            `default:"gorilla" mapstructure:"name"`
	RouterOpts map[string]string `default:"{}" mapstructure:"options"`
}

func (r RouterSpec) Validate() error {
	validationError := RouterValidationError{}
	isValid := true

	if r.Name == "" {
		validationError.NameError = true
		validationError.NameErrorMessage = "name field in router is required"
		isValid = false
	}

	if isValid {
		return nil
	}
	return validationError
}
