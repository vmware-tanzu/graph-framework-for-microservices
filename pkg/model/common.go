package model

import (
	"strings"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type EventType string

const (
	Upsert EventType = "Upsert"
	Delete EventType = "Delete"
)

type NexusAnnotation struct {
	Name                 string                     `json:"name,omitempty"`
	Hierarchy            []string                   `json:"hierarchy,omitempty"`
	Children             map[string]NodeHelperChild `json:"children,omitempty"`
	NexusRestAPIGen      nexus.RestAPISpec          `json:"nexus-rest-api-gen,omitempty"`
	NexusRestAPIMappings map[string]string          `json:"nexus-rest-api-mappings,omitempty"`
}

type NodeHelperChild struct {
	FieldName    string `json:"fieldName"`
	FieldNameGvk string `json:"fieldNameGvk"`
	IsNamed      bool   `json:"isNamed"`
}

type NodeInfo struct {
	Name            string
	ParentHierarchy []string
	Children        map[string]NodeHelperChild
}

func ConstructEchoPathParamURL(uri string) string {
	replacer := strings.NewReplacer("{", ":", "}", "")
	return replacer.Replace(uri)
}
