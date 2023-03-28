package testsupport_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestsupport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testsupport Suite")
}
