package combined_test

import (
	"api-gw/pkg/model"
	"api-gw/pkg/openapi/api"
	"api-gw/pkg/openapi/combined"
	"api-gw/pkg/openapi/declarative"
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	yamlv1 "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Combined OpenAPI tests", func() {
	It("should setup and load openapi file", func() {
		err := declarative.Load(spec)
		Expect(err).To(BeNil())

		Expect(declarative.Paths).To(HaveKey(Uri))
		Expect(declarative.Paths).To(HaveKey(ResourceUri))
	})

	It("should create new datamodel", func() {
		Expect(api.Schemas).To(BeEmpty())
		api.New("vmware.org")
		Expect(api.Schemas["vmware.org"].Info.Title).To(Equal("Nexus API GW APIs"))

		unstructuredObj := unstructured.Unstructured{
			Object: map[string]interface{}{
				"spec": map[string]interface{}{
					"title": "VMWare Datamodel",
				},
			},
		}

		model.ConstructDatamodel(model.Upsert, "vmware2.org", &unstructuredObj)
		api.New("vmware2.org")
		Expect(api.Schemas["vmware2.org"].Info.Title).To(Equal("VMWare Datamodel"))
	})

	It("should add custom description to node", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leader/{orgchart.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}

		crdJson, err := yamlv1.YAMLToJSON([]byte(crdExample))
		Expect(err).NotTo(HaveOccurred())
		var crd apiextensionsv1.CustomResourceDefinition
		err = json.Unmarshal(crdJson, &crd)
		Expect(err).NotTo(HaveOccurred())

		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
			[]string{"roots.orgchart.vmware.org"}, nil, nil, false, "my custom description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		api.New("vmware.org")
		api.AddPath(restUri, "vmware.org")
		Expect(api.Schemas["vmware.org"].Paths[restUri.Uri].Get.Parameters[0].Value.Name).To(Equal("orgchart.Leader"))
		Expect(api.Schemas["vmware.org"].Paths[restUri.Uri].Get.Parameters[0].Value.Description).To(Equal("my custom description"))
	})

	It("should combine openapi specs", func() {
		schema := combined.CombinedSpecs()

		pathItem := schema.Paths.Find("/leader/{orgchart.Leader}")
		Expect(pathItem).ToNot(BeNil())

		pathItem = schema.Paths.Find("/v1alpha1/project/{projectId}/global-namespaces")
		Expect(pathItem).ToNot(BeNil())
	})

	It("should combine openapi specs with additional components", func() {
		s := api.Schemas["vmware.org"]
		s.Components.SecuritySchemes = openapi3.SecuritySchemes{
			"BasicAuth": {
				Value: &openapi3.SecurityScheme{
					Type:   "http",
					Scheme: "basic",
				},
			},
		}
		s.Components.Examples = openapi3.Examples{
			"example": {
				Value: &openapi3.Example{},
			},
		}
		s.Components.Links = openapi3.Links{
			"example": {
				Value: &openapi3.Link{},
			},
		}
		s.Components.Callbacks = openapi3.Callbacks{
			"example": {
				Value: &openapi3.Callback{},
			},
		}
		s.Components.Headers = openapi3.Headers{
			"example": {
				Value: &openapi3.Header{},
			},
		}
		s.Components.Parameters = openapi3.ParametersMap{
			"example": &openapi3.ParameterRef{
				Value: &openapi3.Parameter{Name: "example"},
			},
		}
		api.Schemas["vmware.org"] = s

		schema := combined.CombinedSpecs()

		pathItem := schema.Paths.Find("/leader/{orgchart.Leader}")
		Expect(pathItem).ToNot(BeNil())

		pathItem = schema.Paths.Find("/v1alpha1/project/{projectId}/global-namespaces")
		Expect(pathItem).ToNot(BeNil())
	})
})
