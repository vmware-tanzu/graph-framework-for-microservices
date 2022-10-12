package crd_generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	crdgenerator "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/crd-generator"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

var _ = Describe("Template renderers tests", func() {
	var (
		pkgs       parser.Packages
		parentsMap map[string]parser.NodeHelper
	)

	BeforeEach(func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		graph := parser.ParseDSLNodes(exampleDSLPath, baseGroupName)
		parentsMap = parser.CreateParentsMap(graph)
	})

	It("should resolve graphql vars", func() {
		vars, err := crdgenerator.GenerateGraphqlResolverVars(baseGroupName, crdModulePath, pkgs, parentsMap)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(vars)).To(Equal(12))
		Expect(vars[0].NodeName).To(Equal("Root"))
		Expect(vars[2].PkgName).To(Equal("Config"))
		Expect(vars[2].NodeName).To(Equal("FooType"))
		Expect(vars[2].SchemaName).To(Equal("config_FooType"))
		Expect(vars[2].Alias).To(Equal(""))
		Expect(vars[2].ReturnType).To(Equal(""))

		Expect(vars[2].IsParentNode).To(BeFalse())
		Expect(vars[2].HasParent).To(BeFalse())
		Expect(vars[2].IsSingletonNode).To(BeFalse())
		Expect(vars[2].IsNexusNode).To(BeTrue())
		Expect(vars[2].BaseImportPath).To(Equal("nexustempmodule/"))
		Expect(vars[2].CrdName).To(Equal(""))
	})

	It("should resolve non-singleton root and singleton child node", func() {
		pkgs = parser.ParseDSLPkg("../../example/test-utils/non-singleton-root")
		graph := parser.ParseDSLNodes("../../example/test-utils/non-singleton-root", baseGroupName)
		parentsMap = parser.CreateParentsMap(graph)

		vars, err := crdgenerator.GenerateGraphqlResolverVars(baseGroupName, crdModulePath, pkgs, parentsMap)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(vars)).To(Equal(3))

		Expect(vars[0].NodeName).To(Equal("Root"))
		Expect(vars[1].PkgName).To(Equal("Config"))
		Expect(vars[1].NodeName).To(Equal("Config"))
		Expect(vars[1].SchemaName).To(Equal("config_Config"))

		Expect(vars[1].IsParentNode).To(BeFalse())
		Expect(vars[1].HasParent).To(BeFalse())
		Expect(vars[1].IsSingletonNode).To(BeTrue())
		Expect(vars[1].IsNexusNode).To(BeTrue())
		Expect(vars[1].BaseImportPath).To(Equal("nexustempmodule/"))
		Expect(vars[1].CrdName).To(Equal(""))
	})
})
