package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Node config tests", func() {
	var (
		//err error
		pkg parser.Package
		ok  bool
	)

	BeforeEach(func() {
		pkgs := parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel"]
		Expect(ok).To(BeTrue())
	})

	It("should parse root node config", func() {
		cfg, err := parser.GetNexusNodeConfig(pkg, "Root")
		Expect(err).NotTo(HaveOccurred())

		cfgYaml, err := yaml.Marshal(*cfg)
		Expect(err).NotTo(HaveOccurred())

		expectedCfg := parser.NexusNodeConfig{
			NexusRestAPIGen: parser.NexusRestAPIGen{
				URI:     "/v1alpha1/projects/$PID/global-namespace/$GID",
				Methods: []string{"GET", "PUT", "DELETE"},
				Response: parser.Response{
					Num200: parser.Num200{
						Message: "success message",
					},
					Num400: parser.Num400{
						Message: "not found message",
					},
					Num401: parser.Num401{
						Message: "unauthorized message",
					},
				},
			},
			NexusAPIValidationEndpoint: []parser.NexusAPIValidationEndpoint{
				{
					Service:  "service-name",
					Endpoint: "/foo/bar",
				},
				{
					Service:  "service-name-2",
					Endpoint: "/foo/bar",
				},
			},
			NexusVersion: "v1",
		}

		expectedCfgYaml, err := yaml.Marshal(expectedCfg)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(cfgYaml)).To(Equal(string(expectedCfgYaml)))
	})
})
