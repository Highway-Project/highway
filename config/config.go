package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Global           GlobalConfig              `mapstructure:"global"`
	RouterSpec       RouterSpec                `mapstructure:"router"`
	ServicesSpecs    []ServiceSpec             `mapstructure:"services"`
	RulesSpecs       []RuleSpec                `mapstructure:"rules"`
	MiddlewaresSpecs []MiddlewareSpec          `mapstructure:"middlewares"`
	services         map[string]ServiceSpec    `mapstructure:"-"`
	middlewares      map[string]MiddlewareSpec `mapstructure:"-"`
}

func (c Config) GetServiceSpec(name string) (*ServiceSpec, error) {
	serviceSpec, ok := c.services[name]
	if !ok {
		return nil, ServiceNotFoundError{ServiceName: name}
	}
	return &serviceSpec, nil
}

func (c Config) GetMiddlewareSpec(name string) (*MiddlewareSpec, error) {
	middlewareSpec, ok := c.middlewares[name]
	if !ok {
		return nil, MiddlewareNotFoundError{MiddlewareName: name}
	}
	return &middlewareSpec, nil
}

func (c Config) Validate() error {
	routerErr := c.RouterSpec.Validate()

	if routerErr != nil {
		return routerErr
	}

	var serviceErr error
	for _, service := range c.ServicesSpecs {
		err := service.Validate()
		if err != nil {
			serviceErr = err
			break
		}
	}

	if serviceErr != nil {
		return serviceErr
	}

	var ruleErr error
	for _, rule := range c.RulesSpecs {
		err := rule.Validate()
		if err != nil {
			ruleErr = err
			break
		}
	}

	if ruleErr != nil {
		return ruleErr
	}

	var middlewareErr error
	for _, middleware := range c.MiddlewaresSpecs {
		err := middleware.Validate()
		if err != nil {
			middlewareErr = err
			break
		}
	}

	if middlewareErr != nil {
		return middlewareErr
	}

	return nil
}

func ReadConfig() (*Config, error) {

	//TODO Read configs via providers
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	config.services = make(map[string]ServiceSpec)
	config.middlewares = make(map[string]MiddlewareSpec)

	//TODO Write middleware loading logic

	for _, s := range config.ServicesSpecs {
		if s.Name != "" {
			config.services[s.Name] = s
		}
	}

	return &config, nil
}
