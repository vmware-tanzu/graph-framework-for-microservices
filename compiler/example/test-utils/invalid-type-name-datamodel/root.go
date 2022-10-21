package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/invalid-type-name-datamodel/config"
)

type Root struct {
	nexus.Node
	Config config.Config `nexus:"child"`
}
