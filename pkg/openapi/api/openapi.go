package api

import (
	"api-gw/pkg/model"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"net/http"
	"regexp"
	"strings"
)

var Schema openapi3.T

func New() {
	schema := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:          "Nexus API GW APIs",
			Description:    "",
			TermsOfService: "",
			Contact:        nil,
			License:        nil,
			Version:        "1.0.0",
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "API Gateway",
				URL:         "/",
			},
			&openapi3.Server{
				Description: "Local",
				URL:         "http://localhost:5000",
			},
			&openapi3.Server{
				Description: "Local SSL",
				URL:         "https://localhost:5443",
			},
		},
		Paths: openapi3.Paths{},
		Components: openapi3.Components{
			Schemas:       openapi3.Schemas{},
			RequestBodies: openapi3.RequestBodies{},
			Responses: openapi3.Responses{
				"DefaultResponse": &openapi3.ResponseRef{
					Value: openapi3.NewResponse().
						WithDescription("Default response").
						WithContent(openapi3.NewContentWithJSONSchema(
							openapi3.NewSchema().WithProperty("message", openapi3.NewStringSchema())),
						),
				},
			},
		},
	}
	Schema = schema
}

func AddPath(uri nexus.RestURIs) {
	crdType := model.UriToCRDType[uri.Uri]
	crdInfo := model.CrdTypeToNodeInfo[crdType]
	parseSpec(crdType)

	h := sha1.New()

	params := parseUriParams(uri.Uri, crdInfo.ParentHierarchy)
	pathItem := &openapi3.PathItem{}
	for method, _ := range uri.Methods {
		h.Write([]byte(fmt.Sprintf("%s%s", method, uri.Uri)))
		opId := hex.EncodeToString(h.Sum(nil))
		nameParts := strings.Split(crdInfo.Name, ".")

		switch method {
		case http.MethodGet:
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[0]},
				Parameters:  params,
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/Get" + crdInfo.Name,
					},
				},
			}

			if len(params) == 0 {
				operation = &openapi3.Operation{
					OperationID: opId,
					Tags:        []string{nameParts[0]},
					Parameters:  params,
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Ref: "#/components/responses/List" + crdInfo.Name,
						},
					},
				}
			}
			pathItem.Get = operation
		case http.MethodPut:
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[0]},
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/Create" + crdInfo.Name,
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/DefaultResponse",
					},
				},
				Parameters: params,
			}
			pathItem.Put = operation
		case http.MethodDelete:
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[0]},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Value: openapi3.NewResponse().WithDescription("No content"),
					},
				},
				Parameters: params,
			}
			pathItem.Delete = operation
		}
	}

	Schema.Paths[uri.Uri] = pathItem
}

func parseSpec(crdType string) {
	crdInfo := model.CrdTypeToNodeInfo[crdType]
	crdSpec := model.CrdTypeToSpec[crdType]

	openapiSchema := crdSpec.Versions[0].Schema.OpenAPIV3Schema
	specProps := openapiSchema.Properties["spec"].Properties
	jsonSchema := openapi3.NewObjectSchema()
	parseFields(jsonSchema, specProps)

	Schema.Components.Schemas[crdInfo.Name] = openapi3.NewSchemaRef("", jsonSchema)

	Schema.Components.RequestBodies["Create"+crdInfo.Name] = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().
			WithDescription("Request used to create " + crdInfo.Name).
			WithRequired(true).
			WithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + crdInfo.Name}),
	}

	Schema.Components.Responses["Get"+crdInfo.Name] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Response returned back after getting " + crdInfo.Name + " object").
			WithContent(
				openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + crdInfo.Name}),
			),
	}

	Schema.Components.Responses["List"+crdInfo.Name] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Response returned back after getting " + crdInfo.Name + " objects").
			WithContent(
				openapi3.NewContentWithJSONSchema(openapi3.NewSchema()),
			),
	}
}

func parseFields(jsonSchema *openapi3.Schema, specProps map[string]v1.JSONSchemaProps) {
	for name, prop := range specProps {
		if strings.Contains(name, "Gvk") {
			continue
		}
		// TODO: Support additional properties
		format := prop.Format
		switch prop.Type {
		case "string":
			switch format {
			case "byte":
				jsonSchema.WithProperty(name, openapi3.NewBytesSchema())
			case "date-time":
				jsonSchema.WithProperty(name, openapi3.NewDateTimeSchema())
			default:
				jsonSchema.WithProperty(name, openapi3.NewStringSchema())
			}
		case "boolean":
			jsonSchema.WithProperty(name, openapi3.NewBoolSchema())
		case "object":
			schema := openapi3.NewSchema()
			parseFields(schema, prop.Properties)
			jsonSchema.WithProperty(name, schema)
		case "integer":
			switch format {
			case "int32":
				jsonSchema.WithProperty(name, openapi3.NewInt32Schema())
			case "int64":
				jsonSchema.WithProperty(name, openapi3.NewInt64Schema())
			default:
				jsonSchema.WithProperty(name, openapi3.NewIntegerSchema())
			}
		case "number":
			jsonSchema.WithProperty(name, openapi3.NewFloat64Schema())
		case "array":
			schema := openapi3.NewSchema()
			parseFields(schema, prop.Items.Schema.Properties)
			arraySchema := openapi3.NewArraySchema().WithItems(schema)
			jsonSchema.WithProperty(name, arraySchema)
		default:
			log.Infof("Unknown type %s", prop.Type)
		}
	}
}

func parseUriParams(uri string, hierarchy []string) (parameters []*openapi3.ParameterRef) {
	r := regexp.MustCompile(`{([^{}]+)}`)
	params := r.FindAllStringSubmatch(uri, -1)
	// If there are no params in URI then it's a list
	if len(params) == 0 {
		return
	}

	for _, param := range params {
		description := "Name of the " + param[1] + " node"
		for _, nodeInfo := range model.CrdTypeToNodeInfo {
			if nodeInfo.Name == param[1] {
				if nodeInfo.Description != "" {
					description = nodeInfo.Description
					break
				}
			}
		}
		parameters = append(parameters, &openapi3.ParameterRef{
			Value: openapi3.NewPathParameter(param[1]).
				WithRequired(true).
				WithSchema(openapi3.NewStringSchema()).
				WithDescription(description),
		})
	}

	for _, parent := range hierarchy {
		crdInfo := model.CrdTypeToNodeInfo[parent]
		if crdInfo.IsSingleton {
			continue
		}
		var description string
		if crdInfo.Description != "" {
			description = crdInfo.Description
		} else {
			description = "Name of the " + crdInfo.Name + " node"
		}
		if !paramExist(crdInfo.Name, params) {
			parameters = append(parameters, &openapi3.ParameterRef{
				Value: openapi3.NewQueryParameter(crdInfo.Name).
					WithRequired(true).
					WithSchema(openapi3.NewStringSchema()).
					WithDescription(description),
			})
		}
	}
	return
}

func paramExist(param string, params [][]string) bool {
	for _, p := range params {
		if p[1] == param {
			return true
		}
	}
	return false
}

func Recreate() {
	log.Debug("Recreating openapi spec")
	New()

	for _, uris := range model.CrdTypeToRestUris {
		for _, uri := range uris {
			AddPath(uri)
		}
	}
}
