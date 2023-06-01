package root

import (
	"helloworld/config"
	"helloworld/inventory"
	"helloworld/runtime"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type Root struct {
	nexus.Node
	MyInt     int
	Config    config.Config       `nexus:"child"`
	Runtime   runtime.Runtime     `nexus:"child"`
	Inventory inventory.Inventory `nexus:"child"`
}
