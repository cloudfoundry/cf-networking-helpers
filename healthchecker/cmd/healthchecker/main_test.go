package main_test

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"code.cloudfoundry.org/cf-networking-helpers/healthchecker/config"
	"code.cloudfoundry.org/lager/lagerflags"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("HealthChecker", func() {
	var (
		cfg        config.Config
		configFile *os.File
		binPath    string
		session    *gexec.Session
	)

	BeforeEach(func() {
		cfg = config.Config{
			ComponentName: "healthchecker",
			LagerConfig: lagerflags.LagerConfig{
				LogLevel: "info",
			},
			StartupDelayBuffer:      1 * time.Millisecond,
			HealthCheckPollInterval: 1 * time.Millisecond,
			HealthCheckTimeout:      1 * time.Millisecond,
		}
		var err error
		binPath, err = gexec.Build("code.cloudfoundry.org/cf-networking-helpers/healthchecker/cmd/healthchecker", "-race")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		var err error
		configFile, err = ioutil.TempFile("", "healthchecker.config")
		Expect(err).NotTo(HaveOccurred())

		err = json.NewEncoder(configFile).Encode(cfg)
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
			Expect(session.Err).To(gbytes.Say("Invalid component_name"))
		})
	})

	Context("when there is no server running", func() {
		It("fails", func() {
			Eventually(session, 10*time.Second).Should(gexec.Exit(2))
			Expect(session.Out).To(gbytes.Say("Error running healthcheck"))
		})
	})

	Context("when there is a server running", func() {
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
			Expect(len(server.ReceivedRequests())).To(BeNumerically(">", 0))
		})
	})
})
