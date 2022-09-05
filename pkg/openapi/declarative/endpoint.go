package declarative

import (
	"api-gw/pkg/utils"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

type EndpointContext struct {
	echo.Context

	SpecUri      string
	Method       string     // e.g. PUT
	KindName     string     // e.g. GlobalNamespace
	ResourceName string     // e.g. globalnamespaces
	GroupName    string     // e.g. vmware.org
	CrdName      string     // e.g. globalnamespaces.vmware.org
	Params       [][]string // e.g. [id, projectId, gnsId]
	Identifier   string     // e.g. id or fqdn

	Single bool // used to identify which k8s endpoint we should use (resource/:name or resource/)

	SchemaName string // OpenAPI.components.schema name used to create yaml spec
	ShortName  string
	ShortUri   string
	Uri        string
}

const (
	resourcePattern          = "/apis/%s/v1/%s"
	resourceShortPattern     = "/apis/v1/%s"
	resourceNamePattern      = resourcePattern + "/:name"
	resourceNameShortPattern = resourceShortPattern + "/:name"
)

func SetupContext(uri string, method string, item *openapi3.Operation) *EndpointContext {
	kindName := GetExtensionVal(item, NexusKindName)
	groupName := GetExtensionVal(item, NexusGroupName)
	shortName := GetExtensionVal(item, NexusShortName)
	resourceName := strings.ToLower(utils.ToPlural(kindName))
	crdName := resourceName + "." + groupName
	requiredParams := extractUriParams(uri)
	identifier := GetExtensionVal(item, "x-nexus-identifier")

	path := fmt.Sprintf(resourcePattern, groupName, resourceName)
	shortPath := fmt.Sprintf(resourceShortPattern, shortName)
	single := false
	if identifier != "" && method != http.MethodPut {
		single = true
		path = fmt.Sprintf(resourceNamePattern, groupName, resourceName)
		shortPath = fmt.Sprintf(resourceNameShortPattern, shortName)
	}

	schemaName := ""
	if item.RequestBody != nil && item.RequestBody.Value != nil {
		mediaType := item.RequestBody.Value.Content.Get("application/json")
		if mediaType != nil {
			schemaName = openapi3.DefaultRefNameResolver(mediaType.Schema.Ref)
		}
	}

	if shortName == "" {
		shortPath = ""
	}

	return &EndpointContext{
		SpecUri:      uri,
		KindName:     kindName,
		ResourceName: resourceName,
		GroupName:    groupName,
		CrdName:      crdName,
		Params:       requiredParams,
		Identifier:   identifier,
		Single:       single,
		Uri:          path,
		Method:       method,
		SchemaName:   schemaName,
		ShortName:    shortName,
		ShortUri:     shortPath,
	}
}

func IsArrayResponse(op *openapi3.Operation) bool {
	if op == nil {
		return false
	}

	resp := op.Responses.Get(200)
	if resp == nil {
		return false
	}

	mediaType := resp.Value.Content.Get("application/json")
	if mediaType == nil {
		return false
	}

	if mediaType.Schema.Value.Type == "array" {
		return true
	}

	return false
}

func extractUriParams(uri string) [][]string {
	r := regexp.MustCompile(`{([^{}]+)}`)
	params := r.FindAllStringSubmatch(uri, -1)
	if len(params) == 0 {
		return [][]string{}
	}
	return params
}
