package project

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/non-singleton-root/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/non-singleton-root/nexus"
)

type Project struct {
	nexus.SingletonNode
	Key    string
	Config config.Config `nexus:"child"`
}
