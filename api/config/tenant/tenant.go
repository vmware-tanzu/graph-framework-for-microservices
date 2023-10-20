package tenantconfig

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"

	"github.com/vmware-tanzu/graph-framework-for-microservices/api/common"
)

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Tenant Spec.
type Tenant struct {
	nexus.Node

	Name string `json:"name"`

	// Configuration
	DNSSuffix         string `json:"dns_suffix,omitempty"`
	SkipSaasTlsVerify bool   `json:"skip_saas_tls_verify,omitempty"`
	InstallTenant     bool   `json:"install_tenant,omitempty"`
	InstallClient     bool   `json:"install_client,omitempty"`

	// Order Info
	OrderId      string   `json:"order_id,omitempty"`
	Skus         []string `json:"skus"`
	FeatureFlags []string `json:"feature_flags,omitempty"`

	// Custom labels to be associated with this tenant.
	Labels []Label `json:"labels,omitempty"`

	// Status
	Status TenantStatus `nexus:"status"  json:"status,omitempty"`
}

type Provisioning struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// Tenant Status.
type TenantStatus struct {
	Provisioning Provisioning `json:"provisioning,omitempty"`
}

// Tenant Policy.
type Policy struct {
	nexus.Node

	// Static applications that need to be installed on a tenant.
	// Static applications will be installed on the tenant even if the tenant
	// attributes do not match with application spec.
	StaticApplications []common.Application `json:"static_applications,omitempty"`

	// Applications that need to be pinned to a specified version on a tenant.
	// This will include versions for static and dynamic applications.
	PinApplications []common.Application `json:"pin_applications,omitempty"`

	// Disable dynamic application matching and scheduling on this tenant.
	DynamicAppSchedulingDisable bool `json:"dynamic_app_scheduling_disable,omitempty"`

	// runtime policy
	DisableProvisioning         bool `json:"disable_provisioning,omitempty"`
	DisableAutoScaling          bool `json:"disable_auto_scaling,omitempty"`
	DisableAppClusterOnboarding bool `json:"disable_app_cluster_onboarding,omitempty"`

	// upgrade
	DisableUpgrade     bool `json:"disable_upgrade,omitempty"`
	OnFailureDowngrade bool `json:"on_failure_downgrade,omitempty"`
}
