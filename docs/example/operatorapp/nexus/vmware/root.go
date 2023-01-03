package orgchart

import (
	"vmware/nexus"
	"vmware/pkg/management"
	"vmware/pkg/role"
)

var RootRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/roots",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:RootRestAPISpec
// Datamodel graph root
type Root struct {
	nexus.SingletonNode

	CEO           management.Leader `nexus:"child"`
	ExecutiveRole role.Executive    `nexus:"child"`
	EmployeeRole  role.Employee     `nexus:"child"`
}
