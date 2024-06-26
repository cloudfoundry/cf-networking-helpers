package poller_test

import (
	"errors"
	"os"
	"sync/atomic"
	"time"

	"code.cloudfoundry.org/cf-networking-helpers/poller"
	"code.cloudfoundry.org/lager/v3/lagertest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Poller", func() {
	Describe("Run", func() {
		var (
			logger  *lagertest.TestLogger
			p       *poller.Poller
			signals chan os.Signal
			ready   chan struct{}

			cycleCount uint64
			retChan    chan error
		)

		BeforeEach(func() {
			signals = make(chan os.Signal)
			ready = make(chan struct{})

			cycleCount = 0
			retChan = make(chan error)

			logger = lagertest.NewTestLogger("test")

			p = &poller.Poller{
				Logger:       logger,
				PollInterval: 1 * time.Second,

				SingleCycleFunc: func() error {
					atomic.AddUint64(&cycleCount, 1)
					return nil
				},
			}
		})

		Context("when running", func() {

			Context("when RunBeforeFirstInterval is set", func() {
				BeforeEach(func() {
					p.RunBeforeFirstInterval = true
				})
				It("calls the single cycle func on start", func() {
					go func() {
						retChan <- p.Run(signals, ready)
					}()

					Eventually(ready).Should(BeClosed())
					Expect(atomic.LoadUint64(&cycleCount)).To(Equal(uint64(1)))

					signals <- os.Interrupt
					Eventually(retChan).Should(Receive(BeNil()))
				})
			})
			Context("when RunBeforeFirstInterval is not set", func() {
				BeforeEach(func() {
					p.PollInterval = 2 * time.Second
				})
				It("does not call the single cycle func on start", func() {
					go func() {
						retChan <- p.Run(signals, ready)
					}()

					Eventually(ready).Should(BeClosed())
					Expect(atomic.LoadUint64(&cycleCount)).To(Equal(uint64(0)))

					signals <- os.Interrupt
					Eventually(retChan).Should(Receive())
				})
			})

			It("calls the single cycle func after the poll interval", func() {
				go func() {
					retChan <- p.Run(signals, ready)
				}()

				Eventually(ready).Should(BeClosed())
				Eventually(func() uint64 {
					return atomic.LoadUint64(&cycleCount)
				}, 3*time.Second).Should(BeNumerically(">", 1))

				Consistently(retChan).ShouldNot(Receive())

				signals <- os.Interrupt
				Eventually(retChan).Should(Receive(BeNil()))
			})

		})

		Context("when the cycle func fails with a non-fatal error", func() {
			BeforeEach(func() {
				p.SingleCycleFunc = func() error { return errors.New("banana") }
			})

			It("logs the error and continues", func() {
				go func() {
					retChan <- p.Run(signals, ready)
				}()

				Eventually(ready).Should(BeClosed())

				Eventually(logger, 30*time.Second).Should(gbytes.Say("poll-cycle.*banana"))

				Consistently(retChan).ShouldNot(Receive())

				signals <- os.Interrupt
				Eventually(retChan).Should(Receive(BeNil()))
			})
		})

		Context("when the cycle func fails with a fatal error", func() {
			BeforeEach(func() {
				p.SingleCycleFunc = func() error {
					return poller.FatalError("banana")
				}
			})

			It("logs the error and exits", func() {
				go func() {
					retChan <- p.Run(signals, ready)
				}()

				Eventually(ready).Should(BeClosed())
				Eventually(logger, 30*time.Second).Should(gbytes.Say("poll-cycle.*banana"))
				Eventually(retChan).Should(Receive(MatchError(
					"This cell must be restarted (run \"bosh restart <job>\"): fatal: banana",
				)))
			})
		})
	})
})
