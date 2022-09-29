package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type Root struct {
	nexus.Node
	DisplayName string
	Config      config.Config `nexus:"child"`
	CustomBar   Bar
}

type Bar struct {
	Name string
}

type queryFilters struct {
	StartTime           string
	EndTime             string
	Interval            string
	IsServiceDeployment bool
	StartVal            int
}

var GeneralGraphQLQuerySpec = nexus.GraphQLQuerySpec{
	Queries: []nexus.GraphQLQuery{
		{
			Name: "query",
			ServiceEndpoint: nexus.GraphQLQueryEndpoint{
				Domain: "query-manager",
				Port:   6000,
			},
			Args: queryFilters{},
		},
	},
}
