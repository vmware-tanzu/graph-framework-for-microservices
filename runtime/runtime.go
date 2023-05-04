package runtime

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	tenantruntime "golang-appnet.eng.vmware.com/nexus-sdk/api/runtime/tenant"
)

// Runtime tree.
type Runtime struct {
	nexus.SingletonNode

	// Tenant runtime spec.
	Tenant tenantruntime.Tenant `nexus:"children"  json:"tenant,omitempty"`
}
