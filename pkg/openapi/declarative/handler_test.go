package declarative_test

import (
	"api-gw/pkg/config"
	"api-gw/pkg/openapi/declarative"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
})
