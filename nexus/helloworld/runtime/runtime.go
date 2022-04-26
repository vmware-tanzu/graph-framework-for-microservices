package runtime

import (
	"helloworld/nexus"
)

type Runtime struct {
	nexus.Node
	MyRuntimeInt int
}
