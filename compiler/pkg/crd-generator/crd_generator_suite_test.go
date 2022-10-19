package crd_generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"
)

func TestCRDGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRD generator Suite")
}

func init() {
	conf := &config.Config{
		GroupName:     baseGroupName,
		CrdModulePath: crdModulePath,
	}
	config.ConfigInstance = conf
}
