package build

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

// BuildTestStruct struct is in build directory and should be ignored by parser
type BuildTestStruct struct {
	nexus.Node
	SomeInt int
}
