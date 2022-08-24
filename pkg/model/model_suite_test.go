package model_test

import (
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModel(t *testing.T) {
	log.StandardLogger().ExitFunc = nil
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model Test Suite")
}
