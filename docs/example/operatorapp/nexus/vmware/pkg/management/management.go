package management

import (
	"vmware/nexus"
	"vmware/pkg/engineering"
	"vmware/pkg/hr"
	"vmware/pkg/role"
)

var LeaderRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/leader",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:LeaderRestAPISpec
type Leader struct {
	nexus.SingletonNode

	Designation string
	Name        string
	EmployeeID  int

	EngManagers Mgr               `nexus:"children"`
	HR          hr.HumanResources `nexus:"child"`
	Role        role.Executive    `nexus:"link"`
	Status      LeaderState       `nexus:"status"`
}

type LeaderState struct {
	IsOnVacations            bool
	DaysLeftToEndOfVacations int
}

var MgrRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}/mgr/{management.Mgr}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/mgrs",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:MgrRestAPISpec
type Mgr struct {
	nexus.Node

	Name       string
	EmployeeID int

	Developers engineering.Dev        `nexus:"child"`
	Ops        engineering.Operations `nexus:"child"`

	Role role.Employee `nexus:"link"`
}
