package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

var _ = Describe("Node config tests", func() {
	//var (
	//	//err error
	//	pkg parser.Package
	//	ok  bool
	//)

	BeforeEach(func() {
		pkgs := parser.ParseDSLPkg(exampleDSLPath)
		_, ok := pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel"]
		Expect(ok).To(BeTrue())
	})

	//It("should parse root node config", func() {
	//	cfg, err := parser.GetNexusNodeConfig(pkg, "Root")
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	cfgYaml, err := yaml.Marshal(*cfg)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	expectedCfg := parser.NexusNodeConfig{
	//		NexusRestAPIGen: "nexus-rest-api-gen: GnsRestAPISpec",
	//	}
	//
	//	expectedCfgYaml, err := yaml.Marshal(expectedCfg)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	Expect(string(cfgYaml)).To(Equal(string(expectedCfgYaml)))
	//})
})
