package timeouts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTimeouts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timeouts Suite")
}
