package apigateway

import (
	authentication "golang-appnet.eng.vmware.com/nexus-sdk/api/authn"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type ApiGatewayConfig struct{}

// ApiGateway holds all configuration relevant to a gateway in Nexus runtime.
type ApiGateway struct {
	nexus.Node

	// Configuration.
	Config ApiGatewayConfig

	// Authentication config associated with this Gateway.
	Authn authentication.OIDC `nexus:"child"`
}
