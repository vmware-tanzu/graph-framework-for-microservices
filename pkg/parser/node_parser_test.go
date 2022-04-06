package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

const (
	baseGroupName = "tsm.tanzu.vmware.com"
	crdModulePath = "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/_crd_generated/"
)

var _ = Describe("Node parser tests", func() {
	It("should parse dsl", func() {
		nodes := parser.ParseDSLNodes(exampleDSLPath, baseGroupName)
		_, ok := nodes["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/Root"]
		Expect(ok).To(BeTrue())
	})
})
