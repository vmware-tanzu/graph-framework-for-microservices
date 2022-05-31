package config

import (
	extensions "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/api-extensions"
	authentication "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/authn"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/gateway"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
)

// Config parent's configuration created by user/product.
//
// The configuration is intent-driven and configuration can live
// if it's not ready to be consumed or enabled.
type Config struct {
	nexus.Node

	// Gateway configuration.
	Gateway gateway.Gateway `nexus:"child"`

	// API extensions configuration.
	ApiExtensions extensions.Extension `nexus:"child"`

	// Authentication configuration.
	AuthN map[string]authentication.OIDC `nexus:"child"`
}
