package root

import (
	"helloworld/config"
	"helloworld/inventory"
	"helloworld/nexus"
	"helloworld/runtime"
)

type Root struct {
	nexus.Node
	MyInt     int
	Config    config.Config       `nexus:"child"`
	Runtime   runtime.Runtime     `nexus:"child"`
	Inventory inventory.Inventory `nexus:"child"`
}
