package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/nexus"
)

type Root struct {
	nexus.SingletonNode
	Config config.Config `nexus:"child"`
}

type NonNexusType struct {
	Test int
}
