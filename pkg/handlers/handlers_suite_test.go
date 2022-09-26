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

func GetObject(name, kind, specVal string) *unstructured.Unstructured {
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

func GetParentObject(name, kind string) *unstructured.Unstructured {
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

func GetChildObject(name, kind string) *unstructured.Unstructured {
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

func GetReplicatedObject(name, kind string, status map[string]interface{}) *unstructured.Unstructured {
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

func GetNonHierarchicalDestConfig() utils.ReplicationDestination {
	return utils.ReplicationDestination{
		Hierarchical: false,
	}
}

func GetNonHierarchicalSourceConfig() utils.ReplicationSource {
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

func GetHierarchicalSourceConfig() utils.ReplicationSource {
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

func GetHierarchicalDestConfig() utils.ReplicationDestination {
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

func GetTypeConfig(group, kind string) utils.ReplicationSource {
	return utils.ReplicationSource{
		Kind: utils.Type,
		Type: utils.ObjectType{
			Group:   group,
			Version: "v1",
			Kind:    kind,
		},
	}
}

func GetDifferentTypeDestConfig() utils.ReplicationDestination {
	return utils.ReplicationDestination{
		Hierarchical: false,
		ObjectType: &utils.ObjectType{
			Group:   Group,
			Version: "v1",
			Kind:    AdKind,
		},
	}
}

func GetDefaultResourceObj() *unstructured.Unstructured {
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
