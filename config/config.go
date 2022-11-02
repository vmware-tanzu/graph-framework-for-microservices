package config

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/api/apigateway"
	"golang-appnet.eng.vmware.com/nexus-sdk/api/connect"
	"golang-appnet.eng.vmware.com/nexus-sdk/api/route"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
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
	Connect connect.Connect `nexus:"child"`
}
