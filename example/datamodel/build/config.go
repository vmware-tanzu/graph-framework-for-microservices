package build

import (
	"gitlab.eng.vmware.com/nexus/compiler/example/datamodel/nexus"
)

// BuildTestStruct struct is in build directory and should be ignored by parser
type BuildTestStruct struct {
	nexus.Node
	SomeInt int
}
