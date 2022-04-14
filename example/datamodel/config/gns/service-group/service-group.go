package servicegroup

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

type SvcGroup struct {
	nexus.Node
	DisplayName string
	Description string
	Color       string
	//Services    map[string]core_v1.Service `nexus:"link"`
}
