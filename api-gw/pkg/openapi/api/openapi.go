package api

import (
	"api-gw/pkg/model"
	"api-gw/pkg/utils"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var Schemas = make(map[string]openapi3.T)

func New(datamodel string) {
	// Check if datamodel info is present
	title := "Nexus API GW APIs"
	if info, ok := model.DatamodelToDatamodelInfo[datamodel]; ok {
		title = info.Title
	}
	schema := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:          title,
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
				"NotFoundResponse": &openapi3.ResponseRef{
					Value: openapi3.NewResponse().
						WithDescription("Not Found").
						WithContent(openapi3.NewContentWithJSONSchema(
							openapi3.NewSchema().WithProperty("message", openapi3.NewStringSchema())),
						),
				},
			},
		},
	}
	Schemas[datamodel] = schema
}

func DatamodelUpdateNotification() {
	for {
		select {
		case name := <-model.DatamodelsChan:
			if schema, ok := Schemas[name]; ok {
				model.DatamodelToDatamodelInfoMutex.Lock()
				schema.Info.Title = model.DatamodelToDatamodelInfo[name].Title
				model.DatamodelToDatamodelInfoMutex.Unlock()
				log.Infof("Updated title: %s for %s openapi spec", schema.Info.Title, name)
			}
		}
	}
}

// AddPath creates and adds paths for all the methods of a URI
func AddPath(uri nexus.RestURIs, datamodel string) {
	crdType := model.UriToCRDType[uri.Uri]
	crdInfo := model.CrdTypeToNodeInfo[crdType]
	parseSpec(crdType, datamodel)

	h := sha1.New()

	params := parseUriParams(uri.Uri, crdInfo.ParentHierarchy)
	pathItem := &openapi3.PathItem{}
	for method := range uri.Methods {
		h.Write([]byte(fmt.Sprintf("%s%s", method, uri.Uri)))
		opId := hex.EncodeToString(h.Sum(nil))
		nameParts := strings.Split(crdInfo.Name, ".")

		switch method {
		case "LIST":
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[1]},
				Parameters:  params,
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/List" + crdInfo.Name,
					},
				},
			}
			pathItem.Get = operation
		case http.MethodGet:
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[1]},
				Parameters:  params,
			}
			if uriInfo, ok := model.GetUriInfo(uri.Uri); ok {
				switch uriInfo.TypeOfURI {
				case model.StatusURI:
					operation.Responses = openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Ref: "#/components/responses/Get" + crdInfo.Name + ".Status",
						},
					}
				case model.SingleLinkURI:
					operation.Responses = openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Ref: "#/components/responses/Get" + crdInfo.Name + ".SingleLink",
						},
					}
				case model.NamedLinkURI:
					operation.Responses = openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Ref: "#/components/responses/Get" + crdInfo.Name + ".NamedLink",
						},
					}
				default:
					operation.Responses = openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Ref: "#/components/responses/Get" + crdInfo.Name,
						},
					}
				}
			} else {
				operation.Responses = openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/DefaultResponse",
					},
				}
			}
			pathItem.Get = operation
		case http.MethodPut:
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[1]},
			}
			if uriInfo, ok := model.GetUriInfo(uri.Uri); ok && uriInfo.TypeOfURI == model.StatusURI {
				operation.RequestBody = &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/Create" + crdInfo.Name + ".Status",
				}
				operation.Responses = openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/DefaultResponse",
					},
				}
				operation.Parameters = params
			} else {
				p := constructUpdateParam()
				var putParams []*openapi3.ParameterRef
				putParams = append(putParams, params...)
				putParams = append(putParams, p)
				operation.Parameters = putParams

				operation.RequestBody = &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/Create" + crdInfo.Name,
				}
				operation.Responses = openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/DefaultResponse",
					},
				}
			}
			pathItem.Put = operation
		case http.MethodPatch:
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[1]},
				Parameters:  params,
			}
			operation.Responses = openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Ref: "#/components/responses/DefaultResponse",
				},
				"404": &openapi3.ResponseRef{
					Ref: "#/components/responses/NotFoundResponse",
				},
			}
			if uriInfo, ok := model.GetUriInfo(uri.Uri); ok && uriInfo.TypeOfURI == model.StatusURI {
				operation.RequestBody = &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/Create" + crdInfo.Name + ".Status",
				}
			} else {
				operation.RequestBody = &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/Create" + crdInfo.Name,
				}
			}
			pathItem.Patch = operation
		case http.MethodDelete:
			operation := &openapi3.Operation{
				OperationID: opId,
				Tags:        []string{nameParts[1]},
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

	log.Infof("adding %s path to %s", uri.Uri, datamodel)
	Schemas[datamodel].Paths[uri.Uri] = pathItem
}

// parseSpec parses openapi schema spec and status subresource
func parseSpec(crdType string, datamodel string) {
	crdInfo := model.CrdTypeToNodeInfo[crdType]
	crdSpec := model.CrdTypeToSpec[crdType]

	getKey := makeKey(crdInfo.Name, "Get")
	postKey := makeKey(crdInfo.Name, "Post")
	listKey := makeKey(crdInfo.Name, "List")
	statusKey := makeKey(crdInfo.Name, "Status")
	singleLinkKey := makeKey(crdInfo.Name, "SingleLink")
	namedLinkKey := makeKey(crdInfo.Name, "NamedLink")

	openapiSchema := crdSpec.Versions[0].Schema.OpenAPIV3Schema
	specProps := openapiSchema.Properties["spec"].Properties
	jsonSpecSchema := openapi3.NewObjectSchema()
	parseFields(jsonSpecSchema, specProps)

	statusProps := openapiSchema.Properties["status"].Properties
	delete(statusProps, "nexus")
	jsonStatusSchema := openapi3.NewObjectSchema()
	parseFields(jsonStatusSchema, statusProps)

	Schemas[datamodel].Components.Schemas[statusKey] = openapi3.NewSchemaRef("", jsonStatusSchema)

	jsonSpecAndStatusSchema := openapi3.NewObjectSchema()
	jsonSpecAndStatusSchema.WithProperty("spec", jsonSpecSchema)
	jsonSpecAndStatusSchema.WithProperty("status", jsonStatusSchema)

	Schemas[datamodel].Components.Schemas[postKey] = openapi3.NewSchemaRef("", jsonSpecSchema)
	Schemas[datamodel].Components.Schemas[getKey] = openapi3.NewSchemaRef("", jsonSpecAndStatusSchema)

	jsonListObjectSchema := openapi3.NewObjectSchema()
	jsonListObjectSchema.WithProperty("name", openapi3.NewStringSchema())
	jsonListObjectSchema.WithProperty("spec", jsonSpecSchema)
	jsonListObjectSchema.WithProperty("status", jsonStatusSchema)
	jsonListSchema := openapi3.NewArraySchema().WithItems(jsonListObjectSchema)

	Schemas[datamodel].Components.Schemas[listKey] = openapi3.NewSchemaRef("", jsonListSchema)

	// TODO: Schema for single link and named link need to be generated
	jsonSingleLinkSchema := openapi3.NewObjectSchema()
	jsonNamedLinkSchema := openapi3.NewArraySchema().WithItems(jsonSingleLinkSchema)
	Schemas[datamodel].Components.Schemas[singleLinkKey] = openapi3.NewSchemaRef("", jsonSingleLinkSchema)
	Schemas[datamodel].Components.Schemas[namedLinkKey] = openapi3.NewSchemaRef("", jsonNamedLinkSchema)

	Schemas[datamodel].Components.RequestBodies["Create"+crdInfo.Name] = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().
			WithDescription("Request used to create " + crdInfo.Name).
			WithRequired(true).
			WithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + postKey}),
	}

	Schemas[datamodel].Components.Responses["Get"+crdInfo.Name] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Response returned back after getting " + crdInfo.Name + " object").
			WithContent(
				openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + getKey}),
			),
	}

	Schemas[datamodel].Components.RequestBodies["Create"+statusKey] = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().
			WithDescription("Request used to create Status subresource of " + crdInfo.Name).
			WithRequired(false).
			WithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + statusKey}),
	}

	Schemas[datamodel].Components.Responses["Get"+statusKey] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Response returned back after getting status subresource of " + crdInfo.Name + " object").
			WithContent(
				openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + statusKey}),
			),
	}

	Schemas[datamodel].Components.Responses["List"+crdInfo.Name] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Response returned back after getting " + crdInfo.Name + " objects").
			WithContent(
				openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + listKey}),
			),
	}

	Schemas[datamodel].Components.Responses["Get"+singleLinkKey] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Response returned back after getting " + crdInfo.Name + " objects").
			WithContent(
				openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + singleLinkKey}),
			),
	}

	Schemas[datamodel].Components.Responses["Get"+namedLinkKey] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("Response returned back after getting " + crdInfo.Name + " objects").
			WithContent(
				openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/" + namedLinkKey}),
			),
	}
}

// ParseFields parses openapi schema fields
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

// parseUriParams parses the URI parameters
func parseUriParams(uri string, hierarchy []string) (parameters []*openapi3.ParameterRef) {
	r := regexp.MustCompile(`{([^{}]+)}`)
	params := r.FindAllStringSubmatch(uri, -1)

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
func constructUpdateParam() *openapi3.ParameterRef {
	return &openapi3.ParameterRef{
		Value: openapi3.NewQueryParameter("update_if_exists").
			WithRequired(false).
			WithSchema(openapi3.NewBoolSchema()).
			WithDescription("If set to false, disables update of preexisting object. Default value is true"),
	}
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
	for crdType := range model.CrdTypeToRestUris {
		New(utils.GetDatamodelName(crdType))
	}

	for crdType, uris := range model.CrdTypeToRestUris {
		datamodel := utils.GetDatamodelName(crdType)
		for _, uri := range uris {
			AddPath(uri, datamodel)
		}
	}
}

func makeKey(crd, keyType string) string {
	return crd + "." + keyType
}
