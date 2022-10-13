package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/invalid-singleton-child/nexus"
)

type Root struct {
	nexus.Node
	Id     string
	Config Cfg `nexus:"children"`
}

type Cfg struct {
	nexus.SingletonNode
}
