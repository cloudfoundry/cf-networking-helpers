package httperror_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHttperror(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httperror Suite")
}
