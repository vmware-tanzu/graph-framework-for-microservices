package gns

import (
	"net/http"

	cartv1 "github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"

	service_group "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns/service-group"
	policypkg "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/policy"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

var FooCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodDelete: nexus.HTTPCodesResponse{
		http.StatusOK:              nexus.HTTPResponse{Description: "ok"},
		http.StatusNotFound:        nexus.HTTPResponse{Description: http.StatusText(http.StatusNotFound)},
		nexus.DefaultHTTPErrorCode: nexus.DefaultHTTPError,
	},
}

var BarCustomCodesResponses = nexus.HTTPCodesResponse{
	http.StatusBadRequest: nexus.HTTPResponse{Description: "Bad Request"},
}

var BarCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodPatch: BarCustomCodesResponses,
}

var GNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha2/global-namespace/{Gns.gns}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/v1alpha2/global-namespaces",
			Methods: nexus.HTTPMethodsResponses{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
		{
			Uri:     "/test-foo",
			Methods: FooCustomMethodsResponses,
		},
		{
			Uri:     "/test-bar",
			Methods: BarCustomMethodsResponses,
		},
	},
}

var DNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha2/dns/{Dns.gns}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/v1alpha2/dnses",
			Methods: nexus.HTTPMethodsResponses{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
	},
}

// Gns struct.
// nexus-rest-api-gen:GNSRestAPISpec
// specification of GNS.
type Gns struct {
	nexus.Node
	//nexus-validation: MaxLength=8, MinLength=2
	//nexus-validation: Pattern=abc
	Domain                 string
	UseSharedGateway       bool
	Description            Description
	GnsServiceGroups       map[string]service_group.SvcGroup `nexus:"child"`
	GnsAccessControlPolicy policypkg.AccessControlPolicy     `nexus:"child"`
	Dns                    Dns                               `nexus:"link"`
	State                  GnsState                          `nexus:"status"`

	WorkloadSpec cartv1.WorkloadSpec //external-field
}

// This is Description struct.
type Description struct {
	Color     string
	Version   string
	ProjectId string
}

// nexus-rest-api-gen:DNSRestAPISpec
type Dns struct {
	nexus.Node
}

type GnsState struct {
	Working     bool
	Temperature int
}

type MyStr string
