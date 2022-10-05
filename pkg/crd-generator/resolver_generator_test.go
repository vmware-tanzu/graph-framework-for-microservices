package crd_generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"

	crdgenerator "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/crd-generator"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

var _ = Describe("Template renderers tests", func() {
	var (
		pkgs           parser.Packages
		parentsMap     map[string]parser.NodeHelper
		graphqlQueries map[string]nexus.GraphQLQuerySpec
	)

	BeforeEach(func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		graphqlQueries = parser.ParseGraphqlQuerySpecs(pkgs)
		graph := parser.ParseDSLNodes(exampleDSLPath, baseGroupName, pkgs, graphqlQueries)
		parentsMap = parser.CreateParentsMap(graph)
	})

	It("should resolve graphql vars", func() {
		vars, err := crdgenerator.GenerateGraphqlResolverVars(baseGroupName, crdModulePath, pkgs, parentsMap)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(vars)).To(Equal(42))
		Expect(vars[0].NodeName).To(Equal("Root"))
		Expect(vars[3].PkgName).To(Equal("Config"))
		Expect(vars[3].NodeName).To(Equal("Config"))
		Expect(vars[3].SchemaName).To(Equal("config_Config"))
		Expect(vars[3].Alias).To(Equal(""))
		Expect(vars[3].ReturnType).To(Equal(""))

		Expect(vars[3].IsParentNode).To(BeFalse())
		Expect(vars[3].HasParent).To(BeFalse())
		Expect(vars[3].IsSingletonNode).To(BeFalse())
		Expect(vars[3].IsNexusNode).To(BeTrue())
		Expect(vars[3].BaseImportPath).To(Equal("nexustempmodule/"))
		Expect(vars[3].CrdName).To(Equal(""))
	})

	It("should resolve non-singleton root and singleton child node", func() {
		pkgs = parser.ParseDSLPkg("../../example/test-utils/non-singleton-root")
		graph := parser.ParseDSLNodes("../../example/test-utils/non-singleton-root", baseGroupName, pkgs, graphqlQueries)
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
