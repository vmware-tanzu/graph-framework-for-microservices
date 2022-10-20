package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/duplicated-uris-datamodel/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/duplicated-uris-datamodel/project"
)

type Root struct {
	nexus.Node
	Project project.Project `nexus:"child"`
}
