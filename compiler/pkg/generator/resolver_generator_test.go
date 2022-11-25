package generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var _ = Describe("Template renderers tests", func() {
	var (
		pkgs           parser.Packages
		parentsMap     map[string]parser.NodeHelper
		graphqlQueries map[string]nexus.GraphQLQuerySpec
	)

	BeforeEach(func() {
		config.ConfigInstance.IgnoredDirs = []string{"ignored"}
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		graphqlQueries = parser.ParseGraphqlQuerySpecs(pkgs)
		graph, _, _ := parser.ParseDSLNodes(exampleDSLPath, baseGroupName, pkgs, graphqlQueries)
		parentsMap = parser.CreateParentsMap(graph)
	})

	It("should resolve graphql vars", func() {
		vars, err := generator.GenerateGraphqlResolverVars(baseGroupName, crdModulePath, pkgs, parentsMap)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(vars)).To(Equal(41))
		Expect(vars[0].NodeName).To(Equal("Root"))
		Expect(vars[3].PkgName).To(Equal("Config"))
		Expect(vars[3].NodeName).To(Equal("Config"))
		Expect(vars[3].SchemaName).To(Equal("config_Config"))
		Expect(vars[3].Alias).To(Equal(""))
		Expect(vars[3].ReturnType).To(Equal(""))

		Expect(vars[2].IsParentNode).To(BeFalse())
		Expect(vars[2].HasParent).To(BeFalse())
		Expect(vars[2].IsSingletonNode).To(BeFalse())
		Expect(vars[2].IsNexusNode).To(BeFalse())
		Expect(vars[2].BaseImportPath).To(Equal("nexustempmodule/"))
		Expect(vars[2].CrdName).To(Equal(""))

		// Test that nexus-secret Node should not be present
		Expect(vars).NotTo(ContainElement(HaveField("NodeName", "Foo")))
		Expect(vars).NotTo(ContainElement(HaveField("NodeName", "GnsServiceGroups")))
		Expect(vars).NotTo(ContainElement(HaveField("NodeName", "SourceSvcGroups")))
		Expect(vars).NotTo(ContainElement(HaveField("NodeName", "DestSvcGroups")))
	})

	It("should resolve non-singleton root and singleton child node", func() {
		pkgs = parser.ParseDSLPkg("../../example/test-utils/non-singleton-root")
		graph, _, _ := parser.ParseDSLNodes("../../example/test-utils/non-singleton-root", baseGroupName, pkgs, graphqlQueries)
		parentsMap = parser.CreateParentsMap(graph)

		vars, err := generator.GenerateGraphqlResolverVars(baseGroupName, crdModulePath, pkgs, parentsMap)
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

	It("should validate the import pkg and translate to graphql schema and resolver typeName", func() {
		pkg := pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"]
		schemaTypeName, resolverTypeName := generator.ValidateImportPkg(pkg.Name, "servicegroup.SvcGroup", pkg.GetImportMap(), pkgs)

		Expect(pkg.Name).To(Equal("gns"))
		Expect(schemaTypeName).To(Equal("servicegroup_SvcGroup"))
		Expect(resolverTypeName).To(Equal("ServicegroupSvcGroup"))
	})

	It("should resolve secret spec annotation properly", func() {
		// nexus-secret-spec annotated, different pkg and directory name with no alias name
		pkg := pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/policy"]
		nodeToSkip := generator.SkipSecretSpecAnnotation("DestSvcGroups", "servicegroup.SvcGroup", pkg, pkg.GetImportMap(), pkgs)
		Expect(nodeToSkip).Should(BeTrue())

		// nexus-secret-spec annotated, different pkg and directory name with alias name provided
		pkg = pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"]
		nodeToSkip = generator.SkipSecretSpecAnnotation("GnsServiceGroups", "service_group.SvcGroup", pkg, pkg.GetImportMap(), pkgs)
		Expect(nodeToSkip).Should(BeTrue())

		// nexus-secret-spec annotated, same pkg and directory name
		nodeToSkip = generator.SkipSecretSpecAnnotation("Foo", "Foo", pkg, pkg.GetImportMap(), pkgs)
		Expect(nodeToSkip).Should(BeTrue())

		// Node is not annotated with nexus-secret-spec, should not skip
		nodeToSkip = generator.SkipSecretSpecAnnotation("GnsAccessControlPolicy", "policypkg.AccessControlPolicy", pkg, pkg.GetImportMap(), pkgs)
		Expect(nodeToSkip).Should(BeFalse())
	})
})
