package crd_generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCRDGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRD generator Suite")
}
