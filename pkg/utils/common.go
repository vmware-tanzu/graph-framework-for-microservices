package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Upsert EventType = "Upsert"
	Delete EventType = "Delete"

	RemoteEndpointHost = "REMOTE_ENDPOINT_HOST"
	RemoteEndpointPort = "REMOTE_ENDPOINT_PORT"
	RemoteEndpointCert = "REMOTE_ENDPOINT_CERT"

	DisplayNameKey           = "nexus/display_name"
	DeploymentName           = "DEPLOYMENT_NAME"
	StatusReplication        = "STATUS_REPLICATION"
	StatusEnabled            = "ENABLED"
	NexusReplicationManager  = "nexus-replication-manager"
	NexusReplicationResource = "nexus-replication-resource"
	secretNS                 = "SECRET_NS"
	secretName               = "SECRET_NAME"

	// Nexus-Connect DM CRDs.
	ReplicationConfigCRD = "replicationconfigs.connect.nexus.org"
	NexusEndpointCRD     = "nexusendpoints.connect.nexus.org"

	// CRD Version.
	V1Version = "v1"

	Update = "UPDATE"
	Create = "CREATE"
)

type EventType string

type Children map[string]NodeHelperChild
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

type ResourceAnnotation struct {
	GVR  schema.GroupVersionResource
	Name string
}

func GetCrdType(kind, groupName string) string {
	return GetGroupResourceName(kind) + "." + groupName // eg roots.root.helloworld.com
}

func GetGroupResourceName(nodeName string) string {
	return strings.ToLower(ToPlural(nodeName)) // eg roots
}

func GetGVRFromCrdType(crdType, crdVersion string) schema.GroupVersionResource {
	parts := strings.Split(crdType, ".")
	return schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  crdVersion,
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

func GetReplicationObject(group, kind, name string) ReplicationObject {
	return ReplicationObject{
		Group: group,
		Kind:  kind,
		Name:  name,
	}
}

func DeleteChildGvkFields(fields map[string]interface{}, children map[string]NodeHelperChild) {
	for _, val := range children {
		delete(fields, val.FieldNameGvk)
	}
}

func NexusDatamodelCRDs(group string) bool {
	return strings.Contains(group, "nexus.org")
}

func GenerateAnnotations(annotations map[string]string, gvr schema.GroupVersionResource, name string) map[string]string {
	crInfoInByte, err := json.Marshal(ResourceAnnotation{
		GVR:  gvr,
		Name: name,
	})
	if err != nil {
		log.Errorf("Error marshaling Source CR %v info %v", name, err)
		return nil
	}

	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[NexusReplicationManager] = os.Getenv(DeploymentName)
	annotations[NexusReplicationResource] = string(crInfoInByte)

	return annotations
}

func GetDestinationGvrAndKind(destination ReplicationDestination, gvr schema.GroupVersionResource,
	kind string) (schema.GroupVersionResource, string) {

	destGvr := gvr
	destKind := kind
	if destination.ObjectType != nil && !destination.IsChild {
		destGvr = schema.GroupVersionResource{Group: destination.ObjectType.Group,
			Version:  destination.ObjectType.Version,
			Resource: GetGroupResourceName(destination.ObjectType.Kind)}
		destKind = destination.Kind
	}
	return destGvr, destKind
}
