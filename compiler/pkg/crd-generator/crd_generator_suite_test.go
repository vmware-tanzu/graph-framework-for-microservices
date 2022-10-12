package crd_generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/config"
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
