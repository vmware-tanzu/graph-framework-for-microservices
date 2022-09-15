package echo_server_test

import (
	"api-gw/pkg/config"
	"api-gw/pkg/server/echo_server"
	"net/http"

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

})
