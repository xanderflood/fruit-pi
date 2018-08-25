package chamber_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChamber(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chamber Suite")
}
