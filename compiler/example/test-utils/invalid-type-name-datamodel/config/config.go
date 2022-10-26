package config

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
)

type Config struct {
	nexus.Node
	MyId          string
	FooResource Resource
}

type Resource struct {
	Name string
}
