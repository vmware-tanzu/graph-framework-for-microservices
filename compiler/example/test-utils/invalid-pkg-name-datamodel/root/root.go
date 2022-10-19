package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type Root struct {
	nexus.Node
	Id string
}
