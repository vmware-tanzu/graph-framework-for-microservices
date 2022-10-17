package v1

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator/test_data/proto"
)

// +k8s:openapi-gen=true
type FooWrapper struct {
	Foo proto.Foo `json:"foo"`
}
