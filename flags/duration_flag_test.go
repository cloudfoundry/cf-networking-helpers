package flags_test

import (
	"encoding/json"
	"time"

	"code.cloudfoundry.org/cf-networking-helpers/flags"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DurationFlag", func() {
	type obj struct {
		SomeInterval flags.DurationFlag `json:"some_interval"`
	}

	It("unmarshals from string to time.Duration", func() {
		contents := []byte(`{"some_interval": "10s"}`)

		var o obj
		err := json.Unmarshal(contents, &o)
		Expect(err).NotTo(HaveOccurred())

		Expect(time.Duration(o.SomeInterval)).To(Equal(10 * time.Second))
	})

	It("marshals from time.Duration to string", func() {
		o := obj{SomeInterval: flags.DurationFlag(10 * time.Second)}
		contents, err := json.Marshal(o)
		Expect(err).NotTo(HaveOccurred())

		Expect(contents).To(MatchJSON(`{"some_interval":"10s"}`))
	})
})
