package project

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/duplicated-uris-datamodel/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/duplicated-uris-datamodel/nexus"
)

var ProjectRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{root.Root}/project/{project.Project}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
	},
}

// nexus-rest-api-gen:ProjectRestAPISpec
type Project struct {
	nexus.SingletonNode
	Key    string
	Config config.Config `nexus:"child"`
}
