package gns

import (
	service_group "gitlab.eng.vmware.com/nexus/compiler/example/datamodel/config/gns/service-group"
	"gitlab.eng.vmware.com/nexus/compiler/example/datamodel/config/policy"
	"gitlab.eng.vmware.com/nexus/compiler/example/datamodel/nexus"
)

type Gns struct {
	nexus.Node
	Domain                 string
	UseSharedGateway       bool
	Description            Description
	GnsServiceGroups       map[string]service_group.SvcGroup `nexus:"child"`
	GnsAccessControlPolicy policy.AccessControlPolicy        `nexus:"child"`
	Dns                    Dns                               `nexus:"link"`
	State                  GnsState                          `nexus:"status"`
}

type Description struct {
	Color     string
	Version   string
	ProjectId string
}

type Dns struct {
	nexus.Node
}

type GnsState struct {
	Working     bool
	Temperature int
}
