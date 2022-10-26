package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
)

type Root struct {
	nexus.Node
	MyId Id
}

type Id struct {
	Field1 string
}
