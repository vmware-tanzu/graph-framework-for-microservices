package authentication

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/nexus"
)

type IDPConfig struct {
	ClientId         string
	ClientSecret     string
	OAuthRedirectUrl string
}

// OIDC holds state/config associated with authentication.
//
// Nexus Runtime supports authentication function and the state
// associated with it is rooted on the OIDC node.
type OIDC struct {
	nexus.Node

	// IDP configuration.
	Config IDPConfig
}
