package openapi_generator_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	pkg_generator "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/generator"
	generator "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator/test_data/openapi"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/validation/spec"
)

var _ = Describe("Generator", func() {
	var (
		tmpDir string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "generator-test-")
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
		Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
		compareTmpFileWithExpectedFile(tmpFile, "test_data/00_proto_schema.yaml")
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
		Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
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
			Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
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
			Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
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
			Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
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
			Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
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
			Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
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
			Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
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
		Expect(gen.UpdateYAMLs(tmpDir)).To(Succeed())
		compareTmpFileWithExpectedFile(tmpFile, "test_data/08_kubernetes_flags.yaml")
	})

	Context("checks backward compatibility", func() {
		It("should fail when the spec is changed", func() {
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

			tmpFile := fmt.Sprintf("%s/%s.yaml", tmpDir, "foos")
			err = os.WriteFile(tmpFile, []byte(newFooCRD), 0665)
			Expect(err).NotTo(HaveOccurred())

			// should be unsuccessful due to the incompatibility of the new and old CDs + forcing an upgrade=false
			err = gen.UpdateYAMLs(tmpDir)
			Expect(err).To(BeNil())

			oldCRDDir, err := exampleFileTempTestDir("foos.yaml")
			Expect(err).NotTo(HaveOccurred())

			err = generator.CheckBackwardCompatibility(oldCRDDir, tmpDir, false)
			cleanTempTestDir(oldCRDDir)
			Expect(err).To(HaveOccurred())
		})

		Context("should check nexus annotation and crd name compatibility", func() {
			var (
				crd extensionsv1.CustomResourceDefinition
				f   *os.File
				err error
			)

			AfterEach(func() {
				f.Close()
			})

			BeforeEach(func() {
				tmpFile := fmt.Sprintf("%s/%s.yaml", tmpDir, "foos")
				err := os.WriteFile(tmpFile, []byte(newFooCRD), 0665)
				Expect(err).NotTo(HaveOccurred())

				fooContent, err := os.ReadFile(tmpFile)
				Expect(err).To(BeNil())

				err = yaml.Unmarshal(fooContent, &crd)
				Expect(err).To(BeNil())

				f, err = os.OpenFile(tmpFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				Expect(err).To(BeNil())
			})

			It("should fail when the nexus annotation is changed", func() {
				ann := crd.Annotations["nexus"]
				nexusAnnotation := &pkg_generator.NexusAnnotation{}
				err = json.Unmarshal([]byte(ann), &nexusAnnotation)
				Expect(err).To(BeNil())

				// modify singleton field to `true` which leads to DM incompatible with previous version
				nexusAnnotation.IsSingleton = true
				annotationInByte, err := yaml.Marshal(nexusAnnotation)
				Expect(err).To(BeNil())

				crd.Annotations["nexus"] = string(annotationInByte)
				serialized, err := yaml.Marshal(crd)
				Expect(err).To(BeNil())

				_, err = f.Write(serialized)
				Expect(err).To(BeNil())

				// should be unsuccessful due to the incompatibility of the new and old CDs + forcing an upgrade=false
				oldCRDDir, err := exampleFileTempTestDir("foos.yaml")
				Expect(err).NotTo(HaveOccurred())

				err = generator.CheckBackwardCompatibility(oldCRDDir, tmpDir, false)
				cleanTempTestDir(oldCRDDir)
				Expect(err).To(HaveOccurred())
			})

			It("should fail when the CRD name is not matched", func() {
				crd.Name = "foos.new_test.it"
				serialized, err := yaml.Marshal(crd)
				Expect(err).To(BeNil())

				_, err = f.Write(serialized)
				Expect(err).To(BeNil())

				// should be unsuccessful due to the incompatibility of the new and old CDs + forcing an upgrade=false
				oldCRDDir, err := exampleFileTempTestDir("foos.yaml")
				Expect(err).NotTo(HaveOccurred())

				err = generator.CheckBackwardCompatibility(oldCRDDir, tmpDir, false)
				cleanTempTestDir(oldCRDDir)
				Expect(err.Error()).To(Equal("datamodel upgrade failed due to incompatible datamodel changes: \n \"foos\" is deleted\n"))
			})
		})

		It("should check for incompatible changes if a node is not found in the new CRDs directory", func() {
			oldCRDDir, err := exampleFileTempTestDir("zoos.yaml")
			Expect(err).NotTo(HaveOccurred())

			// should fail when CRD/Node is removed in the new list on force=false
			err = generator.CheckBackwardCompatibility(oldCRDDir, tmpDir, false)
			Expect(err.Error()).To(Equal("datamodel upgrade failed due to incompatible datamodel changes: \n \"foos\" is deleted\n"))
			cleanTempTestDir(oldCRDDir)

			// should not fail when CRD/Node is removed in the new list on force=true
			err = generator.CheckBackwardCompatibility(oldCRDDir, tmpDir, true)
			Expect(err).To(BeNil())
			cleanTempTestDir(oldCRDDir)
		})

		It("should not fail when the existing CRDs directory is empty", func() {
			// shouldn't fail when no crds exists
			emptyDir, err := exampleTestDir()
			Expect(err).NotTo(HaveOccurred())
			err = generator.CheckBackwardCompatibility(emptyDir, tmpDir, false)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func exampleTestDir() (string, error) {
	dir, err := os.MkdirTemp("", "compatibility-test-")
	if err != nil {
		return "", err
	}

	return dir, nil
}

func exampleFileTempTestDir(fileName string) (string, error) {
	dir, err := exampleTestDir()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile("test_data/foos.yaml")
	if err != nil {
		return "", err
	}

	file := filepath.Join(dir, fileName)
	err = os.WriteFile(file, data, 0666)
	if err != nil {
		return "", err
	}
	return dir, err
}

func cleanTempTestDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println(err)
	}
}

var newFooCRD = `
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
