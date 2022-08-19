// This file is used for testing purpose only.

package handlers

import (
	"connector/pkg/utils"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	ApiVersion = "config.mazinger.com/v1"
	Group      = "config.mazinger.com"

	// CRD Types
	Root                  = "roots.apix.mazinger.com"
	Project               = "projects.apix.mazinger.com"
	Config                = "configs.apix.mazinger.com"
	ApiCollaborationSpace = "apicollaborationspaces.config.mazinger.com"
	ApiDevSpace           = "apidevspaces.config.mazinger.com"

	// Object Kind
	AcKind = "ApiCollaborationSpace"
	AdKind = "ApiDevSpace"
)

func GetObject(name, kind string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": ApiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"example": "example",
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

func GetTypeConfig() utils.ReplicationSource {
	return utils.ReplicationSource{
		Kind: utils.Type,
		Type: utils.ObjectType{
			Group:   Group,
			Version: "v1",
			Kind:    AcKind,
		},
	}
}
