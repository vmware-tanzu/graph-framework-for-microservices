package parser_test

import (
	"go/ast"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"

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
		graph = parser.ParseDSLNodes(exampleDSLPath, baseGroupName, nil, nil)
		root, ok = graph["roots.root.tsm.tanzu.vmware.com"]
		Expect(ok).To(BeTrue())
	})

	It("should get all nodes", func() {
		expectedNodes := []string{
			"ACPConfig",
			"AccessControlPolicy",
			"Config",
			"Dns",
			"Gns",
			"Root",
			"SvcGroup",
			"VMpolicy",
		}
		var nodes []string
		root.Walk(func(node *parser.Node) {
			nodes = append(nodes, node.Name)
		})
		sort.Strings(nodes)
		Expect(nodes).To(HaveLen(8))
		Expect(nodes).To(Equal(expectedNodes))
	})

	It("should fail when package names are duplicated.", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		parser.ParseDSLNodes("../../example/test-utils/invalid-pkg-name-datamodel", baseGroupName, nil, nil)
		Expect(fail).To(BeTrue())
	})

	It("should fail when nexus child or link fields is an array.", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		parser.ParseDSLNodes("../../example/test-utils/invalid-type-datamodel", baseGroupName, nil, nil)
		Expect(fail).To(BeTrue())
	})

	It("should fail when nexus child or link fields is a pointer.", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		parser.ParseDSLNodes("../../example/test-utils/pointer-type-datamodel", baseGroupName, nil, nil)
		Expect(fail).To(BeTrue())
	})

	It("should fail when nexus child is singleton node and is named", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		parser.ParseDSLNodes("../../example/test-utils/invalid-singleton-child", baseGroupName, nil, nil)
		Expect(fail).To(BeTrue())
	})

	It("should fail when used type name is reserved", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		parser.ParseDSLNodes("../../example/test-utils/invalid-type-name-datamodel", baseGroupName, nil, nil)
		Expect(fail).To(BeTrue())
	})
	It("should be able to get graphql info from a field", func() {
		graph = parser.ParseDSLNodes(exampleDSLPath, baseGroupName, nil, nil)
		config, ok := graph["roots.root.tsm.tanzu.vmware.com"].SingleChildren["Config"]
		Expect(ok).To(BeTrue())

		if val, ok := config.TypeSpec.Type.(*ast.StructType); ok {
			for _, f := range val.Fields.List {
				isJsonString := parser.IsJsonStringField(f)
				isIgnored := parser.IgnoreField(f)
				name, err := parser.GetFieldName(f)
				Expect(err).NotTo(HaveOccurred())
				switch name {
				case "FooD":
					Expect(isIgnored).To(BeFalse())
					Expect(isJsonString).To(BeTrue())
				case "FooC":
					Expect(isIgnored).To(BeTrue())
					Expect(isJsonString).To(BeFalse())
				case "FooE":
					Expect(isIgnored).To(BeTrue())
					Expect(isJsonString).To(BeFalse())
				case "FooF":
					Expect(isIgnored).To(BeFalse())
					Expect(isJsonString).To(BeTrue())
				default:
					Expect(isIgnored).To(BeFalse())
					Expect(isJsonString).To(BeFalse())
				}
			}
		}
	})
})
