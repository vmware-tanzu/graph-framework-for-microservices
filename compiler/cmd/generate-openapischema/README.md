# generate-openapischema (YAML generator)

* [Overview](#overview)
* [Generation process](#generation-process)
  * [Resolving references](#resolving-references)
  * [Third-party Go structs](#third-party-go-structs)
  * [Go structs from proto enums](#go-structs-from-proto-enums)
  * [Go structs from proto oneOf](#go-structs-from-proto-oneof)
  * [YAMLs generation/update](#yamls-generationupdate)
* [Possible missing schema error messages and how to solve them](#possible-missing-schema-error-messages-and-how-to-solve-them)
* [Things you should know](#things-you-should-know)
* [Generator limitations](#generator-limitations)
* [Possible improvements](#possible-improvements)

## Overview
The YAML generator is used to generate OpenAPI schemas for CustomResourceDefinitions (CRDs).
It is described more in the official Kubernetes documentation
[here](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#specifying-a-structural-schema).

Those schemas are generated from Go structures definition of CRDs. In this repository
they are kept in `pkg/apis` and `pkg/features` paths. The latter ones are generated
automatically from proto files, which makes the generation process not so trivial.
All the challenges and solutions for them are described below.

## Generation process
The YAML generator input are Go schemas generated with `openapi-gen` from
[k8s.io/kube-openapi](https://github.com/kubernetes/kube-openapi/). Those schemas
are not ready to marhsall them into YAMLs out-of-the-box. First we need to solve
some challenges.

### Resolving references
Go schemas are created for each struct separately, which results in something
similar as in an example below.

Input Go structures:
```go
// +k8s:openapi-gen=true
type Foo struct {
    Bar Bar `json:"bar"`
}

// +k8s:openapi-gen=true
type Bar struct {
    Fizz string `json:"fizz"`
    Buzz int    `json:"buzz"`
}
```

Output Go schemas:
```go
func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
   return map[string]common.OpenAPIDefinition{
      "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/example.Bar": schema_KavachSec_policymodel_pkg_example_Bar(ref),
      "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/example.Foo": schema_KavachSec_policymodel_pkg_example_Foo(ref),
   }
}

func schema_KavachSec_policymodel_pkg_example_Bar(ref common.ReferenceCallback) common.OpenAPIDefinition {
   return common.OpenAPIDefinition{
      Schema: spec.Schema{
         SchemaProps: spec.SchemaProps{
            Type: []string{"object"},
            Properties: map[string]spec.Schema{
               "fizz": {
                  SchemaProps: spec.SchemaProps{
                     Default: "",
                     Type:    []string{"string"},
                     Format:  "",
                  },
               },
               "buzz": {
                  SchemaProps: spec.SchemaProps{
                     Default: 0,
                     Type:    []string{"integer"},
                     Format:  "int32",
                  },
               },
            },
            Required: []string{"fizz", "buzz"},
         },
      },
   }
}

func schema_KavachSec_policymodel_pkg_example_Foo(ref common.ReferenceCallback) common.OpenAPIDefinition {
   return common.OpenAPIDefinition{
      Schema: spec.Schema{
         SchemaProps: spec.SchemaProps{
            Type: []string{"object"},
            Properties: map[string]spec.Schema{
               "bar": {
                  SchemaProps: spec.SchemaProps{
                     Default: map[string]interface{}{},
                     Ref:     ref("github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/example.Bar"),
                  },
               },
            },
            Required: []string{"bar"},
         },
      },
      Dependencies: []string{
         "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/example.Bar"},
   }
}

```
Note the `Ref` field in `schema_KavachSec_policymodel_pkg_example_Foo`. As those
are supported in "normal" OpenAPI, they are not supported in Kubernetes. We need
to resolve those references and just put the dependent schema in that place.

Fortunately, each Go schema also contains information about dependent schemas in
`Dependencies` field. This allows us to iterate through the schemas tree recursively
and create schemas without refs.

Unfortunately, this has some edge cases, which are described below.

### Third-party Go structs
Some Go structures are embedding third-party Go structures, for which `openapi-gen`
cannot generate Go schemas, as they do not have `// +k8s:openapi-gen=true` annotation.
In such a case, we have a reference, but we are lacking a schema to replace it.
See the next example below.

Input Go structures:
```go
import (
   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
type Foo struct {
   Bar metav1.ObjectMeta `json:"bar"`
}

```

Output Go schemas:
```go
func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
   return map[string]common.OpenAPIDefinition{
      "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/example.Foo": schema_KavachSec_policymodel_pkg_example_Foo(ref),
   }
}

func schema_KavachSec_policymodel_pkg_example_Foo(ref common.ReferenceCallback) common.OpenAPIDefinition {
   return common.OpenAPIDefinition{
      Schema: spec.Schema{
         SchemaProps: spec.SchemaProps{
            Type: []string{"object"},
            Properties: map[string]spec.Schema{
               "bar": {
                  SchemaProps: spec.SchemaProps{
                     Default: map[string]interface{}{},
                     Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
                  },
               },
            },
            Required: []string{"bar"},
         },
      },
      Dependencies: []string{
         "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
   }
}

```
Note that we have `k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta` both in `Ref`
and `Dependencies`.

To solve that issue, the developer needs to specify the schema by hand in
`pkg/generator/default_schemas.go`. To solve the issue from example above it can
look like this:
```go
defaultSchemas := map[string]common.OpenAPIDefinition{
  "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"): simpleSchema("object", ""),
}

func simpleSchema(schemaType, schemaFormat string) common.OpenAPIDefinition {
    s := common.OpenAPIDefinition{
        Schema: spec.Schema{
            SchemaProps: spec.SchemaProps{
                Type: []string{schemaType},
            },
        },
    }
    if schemaFormat != "" {
        s.Schema.SchemaProps.Format = schemaFormat
    }
    return s
}
```
**Note**: adding empty `type: object` schema (i.e. without any properties specified),
will result in `x-kubernetes-preserve-unknown-fields: true` added to it in order to
disable pruning unknown fields in K8s.

There are many helpers and examples in `pkg/generator/default_schemas.go` so adding
third-party struct schema should be straightforward.

### Go structs from proto enums
When Go structure is generated from proto file and it has some enum embedded:
```proto
enum Enum {
  VALUE1 = 0
  VALUE2 = 1;
}
```

it is translated to something like below (code is simplified for readability):
```go
type Enum int32

const (
    Enum_VALUE1 Enum = 0
    Enum_VALUE2 Enum = 1
)

var Enum_name = map[int32]string{
    0: "VALUE1",
    1: "VALUE2",
}

var Enum_value = map[string]int32{
    "VALUE1": 0,
    "VALUE2": 1,
}

// +k8s:openapi-gen=true
type Foo struct {
  Bar Enum `protobuf:"enum=Enum" json:"bar"`
}
```
**Note:** in Go structures created by us we also have such contructs, but they can
be treated in the same way as those from proto files.

As you can see, the enum is represented as `int32` in the Go structure. But there
is a catch, the `Foo` structure has custom YAML/JSON marshallers which support both,
`int32` value directly or does does silent when receiving a `string`.
In other words, when unmarshalling from YAML it expects the `bar` field to be either
an `int32` or a `string` which is translated into `int32`.

This leads us to our problem. We can not simply require the same type in OpenAPI
schema as it is descibed in Go structure. We need to support both `int32` and `string`.

The solution is to annotate every field using such type with `// Mesh7CodeGenOpenAPIEnum`.
For structs generated from proto, it is done automatically in `./scripts/generate_features_api.sh`.
For structs in `pkg/apis` it needs to by done by hand, for example:
```go
// +k8s:openapi-gen=true
type Segmentation struct {
  // Mesh7CodeGenOpenAPIEnum
  Mode SegmentationMode `json:"mode"`
}

// EgressL7BaselineSpec is Spec part of EgressL7Baseline CRD
// +k8s:openapi-gen=true
type EgressL7BaselineSpec struct {
  // Mesh7CodeGenOpenAPIEnum
  Mode              SegmentationMode              `json:"mode"`
  L4OnlyEnforcement bool                          `json:"l4OnlyEnforcement"`
  Order             *int                          `json:"order,omitempty"`
  Selector          EgressBaselineCrdSelector     `json:"selectors"`
  Application       []EgressL7BaselineApplication `json:"application,omitempty"`
}

type SegmentationMode int

const (
  // Monitoring: We detect Segments, put them in CRD and log all violations
  Monitoring SegmentationMode = iota

  // Disabled: Segmentation doesn't do anything, http logs are passed by
  Disabled

  // Alerting: monitoring + raise events for all violations
  Alerting

  // Enforcing: alerting + load baselines to envoy
  Enforcing
)
```
Annotations are interpreted by the YAML generator, which adds support for both,
`int32` and `string` in such places and removes the annotation from YAML schema.
To be fully supported by K8s API, it also adds `x-kubernetes-int-or-string: true`
flag.

YAML output would look like this:
```yaml
bar:
  anyOf:
  - type: integer
  - type: string
  x-kubernetes-int-or-string: true
```

### Go structs from proto oneOf
When oneOf is used in proto:
```proto
message Foo [
  oneof bar {
    int number = 2;
    string string = 3;
  }
  string fizz = 4;
}
```
it results in a following Go structure in `pkg/features/example/foo.go`:
```go
// +k8s:openapi-gen=true
type Foo struct {
  // Types that are valid to be assigned to Bar:
  //  *Bar_Number
  //  *Bar_String
  Bar isBar_Bar
  Fizz string `json:"fizz"`
}

type isBar_Bar interface{
  isBar_Bar()
}

type Bar_Number {
  Number int `json:"number"`
}

type Bar_String {
  String string `json:"string"`
}

func (*Bar_Number) isBar_Bar() {}
func (*Bar_String) isBar_Bar() {}
```
as `openapi-gen` cannot generate schemas for interfaces, we need to create it by hand:
```go
"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/features/example.bar": common.OpenAPIDefinition{
    Schema: spec.Schema{
        SchemaProps: spec.SchemaProps{
            Type: []string{"object"},
            Properties: map[string]spec.Schema{
                "number": {
                    SchemaProps: spec.SchemaProps{Type: []string{"integer"}},
                },
                "string": {
                    SchemaProps: spec.SchemaProps{Type: []string{"string"}},
                },
            },
        },
    },
}
```
but this will result in a YAML schema like:
```yaml
foo:
    type: object
    properties:
        Bar:
            number:
                type: integer
            string:
                type: string
        fizz:
            type: string
```
Note the `Bar` starting from capital letter. In the Go struct we do not have the
`json` tag defined, so the yaml library uses the field name. And now compare this
to the proto definition - number and string fields should be at the same level as
`fizz`, but they are one level below, inside the `Bar` which is not even described
in the proto file.
To notify YAML generator that those values should be inlined in the parent `Properties`,
Empty property named `Mesh7CodeGenInline` should be added to the `Bar` schema.
Additional property would be removed before updating YAMLs.

The final schema for interface will look like:
```go
"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/features/example.bar": common.OpenAPIDefinition{
    Schema: spec.Schema{
        SchemaProps: spec.SchemaProps{
            Type: []string{"object"},
            Properties: map[string]spec.Schema{
                "Mesh7CodeGenInline": {},
                "number": {
                    SchemaProps: spec.SchemaProps{Type: []string{"integer"}},
                },
                "string": {
                    SchemaProps: spec.SchemaProps{Type: []string{"string"}},
                },
            },
        },
    },
}
```

Which will result in the following YAML schema:
```yaml
foo:
    type: object
    properties:
        number:
            type: integer
        string:
            type: string
        fizz:
            type: string
```

### YAMLs generation/update
After all the schemas are fixed as described above, YAML generator takes the given
path to a directory and recursively opens all the files one by one. Each file content
is split by `---` and all the parts are unmarshalled to `extensionsv1.CustomResourceDefinition`
Go struct. Later, the schema is added (overwriting the previous value) and all the
CRDs are marshalled back to YAML and written to a file with `---` separator.

## Possible missing schema error messages and how to solve them
**NOTE** this section is only about missing schema error messages. All the other
error messages most likely implicate a bug.

All the missing schema error messages look similar, but have small detail which
allow us to identify what is the underlying issue. This detail is the path of the
schema. Refer to the examples below for instructions how to solve any of them.

1. `Missing schema for "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/apis/policy.mesh7.io/v1.tenantconfigspec"`

    The missing schema is from `pkg/apis` path, so the most probable cause is that
    you forgot to annotate the mentioned Go struct (in this case `TenantConfigSpec`)
    with `// +k8s:openapi-gen=true`. Ensure that the is the comment block
    **exactly above** the struct. See examples below.

    ```go
    // This is CORRECT way of annotating

    // +k8s:openapi-gen=true
    type TenantConfigSpec {
        // something
    }
    ```

    ```go
    // This is INCORRECT way of annotating when there is already a comment

    // +k8s:openapi-gen=true

    // TenantConfigSpec contains information about TenantConfig guts
    type TenantConfigSpec {
        // something
    }
    ```

    ```go
    // This is CORRECT, but not preferred, way of annotating when there is already a comment

    // +k8s:openapi-gen=true
    //
    // TenantConfigSpec contains information about TenantConfig guts
    type TenantConfigSpec {
        // something
    }
    ```

    ```go
    // This is CORRECT, and preferred, way of annotating when there is already a comment

    // TenantConfigSpec contains information about TenantConfig guts
    // +k8s:openapi-gen=true
    type TenantConfigSpec {
        // something
    }
    ```

1. `Missing schema for "k8s.io/apimachinery/pkg/api/resource.quantity"`

    The missing schema is neither in `pkg/apis` nor in `pkg/features` path. It means
    that this is a 3rd party schema and needs a manually created entry in
    `pkg/generator/default_schemas.go`. As this is a 3rd party struct, we do not
    need to specify it's content, rather allowing any field in it. To do so,
    schema of type `object` with no specified format of properties needs to be
    created.

    ```go
    // Snippet from `pkg/generator/default_schemas.go`
    k8sAPIMachineryTypeName("api/resource.Quantity"):   schemaForTypeAndFormat("string", noFormat),
    ```

1. `Missing schema for "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/features/plugins/hyper_dlp.ispiiaction_piivariant"`

    The missing schema is in `pkg/features` path. It means that this is a schema
    for a struct generated from `oneof` in proto file and needs a manually created
    entry in `pkg/generator/default_schemas.go`. **NOTE** Remember about adding
    empty `inline` property. For more details refer to
    [Go structs from proto oneOf](#go-structs-from-proto-oneof).

    ```go
    kavachTypeName("features/plugins/hyper_dlp.isPiiAction_PiiVariant"): {
        Schema: spec.Schema{
            SchemaProps: spec.SchemaProps{
                Type: []string{"object"},
                Properties: map[string]spec.Schema{
                    inline: {},
                    "pattern": {
                        SchemaProps: spec.SchemaProps{
                            Type: []string{"string"},
                        },
                    },
                    "patternGroup": {
                        SchemaProps: spec.SchemaProps{
                            Type: []string{"string"},
                        },
                    },
                },
            },
        },
    },
    ```

## Things you should know
1. YAML generator enforces use of `camelCase` style in serialized field naming.
This should be handled properly by unmarshallers for types generated from proto,
but not by Go structs written by hand, withouth custom unmarshaller.
1. YAML generator will notify you when you need to write a schema manually in
`pkg/generator/default_schemas.go`, but **will not** inform you when a modification
is required. Whenever proto file is changes, ensure all the changes are reflected
in the file mentioned above.
1. When a schema has type set to `object`, but no properties are specified,
YAML generator will add `x-kubernetes-preserve-unknown-fields: true` flag to it. This rule
has one exception, which is `.metadata` field which is required by K8s to only
have `type: object` specified.
1. YAML generator marks only `.spec` and `.data` (in schema root) properties as required.

## Generator limitations
1. The YAML generator assumes that there is exactly one version in `crd.spec.versions`.
If there is more or less versions, the script will fail.
1. If some Go struct actually has a field named `Mesh7CodeGenInline`, the generator
will misbehave creating faulty schema.

## Possible improvements
Beside fixing the [generator limitations](#generator-limitations), some items can
be improved:
1. Make use of OpenAPI enums in [Go structs from proto enums](#go-structs-from-proto-enums).
1. Use some of the supported [JSON schema keywords](https://swagger.io/docs/specification/data-models/keywords/)
for better validation.
1. Auto add `+k8s:openapi-gen=true` annotations to all structures in `pkg/apis` automatically.
