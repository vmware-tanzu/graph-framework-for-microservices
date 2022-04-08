package v1

import "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/openapi_generator/test_data/proto"

// +k8s:openapi-gen=true
type FooWrapper struct {
	Foo proto.Foo `json:"foo"`
}
