package shuntyard_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestShuntyard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shuntyard Suite")
}
