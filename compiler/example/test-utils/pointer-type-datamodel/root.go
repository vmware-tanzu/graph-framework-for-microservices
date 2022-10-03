package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/pointer-type-datamodel/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type Root struct {
	nexus.Node
	Config *config.Config `nexus:"child"` // not allowed
}
