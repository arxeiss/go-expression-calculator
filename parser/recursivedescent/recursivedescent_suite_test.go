package recursivedescent_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRecursivedescent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Recursivedescent Suite")
}
