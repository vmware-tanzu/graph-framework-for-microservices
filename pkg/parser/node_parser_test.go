package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
)

var _ = Describe("Node parser tests", func() {
	It("should parse dsl", func() {
		nodes := parser.ParseDSLNodes(exampleDSLPath)
		_, ok := nodes["gitlab.eng.vmware.com/nexus/compiler/example/datamodel/Root"]
		Expect(ok).To(BeTrue())
	})
})
