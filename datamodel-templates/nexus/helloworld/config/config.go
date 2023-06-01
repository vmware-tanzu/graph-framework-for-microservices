package config

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type Config struct {
	nexus.Node
	ExampleStr string
}
