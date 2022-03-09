package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
)

var _ = Describe("Pkg parser tests", func() {
	It("should parse dsl", func() {
		pkgs := parser.ParseDSLPkg(exampleDSLPath)
		_, ok := pkgs["gitlab.eng.vmware.com/nexus/compiler/example/datamodel"]
		Expect(ok).To(BeTrue())
	})
})
