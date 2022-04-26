package rest_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser/rest"
)

var _ = Describe("Rest tests", func() {
	var (
		//err error
		pkgs    map[string]parser.Package
		pkg     parser.Package
		ok      bool
		methods map[string]nexus.HTTPMethodsResponses
		codes   map[string]nexus.HTTPCodesResponse
	)

	BeforeEach(func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel//config/gns"]
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
		apiSpecs := rest.GetRestApiSpecs(pkg, methods, codes)

		expectedDnsApiSpec := nexus.RestAPISpec{
			Uris: []nexus.RestURIs{
				{
					Uri:     "/v1alpha2/projects/{project}/dns/{Dns.gns}",
					Methods: nexus.DefaultHTTPMethodsResponses,
				},
				{
					Uri: "/v1alpha2/projects/{project}/dnses",
					Methods: nexus.HTTPMethodsResponses{
						http.MethodGet: nexus.DefaultHTTPGETResponses,
					},
				},
			},
		}

		dnsRestApiSpec, ok := apiSpecs["DNSRestAPISpec"]
		Expect(ok).To(BeTrue())
		Expect(dnsRestApiSpec).To(Equal(expectedDnsApiSpec))
	})
})
