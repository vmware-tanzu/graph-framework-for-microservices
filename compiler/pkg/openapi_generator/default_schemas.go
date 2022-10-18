package openapi_generator

import (
	"fmt"

	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/validation/spec"
)

const (
	openAPIEnum   = "Mesh7CodeGenOpenAPIEnum"
	inline        = "Mesh7CodeGenInline"
	noFormat      = ""
	defaultSchema = "default"
)

// defaultSchemas returns a map of `schemaName: schema`, which will be used as
// default schemas in the generator.
// Schemas should be added here in 3 cases:
// 1. 3rd party struct is used in Go structure
// 2. Go structis  generated from proto which contains a `oneOf`
// Schemas should be added to the map returned by this function. The key to be used
// is following a format `<importPath>.<structName> e.g. `k8s.io/api/core/v1.Pod`.
// There are helpers below both for `pkg/apis` and `pkg/features` import paths.
// There are also examples below for each of the 2 cases mentioned above.
func defaultSchemas() map[string]common.OpenAPIDefinition {
	//rateLimiterPool, err := spec.NewRef(kavachTypeName("features/rate_limiter.pool"))
	//if err != nil {
	//	// This should not happen, but better safe than sorry.
	//	panic(err)
	//}
	return map[string]common.OpenAPIDefinition{
		//// Case 2. oneOfs. Note the empty `inline` property added. It won't be a part
		//// of output YAML schema, but is needed for internal mechanisms.
		//kavachTypeName("features/common/matcher.isRouteMatch_PathSpecifier"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"prefix": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"path": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"regex": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/plugins/hyper_dlp.isPiiAction_PiiVariant"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"pattern": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"patternGroup": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/plugins/hyper_dlp/defs.isCheck_CheckVariant"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"queryName": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"function": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"object"},
		//						Properties: map[string]spec.Schema{
		//							"name": {
		//								SchemaProps: spec.SchemaProps{
		//									Type: []string{"string"},
		//								},
		//							},
		//							"parameters": {
		//								SchemaProps: spec.SchemaProps{
		//									Type: []string{"object"},
		//								},
		//							},
		//						},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/plugin_interface.isAction_Action"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			// TODO confirm if this is proper
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"ratelimit": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"object"},
		//					},
		//				},
		//				"allow": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"object"},
		//					},
		//				},
		//				"deny": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"object"},
		//					},
		//				},
		//				"skip": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"object"},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/plugins/hyper_dlp/defs.isPatternDatabase_DatabaseFileVariant"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"databaseFile": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"useDefault": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"boolean"},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/plugins/modsecurity.isPolicy_Rule_RuleVariant"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"group": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"id": {
		//					SchemaProps: spec.SchemaProps{
		//						Type:   []string{"integer"},
		//						Format: "int32",
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/plugins/modsecurity.isConfig_Config"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"cfgFile": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"cfgPlain": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/rate_limiter.isPoolSelector_Variant"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"usePool": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"newPool": {
		//					SchemaProps: spec.SchemaProps{
		//						Ref: rateLimiterPool,
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/common/identity_awareness.isIdentityFormat_Jwt_Version"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"claim": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//				"claimPath": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"string"},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//kavachTypeName("features/filters/common/acp.isportsselector_portselector"): {
		//	Schema: spec.Schema{
		//		SchemaProps: spec.SchemaProps{
		//			Type: []string{"object"},
		//			Properties: map[string]spec.Schema{
		//				inline: {},
		//				"port": {
		//					SchemaProps: spec.SchemaProps{
		//						Type:   []string{"integer"},
		//						Format: "int32",
		//					},
		//				},
		//				"range": {
		//					SchemaProps: spec.SchemaProps{
		//						Type: []string{"object"},
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		// Case 1. 3rd party structs. See helpers for proto and k8s related import paths
		protoTypeName("BoolValue"): schemaForTypeAndFormat("boolean", noFormat),
		protoTypeName("Duration"):  schemaForTypeAndFormat("string", noFormat),
		protoTypeName("Int32Value"): schemaForOneOf(
			schemaForTypeAndFormat("integer", noFormat),
			schemaForTypeAndFormat("string", noFormat),
		),
		protoTypeName("StringValue"): schemaForTypeAndFormat("string", noFormat),
		protoTypeName("Struct"):      schemaForTypeAndFormat("object", noFormat),
		protoTypeName("Timestamp"):   schemaForTypeAndFormat("string", noFormat),
		protoTypeName("UInt32Value"): schemaForOneOf(
			schemaForTypeAndFormat("integer", noFormat),
			schemaForTypeAndFormat("string", noFormat),
		),
		protoTypeName("UInt64Value"): schemaForOneOf(
			schemaForTypeAndFormat("integer", noFormat),
			schemaForTypeAndFormat("string", noFormat),
		),

		k8sAPITypeName("apps/v1.DaemonSet"):                   schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.DaemonSetSpec"):               schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.DaemonSetStatus"):             schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.Deployment"):                  schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.DeploymentSpec"):              schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.DeploymentStatus"):            schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.ReplicaSet"):                  schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.ReplicaSetSpec"):              schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.ReplicaSetStatus"):            schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.StatefulSet"):                 schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.StatefulSetSpec"):             schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("apps/v1.StatefulSetStatus"):           schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("batch/v1.Job"):                        schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("batch/v1.JobSpec"):                    schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("batch/v1.JobStatus"):                  schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.ConfigMap"):                   schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.EndpointSubset"):              schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.NodeSpec"):                    schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.NodeStatus"):                  schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.Pod"):                         schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.PodSpec"):                     schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.PodStatus"):                   schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.ReplicationControllerSpec"):   schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.ReplicationControllerStatus"): schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.ServiceSpec"):                 schemaForTypeAndFormat("object", noFormat),
		k8sAPITypeName("core/v1.ServiceStatus"):               schemaForTypeAndFormat("object", noFormat),

		k8sAPIMachineryTypeName("api/resource.Quantity"):   schemaForTypeAndFormat("string", noFormat),
		k8sAPIMachineryTypeName("apis/meta/v1.Duration"):   schemaForTypeAndFormat("string", noFormat),
		k8sAPIMachineryTypeName("apis/meta/v1.ListMeta"):   schemaForTypeAndFormat("object", noFormat),
		k8sAPIMachineryTypeName("apis/meta/v1.ObjectMeta"): schemaForTypeAndFormat("object", noFormat),
		k8sAPIMachineryTypeName("apis/meta/v1.Time"):       schemaForTypeAndFormat("string", noFormat),

		envoyTypeName("config/core/v3.grpcservice"):                                  schemaForTypeAndFormat("object", noFormat),
		envoyTypeName("extensions/access_loggers/grpc/v3.commongrpcaccesslogconfig"): schemaForTypeAndFormat("object", noFormat),

		// The `openapi-gen` cannot generate Go schema for struct builtin. We mark is
		// as an object which allows any properties inside
		"struct{}":     schemaForTypeAndFormat("object", noFormat),
		"struct%7b%7d": schemaForTypeAndFormat("object", noFormat),

		// defaultSchema is used when YAML generator can't find a schema for given key.
		// It will also add the missing schema to `missingDefinitions` which can be accessed
		// with `MissingDefinitions()`.
		defaultSchema: schemaForTypeAndFormat("object", noFormat),
	}
}

func protoTypeName(name string) string {
	return fmt.Sprintf("github.com/gogo/protobuf/types.%v", name)
}

func k8sAPITypeName(name string) string {
	return fmt.Sprintf("k8s.io/api/%v", name)
}

func k8sAPIMachineryTypeName(name string) string {
	return fmt.Sprintf("k8s.io/apimachinery/pkg/%v", name)
}

func kavachTypeName(name string) string { //nolint:deadcode,unused
	return fmt.Sprintf("github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/%v", name)
}

func envoyTypeName(name string) string {
	return fmt.Sprintf("github.com/envoyproxy/go-control-plane/envoy/%v", name)
}

func schemaForOneOf(schemas ...common.OpenAPIDefinition) common.OpenAPIDefinition {
	s := common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				AnyOf: make([]spec.Schema, len(schemas)),
			},
		},
	}
	for i, sch := range schemas {
		s.Schema.SchemaProps.AnyOf[i] = sch.Schema
	}
	return s
}

// schemaForTypeAndFormat returns a simple schema with given type. If given format
// is not empty, it would be filled as well. Given schema does not contain any
// properties.
func schemaForTypeAndFormat(schemaType, schemaFormat string) common.OpenAPIDefinition {
	s := common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{schemaType},
			},
		},
	}
	if schemaFormat != noFormat {
		s.Schema.SchemaProps.Format = schemaFormat
	}
	return s
}
