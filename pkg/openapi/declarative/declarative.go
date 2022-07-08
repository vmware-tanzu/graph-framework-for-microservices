package declarative

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var supportedOperations = []string{"GET", "DELETE", "PUT"}

const NexusKindName = "x-nexus-kind-name"
const NexusGroupName = "x-nexus-group-name"
const NexusListEndpoint = "x-nexus-list-endpoint"

var (
	Paths    = make(map[string]*openapi3.PathItem)
	ApisList = make(map[string]map[string]interface{})
)

func Setup() error {
	_, err := os.Stat("/openapi/openapi.yaml")
	if err == nil {
		f, err := ioutil.ReadFile("/openapi/openapi.yaml")
		if err != nil {
			return err
		}

		return Load(f)
	}
	return nil
}

func Load(data []byte) error {
	doc, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		return err
	}

	for uri, pathInfo := range doc.Paths {
		if !ValidateNexusAnnotations(pathInfo) {
			continue
		}
		Paths[uri] = pathInfo
	}

	return nil
}

func ValidateNexusAnnotations(item *openapi3.PathItem) bool {
	for _, supportedOperation := range supportedOperations {
		op := item.GetOperation(supportedOperation)
		if op != nil {
			if GetExtensionVal(op, NexusKindName) == "" {
				return false
			}

			if GetExtensionVal(op, NexusGroupName) == "" {
				return false
			}
		}
	}

	return true
}

func GetExtensionVal(operation *openapi3.Operation, key string) string {
	val, ok := operation.ExtensionProps.Extensions[key]
	if val == nil || !ok {
		return ""
	}

	out, _ := val.(json.RawMessage).MarshalJSON()
	outStr := string(out)

	if strings.HasPrefix(outStr, `"`) && strings.HasSuffix(outStr, `"`) && len(outStr) > 2 {
		return outStr[1 : len(outStr)-1]
	}

	return outStr
}

func AddApisEndpoint(ec *EndpointContext) {
	if ApisList[ec.Uri] == nil {
		ApisList[ec.Uri] = make(map[string]interface{})
	}

	var params []string
	for _, param := range ec.Params {
		params = append(params, param[1])
	}

	ApisList[ec.Uri][ec.Method] = map[string]interface{}{
		"group":  ec.GroupName,
		"kind":   ec.KindName,
		"params": params,
		"uri":    ec.SpecUri,
	}
}
