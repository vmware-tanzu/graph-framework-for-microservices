package combined

import (
	"api-gw/pkg/config"
	"api-gw/pkg/openapi/api"
	"api-gw/pkg/openapi/declarative"
	"github.com/getkin/kin-openapi/openapi3"
)

func CombinedSpecs() openapi3.T {
	newSchema := openapi3.T{
		OpenAPI:    "3.0.0",
		Components: declarative.Schema.Components,
		Info:       declarative.Schema.Info,
		Paths:      declarative.Schema.Paths,
		Security:   declarative.Schema.Security,
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL: config.Cfg.TenantApiGwDomain + "/tsm/",
			},
			&openapi3.Server{
				URL: config.Cfg.TenantApiGwDomain + "/local/v1/",
			},
			&openapi3.Server{
				URL: "http://localhost:3000/v1/",
			},
			&openapi3.Server{
				URL: "http://localhost:3000/",
			},
		},
		Tags:         declarative.Schema.Tags,
		ExternalDocs: declarative.Schema.ExternalDocs,
	}

	nexusSchemas := api.Schemas

	for _, schema := range nexusSchemas {
		for k, v := range schema.Paths {
			if newSchema.Paths == nil {
				newSchema.Paths = openapi3.Paths{}
			}
			newSchema.Paths[k] = v
		}

		for k, v := range schema.Components.Schemas {
			if newSchema.Components.Schemas == nil {
				newSchema.Components.Schemas = openapi3.Schemas{}
			}
			newSchema.Components.Schemas[k] = v
		}

		for k, v := range schema.Components.Parameters {
			if newSchema.Components.Parameters == nil {
				newSchema.Components.Parameters = openapi3.ParametersMap{}
			}
			newSchema.Components.Parameters[k] = v
		}

		for k, v := range schema.Components.Headers {
			if newSchema.Components.Headers == nil {
				newSchema.Components.Headers = openapi3.Headers{}
			}
			newSchema.Components.Headers[k] = v
		}

		for k, v := range schema.Components.RequestBodies {
			if newSchema.Components.RequestBodies == nil {
				newSchema.Components.RequestBodies = openapi3.RequestBodies{}
			}
			newSchema.Components.RequestBodies[k] = v
		}

		for k, v := range schema.Components.Responses {
			if newSchema.Components.Responses == nil {
				newSchema.Components.Responses = openapi3.Responses{}
			}

			newSchema.Components.Responses[k] = v
		}

		for k, v := range schema.Components.SecuritySchemes {
			if newSchema.Components.SecuritySchemes == nil {
				newSchema.Components.SecuritySchemes = openapi3.SecuritySchemes{}
			}
			newSchema.Components.SecuritySchemes[k] = v
		}

		for k, v := range schema.Components.Examples {
			if newSchema.Components.Examples == nil {
				newSchema.Components.Examples = openapi3.Examples{}
			}
			newSchema.Components.Examples[k] = v
		}

		for k, v := range schema.Components.Links {
			if newSchema.Components.Links == nil {
				newSchema.Components.Links = openapi3.Links{}
			}
			newSchema.Components.Links[k] = v
		}

		for k, v := range schema.Components.Callbacks {
			if newSchema.Components.Callbacks == nil {
				newSchema.Components.Callbacks = openapi3.Callbacks{}
			}
			newSchema.Components.Callbacks[k] = v
		}
	}

	return newSchema
}
