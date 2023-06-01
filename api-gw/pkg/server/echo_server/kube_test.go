package echo_server_test

import (
	"api-gw/pkg/client"
	"api-gw/pkg/config"
	"api-gw/pkg/model"
	"api-gw/pkg/server/echo_server"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
			{
				Group:    "gns.vmware.org",
				Version:  "v1",
				Resource: "globalnamespaces",
			}: "GlobalNamespaceList",
			{
				Group:    "root.vmware.org",
				Version:  "v1",
				Resource: "roots",
			}: "RootList",
			{
				Group:    "orgchart.vmware.org",
				Version:  "v1",
				Resource: "leaders",
			}: "LeaderList",
		})

		//fakeClient = k8sFake.NewSimpleClientset()
		log.SetLevel(log.TraceLevel)
	})

	It("should fail the request with invalid body using kubePost handler", func() {
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

	It("should fail creating object without a spec using kubePost handler", func() {
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

	It("should create root object without a spec using kubePost handler", func() {
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

	It("should create gns object using kubePost handler", func() {
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

	It("should get gns object using kubeGet handler", func() {
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
		c.QueryParams().Add("limit", "1")
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

	It("should update gns object using kubePost handler", func() {
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

	It("should not remove child/link Gvks while update by kubePost handler", func() {
		gvr := schema.GroupVersionResource{
			Group:    "orgchart.vmware.org",
			Version:  "v1",
			Resource: "foos",
		}
		obj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "orgchart.vmware.org/v1",
				"kind":       "Foo",
				"metadata": map[string]interface{}{
					"name": "bb2cbdf1b03e754cea2c9da8e9134c050bc0d547",
				},
				"spec": map[string]interface{}{
					"childGvk": "value_one",
					"linkGvk":  "value_two",
					"name":     "bob",
				},
			},
		}
		_, err := client.Client.Resource(gvr).Create(context.TODO(), obj, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()
		e := echo.New()

		// If the newspec contains new fields, updateResource should add them while retaining the Gvk fields.
		gnsJson := `{
			"apiVersion": "orgchart.vmware.org/v1",
			"kind": "Foo",
			"metadata": {
				"name": "test",
				"labels": {
					"roots.root.vmware.org": "root"
				}
			},
			"spec": {
				"name": "bob",
				"foo": "bar2"
			}
		}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		nc := &echo_server.NexusContext{
			Context:   c,
			CrdType:   "foos.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "foos",
		}
		model.CrdTypeToNodeInfo["foos.orgchart.vmware.org"] = model.NodeInfo{
			Name:            "Foo.foo",
			ParentHierarchy: []string{"roots.root.vmware.org"},
			Children: map[string]model.NodeHelperChild{
				"childGVK": {
					FieldNameGvk: "childGvk",
				},
				"linkGVK": {
					FieldNameGvk: "linkGvk",
				},
			},
		}

		err = echo_server.KubePostHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse := "{\"apiVersion\":\"orgchart.vmware.org/v1\",\"kind\":\"Foo\",\"metadata\":{\"labels\":{\"nexus/display_name\":\"test\",\"nexus/is_name_hashed\":\"true\",\"roots.root.vmware.org\":\"root\"},\"name\":\"bb2cbdf1b03e754cea2c9da8e9134c050bc0d547\"},\"spec\":{\"childGvk\":\"value_one\",\"foo\":\"bar2\",\"linkGvk\":\"value_two\",\"name\":\"bob\"}}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))

		// If the newspec has empty spec, updateResource should remove the existing spec fields while retaining the Gvk fields.
		gnsJson = `{
			"apiVersion": "orgchart.vmware.org/v1",
			"kind": "Foo",
			"metadata": {
				"name": "test",
				"labels": {
					"roots.root.vmware.org": "root"
				}
			},
			"spec": {}
		}`
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(gnsJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		nc = &echo_server.NexusContext{
			Context:   c,
			CrdType:   "foos.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "foos",
		}

		err = echo_server.KubePostHandler(nc)
		Expect(err).To(BeNil())
		expectedResponse = "{\"apiVersion\":\"orgchart.vmware.org/v1\",\"kind\":\"Foo\",\"metadata\":{\"labels\":{\"nexus/display_name\":\"test\",\"nexus/is_name_hashed\":\"true\",\"roots.root.vmware.org\":\"root\"},\"name\":\"bb2cbdf1b03e754cea2c9da8e9134c050bc0d547\"},\"spec\":{\"childGvk\":\"value_one\",\"linkGvk\":\"value_two\"}}\n"
		Expect(rec.Body.String()).To(Equal(expectedResponse))
	})

	It("should get gns object using kubeGetByName handler", func() {
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

	It("should fail kubeGetByName handler by using object with non-existent name", func() {
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

	It("should delete gns object using kubeDelete handler with kubectl user agent", func() {
		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set("User-Agent", "kubectl")
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

	Context("test delete method with label selector", func() {
		It("should create gns object using kubePost handler", func() {
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

		It("should delete gns object using kubeDelete handler", func() {
			//var requestUri string
			server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				//requestUri = req.URL.String()
				res.WriteHeader(200)
				res.Write([]byte(`[]`))
			}))
			defer server.Close()

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/test?labelSelector=roots.root.vmware.org=root", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:name?labelSelector=roots.root.vmware.org=root")
			c.SetParamNames("name")
			c.SetParamValues("test")

			nc := &echo_server.NexusContext{
				Context:   c,
				CrdType:   "globalnamespaces.gns.vmware.org",
				GroupName: "gns.vmware.org",
				Resource:  "globalnamespaces",
			}

			err := echo_server.KubeDeleteHandler(nc)
			Expect(err).To(BeNil())
			expectedResponse := "{\"apiVersion\":\"v1\",\"details\":{\"group\":\"gns.vmware.org\",\"kind\":\"globalnamespaces\",\"name\":\"test\"},\"kind\":\"Status\",\"metadata\":{},\"status\":\"Success\"}\n"
			Expect(rec.Body.String()).To(Equal(expectedResponse))
		})

		It("should delete root object using kubeDeleteHandler", func() {
			//var requestUri string
			server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				//requestUri = req.URL.String()
				res.WriteHeader(200)
				res.Write([]byte(`[]`))
			}))
			defer server.Close()

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/root", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:name")
			c.SetParamNames("name")
			c.SetParamValues("root")

			nc := &echo_server.NexusContext{
				Context:   c,
				CrdType:   "roots.root.vmware.org",
				GroupName: "root.vmware.org",
				Resource:  "roots",
			}

			err := echo_server.KubeDeleteHandler(nc)
			Expect(err).To(BeNil())
			expectedResponse := "{\"apiVersion\":\"v1\",\"details\":{\"group\":\"root.vmware.org\",\"kind\":\"roots\",\"name\":\"root\"},\"kind\":\"Status\",\"metadata\":{},\"status\":\"Success\"}\n"
			log.Debugf(rec.Body.String())
			Expect(rec.Body.String()).To(Equal(expectedResponse))
		})
	})

	It("should fail kubeDelete handler with kubectl user agent by using object with non-existent name", func() {
		//var requestUri string
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//requestUri = req.URL.String()
			res.WriteHeader(200)
			res.Write([]byte(`[]`))
		}))
		defer server.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set("User-Agent", "kubectl")
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

	It("should test updateProxyResponse method when custom not found page cfg is not set", func() {
		if config.Cfg == nil {
			config.Cfg = &config.Config{}
		}
		config.Cfg.CustomNotFoundPage = ""
		err := echo_server.UpdateProxyResponse(&http.Response{})
		Expect(err).To(BeNil())
	})

	It("should test updateProxyResponse method when custom not found page cfg is set", func() {
		if config.Cfg == nil {
			config.Cfg = &config.Config{}
		}

		body := []byte(`404 page`)
		server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write(body)
		}))
		defer server.Close()

		config.Cfg.CustomNotFoundPage = server.URL
		response := &http.Response{
			StatusCode: http.StatusNotFound,
		}
		err := echo_server.UpdateProxyResponse(response)
		Expect(err).To(BeNil())

		bodyRes, err := io.ReadAll(response.Body)
		Expect(err).To(BeNil())

		Expect(bodyRes).To(Equal(body))

	})
})
