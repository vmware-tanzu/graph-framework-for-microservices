package build

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

// BuildTestStruct struct is in build directory and should be ignored by parser
type BuildTestStruct struct {
	nexus.Node
	SomeInt int
}