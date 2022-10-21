package config

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/test-utils/non-singleton-root/nexus"
)

type Config struct {
	nexus.SingletonNode
	Id string
}
