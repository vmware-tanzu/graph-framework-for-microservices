package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/pointer-type-datamodel/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
)

type Root struct {
	nexus.Node
	Config *config.Config `nexus:"child"` // not allowed
}
