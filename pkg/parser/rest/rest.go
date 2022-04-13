package rest

import "net/http"

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

const DefaultHTTPErrorCode ResponseCode = http.StatusNotImplemented

var DefaultHTTPError = HTTPResponse{Description: http.StatusText(http.StatusNotImplemented)}
