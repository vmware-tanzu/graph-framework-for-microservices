package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/invalid-singleton-child/nexus"
)

type Root struct {
	nexus.Node
	Id     string
	Config Cfg `nexus:"children"`
}

type Cfg struct {
	nexus.SingletonNode
}
