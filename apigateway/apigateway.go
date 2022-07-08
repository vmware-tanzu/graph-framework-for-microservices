package apigateway

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/api/admin"
	authentication "golang-appnet.eng.vmware.com/nexus-sdk/api/authn"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

// ApiGateway holds all configuration relevant to a gateway in Nexus runtime.
type ApiGateway struct {
	nexus.Node

	// ProxyRules define a match condition and a corresponding upstream
	ProxyRules map[string]admin.ProxyRule `nexus:"child"`

	// Authentication config associated with this Gateway.
	Authn authentication.OIDC `nexus:"child"`
}
