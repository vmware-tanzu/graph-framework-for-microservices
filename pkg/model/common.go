package model

import (
	"strings"

	authnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.org/v1"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	middlewarenexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/domain.nexus.org/v1"
)

// adding this global variables for CORS to support multiple domain and header configuration
var CorsConfigOrigins = map[string][]string{}
var CorsConfigHeaders = map[string][]string{}

type EventType string

const (
	Upsert EventType = "Upsert"
	Delete EventType = "Delete"
)

type NexusAnnotation struct {
	Name                 string                     `json:"name,omitempty"`
	Hierarchy            []string                   `json:"hierarchy,omitempty"`
	Children             map[string]NodeHelperChild `json:"children,omitempty"`
	Links                map[string]NodeHelperChild `json:"links,omitempty"`
	NexusRestAPIGen      nexus.RestAPISpec          `json:"nexus-rest-api-gen,omitempty"`
	NexusRestAPIMappings map[string]string          `json:"nexus-rest-api-mappings,omitempty"`
	IsSingleton          bool                       `json:"is_singleton,omitempty"`
	Description          string                     `json:"description,omitempty"`
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
	Links           map[string]NodeHelperChild
	IsSingleton     bool
	Description     string
}

type RestUriInfo struct {
	TypeOfURI URIType
}

type URIType int

const (
	DefaultURI URIType = iota
	SingleLinkURI
	NamedLinkURI
	StatusURI
)

func ConstructEchoPathParamURL(uri string) string {
	replacer := strings.NewReplacer("{", ":", "}", "")
	return replacer.Replace(uri)
}

type OidcNodeEvent struct {
	Oidc authnexusv1.OIDC
	Type EventType
}

type DatamodelInfo struct {
	Title string
}

type CorsNodeEvent struct {
	Cors middlewarenexusv1.CORSConfig
	Type EventType
}

// LinkGvk : This model used to carry fully qualified object <gvk> and
// hierarchy information.
type LinkGvk struct {
	Group     string   `json:"group,omitempty" yaml:"group,omitempty"`
	Kind      string   `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name      string   `json:"name,omitempty" yaml:"name,omitempty"`
	Hierarchy []string `json:"hierarchy,omitempty" yaml:"hierarchy,omitempty"`
}
