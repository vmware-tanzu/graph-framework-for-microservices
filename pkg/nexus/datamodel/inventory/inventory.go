package inventory

import (
	"helloworld/nexus"
)

type Inventory struct {
	nexus.Node
	InventoryId int
}
