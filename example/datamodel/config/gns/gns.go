package gns

import (
	"net/http"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
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


type Port uint16

// Host the IP of the endpoint
type Host string

type HostPort struct {
	Host Host
	Port Port
}

type Instance float32
type AliasArr []int

type gnsQueryFilters struct {
	StartTime           string
	EndTime             string
	Interval            string
	IsServiceDeployment bool
	StartVal            int
}

var CloudEndpointGraphQLQuerySpec = nexus.GraphQLQuerySpec{
	Queries: []nexus.GraphQLQuery{
		{
			Name: "queryGns1",
			ServiceEndpoint: nexus.GraphQLQueryEndpoint{
				Domain: "query-manager",
				Port:   15000,
			},
			Args: gnsQueryFilters{},
		},
		{
			Name: "queryGns2",
			ServiceEndpoint: nexus.GraphQLQueryEndpoint{
				Domain: "query-manager2",
				Port:   15002,
			},
			Args: nil,
		},
	},
}

// Gns struct.
// nexus-rest-api-gen:GNSRestAPISpec
// nexus-graphql-query:CloudEndpointGraphQLQuerySpec
// nexus-description: this is my awesome node
// specification of GNS.
type Gns struct {
	nexus.Node
	//nexus-validation: MaxLength=8, MinLength=2
	//nexus-validation: Pattern=abc
	Domain                 string
	UseSharedGateway       bool
	Mydesc           Description
	FooLink          BarLink     `nexus:"link"`
	FooLinks         BarLinks    `nexus:"links"`
	FooChild         BarChild    `nexus:"child"`
	FooChildren      BarChildren `nexus:"children"`
	HostPort         HostPort
	Instance         Instance
	Array1           float32
	Array2           []Description
	Array3           []BarLink
	Array4           []Instance
	Array5           AliasArr
	TestABCLink      ABCLink `nexus:"links"`
}

// This is Description struct.
type Description struct {
	Color     string
	Version   string
	ProjectID []string
	TestAns   []Answer
	Instance  Instance
	HostPort  HostPort
}

type BarLink struct {
	nexus.SingletonNode
	Name string
}

type BarChild struct {
	nexus.SingletonNode
	Name string
}

type BarChildren struct {
	nexus.Node
	Name string
}

type BarLinks struct {
	nexus.Node
	Name string
}

type Answer struct {
	Name string
}

type ABCLink struct {
	nexus.Node
}
