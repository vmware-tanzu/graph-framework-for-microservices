package declarative

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/getkin/kin-openapi/openapi3"
)

var supportedOperations = []string{"GET", "DELETE", "PUT"}

const NexusKindName = "x-nexus-kind-name"
const NexusGroupName = "x-nexus-group-name"
const NexusListEndpoint = "x-nexus-list-endpoint"
const NexusShortName = "x-nexus-short-name"
const OpenApiSpecFile = "/openapi/openapi.yaml"
const OpenApiSpecDir = "/openapi"

var (
	Paths              = make(map[string]*openapi3.PathItem)
	ApisList           = make(map[string]map[string]interface{})
	apisListMutex      = sync.Mutex{}
	Schema             = openapi3.T{}
	Schemas            openapi3.Schemas
	parsedSchemas      = make(map[string]interface{})
	parsedSchemasMutex = sync.Mutex{}
	CrdToSchema        = make(map[string]string)
	crdToSchemaMutex   = sync.Mutex{}
)

func Setup(openApiSpecFile string) error {
	_, err := os.Stat(openApiSpecFile)
	if err == nil {
		f, err := ioutil.ReadFile(openApiSpecFile)
		if err != nil {
			return err
		}

		return Load(f)
	}
	log.Errorln("File", openApiSpecFile, " is not present at setup")
	return nil
}

func Load(data []byte) error {
	doc, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		return err
	}

	Schemas = doc.Components.Schemas
	Schema = *doc

	for uri, pathInfo := range doc.Paths {
		if !ValidateNexusAnnotations(pathInfo) {
			continue
		}
		Paths[uri] = pathInfo
	}

	ParseSchemas()

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
	apisListMutex.Lock()
	crdToSchemaMutex.Lock()
	defer func() {
		apisListMutex.Unlock()
		crdToSchemaMutex.Unlock()
	}()

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

	if ec.SchemaName != "" {
		schema := ConvertSchemaToYaml(ec, params)
		ApisList[ec.Uri]["yaml"] = schema
		CrdToSchema[fmt.Sprintf("%s.%s", ec.ResourceName, ec.GroupName)] = schema
	}

	if ec.ShortUri != "" {
		ApisList[ec.Uri]["short"] = map[string]interface{}{
			"name": ec.ShortName,
			"uri":  ec.ShortUri,
		}
	}
}

func ConvertSchemaToYaml(ec *EndpointContext, params []string) string {
	labels := map[string]interface{}{}
	for _, param := range params {
		if param != ec.Identifier {
			labels[param] = "string"
		}
	}

	obj := map[string]interface{}{
		"apiVersion": ec.GroupName + "/v1",
		"kind":       ec.KindName,
		"metadata": map[string]interface{}{
			"name":   "string",
			"labels": labels,
		},
	}
	obj["spec"] = parsedSchemas[ec.SchemaName]
	yamlObj, err := yaml.Marshal(obj)
	if err != nil {
		log.Warn(err)
	}
	return string(yamlObj)
}

func parseSchema(schemaName string, wg *sync.WaitGroup) {
	parsedSchemasMutex.Lock()
	defer func() {
		parsedSchemasMutex.Unlock()
		wg.Done()
	}()

	spec := make(map[string]interface{})

	for field, val := range Schemas[schemaName].Value.Properties {
		switch val.Value.Type {
		case "string":
			spec[field] = "string"
			if len(val.Value.Enum) > 0 {
				spec[field] = val.Value.Enum[0]
			}
		case "boolean":
			spec[field] = true
		case "number":
			spec[field] = 1.2
		case "integer":
			spec[field] = 1
		case "array":
			if val.Value.Items.Ref != "" {
				ref := openapi3.DefaultRefNameResolver(val.Value.Items.Ref)
				if ref == schemaName {
					spec[field] = "object"
				} else {
					spec[field] = map[string]interface{}{
						"ref":  ref,
						"type": "array",
					}
				}
			} else {
				if val.Value.Items.Value.Type == "string" {
					spec[field] = []string{val.Value.Items.Value.Type}
				}
			}
		case "object":
			spec[field] = "object"
		}

		if val.Ref != "" {
			ref := openapi3.DefaultRefNameResolver(val.Ref)
			if ref == schemaName {
				spec[field] = "object"
			} else {
				spec[field] = map[string]interface{}{
					"ref": ref,
				}
			}
		}
	}

	parsedSchemas[schemaName] = spec
}

func parseSchemaRefs(schemaName string, wg *sync.WaitGroup) {
	parsedSchemasMutex.Lock()
	defer func() {
		parsedSchemasMutex.Unlock()
		wg.Done()
	}()

	for fieldName, fieldVal := range parsedSchemas[schemaName].(map[string]interface{}) {
		if _, ok := fieldVal.(map[string]interface{}); !ok {
			continue
		}

		fv := fieldVal.(map[string]interface{})
		ref := fv["ref"]
		refType := fv["type"]

		if ref == nil || ref == schemaName {
			continue
		}

		refStr := ref.(string)

		if refType == "array" {
			parsedSchemas[schemaName].(map[string]interface{})[fieldName] = []map[string]interface{}{
				parsedSchemas[refStr].(map[string]interface{}),
			}
			continue
		}

		parsedSchemas[schemaName].(map[string]interface{})[fieldName] = parsedSchemas[refStr]
	}
}

func ParseSchemas() {
	wg := &sync.WaitGroup{}
	for schemaName := range Schemas {
		wg.Add(1)
		log.Debugf("Parsing %s schema", schemaName)
		go parseSchema(schemaName, wg)
	}
	wg.Wait()

	for schemaName := range parsedSchemas {
		wg.Add(1)
		log.Debugf("Parsing %s schema refs", schemaName)
		go parseSchemaRefs(schemaName, wg)
	}
	wg.Wait()
	log.Debugf("Finished parsing schemas")
}

func Middleware(endpointContext *EndpointContext, single bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			endpointContext.Context = c
			endpointContext.Single = single
			return next(endpointContext)
		}
	}
}

// ShortNames method creates a list of short names which can be used to access APIs
//func ShortNames(apisList map[string]map[string]interface{}) map[string]string {
//	shortMap := make(map[string]string)
//	for _, methods := range apisList {
//		resources := make(map[string]string)
//
//		for _, val := range methods {
//			if _, ok := val.(map[string]interface{}); !ok {
//				continue
//			}
//			info := val.(map[string]interface{})
//			resources[info["kind"].(string)] = info["group"].(string)
//		}
//
//		for k, g := range resources {
//			resourceName := strings.ToLower(utils.ToPlural(k))
//			if _, ok := shortMap[resourceName]; !ok {
//				shortMap[resourceName] = fmt.Sprintf("%s.%s", resourceName, g)
//			} else {
//				delete(shortMap, resourceName)
//			}
//		}
//	}
//	return shortMap
//}
