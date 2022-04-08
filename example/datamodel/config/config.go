package config

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/policy"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

type Config struct {
	nexus.Node
	GNS    gns.Gns                               `nexus:"child"`
	policy map[string]policy.AccessControlPolicy `nexus:"child"`
}
