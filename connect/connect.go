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

type ReplicationStatusEndpoint string

const (
	Source      ReplicationStatusEndpoint = "Source"
	Destination ReplicationStatusEndpoint = "Destination"
)

type SourceKind string

const (
	Object SourceKind = "Object"
	Type   SourceKind = "Type"
)

type KVP struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Hierarchy identifies a hierarchy of an object in Nexus Runtime.
type Hierarchy struct {
	Labels []KVP `json:"labels"`
}

// ObjectType identifies a type of objects to be replicated
type ObjectType struct {
	Group   string `json:"group"`
	Kind    string `json:"kind"`
	Version string `json:"version"`
}

// ReplicateObject identifies an object in datamodel to be replicated
type SourceObject struct {
	// Object type.
	ObjectType `json:"objectType"`

	// Object name.
	Name string `json:"name"`

	// If true, the hierarchy of the object is relevant for replication.
	Hierarchical bool `json:"hierarchical"`

	// Hierarchical path prefix of the object.
	// This is relevant if the object has to be considered in the context of its hierarchy.
	// Ignored if value of field Hierarchical is false.
	Hierarchy Hierarchy `json:"hierarchy,omitempty"`
}

// ReplicationSource identifies either a single object or all objects of a type
// that are to be replicated to a destination endpoint.
type ReplicationSource struct {
	Kind SourceKind `json:"kind"`

	// Relevant if all objects of a Type are to be replicated.
	Type ObjectType `json:"type,omitempty"`

	// Relevant if a specific object (and optionally its children) are to be replicated.
	Object SourceObject `json:"object,omitempty"`
}

// ReplicationDestination specifies the attributes with which objects are to be created
// in the destination endpoint.
type ReplicationDestination struct {
	// If true, the object will be replicated into the specified hierarchy.
	Hierarchical bool `json:"hierarchical"`

	// Hierarchy into which the object has to be replicated.
	Hierarchy Hierarchy `json:"hierarchy,omitempty"`

	// If specified, the replicated object will be scoped to this namespace.
	Namespace string `json:"namespace,omitempty"`
}

// ReplicationConfig defines a replication request/instance.
type ReplicationConfig struct {
	nexus.Node

	// Pointer to a remote Nexus Runtime.
	RemoteEndpoint NexusEndpoint `nexus:"link"`

	// Credentials to access the remote endpoint.
	AccessToken string

	// Source of the replication.
	Source ReplicationSource

	// Destination of the replication.
	Destination ReplicationDestination

	// Endpoint in which status of replication should be captured.
	// Status can be captured on the corresponding object in source or destination endpoint.
	// This allows for capturing of status at the endpoint where status is being watched for.
	StatusEndpoint ReplicationStatusEndpoint
}
