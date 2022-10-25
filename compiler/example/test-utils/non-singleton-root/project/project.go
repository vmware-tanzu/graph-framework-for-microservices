package project

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/non-singleton-root/config"
)

type Project struct {
	nexus.SingletonNode
	Key    string
	Config config.Config `nexus:"child"`
}
