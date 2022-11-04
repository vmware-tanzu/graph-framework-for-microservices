package config

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
)

var ConfigRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/root/{root.Root}/config/{config.Config}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
	},
}

// nexus-rest-api-gen:ConfigRestAPISpec
type Config struct {
	nexus.SingletonNode
	FieldX        string
	FieldY        int64
	MyStructField MyStruct
}

type MyStruct struct {
	TempFiledA string
	TempFiledB string
}
