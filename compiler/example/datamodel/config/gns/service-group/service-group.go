package servicegroup

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type SvcGroup struct {
	nexus.Node
	DisplayName string
	Description string
	Color       string
	// TODO support links which are not nexus nodes https://jira.eng.vmware.com/browse/NPT-112
	//Services    map[string]core_v1.Service `nexus:"link"`
}
