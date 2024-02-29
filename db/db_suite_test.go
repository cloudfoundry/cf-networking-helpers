package db_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB Suite")
}

var RandomGenerator *rand.Rand
var _ = BeforeSuite(func() {
	randomGenerator = rand.New(rand.NewSource(GinkgoRandomSeed() + int64(GinkgoParallelProcess())))
})
