package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/pointer-type-datamodel/config"
)

type Root struct {
	nexus.Node
	Config *config.Config `nexus:"child"` // not allowed
}
