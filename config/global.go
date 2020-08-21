package config

import "time"

type GlobalConfig struct {
	Port              uint          `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
	MaxHeaderBytes    int           `mapstructure:"maxHeaderBytes"`
}

func (c GlobalConfig) Validate() error {
	validationError := GlobalValidationError{}
	isValid := true

	if c.ReadTimeout < 0 {
		validationError.ReadTimeoutError = true
		validationError.ReadTimeoutErrorMessage = "ReadTimeout should be positive"
		isValid = false
	}

	if c.ReadHeaderTimeout < 0 {
		validationError.ReadHeaderTimeoutError = true
		validationError.ReadHeaderTimeoutErrorMessage = "ReadHeaderTimeout should be positive"
		isValid = false
	}

	if c.WriteTimeout < 0 {
		validationError.WriteTimeoutError = true
		validationError.WriteTimeoutErrorMessage = "WriteTimeout should be positive"
		isValid = false
	}

	if c.IdleTimeout < 0 {
		validationError.IdleTimeoutError = true
		validationError.IdleTimeoutErrorMessage = "IdleTimeout should be positive"
		isValid = false
	}

	if c.MaxHeaderBytes < 0 {
		validationError.MaxHeaderBytesError = true
		validationError.MaxHeaderBytesErrorMessage = "MaxHeaderBytes should be positive"
		isValid = false
	}

	if isValid {
		return nil
	}

	return validationError
}
