package rest

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var HttpCodesResponsesMap = map[string]nexus.HTTPCodesResponse{
	"DefaultHTTPGETResponses":    nexus.DefaultHTTPGETResponses,
	"DefaultHTTPPUTResponses":    nexus.DefaultHTTPPUTResponses,
	"DefaultHTTPDELETEResponses": nexus.DefaultHTTPDELETEResponses,
}

var HttpMethodsResponsesMap = map[string]nexus.HTTPMethodsResponses{
	"DefaultHTTPMethodsResponses": nexus.DefaultHTTPMethodsResponses,
	"HTTPListResponse":            nexus.HTTPListResponse,
}

func ParseResponses(pkgs parser.Packages) (map[string]nexus.HTTPMethodsResponses, map[string]nexus.HTTPCodesResponse) {
	methods := HttpMethodsResponsesMap
	codes := HttpCodesResponsesMap

	// Iterate through packages to get all HTTP Codes
	for _, pkg := range pkgs {
		codes = GetHttpCodesResponses(pkg, codes)
	}

	// Iterate through packages to get all HTTP Methods
	for _, pkg := range pkgs {
		methods = GetHttpMethodsResponses(pkg, methods, codes)
	}

	return methods, codes
}
