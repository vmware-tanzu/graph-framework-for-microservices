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
})
