package utils

import (
	"fmt"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	CRDTypeToParentHierarchy      = make(map[string][]string)
	crdTypeToParentHierarchyMutex = &sync.Mutex{}

	ReplicationEnabledNode      = make(map[ReplicationObject]dynamic.Interface)
	replicationEnabledNodeMutex = &sync.Mutex{}

	ReplicatedNodes      = make(map[string]dynamic.Interface)
	replicatedNodesMutex = &sync.Mutex{}

	CRDTypeToChildren      = make(map[string]Children)
	crdTypeToChildrenMutex = &sync.Mutex{}
)

const (
	Upsert EventType = "Upsert"
	Delete EventType = "Delete"

	DisplayNameKey = "nexus/display_name"

	// Nexus-Connect DM CRDs.
	NexusCRD             = "nexuses.api.nexus.org"
	ConfigCRD            = "configs.config.nexus.org"
	ConnectCRD           = "connects.connect.nexus.org"
	ReplicationConfigCRD = "replicationconfigs.connect.nexus.org"
	NexusEndpointCRD     = "nexusendpoints.connect.nexus.org"

	Update = "UPDATE"
	Create = "CREATE"
)

type EventType string

type Children map[string]NodeHelperChild

type ReplicationConfig struct {
	AccessToken    string            `json:"accessToken"`
	RemoteEndpoint Child             `json:"remoteEndpointGvk"`
	Source         ReplicationObject `json:"source"`
}

type Child struct {
	Group string `json:"group"`
	Kind  string `json:"kind"`
	Name  string `json:"name"`
}

type MultipleChildren map[string]Child

type NexusEndpoint struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Cert string `json:"cert"`
}

type ReplicationObject struct {
	Group string `json:"group"`
	Kind  string `json:"kind"`
	Name  string `json:"name"`
}

type NexusAnnotation struct {
	Name      string                     `json:"name,omitempty"`
	Hierarchy []string                   `json:"hierarchy,omitempty"`
	Children  map[string]NodeHelperChild `json:"children,omitempty"`
}

type NodeHelperChild struct {
	FieldName    string `json:"fieldName"`
	FieldNameGvk string `json:"fieldNameGvk"`
	IsNamed      bool   `json:"isNamed"`
}

func ConstructMapCRDTypeToParentHierarchy(eventType EventType, crdType string, hierarchy []string) {
	crdTypeToParentHierarchyMutex.Lock()
	defer crdTypeToParentHierarchyMutex.Unlock()

	if eventType == Delete {
		delete(CRDTypeToParentHierarchy, crdType)
	}

	CRDTypeToParentHierarchy[crdType] = hierarchy
}

func ConstructMapCRDTypeToChildren(eventType EventType, crdType string, children Children) {
	crdTypeToChildrenMutex.Lock()
	defer crdTypeToChildrenMutex.Unlock()

	if eventType == Delete {
		delete(CRDTypeToChildren, crdType)
	}

	CRDTypeToChildren[crdType] = children
}

func ConstructMapReplicationEnabledNode(repObj ReplicationObject, client dynamic.Interface) {
	replicationEnabledNodeMutex.Lock()
	defer replicationEnabledNodeMutex.Unlock()

	ReplicationEnabledNode[repObj] = client
}

func ConstructMapWithReplicatedNodes(crdType, replicatedNode string, client dynamic.Interface) {
	replicatedNodesMutex.Lock()
	defer replicatedNodesMutex.Unlock()

	rnKey := ConstructReplicatedNodeKey(crdType, replicatedNode)
	ReplicatedNodes[rnKey] = client
}

func ConstructReplicatedNodeKey(crdType, replicatedNode string) string {
	return fmt.Sprintf("%s:%s", crdType, replicatedNode)
}

func GetCrdType(kind, groupName string) string {
	return GetGroupResourceName(kind) + "." + groupName // eg roots.root.helloworld.com
}

func GetGroupResourceName(nodeName string) string {
	return strings.ToLower(ToPlural(nodeName)) // eg roots
}

func GetGVRFromCrdType(crdType string) schema.GroupVersionResource {
	parts := strings.Split(crdType, ".")
	return schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}
}

func GetParentHierarchy(parents []string, labels map[string]string) (hierarchy string) {
	for _, parent := range parents {
		for key, val := range labels {
			if parent == key {
				hierarchy += key + "/" + val + "/"
			}
		}
	}
	return
}

func GetParentLabels(parents []string, labels map[string]string) string {
	var parentLabels string
	immediateParent := parents[len(parents)-1]
	for _, parent := range parents[:len(parents)-1] {
		for key, val := range labels {
			if parent == key {
				parentLabels += key + "=" + val + ","
			}
		}
	}
	return parentLabels + fmt.Sprintf("%s=%s", DisplayNameKey, labels[immediateParent])
}

func GetNodeHierarchy(parents []string, labels map[string]string, node string) string {
	parentPrefix := GetParentHierarchy(parents, labels)
	return parentPrefix + fmt.Sprintf("%s/%s", node, labels[DisplayNameKey])
}

func GetNodeLabels(parents []string, labels map[string]string, node string) string {
	var nodeLabels string
	for _, parent := range parents {
		for key, val := range labels {
			if parent == key {
				nodeLabels += key + "=" + val + ","
			}
		}
	}
	nodeLabels += node + "=" + labels[DisplayNameKey]
	return nodeLabels
}

func DeleteChildGvkFields(fields map[string]interface{}, children map[string]NodeHelperChild) {
	for _, val := range children {
		delete(fields, val.FieldNameGvk)
	}
}

func NexusDatamodelCRDs(crdType string) bool {
	if crdType == NexusCRD ||
		crdType == ConfigCRD ||
		crdType == ConnectCRD ||
		crdType == ReplicationConfigCRD ||
		crdType == NexusEndpointCRD {
		return true
	}
	return false
}
