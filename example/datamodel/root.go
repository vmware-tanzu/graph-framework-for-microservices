package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

type Root struct {
	nexus.SingletonNode
	Config config.Config `nexus:"child"`
}

type NonNexusType struct {
	Test int
}
