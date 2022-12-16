package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/array-type-child/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Root struct {
	nexus.Node
	Config []config.Config `nexus:"child"` // not allowed
}
