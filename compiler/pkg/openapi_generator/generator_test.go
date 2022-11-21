package openapi_generator_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	generator2 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/generator"
	generator "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator/test_data/openapi"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/validation/spec"
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var _ = Describe("Generator", func() {
	var (
		tmpDir string
		oldDir string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "generator-test-")
		oldDir = "test_data"
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		//err := os.RemoveAll(tmpDir)
		//Expect(err).NotTo(HaveOccurred())
	})

	It("00 creates schemas for proto", func() {
		ref := func(pkg string) spec.Ref {
			r, err := spec.NewRef(strings.ToLower(pkg))
			if err != nil {
				panic(err)
			}
			return r
		}
		gen, err := generator.NewGenerator(openapi.GetOpenAPIDefinitions(ref))
		Expect(err).NotTo(HaveOccurred())

		namePrefix := "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator/test_data"
		gen.SetNamePrefix(namePrefix)

		oneOfDefinition := common.OpenAPIDefinition{
			Schema: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"Mesh7CodeGenInline": {},
						"oneof_value_string": {
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
						"oneof_value_int": {
							SchemaProps: spec.SchemaProps{
								Type:   []string{"integer"},
								Format: "int32",
							},
						},
					},
				},
			},
		}
		err = gen.AddDefinition(fmt.Sprintf("%s/proto.isfoo_oneofvalue", namePrefix), oneOfDefinition)
		Expect(err).NotTo(HaveOccurred())
		err = gen.AddDefinition(fmt.Sprintf("%s/proto.isbar_oneofvalue", namePrefix), oneOfDefinition)
		Expect(err).NotTo(HaveOccurred())

		Expect(gen.ResolveRefs()).To(Succeed())

		tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"foowrapper"})
		// should pass the backward compatibility check
		Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
		compareTmpFileWithExpectedFile(tmpFile, "test_data/foowrapper.yaml")
	})

	It("01 creates schemas for simple types", func() {
		rawDefs := map[string]common.OpenAPIDefinition{
			getSchemaName("foo"): fooDefinition(),
			getSchemaName("fizz"): {
				Schema: spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"object"},
						Properties: map[string]spec.Schema{
							"baz": {
								SchemaProps: spec.SchemaProps{
									Type: []string{"string"},
								},
							},
							"baz_bar": {
								SchemaProps: spec.SchemaProps{
									Type: []string{"string"},
								},
							},
						},
						Required: []string{"baz", "baz_bar"},
					},
				},
			},
		}
		gen, err := generator.NewGenerator(rawDefs)
		Expect(err).NotTo(HaveOccurred())

		Expect(gen.ResolveRefs()).To(Succeed())

		tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"foo", "fizz"})
		Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
		compareTmpFileWithExpectedFile(tmpFile, "test_data/01_simple_schema.yaml")
	})

	Context("refs resolution", func() {
		It("02 resolves ref in property", func() {
			fooName := getSchemaName("foo")
			barName := getSchemaName("bar")
			fooRef, err := spec.NewRef(fooName)
			Expect(err).NotTo(HaveOccurred())
			rawDefs := map[string]common.OpenAPIDefinition{
				barName: {
					Schema: spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							Properties: map[string]spec.Schema{
								"foo": {
									SchemaProps: spec.SchemaProps{
										Ref: fooRef,
									},
								},
							},
						},
					},
					Dependencies: []string{fooName},
				},
				fooName: fooDefinition(),
			}
			gen, err := generator.NewGenerator(rawDefs)
			Expect(err).NotTo(HaveOccurred())

			Expect(gen.ResolveRefs()).To(Succeed())

			tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"bar"})
			Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
			compareTmpFileWithExpectedFile(tmpFile, "test_data/02_ref_in_property.yaml")
		})

		It("03 resolves ref in item - single schema", func() {
			fooName := getSchemaName("foo")
			barName := getSchemaName("bar")
			fooRef, err := spec.NewRef(fooName)
			Expect(err).NotTo(HaveOccurred())
			rawDefs := map[string]common.OpenAPIDefinition{
				barName: {
					Schema: spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							Properties: map[string]spec.Schema{
								"foo": {
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Ref: fooRef,
												},
											},
										},
									},
								},
							},
						},
					},
					Dependencies: []string{fooName},
				},
				fooName: fooDefinition(),
			}
			gen, err := generator.NewGenerator(rawDefs)
			Expect(err).NotTo(HaveOccurred())

			Expect(gen.ResolveRefs()).To(Succeed())

			tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"bar"})
			Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
			compareTmpFileWithExpectedFile(tmpFile, "test_data/03_ref_in_items_single.yaml")
		})

		It("04 resolves ref in item - multiple schemas", func() {
			fooName := getSchemaName("foo")
			barName := getSchemaName("bar")
			fooRef, err := spec.NewRef(fooName)
			Expect(err).NotTo(HaveOccurred())
			rawDefs := map[string]common.OpenAPIDefinition{
				barName: {
					Schema: spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							Properties: map[string]spec.Schema{
								"foo": {
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schemas: []spec.Schema{
												{
													SchemaProps: spec.SchemaProps{
														Ref: fooRef,
													},
												},
												{
													SchemaProps: spec.SchemaProps{
														Ref: fooRef,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Dependencies: []string{fooName},
				},
				fooName: fooDefinition(),
			}
			gen, err := generator.NewGenerator(rawDefs)
			Expect(err).NotTo(HaveOccurred())

			Expect(gen.ResolveRefs()).To(Succeed())

			tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"bar"})
			Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
			compareTmpFileWithExpectedFile(tmpFile, "test_data/04_ref_in_items_multiple.yaml")
		})

		It("05 resolves ref in additional property", func() {
			fooName := getSchemaName("foo")
			barName := getSchemaName("bar")
			fooRef, err := spec.NewRef(fooName)
			Expect(err).NotTo(HaveOccurred())
			rawDefs := map[string]common.OpenAPIDefinition{
				barName: {
					Schema: spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"object"},
										Properties: map[string]spec.Schema{
											"foo": {
												SchemaProps: spec.SchemaProps{
													Ref: fooRef,
												},
											},
										},
										AdditionalProperties: &spec.SchemaOrBool{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Type: []string{"object"},
													Properties: map[string]spec.Schema{
														"foo": {
															SchemaProps: spec.SchemaProps{
																Ref: fooRef,
															},
														},
													},
												},
											},
										},
									},
								},
							},
							AdditionalItems: &spec.SchemaOrBool{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"object"},
										Properties: map[string]spec.Schema{
											"foo": {
												SchemaProps: spec.SchemaProps{
													Ref: fooRef,
												},
											},
										},
									},
								},
							},
						},
					},
					Dependencies: []string{fooName},
				},
				fooName: fooDefinition(),
			}
			gen, err := generator.NewGenerator(rawDefs)
			Expect(err).NotTo(HaveOccurred())

			Expect(gen.ResolveRefs()).To(Succeed())

			tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"bar"})
			Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
			compareTmpFileWithExpectedFile(tmpFile, "test_data/05_ref_in_additional_property.yaml")
		})
	})

	Context("enums handling", func() {
		It("06 handles enum in properties", func() {
			barName := getSchemaName("bar")
			rawDefs := map[string]common.OpenAPIDefinition{
				barName: {
					Schema: spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							Properties: map[string]spec.Schema{
								"foo": {
									SchemaProps: spec.SchemaProps{
										Type:        []string{"integer"},
										Format:      "int32",
										Description: "Mesh7CodeGenOpenAPIEnum",
									},
								},
							},
						},
					},
				},
			}
			gen, err := generator.NewGenerator(rawDefs)
			Expect(err).NotTo(HaveOccurred())

			Expect(gen.ResolveRefs()).To(Succeed())

			tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"bar"})
			Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
			compareTmpFileWithExpectedFile(tmpFile, "test_data/06_enum_in_property.yaml")
		})

		It("07 handles enum in array properties", func() {
			barName := getSchemaName("bar")
			rawDefs := map[string]common.OpenAPIDefinition{
				barName: {
					Schema: spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							Properties: map[string]spec.Schema{
								"foo": {
									SchemaProps: spec.SchemaProps{
										Type:        []string{"array"},
										Description: "Mesh7CodeGenOpenAPIEnum",
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Type:   []string{"integer"},
													Format: "int32",
												},
											},
										},
									},
								},
								"bar": {
									SchemaProps: spec.SchemaProps{
										Type:        []string{"array"},
										Description: "Mesh7CodeGenOpenAPIEnum",
										Items: &spec.SchemaOrArray{
											Schemas: []spec.Schema{{
												SchemaProps: spec.SchemaProps{
													Type:   []string{"integer"},
													Format: "int32",
												},
											}},
										},
									},
								},
							},
						},
					},
				},
			}
			gen, err := generator.NewGenerator(rawDefs)
			Expect(err).NotTo(HaveOccurred())

			Expect(gen.ResolveRefs()).To(Succeed())

			tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"bar"})
			Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
			compareTmpFileWithExpectedFile(tmpFile, "test_data/07_enum_in_array_property.yaml")
		})
	})

	It("08 adds kubernetes flags", func() {
		rawDefs := map[string]common.OpenAPIDefinition{
			getSchemaName("fizz"): {
				Schema: spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: []string{"object"},
						Properties: map[string]spec.Schema{
							"emptyObject": {
								SchemaProps: spec.SchemaProps{
									Type: []string{"object"},
								},
							},
							"anyOfProp": {
								SchemaProps: spec.SchemaProps{
									AnyOf: []spec.Schema{
										{SchemaProps: spec.SchemaProps{Type: []string{"integer"}}},
										{SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
									},
								},
							},
							"nested": {
								SchemaProps: spec.SchemaProps{
									Type: []string{"object"},
									Properties: map[string]spec.Schema{
										"emptyObject": {
											SchemaProps: spec.SchemaProps{
												Type: []string{"object"},
											},
										},
										"anyOfProp": {
											SchemaProps: spec.SchemaProps{
												AnyOf: []spec.Schema{
													{SchemaProps: spec.SchemaProps{Type: []string{"integer"}}},
													{SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
												},
											},
										},
									},
								},
							},
							"array": {
								SchemaProps: spec.SchemaProps{
									Type: []string{"array"},
									Items: &spec.SchemaOrArray{
										Schemas: []spec.Schema{
											{
												SchemaProps: spec.SchemaProps{
													Type: []string{"object"},
												},
											},
											{
												SchemaProps: spec.SchemaProps{
													AnyOf: []spec.Schema{
														{SchemaProps: spec.SchemaProps{Type: []string{"integer"}}},
														{SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		gen, err := generator.NewGenerator(rawDefs)
		Expect(err).NotTo(HaveOccurred())

		Expect(gen.ResolveRefs()).To(Succeed())

		tmpFile := createFileWithEmptyYAMLDefinitions(tmpDir, []string{"fizz"})
		Expect(gen.UpdateYAMLs(tmpDir, oldDir, false)).To(Succeed())
		compareTmpFileWithExpectedFile(tmpFile, "test_data/08_kubernetes_flags.yaml")
	})

	Context("checks backward compatibility", func() {
		It("should detect any incompatible changes to crds successfully", func() {
			var inCompatibleCRDBuffer []*bytes.Buffer
			crd := &extensionsv1.CustomResourceDefinition{}
			parts := strings.Split(baseSpec, "---")
			err := yaml.Unmarshal([]byte(parts[1]), crd)
			Expect(err).ToNot(HaveOccurred())

			// change in annotation
			nexusAnnotation := crd.ObjectMeta.Annotations["nexus"]
			n := generator2.NexusAnnotation{}
			err = json.Unmarshal([]byte(nexusAnnotation), &n)
			Expect(err).ToNot(HaveOccurred())

			// modify singleton field to `true`
			n.IsSingleton = true
			b, err := json.Marshal(n)
			Expect(err).ToNot(HaveOccurred())
			crd.ObjectMeta.Annotations["nexus"] = string(b)

			// change in type
			doubleVal := crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["name"]
			doubleVal.Type = "integer"
			doubleVal.Format = "integer"
			crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["name"] = doubleVal

			// adding a new field
			crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["new-field"] = extensionsv1.JSONSchemaProps{Type: "string"}

			// deleting a field
			delete(crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["status"].Properties, "nexus")

			inCompatibleCRDBuffer, err = generator.CheckBackwardCompatibility(inCompatibleCRDBuffer, *crd, []byte(baseSpec))
			Expect(err).ToNot(HaveOccurred())
			for _, c := range inCompatibleCRDBuffer {
				Expect(c.String()).To(Equal("detected changes in model stored in ignorechilds.gns.tsm.tanzu.vmware.com\n\nspec changes: " +
					"\n/spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties/name/type\n  ± value change\n    - string\n    + integer\n  \n\nstatus changes: " +
					"\n/spec/versions/name=v1/schema/openAPIV3Schema/properties/status\n  - one field removed:\n    properties:\n      nexus:\n        type: object\n       " +
					" required:\n        - sourceGeneration\n        - remoteGeneration\n        properties:\n          remoteGeneration:\n            type: integer\n            " +
					"format: int64\n          sourceGeneration:\n            type: integer\n            format: int64\n    \n  " +
					"\n\nnexus annotation changes: \n/is_singleton\n  ± value change\n    - false\n    + true\n  \n\n"))
			}
		})

		It("should check for incompatible changes when forcing flag enable/disable between new and existing CRDs", func() {
			rawDefs := map[string]common.OpenAPIDefinition{
				getSchemaName("foo"): {
					Schema: spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							Properties: map[string]spec.Schema{
								"spec": {
									SchemaProps: spec.SchemaProps{
										Type: []string{"object"},
										Properties: map[string]spec.Schema{
											"changePassword": {
												SchemaProps: spec.SchemaProps{
													Type:   []string{"string"},
													Format: "string",
												},
											},
											"name": {
												SchemaProps: spec.SchemaProps{
													Type:   []string{"integer"},
													Format: "int32",
												},
											},
										},
										Required: []string{"changePassword"},
									},
								},
							},
						},
					},
				},
			}
			gen, err := generator.NewGenerator(rawDefs)
			Expect(err).NotTo(HaveOccurred())

			Expect(gen.ResolveRefs()).To(Succeed())

			// create a new CRD with following changes
			// removing a required field `password`.
			// adding a new field `changePassword`
			// modifying the `name` type from `string` to `int32`
			tmpFile := fmt.Sprintf("%s/%s.yaml", tmpDir, "foos")
			err = os.WriteFile(tmpFile, []byte(newFooCRD), 0665)
			Expect(err).NotTo(HaveOccurred())

			// should be unsuccessful due to the incompatibility of the new and old CDs + forcing an upgrade=false
			err = gen.UpdateYAMLs(tmpDir, oldDir, false)
			Expect(err.Error()).To(Equal("datamodel upgrade failed due to backward incompatible changes:\n " +
				"[detected changes in model stored in foos\n\nspec changes: " +
				"\n/spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties\n  - " +
				"one field removed:\n   " +
				" password:\n      " +
				"type: string\n      " +
				"format: string\n    " +
				"\n  \n\n/spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties/name/format\n  " +
				"± value change\n    " +
				"- string\n    " +
				"+ int32\n  " +
				"\n\n/spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties/name/type\n " +
				" ± value change\n    " +
				"- string\n    " +
				"+ integer\n  " +
				"\n\n/spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/required/0\n " +
				" ± value change\n    " +
				"- password\n    " +
				"+ changePassword\n  \n\n]"))

			// should succeed even if there are incompatible changes when force upgrade is true.
			err = gen.UpdateYAMLs(tmpDir, oldDir, true)
			Expect(err).To(BeNil())
		})

	})
})

var baseSpec = `
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: >
      {"name":"gns.IgnoreChild","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com","gnses.gns.tsm.tanzu.vmware.com"],"is_singleton":false,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: ignorechilds.gns.tsm.tanzu.vmware.com
spec:
  conversion:
    strategy: None
  group: gns.tsm.tanzu.vmware.com
  versions:
    - name: v1
      schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                name:
                  type: string
              required:
                - name
              type: object
            status:
              properties:
                nexus:
                  properties:
                    remoteGeneration:
                      format: int64
                      type: integer
                    sourceGeneration:
                      format: int64
                      type: integer
                  required:
                    - sourceGeneration
                    - remoteGeneration
                  type: object
              type: object
          type: object
`

var newFooCRD = `---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":"gns.Foo","hierarchy":["roots.root.tsm.tanzu.vmware.com","configs.config.tsm.tanzu.vmware.com","gnses.gns.tsm.tanzu.vmware.com"],"is_singleton":false,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: foos
spec:
  conversion:
    strategy: None
  group: test.it
  names:
    kind: Foo
    listKind: FooList
    plural: foos
    shortNames:
    - foo
    singular: foo
  scope: Cluster
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: null
  storedVersions:
  - v1
`
