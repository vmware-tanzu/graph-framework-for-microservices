package config

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/example/test-utils/non-singleton-root/nexus"
)

type Config struct {
	nexus.SingletonNode
	Id string
}
