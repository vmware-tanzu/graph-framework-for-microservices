package generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

var _ = Describe("Graphql Custom query generator tests", func() {
	var (
		pkgs parser.Packages
		gns  parser.Node
	)

	BeforeEach(func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		graphqlQueries := parser.ParseGraphqlQuerySpecs(pkgs)
		graph, _, _ := parser.ParseDSLNodes(exampleDSLPath, baseGroupName, pkgs, graphqlQueries)
		root, ok := graph["roots.root.tsm.tanzu.vmware.com"]
		Expect(ok).To(BeTrue())
		config, ok := root.SingleChildren["Config"]
		Expect(ok).To(BeTrue())
		gns, ok = config.SingleChildren["GNS"]
		Expect(ok).To(BeTrue())
	})

	It("should translate graphql query spec to schema", func() {
		schema := generator.CustomQueryToGraphqlSchema(gns.GraphqlQuerySpec.Queries[0])
		Expect(schema).To(Equal(`    queryGns1(
        StartTime: String
        EndTime: String
        Interval: String
        IsServiceDeployment: Boolean
        StartVal: Int
    ): NexusGraphqlResponse`))
		schema = generator.CustomQueryToGraphqlSchema(gns.GraphqlQuerySpec.Queries[1])
		Expect(schema).To(Equal(`    queryGnsQM1: TimeSeriesData`))
	})
})
