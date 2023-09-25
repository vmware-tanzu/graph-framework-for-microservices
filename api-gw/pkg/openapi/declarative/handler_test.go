package declarative_test

import (
	"api-gw/pkg/config"
	"api-gw/pkg/openapi/declarative"
	"api-gw/pkg/server/echo_server"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	nexus_client "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/nexus-client"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("Handler tests", func() {
	BeforeSuite(func() {
		log.SetLevel(log.DebugLevel)
		err := declarative.Load(spec)
		Expect(err).To(BeNil())
	})

	It("should test ListHandler for gns list url", func() {
		ec := declarative.SetupContext(Uri, http.MethodGet, declarative.Paths[Uri].Get)

		// setup test http server for backend service calls
		var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()
		config.Cfg = &config.Config{BackendService: server.URL}

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ec.Context = c

		err := declarative.ListHandler(ec)
		Expect(err).To(BeNil())
		Expect(rec.Body.String()).To(Equal("[]\n"))
		Expect(requestUri).To(Equal("/v1alpha1/project/default/global-namespaces"))
	})

	It("should test GetHandler for given gns id", func() {
		ec := declarative.SetupContext(ResourceUri, http.MethodGet, declarative.Paths[ResourceUri].Get)

		// setup test http server for backend service calls
		var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`{}`))
		}))
		defer server.Close()
		config.Cfg = &config.Config{BackendService: server.URL}

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("example-gns-id")
		ec.Context = c

		err := declarative.GetHandler(ec)
		Expect(err).To(BeNil())
		Expect(rec.Body.String()).To(Equal("{}\n"))
		Expect(requestUri).To(Equal("/v1alpha1/project/default/global-namespaces/example-gns-id"))
	})

	It("should test PutHandler for given gns id", func() {
		ec := declarative.SetupContext(ResourceUri, http.MethodPut, declarative.Paths[ResourceUri].Put)
		gnsJson := `{
    "metadata": {
        "name": "test"
    },
    "spec": {
        "foo": "bar"
    }
}`

		// setup test http server for backend service calls
		var requestUri string
		var requestBody string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			requestUri = req.URL.String()
			if b, err := io.ReadAll(req.Body); err == nil {
				requestBody = string(b)
			}
			res.WriteHeader(200)
			res.Write([]byte(`{}`))
		}))
		defer server.Close()
		config.Cfg = &config.Config{BackendService: server.URL}

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ec.Context = c

		err := declarative.PutHandler(ec)
		Expect(err).To(BeNil())
		Expect(rec.Body.String()).To(Equal("{}\n"))
		Expect(requestBody).To(Equal("{\"foo\":\"bar\"}"))
		Expect(requestUri).To(Equal("/v1alpha1/project/default/global-namespaces/test"))
	})

	It("should test PutHandler for given gns id with empty spec", func() {
		ec := declarative.SetupContext(ResourceUri, http.MethodPut, declarative.Paths[ResourceUri].Put)
		gnsJson := `{
    "metadata": {
        "name": "test"
    }
}`

		// setup test http server for backend service calls
		var requestUri string
		var requestBody string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			requestUri = req.URL.String()
			if b, err := io.ReadAll(req.Body); err == nil {
				requestBody = string(b)
			}
			res.WriteHeader(200)
			res.Write([]byte(`{}`))
		}))
		defer server.Close()
		config.Cfg = &config.Config{BackendService: server.URL}

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ec.Context = c

		err := declarative.PutHandler(ec)
		Expect(err).To(BeNil())
		Expect(rec.Body.String()).To(Equal("{}\n"))
		Expect(requestBody).To(Equal(""))
		Expect(requestUri).To(Equal("/v1alpha1/project/default/global-namespaces/test"))
	})

	It("should test DeleteHandler for given gns id", func() {
		ec := declarative.SetupContext(ResourceUri, http.MethodDelete, declarative.Paths[ResourceUri].Delete)

		// setup test http server for backend service calls
		var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			requestUri = req.URL.String()
			res.WriteHeader(200)
		}))
		defer server.Close()
		config.Cfg = &config.Config{BackendService: server.URL}

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("example-gns-id")
		ec.Context = c

		err := declarative.DeleteHandler(ec)
		Expect(err).To(BeNil())
		Expect(rec.Code).To(Equal(200))
		Expect(requestUri).To(Equal("/v1alpha1/project/default/global-namespaces/example-gns-id"))
	})

	It("should test buildUrlFromParams method with provided labels", func() {
		config.Cfg.BackendService = ""
		ec := declarative.SetupContext(ResourceUri, http.MethodGet, declarative.Paths[ResourceUri].Get)
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/?labelSelector=projectId=example-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("example-gns-id")
		ec.Context = c
		url, err := declarative.BuildUrlFromParams(ec)
		Expect(err).To(BeNil())
		Expect(url).To(Equal("/v1alpha1/project/example-id/global-namespaces/example-gns-id"))
	})

	It("should test buildUrlFromParams method without labels", func() {
		config.Cfg.BackendService = ""
		ec := declarative.SetupContext(ResourceUri, http.MethodGet, declarative.Paths[ResourceUri].Get)
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("example-gns-id")
		ec.Context = c
		url, err := declarative.BuildUrlFromParams(ec)
		Expect(err).To(BeNil())
		Expect(url).To(Equal("/v1alpha1/project/default/global-namespaces/example-gns-id"))
	})

	It("should test buildUrlFromBody method with provided labels", func() {
		config.Cfg.BackendService = ""
		ec := declarative.SetupContext(ResourceUri, http.MethodPut, declarative.Paths[ResourceUri].Put)
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ec.Context = c
		url, err := declarative.BuildUrlFromBody(ec, map[string]interface{}{
			"name": "test",
			"labels": map[string]interface{}{
				"projectId": "example-id",
			},
		})
		Expect(err).To(BeNil())
		Expect(url).To(Equal("/v1alpha1/project/example-id/global-namespaces/test"))
	})

	It("should test buildUrlFromBody method without labels", func() {
		config.Cfg.BackendService = ""
		ec := declarative.SetupContext(ResourceUri, http.MethodPut, declarative.Paths[ResourceUri].Put)
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ec.Context = c
		url, err := declarative.BuildUrlFromBody(ec, map[string]interface{}{
			"name": "test",
		})
		Expect(err).To(BeNil())
		Expect(url).To(Equal("/v1alpha1/project/default/global-namespaces/test"))
	})

	It("should test Apis handler", func() {
		echoServer := echo_server.NewEchoServer(config.Cfg, &kubernetes.Clientset{}, &nexus_client.Clientset{})
		echoServer.RegisterDeclarativeRouter()

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/declarative/apis", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/declarative/apis")

		err := declarative.ApisHandler(c)
		Expect(err).To(BeNil())

		expectedBody := `{"/apis/gns.vmware.org/v1/globalnamespacelists":{"GET":{"group":"gns.vmware.org","kind":"GlobalNamespaceList","params":null,"uri":"/v1alpha1/global-namespaces/test"},"short":{"name":"gns","uri":"/apis/v1/gns"}},"/apis/gns.vmware.org/v1/globalnamespaces":{"GET":{"group":"gns.vmware.org","kind":"GlobalNamespace","params":["projectId"],"uri":"/v1alpha1/project/{projectId}/global-namespaces"},"PUT":{"group":"gns.vmware.org","kind":"GlobalNamespace","params":["projectId","id"],"uri":"/v1alpha1/project/{projectId}/global-namespaces/{id}"},"short":{"name":"gns","uri":"/apis/v1/gns"},"yaml":"apiVersion: gns.vmware.org/v1\nkind: GlobalNamespace\nmetadata:\n  labels:\n    projectId: string\n  name: string\nspec:\n  api_discovery_enabled: true\n  ca: string\n  ca_type: PreExistingCA\n  color: string\n  description: string\n  display_name: string\n  domain_name: string\n  match_conditions:\n  - cluster:\n      match: string\n      type: string\n    namespace:\n      match: string\n      type: string\n    service: object\n  mtls_enforced: true\n  name: string\n  use_shared_gateway: true\n  version: string\n"},"/apis/gns.vmware.org/v1/globalnamespaces/:name":{"DELETE":{"group":"gns.vmware.org","kind":"GlobalNamespace","params":["projectId","id"],"uri":"/v1alpha1/project/{projectId}/global-namespaces/{id}"},"GET":{"group":"gns.vmware.org","kind":"GlobalNamespace","params":["projectId","id"],"uri":"/v1alpha1/project/{projectId}/global-namespaces/{id}"},"short":{"name":"gns","uri":"/apis/v1/gns/:name"}}}
`
		Expect(rec.Body.String()).To(Equal(expectedBody))
	})

	It("should test Apis handler with globalnamespaces.gns.vmware.org crd", func() {
		echoServer := echo_server.NewEchoServer(config.Cfg, &kubernetes.Clientset{}, &nexus_client.Clientset{})
		echoServer.RegisterDeclarativeRouter()

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/declarative/apis?crd=globalnamespaces.gns.vmware.org", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/declarative/apis?crd=globalnamespaces.gns.vmware.org")

		err := declarative.ApisHandler(c)
		Expect(err).To(BeNil())

		expectedBody := `apiVersion: gns.vmware.org/v1
kind: GlobalNamespace
metadata:
  labels:
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
    service: object
  mtls_enforced: true
  name: string
  use_shared_gateway: true
  version: string
`
		Expect(rec.Body.String()).To(Equal(expectedBody))
	})

	It("should test Apis handler with non-existent crd", func() {
		echoServer := echo_server.NewEchoServer(config.Cfg, &kubernetes.Clientset{}, &nexus_client.Clientset{})
		echoServer.RegisterDeclarativeRouter()

		// setup echo test
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/declarative/apis?crd=non-existent-crd", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/declarative/apis?crd=non-existent-crd")

		err := declarative.ApisHandler(c)
		Expect(err).To(BeNil())
		Expect(rec.Code).To(Equal(http.StatusNotFound))
	})
})
