package apigateway

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/api/admin"
	authentication "github.com/vmware-tanzu/graph-framework-for-microservices/api/authn"
	domain "github.com/vmware-tanzu/graph-framework-for-microservices/api/domain"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

// ApiGateway holds all configuration relevant to a gateway in Nexus runtime.
type ApiGateway struct {
	nexus.Node

	// ProxyRules define a match condition and a corresponding upstream
	ProxyRules admin.ProxyRule `nexus:"children"`

	// Authentication config associated with this Gateway.
	Authn authentication.OIDC `nexus:"child"`

	//Domain objects
	Cors domain.CORSConfig `nexus:"children"`
}
