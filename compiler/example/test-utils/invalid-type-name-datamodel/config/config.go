package config

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Config struct {
	nexus.Node
	Id          string
	FooResource Resource
}

type Resource struct {
	Name string
}
