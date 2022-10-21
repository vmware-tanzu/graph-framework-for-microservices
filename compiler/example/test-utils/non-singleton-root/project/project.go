package project

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/non-singleton-root/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/non-singleton-root/nexus"
)

type Project struct {
	nexus.SingletonNode
	Key    string
	Config config.Config `nexus:"child"`
}
