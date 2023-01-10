package engineering

import (
	"net/http"
	"vmware/nexus"
	"vmware/pkg/role"
)

var DevRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}/mgr/{management.Mgr}/dev/{engineering.Dev}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/devs",
			QueryParams: []string{
				"management.Mgr",
			},
			Methods: nexus.HTTPListResponse,
		},
		{
			Uri: "/root/{orgchart.Root}/dev/{engineering.Dev}",
			QueryParams: []string{
				"management.Mgr",
			},
			Methods: nexus.HTTPMethodsResponses{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
				http.MethodPut: nexus.HTTPCodesResponse{
					http.StatusOK: nexus.HTTPResponse{Description: "Example response"},
				},
			},
		},
	},
}

// nexus-rest-api-gen:DevRestAPISpec
type Dev struct {
	nexus.Node

	Name       string
	EmployeeID int

	Role role.Employee `nexus:"link"`
}

var OperationsRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}/mgr/{management.Mgr}/operations/{engineering.Operations}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/operations",
			QueryParams: []string{
				"management.Mgr",
			},
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:OperationsRestAPISpec
type Operations struct {
	nexus.Node

	Name       string
	EmployeeID int

	Role role.Employee `nexus:"link"`
}
