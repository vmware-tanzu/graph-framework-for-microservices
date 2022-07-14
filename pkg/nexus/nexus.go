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
	QueryParams []string             `json:"queryParams"`
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
	DefaultHTTPErrorCode: DefaultHTTPError,
}

// Default HTTP methods and responses.
var DefaultHTTPMethodsResponses = HTTPMethodsResponses{
	http.MethodGet:    DefaultHTTPGETResponses,
	http.MethodPut:    DefaultHTTPPUTResponses,
	http.MethodDelete: DefaultHTTPDELETEResponses,
}
