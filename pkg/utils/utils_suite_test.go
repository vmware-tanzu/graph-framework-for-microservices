package utils_test

import (
	"testing"

	log "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	log.StandardLogger().ExitFunc = nil
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Test Suite")
}
