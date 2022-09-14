package openapi_generator_test

import (
	"fmt"
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	generator "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/openapi_generator"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/openapi_generator/test_data/openapi"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/pkg/validation/spec"
)

var _ = Describe("Generator", func() {
	var (
		tmpDir string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "generator-test-")
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

		namePrefix := "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/openapi_generator/test_data"
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
})
