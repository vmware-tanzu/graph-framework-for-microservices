package nexus

import "net/http"

const (
	tagName              = "nexus.vmware"
	BaseGroupName string = "tanzu.tsm.vmware.com"
)

type ID struct {
	Id string `nexus.vmware:"id"`
}

type Node struct {
	ID
}

type SingletonNode struct {
	ID
}

// HTTPMethod type.
type HTTPMethod string

// ResponseCode type.
type ResponseCode int

// HTTPResponse type.
type HTTPResponse struct {
	Description string `json:"description"`
}

// HTTPCodesResponse code to response type.
type HTTPCodesResponse map[ResponseCode]HTTPResponse

// HTTPMethodsResponses to response mapping.
type HTTPMethodsResponses map[HTTPMethod]HTTPCodesResponse

// RestURIs and associated data.
type RestURIs struct {
	Uri         string               `json:"uri"`
	QueryParams []string             `json:"query_params,omitempty"`
	Methods     HTTPMethodsResponses `json:"methods"`
}

type RestAPISpec struct {
	Uris []RestURIs `json:"uris"`
}

// Default HTTP error code and description.
const DefaultHTTPErrorCode ResponseCode = http.StatusNotImplemented

var DefaultHTTPError = HTTPResponse{Description: http.StatusText(http.StatusNotImplemented)}

// Default HTTP GET Response mappings.
var DefaultHTTPGETResponses = HTTPCodesResponse{
	http.StatusOK:        HTTPResponse{Description: http.StatusText(http.StatusOK)},
	http.StatusNotFound:  HTTPResponse{Description: http.StatusText(http.StatusNotFound)},
	DefaultHTTPErrorCode: DefaultHTTPError,
}

// Default HTTP PUT Response mappings.
var DefaultHTTPPUTResponses = HTTPCodesResponse{
	http.StatusOK:        HTTPResponse{Description: http.StatusText(http.StatusOK)},
	http.StatusCreated:   HTTPResponse{Description: http.StatusText(http.StatusCreated)},
	DefaultHTTPErrorCode: DefaultHTTPError,
}

// Default HTTP DELETE Response mappings.
var DefaultHTTPDELETEResponses = HTTPCodesResponse{
	http.StatusOK:        HTTPResponse{Description: http.StatusText(http.StatusOK)},
	http.StatusNotFound:  HTTPResponse{Description: http.StatusText(http.StatusNotFound)},
	DefaultHTTPErrorCode: DefaultHTTPError,
}

// Default HTTP methods and responses.
var DefaultHTTPMethodsResponses = HTTPMethodsResponses{
	http.MethodGet:    DefaultHTTPGETResponses,
	http.MethodPut:    DefaultHTTPPUTResponses,
	http.MethodDelete: DefaultHTTPDELETEResponses,
}

// HTTP Response for List operation
var HTTPListResponse = HTTPMethodsResponses{
	"LIST": DefaultHTTPGETResponses,
}

// GraphQL Types.

// A GraphQLQueryEndpoint specifies the network endpoint that serves a GraphQL query.
type GraphQLQueryEndpoint struct {
	Domain string `json:"domain"` // fully qualified domain name of the network endpoint
	Port   int    `json:"port"`   // service port
}

// A GraphQLQuery specifies a custom query available via GraphQL API.
// Each GraphQLQuery is self contained unit of the exposed custom query.
type GraphQLQuery struct {
	Name            string               `json:"name,omitempty"`            // query identifier
	ServiceEndpoint GraphQLQueryEndpoint `json:"servce_endpoint,omitempty"` // endpoint that serves this query
	Args            interface{}          `json:"args,omitempty"`            // custom graphql filters and arguments
	ApiType         GraphQlApiType       `json:"api_type,omitempty"`
}

// A GraphQLQuerySpec is a collection of GraphQLQuery.
// GraphQLQuerySpec provides a handle to represent and refer a collection of GraphQLQuery.
type GraphQLQuerySpec struct {
	Queries []GraphQLQuery `json:"queries"`
}

type GraphQlApiType int

const (
	GraphQLQueryApi GraphQlApiType = iota
	GetMetricsApi
)
