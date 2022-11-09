package config

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var ConfigRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{root.Root}/project/{config.Config}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
	},
}

// nexus-rest-api-gen:ConfigRestAPISpec
type Config struct {
	nexus.SingletonNode
	Id string
}
