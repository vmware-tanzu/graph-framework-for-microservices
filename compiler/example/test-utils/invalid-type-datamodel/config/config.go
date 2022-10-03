package config

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type Config struct {
	nexus.Node
	Id string
}
