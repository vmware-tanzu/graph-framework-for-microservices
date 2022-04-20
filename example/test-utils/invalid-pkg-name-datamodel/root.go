package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
)

type Root struct {
	nexus.Node
	Name string
}
