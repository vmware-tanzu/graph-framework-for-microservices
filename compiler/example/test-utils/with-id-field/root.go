package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
)

type Root struct {
	nexus.Node
	MyField MyField
}

type MyField struct {
	ID string
}
