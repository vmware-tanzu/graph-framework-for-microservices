package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

var _ = Describe("Graphql parsing tests", func() {
	var (
		//err error
		pkgs map[string]parser.Package
		//pkg  parser.Package
		graph map[string]parser.Node
	)

	FIt("should parse graphql query specs", func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		parser.ParseGraphqlQuerySpecs(pkgs)
		graph = parser.ParseDSLNodes(exampleDSLPath, baseGroupName, pkgs)
		root, ok := graph["roots.root.tsm.tanzu.vmware.com"]
		Expect(ok).To(BeTrue())
		config, ok := root.SingleChildren["Config"]
		Expect(ok).To(BeTrue())
		gns, ok := config.SingleChildren["GNS"]
		Expect(ok).To(BeTrue())

		Expect(gns.GraphqlSpec.Queries).To(HaveLen(2))
		Expect(gns.GraphqlSpec.Queries[0].Name).To(Equal("queryGns1"))
		Expect(gns.GraphqlSpec.Queries[0].ServiceEndpoint.Domain).To(Equal("query-manager"))
		Expect(gns.GraphqlSpec.Queries[0].ServiceEndpoint.Port).To(Equal(15000))
		args := gns.GraphqlSpec.Queries[0].Args.([]parser.GraphQlArg)
		Expect(len(args)).To(Equal(5))

		Expect(args[0].Name).To(Equal("StartTime"))
		Expect(args[0].Type).To(Equal("string"))
		Expect(args[1].Name).To(Equal("EndTime"))
		Expect(args[1].Type).To(Equal("string"))
		Expect(args[2].Name).To(Equal("Interval"))
		Expect(args[2].Type).To(Equal("string"))
		Expect(args[3].Name).To(Equal("IsServiceDeployment"))
		Expect(args[3].Type).To(Equal("bool"))
		Expect(args[4].Name).To(Equal("StartVal"))
		Expect(args[4].Type).To(Equal("int"))

		Expect(gns.GraphqlSpec.Queries[1].Name).To(Equal("queryGns2"))
		Expect(gns.GraphqlSpec.Queries[1].ServiceEndpoint.Domain).To(Equal("query-manager2"))
		Expect(gns.GraphqlSpec.Queries[1].ServiceEndpoint.Port).To(Equal(15002))
		Expect(gns.GraphqlSpec.Queries[1].Args).To(BeNil())
	})

	FIt("should match graphql query specs from other packages", func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		parser.ParseGraphqlQuerySpecs(pkgs)
		graph = parser.ParseDSLNodes(exampleDSLPath, baseGroupName, pkgs)
		root, ok := graph["roots.root.tsm.tanzu.vmware.com"]
		Expect(ok).To(BeTrue())
		config, ok := root.SingleChildren["Config"]
		Expect(ok).To(BeTrue())

		Expect(config.GraphqlSpec.Queries).To(HaveLen(1))
		Expect(config.GraphqlSpec.Queries[0].Name).To(Equal("query"))
		Expect(config.GraphqlSpec.Queries[0].ServiceEndpoint.Domain).To(Equal("query-manager"))
		Expect(config.GraphqlSpec.Queries[0].ServiceEndpoint.Port).To(Equal(6000))
		args := config.GraphqlSpec.Queries[0].Args.([]parser.GraphQlArg)
		Expect(len(args)).To(Equal(5))
	})

})
