package watchdog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	"code.cloudfoundry.org/lager"
)

const (
	numRetries = 3
)

type Watchdog struct {
	url           *url.URL
	componentName string
	pollInterval  time.Duration
	client        http.Client
	logger        lager.Logger
}

func NewWatchdog(u *url.URL, componentName string, pollInterval time.Duration, healthcheckTimeout time.Duration, logger lager.Logger) *Watchdog {
	client := http.Client{
		Timeout: healthcheckTimeout,
	}
	if strings.HasPrefix(u.Host, "unix") {
		socket := strings.TrimPrefix(u.Host, "unix")
		client.Transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socket)
			},
		}
	}
	return &Watchdog{
		url:           u,
		componentName: componentName,
		pollInterval:  pollInterval,
		client:        client,
		logger:        logger,
	}
}

func (w *Watchdog) WatchHealthcheckEndpoint(ctx context.Context, signals <-chan os.Signal) error {
	pollTimer := time.NewTimer(w.pollInterval)
	errCounter := 0
	defer pollTimer.Stop()
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Context done, exiting")
			return nil
		case sig := <-signals:
			if sig == syscall.SIGUSR1 {
				w.logger.Info("Received USR1 signal, exiting")
				return nil
			}
		case <-pollTimer.C:
			w.logger.Debug("Verifying endpoint", lager.Data{"component": w.componentName, "poll-interval": w.pollInterval})
			err := w.HitHealthcheckEndpoint()
			if err != nil {
				errCounter += 1
				if errCounter >= numRetries {
					select {
					case sig := <-signals:
						if sig == syscall.SIGUSR1 {
							w.logger.Info("Received USR1 signal, exiting")
							return nil
						}
					default:
						return err
					}
				} else {
					w.logger.Debug("Received error", lager.Data{"error": err.Error(), "attempt": errCounter})
				}
			} else {
				errCounter = 0
			}
			pollTimer.Reset(w.pollInterval)
		}
	}
}

func (w *Watchdog) HitHealthcheckEndpoint() error {
	req, err := http.NewRequest("GET", w.url.String(), nil)
	if err != nil {
		return err
	}
	if req.URL.Host == "" {
		req.URL.Host = w.url.Host
	}

	response, err := w.client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf(
			"%v received from healthcheck endpoint (200 expected)",
			response.StatusCode))
	}
	return nil
}
