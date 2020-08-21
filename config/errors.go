package config

import "fmt"

type ServiceNotFoundError struct {
	ServiceName string
}

func (s ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service \" %s \" not found", s.ServiceName)
}

type MiddlewareNotFoundError struct {
	MiddlewareName string
}

func (m MiddlewareNotFoundError) Error() string {
	return fmt.Sprintf("middleware \" %s \" not found", m.MiddlewareName)
}

type ServiceValidationError struct {
	NameError                bool
	NameErrorMessage         string
	LoadBalancerError        bool
	LoadBalancerErrorMessage string
	BackendError             bool
	BackendErrorMessage      string
}

func (s ServiceValidationError) Error() string {
	message := "ServiceValidationError - "
	if s.NameError {
		message += fmt.Sprintf("NameError: %s ,", s.NameErrorMessage)
	}

	if s.LoadBalancerError {
		message += fmt.Sprintf("LoadBalancerError: %s ,", s.LoadBalancerErrorMessage)
	}

	if s.BackendError {
		message += fmt.Sprintf("BackendError: %s ,", s.BackendErrorMessage)
	}

	return message
}

type BackendValidationError struct {
	NameError           bool
	NameErrorMessage    string
	WeightError         bool
	WeightErrorMessage  string
	AddressError        bool
	AddressErrorMessage string
}

func (b BackendValidationError) Error() string {
	message := "BackendValidationError - "
	if b.NameError {
		message += fmt.Sprintf("NameError: %s ,", b.NameErrorMessage)
	}

	if b.WeightError {
		message += fmt.Sprintf("LoadBalancerError: %s ,", b.WeightErrorMessage)
	}

	if b.AddressError {
		message += fmt.Sprintf("BackendError: %s ,", b.AddressErrorMessage)
	}

	return message
}

type MiddlewareValidationError struct {
	NameError              bool
	NameErrorMessage       string
	PluginPathError        bool
	PluginPathErrorMessage string
}

func (m MiddlewareValidationError) Error() string {
	message := "MiddlewareValidationError - "
	if m.NameError {
		message += fmt.Sprintf("NameError: %s ,", m.NameErrorMessage)
	}

	if m.PluginPathError {
		message += fmt.Sprintf("PluginPathError: %s ,", m.PluginPathErrorMessage)
	}

	return message
}

type RouterValidationError struct {
	NameError        bool
	NameErrorMessage string
}

func (r RouterValidationError) Error() string {
	return fmt.Sprintf("RouterValidationError - NameError: %s", r.NameErrorMessage)
}

type RuleValidationError struct {
	NameError                   bool
	NameErrorMessage            string
	SchemaError                 bool
	SchemaErrorMessage          string
	PathPrefixError             bool
	PathPrefixErrorMessage      string
	HostsError                  bool
	HostsErrorMessage           string
	MethodsError                bool
	MethodsErrorMessage         string
	HeadersError                bool
	HeadersErrorMessage         string
	QueriesError                bool
	QueriesErrorMessage         string
	MiddlewareNamesError        bool
	MiddlewareNamesErrorMessage string
}

func (r RuleValidationError) Error() string {
	message := "RuleValidationError - "

	if r.NameError {
		message += fmt.Sprintf("NameError: %s ,", r.NameErrorMessage)
	}

	if r.SchemaError {
		message += fmt.Sprintf("SchemaError: %s ,", r.SchemaErrorMessage)
	}

	if r.PathPrefixError {
		message += fmt.Sprintf("PathPrefixError: %s ,", r.PathPrefixErrorMessage)
	}

	if r.HostsError {
		message += fmt.Sprintf("HostsError: %s ,", r.HostsErrorMessage)
	}

	if r.MethodsError {
		message += fmt.Sprintf("MethodsError: %s ,", r.MethodsErrorMessage)
	}

	if r.HeadersError {
		message += fmt.Sprintf("HeadersError: %s ,", r.HeadersErrorMessage)
	}

	if r.QueriesError {
		message += fmt.Sprintf("QueriesError: %s ,", r.QueriesErrorMessage)
	}

	if r.MiddlewareNamesError {
		message += fmt.Sprintf("MiddlewareNamesError: %s ,", r.MiddlewareNamesErrorMessage)
	}

	return message
}

type GlobalValidationError struct {
	PortError                     bool
	PortErrorMessage              string
	ReadTimeoutError              bool
	ReadTimeoutErrorMessage       string
	ReadHeaderTimeoutError        bool
	ReadHeaderTimeoutErrorMessage string
	WriteTimeoutError             bool
	WriteTimeoutErrorMessage      string
	IdleTimeoutError              bool
	IdleTimeoutErrorMessage       string
	MaxHeaderBytesError           bool
	MaxHeaderBytesErrorMessage    string
}

func (g GlobalValidationError) Error() string {
	message := "GlobalValidationError - "
	if g.PortError {
		message += fmt.Sprintf("PortError: %s, ", g.PortErrorMessage)
	}

	if g.ReadTimeoutError {
		message += fmt.Sprintf("ReadTimeoutError: %s, ", g.ReadTimeoutErrorMessage)
	}

	if g.ReadHeaderTimeoutError {
		message += fmt.Sprintf("ReadHeaderTimeoutError: %s, ", g.ReadHeaderTimeoutErrorMessage)
	}

	if g.WriteTimeoutError {
		message += fmt.Sprintf("WriteTimeoutError: %s, ", g.WriteTimeoutErrorMessage)
	}

	if g.IdleTimeoutError {
		message += fmt.Sprintf("IdleTimeoutError: %s, ", g.IdleTimeoutErrorMessage)
	}

	if g.MaxHeaderBytesError {
		message += fmt.Sprintf("MaxHeaderBytesError: %s, ", g.MaxHeaderBytesErrorMessage)
	}

	return message
}
