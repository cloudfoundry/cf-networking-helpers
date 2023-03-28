package json_client_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJsonClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JsonClient Suite")
}
