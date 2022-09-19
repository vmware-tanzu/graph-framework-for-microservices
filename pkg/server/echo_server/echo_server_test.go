package echo_server_test

import (
	"api-gw/pkg/config"
	"api-gw/pkg/server/echo_server"
	"github.com/labstack/echo/v4"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Echo server tests", func() {
	var e *echo_server.EchoServer

	It("should init echo server", func() {
		config.Cfg = &config.Config{
			Server:             config.ServerConfig{},
			EnableNexusRuntime: true,
			BackendService:     "",
		}
		e = echo_server.NewEchoServer(config.Cfg)
	})

	It("should register CrdRouter for globalnamespaces.gns.vmware.org", func() {
		e.RegisterCrdRouter("globalnamespaces.gns.vmware.org")

		c := e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/gns.vmware.org/v1/globalnamespaces/:name", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces/:name"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/gns.vmware.org/v1/globalnamespaces", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodPost, "/apis/gns.vmware.org/v1/globalnamespaces", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodDelete, "/apis/gns.vmware.org/v1/globalnamespaces/:name", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces/:name"))
	})

	It("should start echo server", func() {
		stopCh := make(chan struct{})
		e := echo_server.InitEcho(stopCh, &config.Config{
			Server: config.ServerConfig{
				HttpPort: "0",
			},
			EnableNexusRuntime: true,
			BackendService:     "http://localhost",
		})
		e.StopServer()
	})

	It("should start echo server and restart through channel", func() {
		stopCh := make(chan struct{})
		e := echo_server.InitEcho(stopCh, &config.Config{
			Server: config.ServerConfig{
				HttpPort: "0",
			},
			EnableNexusRuntime: true,
			BackendService:     "http://localhost",
		})

		stopCh <- struct{}{}

		e.StopServer()
	})

	It("should get nexus crd context", func() {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		nexusCrdContext := e.GetNexusCrdContext("test.vmware.org", "vmware.org", "tests")

		var actualContext *echo_server.NexusContext
		c.SetHandler(
			nexusCrdContext(
				func(c echo.Context) error {
					actualContext = c.(*echo_server.NexusContext)
					return c.NoContent(200)
				},
			),
		)
		err := c.Handler()(c)
		Expect(err).NotTo(HaveOccurred())

		Expect(actualContext.Resource).To(Equal("tests"))
		Expect(actualContext.GroupName).To(Equal("vmware.org"))
		Expect(actualContext.CrdType).To(Equal("test.vmware.org"))
	})

	It("should get nexus context", func() {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		restUri := nexus.RestURIs{
			Uri:     "/test",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		codes := map[nexus.ResponseCode]nexus.HTTPResponse{
			http.StatusOK: {
				Description: "description",
			},
		}

		nexusCrdContext := e.GetNexusContext(restUri, codes)
		var actualContext *echo_server.NexusContext
		c.SetHandler(
			nexusCrdContext(
				func(c echo.Context) error {
					actualContext = c.(*echo_server.NexusContext)
					return c.NoContent(200)
				},
			),
		)
		err := c.Handler()(c)
		Expect(err).NotTo(HaveOccurred())

		Expect(actualContext.NexusURI).To(Equal("/test"))
		Expect(actualContext.Codes[http.StatusOK].Description).To(Equal("description"))
	})

})
