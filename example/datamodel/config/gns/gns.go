package gns

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config"
	"net/http"

	service_group "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns/service-group"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/policy"
        "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

var FooCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodDelete: nexus.HTTPCodesResponse{
		http.StatusOK:              nexus.HTTPResponse{Description: "ok"},
		http.StatusNotFound:        nexus.HTTPResponse{Description: http.StatusText(http.StatusNotFound)},
		nexus.DefaultHTTPErrorCode: nexus.DefaultHTTPError,
	},
}

var GNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/v1alpha2/projects/{project}/global-namespace/{Gns.gns}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/v1alpha2/projects/{project}/global-namespaces",
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
			Methods: config.BarCustomMethodsResponses,
		},
	},
}

// Gns struct.
// nexus-rest-api-gen:GNSRestAPISpec
// specification of GNS.
type Gns struct {
	nexus.Node
	Domain                 string
	UseSharedGateway       bool
	Description            Description
	GnsServiceGroups       map[string]service_group.SvcGroup `nexus:"child"`
	GnsAccessControlPolicy policy.AccessControlPolicy        `nexus:"child"`
	Dns                    Dns                               `nexus:"link"`
	State                  GnsState                          `nexus:"status"`
	// TODO support external fields https://jira.eng.vmware.com/browse/NPT-124
	//Box                    text.Box
}

// This is Description struct.
type Description struct {
	Color     string
	Version   string
	ProjectId string
}

var DNSRestAPISpec = nexus.RestAPISpec{
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

// nexus-rest-api-gen:DNSRestAPISpec
type Dns struct {
	nexus.Node
}

type GnsState struct {
	Working     bool
	Temperature int
}
