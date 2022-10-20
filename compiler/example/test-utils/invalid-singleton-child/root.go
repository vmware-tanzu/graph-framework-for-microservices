package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
)

type Root struct {
	nexus.Node
	Id     string
	Config Cfg `nexus:"children"`
}

type Cfg struct {
	nexus.SingletonNode
}
