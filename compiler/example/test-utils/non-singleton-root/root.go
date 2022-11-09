package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/non-singleton-root/project"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Root struct {
	nexus.Node
	Project project.Project `nexus:"child"`
	IsRoot  IsRoot          // <--- to verify alias type
}

type IsRoot bool
