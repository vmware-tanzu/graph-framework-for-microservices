package declarative_test

import (
	"api-gw/pkg/openapi/declarative"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("OpenAPI tests", func() {

	const (
		Uri         = "/v1alpha1/global-namespaces"
		ResourceUri = "/v1alpha1/global-namespaces/{id}"
		ListUri     = "/v1alpha1/global-namespaces/test"
	)

	It("should setup and load openapi file", func() {
		err := declarative.Load(spec)
		Expect(err).To(BeNil())

		Expect(declarative.Paths).To(HaveKey(Uri))
		Expect(declarative.Paths).To(HaveKey(ResourceUri))
	})

	It("should get extension value for kind and group", func() {
		kind := declarative.GetExtensionVal(declarative.Paths[Uri].Get, declarative.NexusKindName)
		Expect(kind).To(Equal("GlobalNamespace"))

		group := declarative.GetExtensionVal(declarative.Paths[Uri].Get, declarative.NexusGroupName)
		Expect(group).To(Equal("gns.vmware.org"))

		list := declarative.GetExtensionVal(declarative.Paths[ListUri].Get, declarative.NexusListEndpoint)
		Expect(list).To(Equal("true"))
	})

	It("should setup context for resource list operation", func() {
		ec := declarative.SetupContext(Uri, http.MethodGet, declarative.Paths[Uri].Get)

		expectedEc := declarative.EndpointContext{
			Context:      nil,
			SpecUri:      Uri,
			Method:       http.MethodGet,
			KindName:     "GlobalNamespace",
			ResourceName: "globalnamespaces",
			GroupName:    "gns.vmware.org",
			CrdName:      "globalnamespaces.gns.vmware.org",
			Params:       [][]string{},
			IdParam:      "",
			Single:       false,
			Uri:          "/apis/gns.vmware.org/v1/globalnamespaces",
		}

		Expect(ec).To(Equal(&expectedEc))
	})

	It("should setup context for resource get operation", func() {
		ec := declarative.SetupContext(ResourceUri, http.MethodGet, declarative.Paths[ResourceUri].Get)

		expectedEc := declarative.EndpointContext{
			Context:      nil,
			SpecUri:      ResourceUri,
			Method:       http.MethodGet,
			KindName:     "GlobalNamespace",
			ResourceName: "globalnamespaces",
			GroupName:    "gns.vmware.org",
			CrdName:      "globalnamespaces.gns.vmware.org",
			Params:       [][]string{{"{id}", "id"}},
			IdParam:      "{id}",
			Single:       true,
			Uri:          "/apis/gns.vmware.org/v1/globalnamespaces/:name",
		}

		Expect(ec).To(Equal(&expectedEc))
	})

	It("should check if resource get operation have an array response", func() {
		isArray := declarative.IsArrayResponse(declarative.Paths[Uri].Get)
		Expect(isArray).To(BeTrue())
	})

	It("should fail on nil operation when checking if response is array", func() {
		isArray := declarative.IsArrayResponse(nil)
		Expect(isArray).To(BeFalse())
	})
})
