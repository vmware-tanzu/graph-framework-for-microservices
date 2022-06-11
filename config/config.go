package config

import (
	extensions "golang-appnet.eng.vmware.com/nexus-sdk/api/api-extensions"
	authentication "golang-appnet.eng.vmware.com/nexus-sdk/api/authn"
	"golang-appnet.eng.vmware.com/nexus-sdk/api/connect"
	"golang-appnet.eng.vmware.com/nexus-sdk/api/gateway"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
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

	// Nexus Connect configuration.
	Connect connect.Connect `nexus:"child"`
}
