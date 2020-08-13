package config

import "github.com/asaskevich/govalidator"

type ServiceSpec struct {
	Name             string        `mapstructure:"name"`
	LoadBalancerName string        `default:"round-robin" mapstructure:"loadbalancer"`
	BackendsSpecs    []BackendSpec `mapstructure:"backends"`
}

type BackendSpec struct {
	BackendName string `mapstructure:"name"`
	Weight      int8   `default:"1" mapstructure:"weight"`
	Address     string `mapstructure:"address"`
}

func (s ServiceSpec) Validate() error {
	validationError := ServiceValidationError{}
	isValid := true

	//Validate name
	if s.Name == "" {
		validationError.NameError = true
		validationError.NameErrorMessage = "name field for service is required"
		isValid = false
	}

	//Validate Load Balancer
	if s.LoadBalancerName == "" {
		validationError.LoadBalancerError = true
		validationError.LoadBalancerErrorMessage = "loadbalancer field for service is required"
		isValid = false
	}

	//TODO Validate LoadBalancer algorithm existence in app

	//Validate Backends
	if len(s.BackendsSpecs) == 0 {
		validationError.BackendError = true
		validationError.LoadBalancerErrorMessage = "each service should have at least one backend"
		isValid = false
	} else {
		backendMessage := ""
		for _, backend := range s.BackendsSpecs {
			err := backend.Validate()
			if err != nil {
				backendMessage += err.Error()
			}
		}
		if backendMessage != "" {
			validationError.LoadBalancerError = true
			validationError.BackendErrorMessage = backendMessage
			isValid = false
		}
	}

	if isValid {
		return nil
	}

	return validationError
}

func (b BackendSpec) Validate() error {
	validationError := BackendValidationError{}
	isValid := true

	// Validate Name
	if b.BackendName == "" {
		validationError.NameError = true
		validationError.NameErrorMessage = "name field for backend is required"
		isValid = false
	}

	// Validate Weight
	if b.Weight <= 0 {
		validationError.WeightError = true
		validationError.WeightErrorMessage = "weight field for backend should be above zero"
		isValid = false
	}

	// Validate Address
	if !govalidator.IsURL(b.Address) {
		validationError.AddressError = true
		validationError.AddressErrorMessage = "address field for backend should be in URL format"
		isValid = false
	}

	if isValid {
		return nil
	}
	return validationError
}
