package servicegroup

import (
        "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type SvcGroup struct {
	nexus.Node
	DisplayName string
	Description string
	Color       string
	// TODO support links which are not nexus nodes https://jira.eng.vmware.com/browse/NPT-112
	//Services    map[string]core_v1.Service `nexus:"link"`
}
