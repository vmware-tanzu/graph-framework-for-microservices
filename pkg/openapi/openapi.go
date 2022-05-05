package openapi

import (
	"api-gw/controllers"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"net/http"
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
				Description: "Local",
				URL:         "http://localhost:5000",
			},
		},
		Paths: openapi3.Paths{},
		Components: openapi3.Components{
			RequestBodies: openapi3.RequestBodies{},
		},
	}
	Schema = schema
}

func AddPath(uri nexus.RestURIs) {
	crdType := controllers.GlobalURIToCRDTypes[uri.Uri]
	crd := controllers.GlobalCRDTypeToNodes[crdType]
	parseSpec(crdType)

	pathItem := &openapi3.PathItem{}
	for method, _ := range uri.Methods {
		switch method {
		case http.MethodGet:
			operation := &openapi3.Operation{
				OperationID: "Get" + crd.Name,
				Tags:        []string{crd.Name},
				Responses:   openapi3.Responses{},
			}
			pathItem.Get = operation
		case http.MethodPut:
			operation := &openapi3.Operation{
				OperationID: "Put" + crd.Name,
				Tags:        []string{crd.Name},
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/Create" + crd.Name,
				},
				Responses: openapi3.Responses{},
			}
			pathItem.Put = operation
		case http.MethodDelete:
			operation := &openapi3.Operation{
				OperationID: "Delete" + crd.Name,
				Tags:        []string{crd.Name},
				Responses:   openapi3.Responses{},
			}
			pathItem.Delete = operation
		}
	}

	Schema.Paths[uri.Uri] = pathItem
}

func parseSpec(crdType string) {
	crd := controllers.GlobalCRDTypeToNodes[crdType]
	crdSpec := controllers.GlobalCRDTypeToSpec[crdType]

	openapiSchema := crdSpec.Versions[0].Schema.OpenAPIV3Schema
	specProps := openapiSchema.Properties["spec"].Properties
	jsonSchema := openapi3.NewSchema()
	parseFields(jsonSchema, specProps)

	Schema.Components.RequestBodies["Create"+crd.Name] = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().
			WithDescription("Request used to create " + crd.Name).
			WithRequired(true).
			WithJSONSchema(jsonSchema),
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
