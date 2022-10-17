package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

var _ = Describe("Pkg parser tests", func() {
	It("should parse dsl", func() {
		pkgs := parser.ParseDSLPkg(exampleDSLPath)
		_, ok := pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel"]
		Expect(ok).To(BeTrue())
	})
})
