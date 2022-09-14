package gns

import (
	"net/http"

	cartv1 "github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"

	service_group "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns/service-group"
	policypkg "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/policy"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/nexus"
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
			Uri: "/v1alpha2/global-namespace/{gns.Gns}",
			QueryParams: []string{
				"config.Config",
			},
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/v1alpha2/global-namespaces",
			QueryParams: []string{
				"config.Config",
			},
			Methods: nexus.HTTPListResponse,
		},
		{
			Uri: "/test-foo",
			QueryParams: []string{
				"config.Config",
			},
			Methods: FooCustomMethodsResponses,
		},
		{
			Uri: "/test-bar",
			QueryParams: []string{
				"config.Config",
			},
			Methods: BarCustomMethodsResponses,
		},
	},
}

var DNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri: "/v1alpha2/dns",
			QueryParams: []string{
				"config.Config",
				"gns.Dns",
			},
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/v1alpha2/dnses",
			QueryParams: []string{
				"config.Config",
			},
			Methods: nexus.HTTPListResponse,
		},
	},
}

type MyConst string
type SourceKind string

const (
	Object SourceKind = "Object"
	Type   SourceKind = "Type"
	XYZ    MyConst    = "xyz"
)

type ReplicationSource struct {
	Kind SourceKind
}

// Gns struct.
// nexus-rest-api-gen:GNSRestAPISpec
// nexus-description: this is my awesome node
// specification of GNS.
type Gns struct {
	nexus.Node
	//nexus-validation: MaxLength=8, MinLength=2
	//nexus-validation: Pattern=abc
	Domain                 string
	UseSharedGateway       bool
	Description            Description
	GnsServiceGroups       service_group.SvcGroup        `nexus:"children"`
	GnsAccessControlPolicy policypkg.AccessControlPolicy `nexus:"child"`
	Dns                    Dns                           `nexus:"link"`
	State                  GnsState                      `nexus:"status"`
	Meta                   string

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
	nexus.SingletonNode
}

type GnsState struct {
	Working     bool
	Temperature int
}

type MyStr string
