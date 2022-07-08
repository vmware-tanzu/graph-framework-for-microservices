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
	)

	It("should setup and load openapi file", func() {
		err := declarative.Load(spec)
		Expect(err).To(BeNil())

		Expect(declarative.Paths).To(HaveKey(Uri))
		Expect(declarative.Paths).To(HaveKey(ResourceUri))
	})

	It("should add resource get operation uri to apis list", func() {
		ec := declarative.SetupContext(Uri, http.MethodGet, declarative.Paths[Uri].Get)
		declarative.AddApisEndpoint(ec)

		var params []string
		Expect(declarative.ApisList).To(HaveKey(ec.Uri))
		Expect(declarative.ApisList[ec.Uri]).To(HaveKey(http.MethodGet))
		Expect(declarative.ApisList[ec.Uri]).ToNot(HaveKey(http.MethodPost))
		Expect(declarative.ApisList[ec.Uri][http.MethodGet]).To(BeEquivalentTo(map[string]interface{}{
			"group":  ec.GroupName,
			"kind":   ec.KindName,
			"params": params,
			"uri":    ec.SpecUri,
		}))
	})
})
