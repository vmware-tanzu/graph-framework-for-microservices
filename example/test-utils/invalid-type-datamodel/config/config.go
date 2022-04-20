package config

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
)

type Config struct {
	nexus.Node
	Id string
}
