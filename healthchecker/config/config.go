package config

import (
	"time"
)

type HealthCheckEndpoint struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Path     string `yaml:"path"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	ComponentName              string              `yaml:"component_name"`
	HealthCheckEndpoint        HealthCheckEndpoint `yaml:"healthcheck_endpoint"`
	HealthCheckPollInterval    time.Duration       `yaml:"healthcheck_poll_interval",default:"10s"`
	HealthCheckTimeout         time.Duration       `yaml:"healthcheck_timeout",default:"5s"`
	StartResponseDelayInterval time.Duration       `yaml:"start_response_delay_interval,omitempty",default:"5s"`
	StartupDelayBuffer         time.Duration       `yaml:"startup_delay_buffer",default:"5s"`
	LogLevel                   string              `yaml:"log_level",default:"info"`
}
