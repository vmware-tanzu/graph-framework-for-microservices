package project

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/group-name-with-hyphen-datamodel/config"
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
	Field1 string
	Field2 int64
	Config config.Config `nexus:"child"`
}
