package rest_test

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser/rest"
)

var _ = Describe("Rest tests", func() {
	var (
		//err error
		pkgs       map[string]parser.Package
		pkg        parser.Package
		ok         bool
		methods    map[string]nexus.HTTPMethodsResponses
		codes      map[string]nexus.HTTPCodesResponse
		parentsMap map[string]parser.NodeHelper
	)

	BeforeEach(func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns"]
		graph := parser.ParseDSLNodes(exampleDSLPath, baseGroupName)
		parentsMap = parser.CreateParentsMap(graph)
		Expect(ok).To(BeTrue())
	})

	It("should parse responses", func() {
		methods, codes = rest.ParseResponses(pkgs)

		expectedFooMethods := nexus.HTTPMethodsResponses{
			http.MethodDelete: nexus.HTTPCodesResponse{
				http.StatusOK:              nexus.HTTPResponse{Description: "ok"},
				http.StatusNotFound:        nexus.HTTPResponse{Description: http.StatusText(http.StatusNotFound)},
				nexus.DefaultHTTPErrorCode: nexus.DefaultHTTPError,
			},
		}

		fooMethods, ok := methods["FooCustomMethodsResponses"]
		Expect(ok).To(BeTrue())
		Expect(fooMethods).To(Equal(expectedFooMethods))

		expectedBarCodes := nexus.HTTPCodesResponse{
			http.StatusBadRequest: nexus.HTTPResponse{Description: "Bad Request"},
		}

		barCodes, ok := codes["BarCustomCodesResponses"]
		Expect(ok).To(BeTrue())
		Expect(barCodes).To(Equal(expectedBarCodes))
	})

	It("should get rest api specs for gns package", func() {
		apiSpecs := rest.GetRestApiSpecs(pkg, methods, codes, parentsMap)

		expectedDnsApiSpec := nexus.RestAPISpec{
			Uris: []nexus.RestURIs{
				{
					Uri:         "/v1alpha2/dns",
					QueryParams: []string{"config.Config", "gns.Dns"},
					Methods:     nexus.DefaultHTTPMethodsResponses,
				},
				{
					Uri:         "/v1alpha2/dnses",
					QueryParams: []string{"config.Config"},
					Methods:     nexus.HTTPListResponse,
				},
			},
		}

		dnsRestApiSpec, ok := apiSpecs["DNSRestAPISpec"]
		Expect(ok).To(BeTrue())
		Expect(dnsRestApiSpec).To(Equal(expectedDnsApiSpec))
	})

	It("should validate RestAPISpec for list endpoint", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		restApiSpec := nexus.RestAPISpec{
			Uris: []nexus.RestURIs{
				{
					Uri: "/v1alpha2/dnses",
					QueryParams: []string{
						"config.Config",
					},
					Methods: nexus.HTTPListResponse,
				},
			},
		}
		rest.ValidateRestApiSpec(restApiSpec, parentsMap, "dnses.gns.tsm.tanzu.vmware.com")
		Expect(fail).To(BeFalse())
	})

	It("should fail validation of RestAPISpec for list endpoint with node name in URI", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		restApiSpec := nexus.RestAPISpec{
			Uris: []nexus.RestURIs{
				{
					Uri: "/v1alpha2/dnses/{gns.Dns}",
					QueryParams: []string{
						"config.Config",
					},
					Methods: nexus.HTTPListResponse,
				},
			},
		}
		rest.ValidateRestApiSpec(restApiSpec, parentsMap, "dnses.gns.tsm.tanzu.vmware.com")
		Expect(fail).To(BeTrue())
	})
})
