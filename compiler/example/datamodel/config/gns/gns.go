package gns

import (
	"net/http"

	cartv1 "github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"

	service_group "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns/service-group"
	policypkg "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/policy"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
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

type Port uint16
type Host string

type HostPort struct {
	Host Host
	Port Port
}

type Instance float32
type AliasArr []int

const (
	Object SourceKind = "Object"
	Type   SourceKind = "Type"
	XYZ    MyConst    = "xyz"
)

type ReplicationSource struct {
	Kind SourceKind
}

type gnsQueryFilters struct {
	StartTime           string
	EndTime             string
	Interval            string
	IsServiceDeployment bool
	StartVal            int
}

type metricsFilers struct {
	StartTime    string
	EndTime      string
	TimeInterval string
	SomeUserArg1 string
	SomeUserArg2 int
	SomeUserArg3 bool
}

var CloudEndpointGraphQLQuerySpec = nexus.GraphQLQuerySpec{
	Queries: []nexus.GraphQLQuery{
		{
			Name: "queryGns1",
			ServiceEndpoint: nexus.GraphQLQueryEndpoint{
				Domain: "nexus-query-responder",
				Port:   15000,
			},
			Args:    gnsQueryFilters{},
			ApiType: nexus.GraphQLQueryApi,
		},
		{
			Name: "queryGnsQM1",
			ServiceEndpoint: nexus.GraphQLQueryEndpoint{
				Domain: "query-manager",
				Port:   15002,
			},
			Args:    nil,
			ApiType: nexus.GetMetricsApi,
		},
		{
			Name: "queryGnsQM",
			ServiceEndpoint: nexus.GraphQLQueryEndpoint{
				Domain: "query-manager",
				Port:   15003,
			},
			Args:    metricsFilers{},
			ApiType: nexus.GetMetricsApi,
		},
	},
}

// Gns struct.
// nexus-graphql-query:CloudEndpointGraphQLQuerySpec
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
	GnsAccessControlPolicy policypkg.AccessControlPolicy `nexus:"child" nexus-graphql:"type:string"`
	Dns                    Dns                           `nexus:"link" nexus-graphql:"ignore:true"`
	State                  GnsState                      `nexus:"status"`
	FooChild               BarChild                      `nexus:"child" nexus-graphql:"type:string"`
	IgnoreChild            IgnoreChild                   `nexus:"child" nexus-graphql:"ignore:true"`
	Meta                   string

	Port             *int         // pointer test
	OtherDescription *Description // pointer test - struct
	MapPointer       *map[string]string
	SlicePointer     *[]string

	WorkloadSpec  cartv1.WorkloadSpec  //external-field
	DifferentSpec *cartv1.WorkloadSpec // external-field - pointer

	ServiceSegmentRef        ServiceSegmentRef            `json:"serviceSegmentRef,omitempty"`
	ServiceSegmentRefPointer *ServiceSegmentRef           `json:"serviceSegmentRefPointer,omitempty"`
	ServiceSegmentRefs       []ServiceSegmentRef          `json:"serviceSegmentRefs,omitempty"`
	ServiceSegmentRefMap     map[string]ServiceSegmentRef `json:"serviceSegmentRefMap,omitempty"`
}

type ServiceSegmentRef struct {
	Field1 string
	Field2 string
}

// This is Description struct.
type Description struct {
	Color     string
	Version   string
	ProjectId string
	TestAns   []Answer
	Instance  Instance
	HostPort  HostPort
}

type BarChild struct {
	nexus.SingletonNode
	Name string
}

type IgnoreChild struct {
	nexus.Node
	Name string
}

// nexus-rest-api-gen:DNSRestAPISpec
type Dns struct {
	nexus.SingletonNode
}

type Answer struct {
	Name string
}

type GnsState struct {
	Working     bool
	Temperature int
}

type MyStr string
