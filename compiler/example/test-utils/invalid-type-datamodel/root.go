package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/invalid-type-datamodel/config"
)

type Root struct {
	nexus.Node
	Config []config.Config `nexus:"child"`
}
