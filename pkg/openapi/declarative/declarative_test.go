package declarative_test

import (
	"api-gw/pkg/config"
	"api-gw/pkg/openapi/declarative"
	"api-gw/pkg/server/echo_server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("OpenAPI tests", func() {
	It("should setup and load openapi file", func() {
		err := declarative.Load(spec)
		Expect(err).To(BeNil())

		Expect(declarative.Paths).To(HaveKey(Uri))
		Expect(declarative.Paths).To(HaveKey(ResourceUri))
	})

	It("should add resource get operation uri to apis list", func() {
		ec := declarative.SetupContext(Uri, http.MethodGet, declarative.Paths[Uri].Get)
		declarative.AddApisEndpoint(ec)

		Expect(declarative.ApisList).To(HaveKey(ec.Uri))
		Expect(declarative.ApisList[ec.Uri]).To(HaveKey(http.MethodGet))
		Expect(declarative.ApisList[ec.Uri]).ToNot(HaveKey(http.MethodPost))
		Expect(declarative.ApisList[ec.Uri][http.MethodGet]).To(BeEquivalentTo(map[string]interface{}{
			"group":  ec.GroupName,
			"kind":   ec.KindName,
			"params": []string{"projectId"},
			"uri":    ec.SpecUri,
		}))
	})

	It("should register declarative router", func() {
		config.Cfg = &config.Config{
			Server:             config.ServerConfig{},
			EnableNexusRuntime: true,
			BackendService:     "",
		}
		e := echo_server.NewEchoServer(config.Cfg)
		e.RegisterDeclarativeRouter()

		c := e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/gns.vmware.org/v1/globalnamespaces", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/gns.vmware.org/v1/globalnamespaces/:name", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces/:name"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodPut, "/apis/gns.vmware.org/v1/globalnamespaces", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodDelete, "/apis/gns.vmware.org/v1/globalnamespaces/:name", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces/:name"))
	})

	It("should parse schema for GlobalNamespace", func() {
		ec := declarative.SetupContext(Uri, http.MethodGet, declarative.Paths[Uri].Get)
		declarative.AddApisEndpoint(ec)

		Expect(declarative.ApisList).To(HaveKey("/apis/gns.vmware.org/v1/globalnamespaces"))
		Expect(declarative.ApisList["/apis/gns.vmware.org/v1/globalnamespaces"]).To(HaveKey("yaml"))

		expectedYaml := `apiVersion: gns.vmware.org/v1
kind: GlobalNamespace
metadata:
  labels:
    id: string
    projectId: string
  name: string
spec:
  api_discovery_enabled: true
  ca: string
  ca_type: PreExistingCA
  color: string
  description: string
  display_name: string
  domain_name: string
  match_conditions:
  - cluster:
      match: string
      type: string
    namespace:
      match: string
      type: string
  mtls_enforced: true
  name: string
  use_shared_gateway: true
  version: string
`
		Expect(declarative.ApisList["/apis/gns.vmware.org/v1/globalnamespaces"]["yaml"]).To(Equal(expectedYaml))
	})
})
