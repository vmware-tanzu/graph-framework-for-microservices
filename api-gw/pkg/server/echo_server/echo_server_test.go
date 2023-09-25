package echo_server_test

import (
	"api-gw/pkg/client"
	"api-gw/pkg/config"
	"api-gw/pkg/model"
	"api-gw/pkg/server/echo_server"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apinexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/api.nexus.vmware.com/v1"
	confignexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/config.nexus.vmware.com/v1"
	domain_nexus_org "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/domain.nexus.vmware.com/v1"
	nexus_client "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/nexus-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("Echo server tests", func() {
	var e *echo_server.EchoServer

	It("should watch for OpenApi spec file creation", func() {
		openApiSpecDir := "testDir"
		openApiSpecFile := "testDir/testFile"
		echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter = 0

		stopCh := make(chan struct{})
		receivedSignalFromStopCh := false
		go func() {
			for {
				select {
				case <-stopCh:
					receivedSignalFromStopCh = true
					return
				}
			}
		}()
		err := os.Mkdir(openApiSpecDir, 0755)
		defer os.RemoveAll(openApiSpecDir)
		Expect(err).To(BeNil())
		echo_server.WatchForOpenApiSpecChanges(stopCh, openApiSpecDir, openApiSpecFile)
		time.Sleep(5 * time.Second)
		f, err := os.Create(openApiSpecFile)
		Expect(err).To(BeNil())
		f.Sync()
		defer f.Close()

		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			if receivedSignalFromStopCh && echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter == 1 {
				break
			}
		}
		Expect(echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter).To(Equal(1))
		Expect(receivedSignalFromStopCh).To(Equal(true))
	})
	It("should watch for OpenApi spec file update", func() {
		openApiSpecDir := "testDir"
		openApiSpecFile := "testDir/testFile"
		echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter = 0

		stopCh := make(chan struct{})
		receivedSignalFromStopCh := false
		go func() {
			for {
				select {
				case <-stopCh:
					receivedSignalFromStopCh = true
					return
				}
			}
		}()

		err := os.Mkdir(openApiSpecDir, 0755)
		defer os.RemoveAll(openApiSpecDir)
		Expect(err).To(BeNil())
		f, err := os.Create(openApiSpecFile)
		Expect(err).To(BeNil())
		f.Sync()
		defer f.Close()

		echo_server.WatchForOpenApiSpecChanges(stopCh, openApiSpecDir, openApiSpecFile)
		time.Sleep(5 * time.Second)
		bytesWritten, err := f.WriteString("writes\n")
		Expect(err).To(BeNil())
		Expect(bytesWritten).ToNot(Equal(0))
		f.Sync()
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			if receivedSignalFromStopCh && echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter == 1 {
				break
			}
		}
		Expect(echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter).To(Equal(1))
		Expect(receivedSignalFromStopCh).To(Equal(true))
	})
	It("should watch for OpenApi spec file deletion", func() {
		openApiSpecDir := "testDir"
		openApiSpecFile := "testDir/testFile"
		echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter = 0

		stopCh := make(chan struct{})
		receivedSignalFromStopCh := false
		go func() {
			for {
				select {
				case <-stopCh:
					receivedSignalFromStopCh = true
					return
				}
			}
		}()

		err := os.Mkdir(openApiSpecDir, 0755)
		defer os.RemoveAll(openApiSpecDir)
		Expect(err).To(BeNil())
		f, err := os.Create(openApiSpecFile)
		Expect(err).To(BeNil())
		f.Sync()
		f.Close()

		echo_server.WatchForOpenApiSpecChanges(stopCh, openApiSpecDir, openApiSpecFile)
		time.Sleep(5 * time.Second)
		os.Remove(openApiSpecFile)

		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			if receivedSignalFromStopCh && echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter == 1 {
				break
			}
		}
		Expect(echo_server.HttpServerRestartFromOpenApiSpecUpdateCounter).To(Equal(1))
		Expect(receivedSignalFromStopCh).To(Equal(true))
	})

	It("should init echo server", func() {

		client.NexusClient = nexus_client.NewFakeClient()
		_, err := client.NexusClient.Api().CreateNexusByName(context.TODO(), &apinexusv1.Nexus{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = client.NexusClient.Config().CreateConfigByName(context.TODO(), &confignexusv1.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "943ea6107388dc0d02a4c4d861295cd2ce24d551",
				Labels: map[string]string{
					"nexus/display_name": "default",
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		config.Cfg = &config.Config{
			Server:             config.ServerConfig{},
			EnableNexusRuntime: true,
			BackendService:     "",
		}
		e = echo_server.NewEchoServer(config.Cfg, &kubernetes.Clientset{}, &nexus_client.Clientset{})

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

	It("should return true when origin is present", func() {
		model.CorsConfigOrigins["default"] = []string{"http://testdomain"}

		res, err := echo_server.CheckCorsOrigin("http://localhost")
		Expect(res).To(BeFalse())
		Expect(err).NotTo(HaveOccurred())

		res, err = echo_server.CheckCorsOrigin("http://testdomain")
		Expect(res).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
		//removing the added object
		delete(model.CorsConfigOrigins, "default")

		res, err = echo_server.CheckCorsOrigin("http://testdomain")
		Expect(res).To(BeFalse())
		Expect(err).NotTo(HaveOccurred())

	})

	It("should register cors middleware when event is recieved", func() {

		e.RegisterCrdRouter("globalnamespaces.gns.vmware.org")
		corsObject := domain_nexus_org.CORSConfig{
			Spec: domain_nexus_org.CORSConfigSpec{
				Origins: []string{"http://tester"},
				Headers: []string{""},
			},
			ObjectMeta: v1.ObjectMeta{
				Name: "test",
			},
		}

		//check if corsOrigin is setup
		corsUpdate := &model.CorsNodeEvent{Cors: corsObject, Type: model.Upsert}
		echo_server.HandleCorsNodeUpdate(corsUpdate, e.Echo)
		Eventually(func() bool {
			res, _ := echo_server.CheckCorsOrigin("http://tester")
			return res
		})

		//delete should remove the node
		corsUpdate = &model.CorsNodeEvent{Cors: corsObject, Type: model.Delete}
		delete(model.CorsConfigOrigins, "test")
		echo_server.HandleCorsNodeUpdate(corsUpdate, e.Echo)
		Eventually(func() bool {
			res, _ := echo_server.CheckCorsOrigin("http://tester")
			if res {
				return false
			}
			return true
		})

		err := echo_server.HandleCorsNodeUpdate(nil, e.Echo)
		Expect(err).To(HaveOccurred())

	})

	It("should start echo server", func() {
		stopCh := make(chan struct{})
		e := echo_server.InitEcho(stopCh, &config.Config{
			Server: config.ServerConfig{
				HttpPort: "0",
			},
			EnableNexusRuntime: true,
			BackendService:     "http://localhost",
		}, &kubernetes.Clientset{},
			&nexus_client.Clientset{})
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
		}, &kubernetes.Clientset{},
			&nexus_client.Clientset{})

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

	It("should register TSM Routes", func() {
		e.RegisterCosmosAdminRoutes()

		c := e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodPut, "/v0/tenants/instance", c)
		Expect(c.Path()).To(Equal("/v0/tenants/instance"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/v0/version", c)
		Expect(c.Path()).To(Equal("/v0/version"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodPost, "/v0/users/login", c)
		Expect(c.Path()).To(Equal("/v0/users/login"))
	})

	It("should start echo server and timeout on port checking", func() {
		log.StandardLogger().ExitFunc = func(i int) {
			Expect(i).To(Equal(1))
		}

		stopCh := make(chan struct{})
		e := echo_server.InitEcho(stopCh, &config.Config{
			Server: config.ServerConfig{
				HttpPort: "0",
			},
			EnableNexusRuntime: true,
			BackendService:     "http://localhost",
		}, &kubernetes.Clientset{}, &nexus_client.Clientset{})

		listen, err := net.Listen("tcp", ":0")
		Expect(err).To(BeNil())
		e.Config.Server.HttpPort = fmt.Sprintf("%d", listen.Addr().(*net.TCPAddr).Port)
		e.StopServer()

		err = listen.Close()
		Expect(err).To(BeNil())
	})

})
