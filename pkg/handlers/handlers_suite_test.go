package handlers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"connector/pkg/utils"
)

const (
	ApiVersion = "config.mazinger.com/v1"
	Group      = "config.mazinger.com"

	// CRD Types
	Root                  = "roots.config.mazinger.com"
	Project               = "projects.config.mazinger.com"
	Config                = "configs.config.mazinger.com"
	ApiCollaborationSpace = "apicollaborationspaces.config.mazinger.com"
	ApiDevSpace           = "apidevspaces.config.mazinger.com"

	// Object Kind
	RootKind    = "Root"
	ProjectKind = "Project"
	ConfigKind  = "Config"
	AcKind      = "ApiCollaborationSpace"
	AdKind      = "ApiDevSpace"
)

var (
	apicollaborationspace = schema.GroupVersionResource{Group: "config.mazinger.com", Version: "v1", Resource: "apicollaborationspaces"}
	apidevspace           = schema.GroupVersionResource{Group: "config.mazinger.com", Version: "v1", Resource: "apidevspaces"}
	deployment            = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
)

func getObject(name, kind, specVal string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": ApiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"example": specVal,
			},
		},
	}
}

func getParentObject(name, kind string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": ApiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
				"labels": map[string]interface{}{
					Root:                 "root",
					Project:              "project",
					Config:               "config",
					"nexus/display_name": name,
				},
			},
			"spec": map[string]interface{}{
				"example":        "example",
				"apiDevSpaceGvk": "value",
			},
		},
	}
}

func getChildObject(name, kind string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": ApiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
				"labels": map[string]interface{}{
					Root:                  "root",
					Project:               "project",
					Config:                "config",
					ApiCollaborationSpace: "foo",
					"nexus/display_name":  name,
				},
			},
			"spec": map[string]interface{}{
				"example": "example",
			},
		},
	}
}

func getReplicatedObject(name, kind string, status map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": ApiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
				"annotations": map[string]interface{}{
					utils.NexusReplicationManager:  "connector",
					utils.NexusReplicationResource: `{"GVR":{"Group":"config.mazinger.com","Version":"v1","Resource":"configs"},"Name":"New"}`,
				},
			},
			"spec": map[string]interface{}{
				"example": "example",
			},
			"status": status,
		},
	}
}

func getNonHierarchicalDestConfig() utils.ReplicationDestination {
	return utils.ReplicationDestination{
		Hierarchical: false,
	}
}

func getNonHierarchicalSourceConfig() utils.ReplicationSource {
	return utils.ReplicationSource{
		Kind: utils.Object,
		Object: utils.SourceObject{
			ObjectType: utils.ObjectType{
				Group:   Group,
				Version: "v1",
				Kind:    AcKind,
			},
			Hierarchical: false,
		},
	}
}

func getHierarchicalSourceConfig() utils.ReplicationSource {
	return utils.ReplicationSource{
		Kind: utils.Object,
		Object: utils.SourceObject{
			ObjectType: utils.ObjectType{
				Group:   Group,
				Version: "v1",
				Kind:    AcKind,
			},
			Hierarchical: true,
			Hierarchy: utils.Hierarchy{
				Labels: []utils.KVP{
					{
						Key:   Root,
						Value: "root",
					},
					{
						Key:   Project,
						Value: "project",
					},
					{
						Key:   Config,
						Value: "config",
					},
				},
			},
		},
	}
}

func getHierarchicalDestConfig() utils.ReplicationDestination {
	return utils.ReplicationDestination{
		Hierarchical: true,
		Hierarchy: utils.Hierarchy{
			Labels: []utils.KVP{
				{
					Key:   Root,
					Value: "root",
				},
				{
					Key:   Project,
					Value: "project",
				},
				{
					Key:   Config,
					Value: "config",
				},
				{
					Key:   utils.DisplayNameKey,
					Value: "bar",
				},
			},
		},
	}
}

func getTypeConfig(group, kind string) utils.ReplicationSource {
	return utils.ReplicationSource{
		Kind: utils.Type,
		Type: utils.ObjectType{
			Group:   group,
			Version: "v1",
			Kind:    kind,
		},
	}
}

func getDifferentTypeDestConfig() utils.ReplicationDestination {
	return utils.ReplicationDestination{
		Hierarchical: false,
		ObjectType: &utils.ObjectType{
			Group:   Group,
			Version: "v1",
			Kind:    AdKind,
		},
	}
}

func getNexusEndpointObject(name, host, port, cert interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"host": host,
				"port": port,
				"cert": cert,
			},
		},
	}
}

func getReplicationConfigObject() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec": map[string]interface{}{
				"source": map[string]interface{}{
					"kind": "Type",
					"type": map[string]interface{}{
						"group":   "config.mazinger.com",
						"version": "v1",
						"kind":    "ApiCollaborationSpace",
					},
				},
				"destination": map[string]interface{}{
					"hierarchical": false,
				},
				"remoteEndpointGvk": map[string]interface{}{
					"group": "connect.nexus.org",
					"kind":  "NexusEndpoint",
					"name":  "default",
				},
			},
		},
	}
}

func getDefaultResourceObj() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "zoo",
			},
			"spec": map[string]interface{}{
				"example": "example",
			},
		},
	}
}

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}
