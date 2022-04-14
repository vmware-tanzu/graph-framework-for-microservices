package rest

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
)

var HttpCodesResponsesMap = map[string]nexus.HTTPCodesResponse{
	"DefaultHTTPGETResponses":    nexus.DefaultHTTPGETResponses,
	"DefaultHTTPPUTResponses":    nexus.DefaultHTTPPUTResponses,
	"DefaultHTTPDELETEResponses": nexus.DefaultHTTPDELETEResponses,
}

var HttpMethodsResponsesMap = map[string]nexus.HTTPMethodsResponses{
	"DefaultHTTPMethodsResponses": nexus.DefaultHTTPMethodsResponses,
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
