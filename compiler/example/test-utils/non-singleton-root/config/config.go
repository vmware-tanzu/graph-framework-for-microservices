package config

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
)

type Config struct {
	nexus.SingletonNode
	MyId string
}
