package config

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/api/apigateway"
	tenantconfig "github.com/vmware-tanzu/graph-framework-for-microservices/api/config/tenant"
	"github.com/vmware-tanzu/graph-framework-for-microservices/api/config/user"
	"github.com/vmware-tanzu/graph-framework-for-microservices/api/connect"
	"github.com/vmware-tanzu/graph-framework-for-microservices/api/route"
)

// Config holds the Nexus configuration.
// Configuration in Nexus is intent-driven.
type Config struct {
	nexus.Node

	// Gateway configuration.
	ApiGateway apigateway.ApiGateway `nexus:"child"`

	// API extensions configuration.
	Routes route.Route `nexus:"children"`

	// Nexus Connect configuration.
	Connect      connect.Connect     `nexus:"child"`
	Tenant       tenantconfig.Tenant `nexus:"children" json:"tenant,omitempty"`
	TenantPolicy tenantconfig.Policy `nexus:"children" json:"tenant_policy,omitempty"`

	User user.User `nexus:"children" json:"user,omitempty"`
}
