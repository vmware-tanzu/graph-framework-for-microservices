package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

var _ = Describe("Node config tests", func() {
	var (
		//err error
		pkg parser.Package
		ok  bool
	)

	BeforeEach(func() {
		pkgs := parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"]
		Expect(ok).To(BeTrue())
	})

	It("should parse gns node annotation", func() {
		annotation, ok := parser.GetNexusRestAPIGenAnnotation(pkg, "Gns")
		Expect(ok).To(BeTrue())
		Expect(annotation).To(Equal("GNSRestAPISpec"))
	})
})
