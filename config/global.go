package config

import "time"

type GlobalConfig struct {
	Port              string        `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
	MaxHeaderBytes    int           `mapstructure:"maxHeaderBytes"`
}
