package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

var _ = Describe("Node parser tests", func() {
	var (
		//err error
		graph map[string]parser.Node
		root  parser.Node
		ok    bool
	)

	BeforeEach(func() {
		graph = parser.ParseDSLNodes(exampleDSLPath, baseGroupName)
		root, ok = graph["roots.root.tsm.tanzu.vmware.com"]
		Expect(ok).To(BeTrue())
	})

	It("should get all nodes", func() {
		expectedNodes := []string{"Root", "Config", "Gns", "AccessControlPolicy", "ACPConfig", "SvcGroup"}
		var nodes []string
		root.Walk(func(node *parser.Node) {
			nodes = append(nodes, node.Name)
		})
		Expect(nodes).To(HaveLen(6))
		Expect(nodes).To(Equal(expectedNodes))
	})
})
