package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/map-type-child/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Root struct {
	nexus.Node
	Config map[string]config.Config `nexus:"child"` // not allowed
}
