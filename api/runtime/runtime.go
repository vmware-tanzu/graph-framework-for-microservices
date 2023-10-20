package runtime

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	tenantruntime "github.com/vmware-tanzu/graph-framework-for-microservices/api/runtime/tenant"
)

// Runtime tree.
type Runtime struct {
	nexus.SingletonNode

	// Tenant runtime spec.
	Tenant tenantruntime.Tenant `nexus:"children"  json:"tenant,omitempty"`
}
