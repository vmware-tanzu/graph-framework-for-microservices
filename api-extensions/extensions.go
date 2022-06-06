package extensions

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
)

type Service struct {
	Name   string
	Port   int32 `yaml:"omitempty"`
	Scheme string
}

type ResourceConfig struct {
	Name string
}

// Extension speicifies configuration to extend the API gateway with
// custom API's.
type Extension struct {
	nexus.Node

	Uri      string
	Service  Service
	Resource ResourceConfig
}
