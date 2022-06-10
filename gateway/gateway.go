package gateway

import (
	authentication "golang-appnet.eng.vmware.com/nexus-sdk/api/authn"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type GatewayConfig struct{}

// Gateway holds all configuration relevant to a gateway in Nexus runtime.
type Gateway struct {
	nexus.Node

	// Configuration.
	Config GatewayConfig

	// Authentication config assocciated with this Gateway.
	Authn authentication.OIDC `nexus:"link"`
}
