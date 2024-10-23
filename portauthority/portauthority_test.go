package portauthority_test

import (
	"code.cloudfoundry.org/cf-networking-helpers/portauthority"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Portallocator", func() {
	var (
		allocator portauthority.PortAllocator
		port      uint16
		err       error
	)

	BeforeEach(func() {
		allocator, err = portauthority.New(30, 65355)
		Expect(err).NotTo(HaveOccurred())
	})

	It("starts allocating at the beginning of the range, inclusive", func() {
		port, err = allocator.ClaimPorts(1)
		Expect(err).NotTo(HaveOccurred())
		Expect(port).To(Equal(uint16(30)))
	})

	It("returns a different int each time NextPort is called", func() {
		Expect(allocator.ClaimPorts(1)).Should(BeEquivalentTo(30))
		Expect(allocator.ClaimPorts(1)).Should(BeEquivalentTo(31))
		Expect(allocator.ClaimPorts(1)).Should(BeEquivalentTo(32))
	})

	Context("when the allocator runs out of ports in the range", func() {
		BeforeEach(func() {
			allocator, err = portauthority.New(30, 30)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error", func() {
			port, err := allocator.ClaimPorts(1)
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(BeEquivalentTo(30))

			port, err = allocator.ClaimPorts(1)
			Expect(port).To(BeZero())
			Expect(err).To(MatchError("insufficient ports available"))
		})
	})
	Context("bounds checking for ClaimPort()", func() {
		BeforeEach(func() {
			allocator, err = portauthority.New(65530, 65535)
			Expect(err).NotTo(HaveOccurred())
		})
		Context("when the allocator is requested more ports than can be possible", func() {
			It("errors", func() {
				_, err := allocator.ClaimPorts(65536)
				Expect(err).To(MatchError("number of ports requested must be between 1-65535"))
			})

		})
		Context("when the allocator is requested for 0 ports", func() {
			It("errors", func() {
				_, err := allocator.ClaimPorts(0)
				Expect(err).To(MatchError("number of ports requested must be between 1-65535"))
			})

		})
		Context("when the allocator is requested for negative ports", func() {
			It("errors", func() {
				_, err := allocator.ClaimPorts(-1)
				Expect(err).To(MatchError("number of ports requested must be between 1-65535"))
			})
		})

		Context("when the allocator is requested for a numbef of ports that would wrap-around", func() {
			It("errors", func() {
				_, err := allocator.ClaimPorts(7)
				Expect(err).To(MatchError("too many ports requested, will exceed maximum port of 65535"))
			})
		})
	})

	Context("when ClaimPorts is asked for more than one port", func() {
		BeforeEach(func() {
			port, err = allocator.ClaimPorts(4)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the first of those ports", func() {
			Expect(port).To(BeEquivalentTo(30))
		})

		It("skips those P ports the next time it is called", func() {
			port, err = allocator.ClaimPorts(4)
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(BeEquivalentTo(34))
		})

		Context("and there aren't enough ports available", func() {
			It("errors and returns 0 as a port", func() {
				port, err = allocator.ClaimPorts(65355)
				Expect(err).To(HaveOccurred())
				Expect(port).To(BeZero())
			})
		})
	})

	Context("when a start port is too high", func() {
		It("errors", func() {
			allocator, err = portauthority.New(65536, 30)
			Expect(err).To(MatchError("Invalid port range requested. Ports can only be numbers between 1-65535"))
		})
	})
	Context("when a end port is too high", func() {
		It("errors", func() {
			allocator, err = portauthority.New(30, 65536)
			Expect(err).To(MatchError("Invalid port range requested. Ports can only be numbers between 1-65535"))
		})
	})
	Context("when a negative start port is requested", func() {
		It("errors", func() {
			allocator, err = portauthority.New(-1, 30)
			Expect(err).To(MatchError("Invalid port range requested. Ports can only be numbers between 1-65535"))
		})
	})
	Context("when a negative end port is requested", func() {
		It("errors", func() {
			allocator, err = portauthority.New(1, -1)
			Expect(err).To(MatchError("Invalid port range requested. Ports can only be numbers between 1-65535"))
		})
	})
	Context("when a 0 start port is requested", func() {
		It("errors", func() {
			allocator, err = portauthority.New(0, 1)
			Expect(err).To(MatchError("Invalid port range requested. Ports can only be numbers between 1-65535"))
		})
	})
	Context("when a 0 end port is requested", func() {
		It("errors", func() {
			allocator, err = portauthority.New(1, 0)
			Expect(err).To(MatchError("Invalid port range requested. Ports can only be numbers between 1-65535"))
		})
	})
	Context("when an end port is lower than start port", func() {
		It("errors", func() {
			allocator, err = portauthority.New(30, 10)
			Expect(err).To(MatchError("Invalid port range requested. Starting port must be < ending port"))
		})
	})
})
