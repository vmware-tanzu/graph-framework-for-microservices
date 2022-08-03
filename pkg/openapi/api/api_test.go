package api_test

import (
	"encoding/json"
	yamlv1 "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"api-gw/pkg/model"
	"api-gw/pkg/openapi/api"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

var _ = Describe("OpenAPI tests", func() {
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
			[]string{"roots.orgchart.vmware.org"}, nil, false, "my custom description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		api.New()
		api.AddPath(restUri)
		Expect(api.Schema.Paths[restUri.Uri].Get.Parameters[0].Value.Name).To(Equal("orgchart.Leader"))
		Expect(api.Schema.Paths[restUri.Uri].Get.Parameters[0].Value.Description).To(Equal("my custom description"))
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
			[]string{}, nil, false, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		api.New()
		api.AddPath(restUri)
		Expect(api.Schema.Paths[restUri.Uri].Get.Parameters[0].Value.Name).To(Equal("orgchart.Leader"))
		Expect(api.Schema.Paths[restUri.Uri].Get.Parameters[0].Value.Description).
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
			[]string{}, nil, false, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		model.ConstructMapCRDTypeToSpec(model.Upsert, "leaders.orgchart.vmware.org", crd.Spec)
		api.New()
		api.AddPath(restUri)
		Expect(api.Schema.Paths[restUri.Uri].Get).To(Not(BeNil()))
	})
})
