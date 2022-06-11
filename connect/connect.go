package connect

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

// Nexus Connect configuration.
type Connect struct {
	nexus.Node

	Endpoints map[string]NexusEndpoint `nexus:"child"`

	ReplicationConfig map[string]ReplicationConfig `nexus:"child"`
}

// NexusEndpoint identifies a Nexus Runtime endpoint.
type NexusEndpoint struct {
	nexus.Node

	Host string
	Port string
	Cert string
}

// ReplicationConfig defines a replication request/instance.
type ReplicationConfig struct {
	nexus.Node

	// Pointer to a remote Nexus Runtime.
	RemoteEndpoint NexusEndpoint `nexus:"link"`

	// Credentials to access the remote endpoint.
	AccessToken string

	// Source of the replication.
	Source ReplicationObject `nexus:"child"`

	// Destination of the replication.
	Destination ReplicationObject `nexus:"child"`
}

// ReplicationObject identifies a resource in the Nexus Runtime.
type ReplicationObject struct {
	nexus.Node

	// Object identity.
	Group string `json:"group" yaml:"group"`
	Kind  string `json:"kind" yaml:"kind"`
	Name  string `json:"name" yaml:"name"`

	// Identifies if the object is available in local runtime.
	// If false, the object is available in the remote runtime.
	LocalRuntime bool

	// If true, the hierarchy of the object is relevant for replication.
	Hierarchical bool

	// Hierarchical path prefix of the object.
	// This is relevant is the object has to be considered in the context of its hierarchy.
	Hierarchy Hierarchy
}

// Hierarchy identifies a hierarchy of an object in Nexus Runtime.
type Hierarchy struct {
	Path string
}
