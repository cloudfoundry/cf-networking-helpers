package lagerlevel_test

import (
	"bufio"
	"net"

	. "code.cloudfoundry.org/cf-networking-helpers/lagerlevel"

	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"code.cloudfoundry.org/cf-networking-helpers/testsupport/ports"
	"code.cloudfoundry.org/lager/v3"
	"code.cloudfoundry.org/lager/v3/lagertest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
)

var _ = Describe("Server", func() {
	var (
		serverProc ifrit.Process
		testLogger *lagertest.TestLogger
		sink       *lager.ReconfigurableSink
		server     *Server
		port       int
	)

	BeforeEach(func() {
		testLogger = lagertest.NewTestLogger("test")
		port = ports.PickAPort()
		sink = lager.NewReconfigurableSink(testLogger, lager.ERROR)

		server = NewServer("localhost", port, 200*time.Millisecond, sink, testLogger)
	})

	Context("when it starts", func() {
		BeforeEach(func() {
			serverProc = ifrit.Invoke(server)
		})

		AfterEach(func() {
			serverProc.Signal(os.Interrupt)
			Eventually(serverProc.Wait()).Should(Receive())
		})

		It("sets the debug level", func() {
			response := requestLogChange("debug", port)
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
			Expect(testLogger.LogMessages()).To(ContainElement("test.Log level set to DEBUG"))
			Expect(sink.GetMinLevel()).To(Equal(lager.DEBUG))
		})

		It("sets the info level", func() {
			response := requestLogChange("info", port)
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
			Expect(testLogger.LogMessages()).To(ContainElement("test.Log level set to INFO"))
			Expect(sink.GetMinLevel()).To(Equal(lager.INFO))
		})

		It("rejects other levels", func() {
			response := requestLogChange("error", port)
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(sink.GetMinLevel()).To(Equal(lager.ERROR))
			Expect(testLogger.LogMessages()).To(ContainElement("test.Invalid log level requested: `error`. Skipping."))
		})

		Context("when request header write times out", func() {
			It("closes the request", func() {
				var conn net.Conn
				var err error
				Eventually(func() error {
					conn, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
					return err
				}).Should(BeNil())
				defer conn.Close()

				writer := bufio.NewWriter(conn)

				fmt.Fprintf(writer, "POST /log-level HTTP/1.1\r\n")

				// started writing headers
				fmt.Fprintf(writer, "Host: localhost\r\n")
				writer.Flush()

				time.Sleep(300 * time.Millisecond)

				// done
				fmt.Fprintf(writer, "\r\n")
				writer.Flush()

				resp := bufio.NewReader(conn)
				_, err = resp.ReadString('\n')
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when it is unable to start", func() {
		var conflictingServer *http.Server

		BeforeEach(func() {
			conflictingServer = launchConflictingServer(port)
		})

		AfterEach(func() {
			conflictingServer.Close()
			serverProc.Signal(os.Interrupt)
		})

		It("logs and quits", func() {
			serverProc = ifrit.Invoke(server)
			Eventually(serverProc.Wait()).Should(Receive())
			Eventually(testLogger.LogMessages(), 5*time.Second).Should(ContainElement("test.Listen and serve exited with error:"))
		})
	})
})

func launchConflictingServer(port int) *http.Server {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	conflictingServer := &http.Server{Addr: address}
	go func() { conflictingServer.ListenAndServe() }()
	Eventually(func() bool {
		client := &http.Client{}
		resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
		if err != nil {
			return false
		}
		return resp.StatusCode == 404
	}).Should(BeTrue())
	return conflictingServer
}

func requestLogChange(logLevel string, port int) *http.Response {
	client := &http.Client{}
	postBody := strings.NewReader(logLevel)
	response, err := client.Post(fmt.Sprintf("http://localhost:%d/log-level", port), "text/plain", postBody)
	Expect(err).ToNot(HaveOccurred())
	return response
}
