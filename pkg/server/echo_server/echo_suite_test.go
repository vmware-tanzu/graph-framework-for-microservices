package echo_server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEcho(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Echo Suite")
}
