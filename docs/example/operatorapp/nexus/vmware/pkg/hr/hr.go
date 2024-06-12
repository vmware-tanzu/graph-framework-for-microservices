package hr

import (
	"vmware/nexus"
	"vmware/pkg/role"
)

var HumanResourcesRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}/hr/{hr.HumanResources}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/hr/{hr.HumanResources}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/humanresources",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:HumanResourcesRestAPISpec
type HumanResources struct {
	nexus.Node

	Name       string
	EmployeeID int

	Role role.Employee `nexus:"link"`
}
