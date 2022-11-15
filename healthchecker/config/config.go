package config

import (
	"time"

	"code.cloudfoundry.org/lager/lagerflags"
)

type HealthCheckEndpoint struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	ComponentName              string              `json:"component_name"`
	HealthCheckEndpoint        HealthCheckEndpoint `json:"healthcheck_endpoint"`
	HealthCheckPollInterval    time.Duration       `json:"healthcheck_poll_interval",default:"10s"`
	HealthCheckTimeout         time.Duration       `json:"healthcheck_timeout",default:"5s"`
	StartResponseDelayInterval time.Duration       `json:"start_response_delay_interval,omitempty",default:"5s"`
	StartupDelayBuffer         time.Duration       `json:"startup_delay_buffer",default:"5s"`
	lagerflags.LagerConfig
}
