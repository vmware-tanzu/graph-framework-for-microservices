package api_test

import (
	"encoding/json"

	yamlv1 "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"api-gw/pkg/model"
	"api-gw/pkg/openapi/api"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

var _ = Describe("OpenAPI tests", func() {
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

	It("should add default description to node if custom is not present", func() {
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
			[]string{}, nil, nil, false, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		api.New("vmware.org")
		api.AddPath(restUri, "vmware.org")
		Expect(api.Schemas["vmware.org"].Paths[restUri.Uri].Get.Parameters[0].Value.Name).To(Equal("orgchart.Leader"))
		Expect(api.Schemas["vmware.org"].Paths[restUri.Uri].Get.Parameters[0].Value.Description).
			To(Equal("Name of the orgchart.Leader node"))
	})

	It("should add list endpoint", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		}

		crdJson, err := yamlv1.YAMLToJSON([]byte(crdExample))
		Expect(err).NotTo(HaveOccurred())
		var crd apiextensionsv1.CustomResourceDefinition
		err = json.Unmarshal(crdJson, &crd)
		Expect(err).NotTo(HaveOccurred())

		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
			[]string{}, nil, nil, false, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		api.New("vmware.org")
		api.AddPath(restUri, "vmware.org")
		Expect(api.Schemas["vmware.org"].Paths[restUri.Uri].Get).To(Not(BeNil()))
	})

	It("should test Recreate func", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		}

		crdJson, err := yamlv1.YAMLToJSON([]byte(crdExample))
		Expect(err).NotTo(HaveOccurred())
		var crd apiextensionsv1.CustomResourceDefinition
		err = json.Unmarshal(crdJson, &crd)
		Expect(err).NotTo(HaveOccurred())

		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
			[]string{}, nil, nil, false, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		model.ConstructMapCRDTypeToRestUris(model.Upsert, "leaders.orgchart.vmware.org", nexus.RestAPISpec{
			Uris: []nexus.RestURIs{
				restUri,
			},
		})
		api.Recreate()
		Expect(api.Schemas).To(HaveKey("vmware.org"))
		Expect(api.Schemas["vmware.org"].Components.Responses).To(HaveKey("Listorgchart.Leader"))
	})

	It("should test update notification for new crd", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		}

		crdJson, err := yamlv1.YAMLToJSON([]byte(crdExample))
		Expect(err).NotTo(HaveOccurred())
		var crd apiextensionsv1.CustomResourceDefinition
		err = json.Unmarshal(crdJson, &crd)
		Expect(err).NotTo(HaveOccurred())

		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
			[]string{}, nil, nil, false, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		model.ConstructMapCRDTypeToRestUris(model.Upsert, "leaders.orgchart.vmware.org", nexus.RestAPISpec{
			Uris: []nexus.RestURIs{
				restUri,
			},
		})

		api.Recreate()

		unstructuredObj := unstructured.Unstructured{
			Object: map[string]interface{}{
				"spec": map[string]interface{}{
					"title": "VMWare Datamodel",
				},
			},
		}

		go api.DatamodelUpdateNotification()
		model.ConstructDatamodel(model.Delete, "vmware.org", &unstructuredObj)

		Eventually(func() bool {
			if api.Schemas["vmware.org"].Info.Title == "VMWare Datamodel" {
				return true
			}

			return false
		})
	})
})
