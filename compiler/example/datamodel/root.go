package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type Root struct {
	nexus.SingletonNode
	Config config.Config `nexus:"child"`
}

type NonNexusType struct {
	Test int
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
			Name: "QueryExample",
			ServiceEndpoint: nexus.GraphQLQueryEndpoint{
				Domain: "query-manager",
				Port:   6000,
			},
			Args: queryFilters{},
		},
	},
}
