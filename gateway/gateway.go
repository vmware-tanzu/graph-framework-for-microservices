package gateway

import (
	authentication "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/authn"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
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
