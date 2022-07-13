package inventory

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type Inventory struct {
	nexus.Node
	InventoryId int
}
