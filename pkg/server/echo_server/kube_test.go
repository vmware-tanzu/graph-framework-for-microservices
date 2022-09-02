package echo_server_test

import (
	"api-gw/pkg/client"
	"api-gw/pkg/model"
	"api-gw/pkg/server/echo_server"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

var _ = Describe("Kube tests", func() {
	var (
	//fakeClient kubernetes.Interface
	)

	BeforeSuite(func() {
		scheme := runtime.NewScheme()
		//client.Client = fake.NewSimpleDynamicClient(scheme)
		client.Client = fake.NewSimpleDynamicClientWithCustomListKinds(scheme, map[schema.GroupVersionResource]string{
			schema.GroupVersionResource{
				Group:    "gns.vmware.org",
				Version:  "v1",
				Resource: "globalnamespaces",
			}: "GlobalNamespaceList",
			schema.GroupVersionResource{
				Group:    "root.vmware.org",
				Version:  "v1",
				Resource: "roots",
			}: "RootList",
			schema.GroupVersionResource{
				Group:    "orgchart.vmware.org",
				Version:  "v1",
				Resource: "leaders",
			}: "LeaderList",
		})

		//fakeClient = k8sFake.NewSimpleClientset()
		log.SetLevel(log.TraceLevel)
	})

	It("should test kubePost handler with invalid body", func() {
		gnsJson := `invalid`

		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "roots.root.vmware.org",
			GroupName: "root.vmware.org",
			Resource:  "roots",
		}

		err := echo_server.KubePostHandler(nc)
		Expect(err).To(Not(BeNil()))
	})

	It("should test kubePost handler without k8s spec", func() {
		gnsJson := `{}`

		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "roots.root.vmware.org",
			GroupName: "root.vmware.org",
			Resource:  "roots",
		}

		model.CrdTypeToNodeInfo["roots.root.vmware.org"] = model.NodeInfo{
			Name:            "Root.root",
			ParentHierarchy: []string{},
			Children: map[string]model.NodeHelperChild{
				"globalnamespaces.gns.vmware.org": {
					FieldName:    "gns",
					FieldNameGvk: "gnsGvk",
					IsNamed:      false,
				},
			},
		}

		err := echo_server.KubePostHandler(nc)
		Expect(err).To(Not(BeNil()))
	})

	It("should test kubePost handler and add root object without a spec", func() {
		gnsJson := `{
	"apiVersion": "root.vmware.org/v1",
	"kind": "Root",
    "metadata": {
        "name": "root"
    }
}`

		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "roots.root.vmware.org",
			GroupName: "root.vmware.org",
			Resource:  "roots",
		}

		model.CrdTypeToNodeInfo["roots.root.vmware.org"] = model.NodeInfo{
			Name:            "Root.root",
			ParentHierarchy: []string{},
			Children: map[string]model.NodeHelperChild{
				"globalnamespaces.gns.vmware.org": {
					FieldName:    "gns",
					FieldNameGvk: "gnsGvk",
					IsNamed:      false,
				},
			},
		}

		err := echo_server.KubePostHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"apiVersion\":\"root.vmware.org/v1\",\"kind\":\"Root\",\"metadata\":{\"labels\":{\"nexus/display_name\":\"root\",\"nexus/is_name_hashed\":\"true\"},\"name\":\"de3f9fe476b35572145d6b4031712249619efdae\"},\"spec\":{}}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should test kubePost handler and add gns object", func() {
		gnsJson := `{
	"apiVersion": "gns.vmware.org/v1",
	"kind": "GlobalNamespace",
    "metadata": {
        "name": "test",
		"labels": {
			"roots.root.vmware.org": "root"
		}
    },
    "spec": {
        "foo": "bar"
    }
}`

		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "globalnamespaces.gns.vmware.org",
			GroupName: "gns.vmware.org",
			Resource:  "globalnamespaces",
		}

		model.CrdTypeToNodeInfo["globalnamespaces.gns.vmware.org"] = model.NodeInfo{
			Name:            "Gns.gns",
			ParentHierarchy: []string{"roots.root.vmware.org"},
		}

		err := echo_server.KubePostHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"apiVersion\":\"gns.vmware.org/v1\",\"kind\":\"GlobalNamespace\",\"metadata\":{\"labels\":{\"nexus/display_name\":\"test\",\"nexus/is_name_hashed\":\"true\",\"roots.root.vmware.org\":\"root\"},\"name\":\"2587591c2e1023ff9498b1b70ac5cbcb84504352\"},\"spec\":{\"foo\":\"bar\"}}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should test kubeGet handler", func() {
		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "globalnamespaces.gns.vmware.org",
			GroupName: "gns.vmware.org",
			Resource:  "globalnamespaces",
		}

		err := echo_server.KubeGetHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"apiVersion\":\"gns.vmware.org/v1\",\"items\":[{\"apiVersion\":\"gns.vmware.org/v1\",\"kind\":\"GlobalNamespace\",\"metadata\":{\"labels\":{\"nexus/display_name\":\"test\",\"nexus/is_name_hashed\":\"true\",\"roots.root.vmware.org\":\"root\"},\"name\":\"2587591c2e1023ff9498b1b70ac5cbcb84504352\"},\"spec\":{\"foo\":\"bar\"}}],\"kind\":\"GlobalNamespaceList\",\"metadata\":{\"resourceVersion\":\"\"}}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should test kubePost handler and add update object", func() {
		gnsJson := `{
	"apiVersion": "gns.vmware.org/v1",
	"kind": "GlobalNamespace",
    "metadata": {
        "name": "test",
		"labels": {
			"roots.root.vmware.org": "root"
		}
    },
    "spec": {
        "foo": "bar2"
    }
}`

		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "globalnamespaces.gns.vmware.org",
			GroupName: "gns.vmware.org",
			Resource:  "globalnamespaces",
		}

		model.CrdTypeToNodeInfo["globalnamespaces.gns.vmware.org"] = model.NodeInfo{
			Name:            "Gns.gns",
			ParentHierarchy: []string{"roots.root.vmware.org"},
		}

		err := echo_server.KubePostHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"apiVersion\":\"gns.vmware.org/v1\",\"kind\":\"GlobalNamespace\",\"metadata\":{\"labels\":{\"nexus/display_name\":\"test\",\"nexus/is_name_hashed\":\"true\",\"roots.root.vmware.org\":\"root\"},\"name\":\"2587591c2e1023ff9498b1b70ac5cbcb84504352\"},\"spec\":{\"foo\":\"bar2\"}}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should test kubeGetByName handler", func() {
		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("2587591c2e1023ff9498b1b70ac5cbcb84504352")

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "globalnamespaces.gns.vmware.org",
			GroupName: "gns.vmware.org",
			Resource:  "globalnamespaces",
		}

		err := echo_server.KubeGetByNameHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"apiVersion\":\"gns.vmware.org/v1\",\"kind\":\"GlobalNamespace\",\"metadata\":{\"labels\":{\"nexus/display_name\":\"test\",\"nexus/is_name_hashed\":\"true\",\"roots.root.vmware.org\":\"root\"},\"name\":\"2587591c2e1023ff9498b1b70ac5cbcb84504352\"},\"spec\":{\"foo\":\"bar2\"}}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should test kubeGetByName handler with non-existent name and get an error", func() {
		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("non-existent-id")

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "globalnamespaces.gns.vmware.org",
			GroupName: "gns.vmware.org",
			Resource:  "globalnamespaces",
		}

		err := echo_server.KubeGetByNameHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"metadata\":{},\"status\":\"Failure\",\"message\":\"globalnamespaces.gns.vmware.org \\\"non-existent-id\\\" not found\",\"reason\":\"NotFound\",\"details\":{\"name\":\"non-existent-id\",\"group\":\"gns.vmware.org\",\"kind\":\"globalnamespaces\"},\"code\":404}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should test kubeDelete handler", func() {
		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("2587591c2e1023ff9498b1b70ac5cbcb84504352")

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "globalnamespaces.gns.vmware.org",
			GroupName: "gns.vmware.org",
			Resource:  "globalnamespaces",
		}

		err := echo_server.KubeDeleteHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"apiVersion\":\"v1\",\"details\":{\"group\":\"gns.vmware.org\",\"kind\":\"globalnamespaces\",\"name\":\"2587591c2e1023ff9498b1b70ac5cbcb84504352\"},\"kind\":\"Status\",\"metadata\":{},\"status\":\"Success\"}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should test kubeDelete handler with non-existent name and get an error", func() {
		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:name")
		c.SetParamNames("name")
		c.SetParamValues("non-existent-id")

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "globalnamespaces.gns.vmware.org",
			GroupName: "gns.vmware.org",
			Resource:  "globalnamespaces",
		}

		err := echo_server.KubeDeleteHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"metadata\":{},\"status\":\"Failure\",\"message\":\"globalnamespaces.gns.vmware.org \\\"non-existent-id\\\" not found\",\"reason\":\"NotFound\",\"details\":{\"name\":\"non-existent-id\",\"group\":\"gns.vmware.org\",\"kind\":\"globalnamespaces\"},\"code\":404}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})
})
