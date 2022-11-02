package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/pointer-type-datamodel/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Root struct {
	nexus.Node
	Config *config.Config `nexus:"child"` // not allowed
}
