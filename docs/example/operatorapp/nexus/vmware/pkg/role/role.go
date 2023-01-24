package role

import (
	"vmware/nexus"
)

var ExecutiveRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}/executive/{role.Executive}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/executives",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:ExecutiveRestAPISpec
type Executive struct {
	nexus.Node
}

var EmployeeRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}/employee/{role.Employee}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/employees",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:EmployeeRestAPISpec
type Employee struct {
	nexus.Node
}
