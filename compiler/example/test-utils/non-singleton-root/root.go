package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/non-singleton-root/nexus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/non-singleton-root/project"
)

type Root struct {
	nexus.Node
	Project project.Project `nexus:"child"`
	IsRoot  IsRoot          // <--- to verify alias type
}

type IsRoot bool
