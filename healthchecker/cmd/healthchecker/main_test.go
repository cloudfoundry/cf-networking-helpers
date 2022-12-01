package main_test

import (
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"code.cloudfoundry.org/cf-networking-helpers/healthchecker/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
	"gopkg.in/yaml.v2"
)

var _ = Describe("HealthChecker", func() {
	var (
		cfg        config.Config
		configFile *os.File
		binPath    string
		session    *gexec.Session
	)

	BeforeEach(func() {
		failureCounterFile, err := os.CreateTemp("", "ginkgoWatchdogFailureCountFile.*")
		Expect(err).NotTo(HaveOccurred())

		cfg = config.Config{
			ComponentName:              "healthchecker",
			FailureCounterFile:         failureCounterFile.Name(),
			LogLevel:                   "info",
			StartupDelayBuffer:         1 * time.Millisecond,
			StartResponseDelayInterval: 1 * time.Millisecond,
			HealthCheckPollInterval:    1 * time.Millisecond,
			HealthCheckTimeout:         1 * time.Millisecond,
		}
		binPath, err = gexec.Build("code.cloudfoundry.org/cf-networking-helpers/healthchecker/cmd/healthchecker", "-race")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		var err error
		configFile, err = ioutil.TempFile("", "healthchecker.config")
		Expect(err).NotTo(HaveOccurred())

		cfgBytes, err := yaml.Marshal(cfg)
		Expect(err).NotTo(HaveOccurred())

		_, err = configFile.Write(cfgBytes)
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command(binPath, "-c", configFile.Name())
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(configFile.Name())
		os.RemoveAll(binPath)
	})

	Context("when there is no component name in config", func() {
		BeforeEach(func() {
			cfg = config.Config{}
		})

		It("fails with error", func() {
			Eventually(session).Should(gexec.Exit(2))
			Expect(session.Err).To(gbytes.Say("Missing component_name"))
		})
	})

	Context("when there is no server running", func() {
		BeforeEach(func() {
			cfg.HealthCheckEndpoint.Host = "invalid-host"
			cfg.HealthCheckEndpoint.Port = 4444
		})

		It("fails", func() {
			Eventually(session, 10*time.Second).Should(gexec.Exit(2))
			Expect(session.Out).To(gbytes.Say("Error running healthcheck"))
		})
	})

	Context("when there is a tcp server running", func() {
		var server *ghttp.Server
		BeforeEach(func() {
			server = ghttp.NewServer()
			server.RouteToHandler(
				"GET", "/some-path",
				ghttp.RespondWith(200, "ok"),
			)
			u, err := url.Parse(server.URL())
			Expect(err).NotTo(HaveOccurred())

			cfg.HealthCheckEndpoint.Host = u.Hostname()
			port, err := strconv.Atoi(u.Port())
			Expect(err).NotTo(HaveOccurred())
			cfg.HealthCheckEndpoint.Port = port
			cfg.LogLevel = "debug"
			cfg.HealthCheckEndpoint.Path = "/some-path"
			cfg.StartupDelayBuffer = 5 * time.Second
			cfg.HealthCheckPollInterval = 500 * time.Millisecond
			cfg.HealthCheckTimeout = 5 * time.Second
		})

		AfterEach(func() {
			server.Close()
		})

		It("works", func() {
			Eventually(session.Out, 10*time.Second).Should(gbytes.Say("Verifying endpoint"))
			Eventually(func() int { return len(server.ReceivedRequests()) }, 10*time.Second).Should(BeNumerically(">", 0))
		})
	})
	Context("when there is a unix socket server running", func() {
		var server *ghttp.Server
		BeforeEach(func() {
			unixSocket, err := os.CreateTemp("", "ghttpUnixSocket.*")
			Expect(err).NotTo(HaveOccurred())
			err = os.Remove(unixSocket.Name())
			Expect(err).NotTo(HaveOccurred())

			unixListener, err := net.Listen("unix", unixSocket.Name())
			Expect(err).NotTo(HaveOccurred())

			server = ghttp.NewUnstartedServer()
			server.HTTPTestServer.Listener = unixListener
			server.RouteToHandler(
				"GET", "/some-path",
				ghttp.RespondWith(200, "ok"),
			)
			server.Start()

			cfg.HealthCheckEndpoint.Socket = unixSocket.Name()
			cfg.LogLevel = "debug"
			cfg.HealthCheckEndpoint.Path = "/some-path"
			cfg.StartupDelayBuffer = 5 * time.Second
			cfg.HealthCheckPollInterval = 500 * time.Millisecond
			cfg.HealthCheckTimeout = 5 * time.Second
		})

		AfterEach(func() {
			server.Close()
		})

		It("works", func() {
			Eventually(session.Out, 10*time.Second).Should(gbytes.Say("Verifying endpoint"))
			Eventually(func() int { return len(server.ReceivedRequests()) }, 10*time.Second).Should(BeNumerically(">", 0))
		})
	})
})
