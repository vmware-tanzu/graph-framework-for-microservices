package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var _ = Describe("Graphql parsing tests", func() {
	var (
		//err error
		pkgs map[string]parser.Package
		//pkg  parser.Package
		graph map[string]parser.Node
	)

	It("should parse graphql query specs", func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		graphqlQueries := parser.ParseGraphqlQuerySpecs(pkgs)
		graph, _, _ = parser.ParseDSLNodes(exampleDSLPath, baseGroupName, pkgs, graphqlQueries)
		root, ok := graph["roots.root.tsm.tanzu.vmware.com"]
		Expect(ok).To(BeTrue())
		config, ok := root.SingleChildren["Config"]
		Expect(ok).To(BeTrue())
		gns, ok := config.SingleChildren["GNS"]
		Expect(ok).To(BeTrue())

		Expect(gns.GraphqlQuerySpec.Queries).To(HaveLen(3))
		Expect(gns.GraphqlQuerySpec.Queries[0].Name).To(Equal("queryGns1"))
		Expect(gns.GraphqlQuerySpec.Queries[0].ServiceEndpoint.Domain).To(Equal("nexus-query-responder"))
		Expect(gns.GraphqlQuerySpec.Queries[0].ServiceEndpoint.Port).To(Equal(15000))
		args := gns.GraphqlQuerySpec.Queries[0].Args.([]parser.GraphQlArg)
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

		Expect(gns.GraphqlQuerySpec.Queries[1].Name).To(Equal("queryGnsQM1"))
		Expect(gns.GraphqlQuerySpec.Queries[1].ServiceEndpoint.Domain).To(Equal("query-manager"))
		Expect(gns.GraphqlQuerySpec.Queries[1].ServiceEndpoint.Port).To(Equal(15002))
		Expect(gns.GraphqlQuerySpec.Queries[1].Args).To(BeNil())
	})

	It("should match graphql query specs from other packages", func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		graphqlQueries := parser.ParseGraphqlQuerySpecs(pkgs)
		graph, _, _ = parser.ParseDSLNodes(exampleDSLPath, baseGroupName, pkgs, graphqlQueries)
		root, ok := graph["roots.root.tsm.tanzu.vmware.com"]
		Expect(ok).To(BeTrue())
		config, ok := root.SingleChildren["Config"]
		Expect(ok).To(BeTrue())

		Expect(config.GraphqlQuerySpec.Queries).To(HaveLen(1))
		Expect(config.GraphqlQuerySpec.Queries[0].Name).To(Equal("QueryExample"))
		Expect(config.GraphqlQuerySpec.Queries[0].ServiceEndpoint.Domain).To(Equal("query-manager"))
		Expect(config.GraphqlQuerySpec.Queries[0].ServiceEndpoint.Port).To(Equal(6000))
		args := config.GraphqlQuerySpec.Queries[0].Args.([]parser.GraphQlArg)
		Expect(len(args)).To(Equal(5))
	})

	It("should parse graphql files", func() {
		files := parser.ParseGraphQLFiles(exampleDSLPath)
		Expect(files["../../example/datamodel/example.graphql"]).ToNot(BeNil())
	})

	It("should parse graphql specs", func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		specs := parser.ParseGraphqlSpecs(pkgs)

		Expect(specs).To(Equal(map[string]nexus.GraphQLSpec{
			"root.GraphQLIdNameSpec": {
				IdName:     "name",
				IdNullable: true,
			},
			"root.GraphQLIdNullable": {
				IdNullable: false,
			},
		}))
	})

})
