package route

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type Service struct {
	Name   string
	Port   int32 `yaml:"omitempty"`
	Scheme string
}

type ResourceConfig struct {
	Name string
}

// Route specifies configuration to extend the API gateway with
// custom APIs.
type Route struct {
	nexus.Node

	Uri      string
	Service  Service
	Resource ResourceConfig
}
