package runtime

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type Runtime struct {
	nexus.Node
	MyRuntimeInt int
}
