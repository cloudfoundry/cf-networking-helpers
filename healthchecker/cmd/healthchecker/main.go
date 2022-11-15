package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.cloudfoundry.org/cf-networking-helpers/healthchecker/config"
	"code.cloudfoundry.org/cf-networking-helpers/healthchecker/watchdog"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerflags"
)

const (
	SIGNAL_BUFFER_SIZE = 1024
)

func main() {
	var configFile string
	var c config.Config
	flag.StringVar(&configFile, "c", "", "Configuration File")
	flag.Parse()

	if configFile != "" {
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			panic(fmt.Sprintf("Could not read config file: %s, err: %s", configFile, err.Error()))
		}
		err = json.Unmarshal(b, &c)
		if err != nil {
			panic(fmt.Sprintf("Could not unmarshal config file: %s, err: %s", configFile, err.Error()))
		}
	}

	if c.ComponentName == "" {
		panic(fmt.Sprintf("Invalid component_name in config: %s", configFile))
	}
	logger, _ := lagerflags.NewFromConfig(c.ComponentName, c.LagerConfig)

	startupDelay := c.StartResponseDelayInterval + c.StartupDelayBuffer
	logger.Debug("Sleeping before gorouter responds to /health endpoint on startup", lager.Data{"sleep_time_seconds": startupDelay.Seconds()})
	time.Sleep(startupDelay)

	logger.Info("Starting")

	u := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", c.HealthCheckEndpoint.Host, c.HealthCheckEndpoint.Port),
		User:   url.UserPassword(c.HealthCheckEndpoint.User, c.HealthCheckEndpoint.Password),
		Path:   c.HealthCheckEndpoint.Path,
	}
	host := u.String()

	w := watchdog.NewWatchdog(host, c.ComponentName, c.HealthCheckPollInterval, c.HealthCheckTimeout, logger)
	signals := make(chan os.Signal, SIGNAL_BUFFER_SIZE)
	signal.Notify(signals, syscall.SIGUSR1)

	err := w.WatchHealthcheckEndpoint(context.Background(), signals)
	if err != nil {
		logger.Fatal("Error running healthcheck:", err)
	}
}
