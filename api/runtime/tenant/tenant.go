package tenantruntime

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

// Tenant runtime spec.
type Tenant struct {
	nexus.Node

	// K8s Namespace associated with the Tenant.
	Namespace  string `json:"namespace,omitempty"`
	TenantName string `json:"tenantName"`

	// Dynamic / runtime attributes of the Tenant.
	Attributes Attributes `json:"attributes,omitempty"`

	// Domains associated with the Tenant.
	SaasDomainName    string `json:"saasDomainName,omitempty"`
	SaasApiDomainName string `json:"saasApiDomainName,omitempty"`
	M7Enabled         string `json:"m7Enabled,omitempty"`
	LicenseType       string `json:"licenseType,omitempty"`
	// Infrastructure associated with the Tenant.
	StreamName              string `json:"streamName,omitempty"`
	AwsS3Bucket             string `json:"awsS3Bucket,omitempty"`
	AwsKmsKeyId             string `json:"awsKmsKeyId,omitempty"`
	M7InstallationScheduled string `json:"m7InstallationScheduled,omitempty"`
	//removing cloud because it would be cluster object
	//Cloud       Cloud  `json:"cloud,omitempty"`

	// cosmos-release version
	ReleaseVersion string `json:"releaseVersion,omitempty"`
	// Runtime status of the Tenant.
	AppStatus TenantStatus `nexus:"status"`
}

// Cloud infra associated with the Tenant.
type Cloud struct {
	Provider string `json:"provider,omitempty"`
	Region   string `json:"region,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

// Runtime attributes associated with the Tenant.
type Attributes struct {
	Skus []string `json:"skus"`
}

// Tenant health spec.
type Health struct {
	Healthy                  bool   `json:"healthy,omitempty"`
	LasHealthCheckTimestamp  string `json:"last_health_check_timestamp,omitempty"`
	NextHealthCheckTimestamp string `json:"next_health_check_timestamp,omitempty"`
	Message                  string `json:"message,omitempty"`
}

// Runtime status associated with the Tenant.
type TenantStatus struct {

	// Applications currently installed in the Tenant.
	InstalledApplications ApplicationStatus `json:"installedApplications,omitempty"`
	ReleaseVersion        string                   `json:"releaseVersion,omitempty"`
	ReleaseStatus         string                   `json:"releaseStatus,omitempty"`
	PreviousRelease       string                   `json:"previousRelease,omitempty"`
}

type ApplicationStatus struct {
        NexusApps map[string]NexusApp `json:"nexusApps, omitempty"`
}

type NexusApp struct {
        OamApp      OamApp `json:"oamApp, omitempty"`
        State       string `json:"state, omitempty"`
        StateReason string `json:"stateReason, omitempty"`
}

type OamApp struct {
        Components map[string]ComponentDefinition `json:"components, omitempty"`
}

type ComponentDefinition struct {
        Name   string `json:"name, omitempty"`
        Sync   string `json:"sync, omitempty"`
        Health string `json:"health, omitempty"`
}
