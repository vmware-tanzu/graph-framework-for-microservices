package envoy_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnvoy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envoy Suite")
}
