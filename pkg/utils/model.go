package utils

import "k8s.io/client-go/dynamic"

type SourceKind string

const (
	Object SourceKind = "Object"
	Type   SourceKind = "Type"
)

type Link struct {
	Group string `json:"group"`
	Kind  string `json:"kind"`
	Name  string `json:"name"`
}

type ReplicationConfig struct {
	AccessToken    string                 `json:"accessToken"`
	Source         ReplicationSource      `json:"source"`
	Destination    ReplicationDestination `json:"destination"`
	RemoteEndpoint Link                   `json:"remoteEndpointGvk"`
}

type ReplicationConfigSpec struct {
	Client      dynamic.Interface
	Source      ReplicationSource
	Destination ReplicationDestination
}

type ReplicationSource struct {
	Kind   SourceKind   `json:"kind"`
	Type   ObjectType   `json:"type"`
	Object SourceObject `json:"object"`
}

type ObjectType struct {
	Group   string `json:"group"`
	Kind    string `json:"kind"`
	Version string `json:"version"`
}

type NexusEndpoint struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Cert string `json:"cert"`
}

type SourceObject struct {
	ObjectType   `json:"objectType"`
	Name         string    `json:"name"`
	Hierarchical bool      `json:"hierarchical"`
	Hierarchy    Hierarchy `json:"hierarchy"`
}

type ReplicationDestination struct {
	Hierarchical bool      `json:"hierarchical"`
	Hierarchy    Hierarchy `json:"hierarchy"`
	Namespace    string    `json:"namespace"`
	*ObjectType  `json:"objectType"`
	IsChild      bool `json:"-"`
}

type KVP struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Hierarchy struct {
	Labels []KVP `json:"labels"`
}
