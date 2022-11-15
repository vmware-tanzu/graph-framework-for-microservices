package echo_server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"api-gw/controllers"
	"api-gw/pkg/client"
	"api-gw/pkg/config"
	"api-gw/pkg/model"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var _ = Describe("Echo server tests", func() {
	var e *EchoServer

	BeforeEach(func() {
		config.Cfg = &config.Config{
			Server:             config.ServerConfig{},
			EnableNexusRuntime: true,
			BackendService:     "",
		}
		e = NewEchoServer(config.Cfg)
	})

	It("should handle list query for empty list", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodGet, "/:orgchart.Leader", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		nc := &NexusContext{
			NexusURI: "/leaders",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := listHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
		Expect(rec.Body.String()).Should(Equal("[]\n"))
	})

	It("shouldn't handle get query for singleton object if nexus object name is empty string", func() {
		restUri := nexus.RestURIs{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodGet, "/root/:orgchart.Root/leader/:management.Leader", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetParamNames("management.Leader")
		c.SetParamValues("")
		nc := &NexusContext{
			NexusURI: "/root/{orgchart.Root}/leader/{management.Leader}",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "management.vmware.org",
			Resource:  "leaders",
		}

		err := getHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).ToNot(Equal(200))
	})

	It("should handle put query for singleton object without passing name parameters in request", func() {
		leaderJson := `{
			"designation": "abc",
			"employeeID": 100,
			"name": "xyz"
		}`

		restUri := nexus.RestURIs{
			Uri:     "/leader",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodPost, "/:orgchart.Leader", strings.NewReader(leaderJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		nc := &NexusContext{
			NexusURI: "/leader",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	It("should handle put query for singleton object with default as name", func() {
		leaderJson := `{
			"designation": "abc",
			"employeeID": 100,
			"name": "xyz"
		}`

		restUri := nexus.RestURIs{
			Uri:     "/leader",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodPost, "/:orgchart.Leader", strings.NewReader(leaderJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetParamNames("orgchart.Leader")
		c.SetParamValues("default")
		nc := &NexusContext{
			NexusURI: "/leader",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	It("should not remove child/link GVKs while update by putHandler", func() {
		req := httptest.NewRequest(http.MethodPost, "/:orgchart.Leader", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		nc := &NexusContext{
			Context: e.Echo.NewContext(req, rec),
		}
		obj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "orgchart.vmware.org/v1",
				"kind":       "Foo",
				"metadata": map[string]interface{}{
					"name": "zoo",
				},
				"spec": map[string]interface{}{
					"childGvk": "value_one",
					"linkGvk":  "value_two",
				},
			},
		}
		gvr := schema.GroupVersionResource{
			Group:    "orgchart.vmware.org",
			Version:  "v1",
			Resource: "foos",
		}
		_, err := client.Client.Resource(gvr).Create(context.TODO(), obj, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		crdInfo := model.NodeInfo{
			Children: map[string]model.NodeHelperChild{
				"childGVK": {
					FieldNameGvk: "childGvk",
				},
				"linkGVK": {
					FieldNameGvk: "linkGvk",
				},
			},
		}

		// If the newspec contains new fields, updateResource should add them while retaining the Gvk fields.
		body := map[string]interface{}{
			"employeeID": "100",
			"name":       "bob",
		}
		err = updateResource(nc, gvr, obj, body, crdInfo)
		Expect(err).NotTo(HaveOccurred())

		obj, err = client.Client.Resource(gvr).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		spec := obj.Object["spec"].(map[string]interface{})
		Expect(spec).To(HaveKey("childGvk"))
		Expect(spec).To(HaveKey("linkGvk"))
		Expect(spec).To(HaveLen(4))

		// If the newspec has empty spec, updateResource should remove the existing spec fields while retaining the Gvk fields.
		body = map[string]interface{}{}
		err = updateResource(nc, gvr, obj, body, crdInfo)
		Expect(err).NotTo(HaveOccurred())

		obj, err = client.Client.Resource(gvr).Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		spec = obj.Object["spec"].(map[string]interface{})
		Expect(spec).To(HaveKey("childGvk"))
		Expect(spec).To(HaveKey("linkGvk"))
		Expect(spec).To(HaveLen(2))
	})

	It("should fail update in put query if query param update_if_exists=false", func() {
		leaderJson := `{
			"designation": "abc",
			"employeeID": 100,
			"name": "xyz"
		}`

		restUri := nexus.RestURIs{
			Uri:     "/leader",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodPost, "/:orgchart.Leader", strings.NewReader(leaderJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetParamNames("orgchart.Leader")
		c.SetParamValues("default")
		c.QueryParams().Add("update_if_exists", "false")
		nc := &NexusContext{
			NexusURI: "/leader",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(403))
	})

	It("shouldn't handle put query for singleton object with not default as name", func() {
		leaderJson := `{
			"designation": "abc",
			"employeeID": 100,
			"name": "xyz"
		}`

		restUri := nexus.RestURIs{
			Uri:     "/leader/{orgchart.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
			[]string{}, nil, nil, true, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(leaderJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetPath("/:orgchart.Leader")
		c.SetParamNames("orgchart.Leader")
		c.SetParamValues("notdefault")
		nc := &NexusContext{
			NexusURI:  "/leader/{orgchart.Leader}",
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(400))

		// should not handle put query for singleton object with not default as name
		patchJson := `{
			"designation": "CEO",
			"new-field": "xyz"
		}`

		req = httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(patchJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = e.Echo.NewContext(req, rec)
		c.SetPath("/:orgchart.Leader")
		c.SetParamNames("orgchart.Leader")
		c.SetParamValues("notdefault")
		nc = &NexusContext{
			NexusURI:  "/leader/{orgchart.Leader}",
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
			Codes: map[nexus.ResponseCode]nexus.HTTPResponse{
				http.StatusBadRequest: {Description: http.StatusText(http.StatusBadRequest)},
			},
		}

		err = patchHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(400))
	})

	It("shouldn't handle put query for object without a name", func() {
		leaderJson := `{
			"designation": "abc",
			"employeeID": 100,
			"name": "xyz"
		}`

		restUri := nexus.RestURIs{
			Uri:     "/leader/{orgchart.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
			[]string{}, nil, nil, false, "")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(leaderJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetPath("/:orgchart.Leader")
		nc := &NexusContext{
			NexusURI:  "/leader/{orgchart.Leader}",
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(400))
	})

	It("should handle list query", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodGet, "/:orgchart.Leader", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		nc := &NexusContext{
			NexusURI: "/leaders",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := listHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
		Expect(rec.Body.String()).Should(Equal("[{\"name\":\"default\",\"spec\":{\"designation\":\"abc\",\"employeeID\":100,\"name\":\"xyz\"},\"status\":{}}]\n"))
	})

	It("should handle get query", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leader/{management.Leader}",
			Methods: nexus.HTTPListResponse,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodGet, "/leader", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetPath("/leader/:management.Leader")
		c.SetParamNames("management.Leader")
		c.SetParamValues("default")
		nc := &NexusContext{
			NexusURI: "/leader/{management.Leader}",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := getHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
		Expect(rec.Body.String()).Should(Equal("{\"spec\":{\"designation\":\"abc\",\"employeeID\":100,\"name\":\"xyz\"},\"status\":{}}\n"))
	})

	It("should handle delete query for singleton object", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leader/{management.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodDelete, "/leader", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetPath("/leader")
		nc := &NexusContext{
			NexusURI:  "/leader/{management.Leader}",
			Codes:     nexus.DefaultHTTPDELETEResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := deleteHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	It("should fail delete query for singleton object with provided name", func() {
		restUri := nexus.RestURIs{
			Uri:     "/leader/{management.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodDelete, "/leader", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetPath("/leader/:management.Leader")
		c.SetParamNames("management.Leader")
		c.SetParamValues("some")
		nc := &NexusContext{
			NexusURI:  "/leader/{management.Leader}",
			Codes:     nexus.DefaultHTTPDELETEResponses,
			Context:   c,
			CrdType:   "leaders.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "leaders",
		}

		err := deleteHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(400))
	})

	It("should handle put query for non-singleton object", func() {
		mgrJson := `{
			"designation": "abc",
			"employeeID": 100,
			"name": "xyz"
		}`

		restUri := nexus.RestURIs{
			Uri:     "/mgr/{management.Mgr}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "mgrs.orgchart.vmware.org", "management.Mgr",
			[]string{}, nil, nil, false, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "mgrs.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodPost, "/mgr/:management.Mgr", strings.NewReader(mgrJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetParamNames("management.Mgr")
		c.SetParamValues("mgr1")
		nc := &NexusContext{
			NexusURI: "/mgr/{management.Mgr}",
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c,
			CrdType:   "mgrs.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "mgrs",
		}

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	It("should handle delete query for non-singleton object with provided name", func() {
		restUri := nexus.RestURIs{
			Uri:     "/mgr/{management.Mgr}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "mgrs.orgchart.vmware.org", "management.Mgr",
			[]string{}, nil, nil, false, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "mgrs.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodDelete, "/leader", nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)
		c.SetPath("/mgr/:management.Mgr")
		c.SetParamNames("management.Mgr")
		c.SetParamValues("mgr1")
		nc := &NexusContext{
			NexusURI:  "/mgr/{management.Mgr}",
			Codes:     nexus.DefaultHTTPDELETEResponses,
			Context:   c,
			CrdType:   "mgrs.orgchart.vmware.org",
			GroupName: "orgchart.vmware.org",
			Resource:  "mgrs",
		}

		err := deleteHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	Context("should GET child from Parent object", func() {
		BeforeEach(func() {
			// Create `Leader` object
			leaderChildrenJson := getLeaderChildrenJson("", "")
			restUri := nexus.RestURIs{
				Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}

			nc, rec := initNode(e, "leaders.management.vmware.org", "management.vmware.org",
				"leaders", "management.Leader", http.MethodPost, leaderChildrenJson,
				"/root/:orgchart.Root/leader/:management.Leader", restUri)

			err := putHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))
		})

		It("should process single child from parent", func() {
			// create HR child object
			hrJson := createTestNode("hr.vmware.org/v1", "HumanResources", "default")
			restUri := nexus.RestURIs{
				Uri:     "/root/{orgchart.Root}/hr/{hr.HumanResources}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}
			hrCtx, hrRec := initNode(e, "humanresourceses.hr.vmware.org", "hr.vmware.org",
				"humanresourceses", "hr.HumanResources", http.MethodPost, hrJson,
				"/root/:orgchart.Root/hr/:hr.HumanResources", restUri)

			err := putHandler(hrCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(hrRec.Code).To(Equal(200))

			// construct annotation
			n := constructTestChildrenAnnotation()

			urisMap := make(map[string]model.RestUriInfo)
			// add child, link and status URIs for each GET method
			var newUris []nexus.RestURIs
			controllers.ConstructNewURIs(n, urisMap, &newUris)

			// should contain all the URIs constructed with child/link only with GET rest API.
			Expect(newUris).Should(ConsistOf(expectedHRRestURIs()))
			n.NexusRestAPIGen.Uris = append(n.NexusRestAPIGen.Uris, newUris...)

			for _, restUri := range n.NexusRestAPIGen.Uris {
				e.RegisterRouter(restUri)
			}

			model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
				n.Hierarchy, n.Children, n.Links, true, "some description")
			model.ConstructMapURIToCRDType(model.Upsert, "leaders.management.vmware.org", n.NexusRestAPIGen.Uris)
			model.ConstructMapUriToUriInfo(model.Upsert, urisMap)

			nc, rec := createSampleHRRequest(e)
			// should GET `HR` child object successfully from Parent object `Leader`
			err = getHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))

			Expect(rec.Body.String()).Should(Equal("{\"group\":\"hr.vmware.org/v1\",\"kind\":\"HumanResources\",\"name\":\"default\"}\n"))

			By("deleting the HR object, should fail to find the object on GET")
			err = deleteHandler(hrCtx)
			Expect(err).NotTo(HaveOccurred())

			nc, rec = createSampleHRRequest(e)
			// should fail when child object not exists in DB
			err = getHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(404))
			Expect(rec.Body.String()).To(Equal("{\"message\":\"Couldn't find object\"}\n"))
		})

		It("should GET multiple child from parent", func() {
			// create engineering manager `default` child object
			mgrJson1 := createTestNode("management.vmware.org/v1", "Mgr", "default")
			restUri := nexus.RestURIs{
				Uri:     "/root/{orgchart.Root}/mgr/{management.Mgr}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}
			engCtx, engRec := initNode(e, "mgrs.management.vmware.org", "management.vmware.org",
				"mgrs", "management.Mgr", http.MethodPost, mgrJson1,
				"/root/:orgchart.Root/mgr/:management.Mgr", restUri)

			err := putHandler(engCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(engRec.Code).To(Equal(200))

			// create engineering manager `foo` child object
			mgrJson2 := createTestNode("management.vmware.org/v1", "Mgr", "foo")
			restUri = nexus.RestURIs{
				Uri:     "/root/{orgchart.Root}/mgr/{management.Mgr}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}
			engCtx, engRec = initNode(e, "mgrs.management.vmware.org", "management.vmware.org",
				"mgrs", "management.Mgr", http.MethodPost, mgrJson2,
				"/root/:orgchart.Root/mgr/:management.Mgr", restUri)

			err = putHandler(engCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(engRec.Code).To(Equal(200))

			// construct annotation
			n := constructTestLinkAnnotation()

			urisMap := make(map[string]model.RestUriInfo)
			// add child, link and status URIs for each GET method
			var newUris []nexus.RestURIs
			controllers.ConstructNewURIs(n, urisMap, &newUris)

			// should contain all the URIs constructed with child/link only with GET rest API.
			Expect(newUris).Should(ConsistOf(expectedEngManagersRestURIs()))

			n.NexusRestAPIGen.Uris = append(n.NexusRestAPIGen.Uris, newUris...)
			for _, restUri := range n.NexusRestAPIGen.Uris {
				e.RegisterRouter(restUri)
			}

			model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
				n.Hierarchy, n.Children, n.Links, true, "some description")
			model.ConstructMapURIToCRDType(model.Upsert, "leaders.management.vmware.org", n.NexusRestAPIGen.Uris)
			model.ConstructMapUriToUriInfo(model.Upsert, urisMap)

			nc, rec := createSampleEngManagerRequest(e)
			// should GET `Mgr` children object successfully from Parent object `Leader`
			err = getHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))

			Expect(rec.Body.String()).Should(Equal("[{\"group\":\"management.vmware.org/v1\"," +
				"\"kind\":\"Mgr\",\"name\":\"default\"},{\"group\":\"management.vmware.org/v1\",\"kind\":\"Mgr\",\"name\":\"foo\"}]\n"))

			By("deleting the EngManagers object, should fail to find the object on GET")
			err = deleteHandler(engCtx)
			Expect(err).NotTo(HaveOccurred())

			nc, rec = createSampleEngManagerRequest(e)
			// should throw empty list when not exists in DB
			err = getHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Body.String()).To(Equal("[{},{}]\n"))
		})

		It("should fail to determine object when invalid gvk", func() {
			// construct annotation
			n := constructTestLinkAnnotation()

			urisMap := make(map[string]model.RestUriInfo)
			newUris := []nexus.RestURIs{
				{
					Uri: "/root/{orgchart.Root}/leader/{management.Leader}/Role",
					Methods: map[nexus.HTTPMethod]nexus.HTTPCodesResponse{
						http.MethodGet: nexus.DefaultHTTPGETResponses,
					},
				}}
			controllers.ConstructNewURIs(n, urisMap, &newUris)
			urisMap["/root/{orgchart.Root}/leader/{management.Leader}/Role"] = model.RestUriInfo{
				TypeOfURI: model.SingleLinkURI,
			}
			n.NexusRestAPIGen.Uris = append(n.NexusRestAPIGen.Uris, newUris...)

			for _, restUri := range n.NexusRestAPIGen.Uris {
				e.RegisterRouter(restUri)
			}

			model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
				n.Hierarchy, n.Children, n.Links, true, "some description")
			model.ConstructMapURIToCRDType(model.Upsert, "leaders.management.vmware.org", n.NexusRestAPIGen.Uris)
			model.ConstructMapUriToUriInfo(model.Upsert, urisMap)

			// try to get object with empty gvk
			nc, rec := createSampleRoleRequest(e)
			// should not GET `Role` link
			err := getHandler(nc)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(500))
			Expect(rec.Body.String()).To(Equal("{\"message\":\"Couldn't determine gvk of link\"}\n"))

			// Update `Leader` spec with roleGvk
			leaderCtx, leaderRec := initNode(e, "leaders.management.vmware.org", "management.vmware.org",
				"leaders", "management.Leader", http.MethodPost, getLeaderChildrenJson("roleGvk", "1"),
				"/root/:orgchart.Root/leader/:management.Leader", nexus.RestURIs{})

			model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
				n.Hierarchy, n.Children, map[string]model.NodeHelperChild{
					"role.executive.vmware.org": {
						FieldName:    "Role",
						FieldNameGvk: "roleGvk",
						IsNamed:      false,
					},
				}, true, "some description")

			err = putHandler(leaderCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(leaderRec.Code).To(Equal(200))

			//should fail to unmarshal the object when parent has invalid spec.
			nc, rec = createSampleRoleRequest(e)
			// should not GET `Role` link
			err = getHandler(nc)
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(500))
			Expect(rec.Body.String()).To(Equal("{\"message\":\"Couldn't unmarshal gvk of link\"}\n"))
		})

		It("shouldn't show Child and Links GVK when doing object GET", func() {
			// create HR child object
			hrJson := createTestNode("hr.vmware.org/v1", "HumanResources", "default")
			restUri := nexus.RestURIs{
				Uri:     "/root/{orgchart.Root}/hr/{hr.HumanResources}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}
			hrCtx, hrRec := initNode(e, "humanresourceses.hr.vmware.org", "hr.vmware.org",
				"humanresourceses", "hr.HumanResources", http.MethodPost, hrJson,
				"/root/:orgchart.Root/hr/:hr.HumanResources", restUri)

			err := putHandler(hrCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(hrRec.Code).To(Equal(200))

			// construct annotation
			n := constructTestAnnotation()

			urisMap := make(map[string]model.RestUriInfo)
			// add child, link and status URIs for each GET method
			var newUris []nexus.RestURIs
			controllers.ConstructNewURIs(n, urisMap, &newUris)

			n.NexusRestAPIGen.Uris = append(n.NexusRestAPIGen.Uris, newUris...)

			for _, restUri := range n.NexusRestAPIGen.Uris {
				e.RegisterRouter(restUri)
			}

			model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
				n.Hierarchy, n.Children, n.Links, true, "some description")
			model.ConstructMapURIToCRDType(model.Upsert, "leaders.management.vmware.org", n.NexusRestAPIGen.Uris)
			model.ConstructMapUriToUriInfo(model.Upsert, urisMap)

			nc, rec := createSampleLeaderRequest(e)
			err = getHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))
			Expect(rec.Body.String()).Should(Equal(
				"{\"spec\":{\"designation\":\"NexusLead\",\"employeeID\":1},\"status\":{}}\n"))
		})
	})

	It("should handle PUT and GET status subresource", func() {
		// Create `Leader` object
		leaderJson := `{
				"designation": "abc",
				"employeeID": 100,
				"name": "xyz"
			  }`
		restUri := nexus.RestURIs{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}

		nc, rec := initNode(e, "leaders.management.vmware.org", "management.vmware.org",
			"leaders", "management.Leader", http.MethodPost, leaderJson,
			"/root/:orgchart.Root/leader/:management.Leader", restUri)

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))

		// =========== Status subresource
		leaderStatusJson := `{
				"status": {
				  "DaysLeftToEndOfVacations": 139,
				  "IsOnVacations": true
				}
			  }`
		statusUriPath := "/root/{orgchart.Root}/leader/{management.Leader}/status"
		targetUri := "/root/:orgchart.Root/leader/:management.Leader/status"
		model.UriToUriInfo[statusUriPath] = model.RestUriInfo{TypeOfURI: model.StatusURI}
		restUriForStatus := nexus.RestURIs{
			Uri: statusUriPath,
			// Methods: nexus.DefaultHTTPMethodsResponses,
			Methods: nexus.HTTPMethodsResponses{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
				http.MethodPut: nexus.DefaultHTTPPUTResponses,
			},
		}
		urisMap := map[string]model.RestUriInfo{
			statusUriPath: {
				TypeOfURI: model.StatusURI,
			},
		}
		model.ConstructMapUriToUriInfo(model.Upsert, urisMap)

		// =========== status PUT
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.management.vmware.org", []nexus.RestURIs{restUriForStatus})

		req1 := httptest.NewRequest(http.MethodPost, targetUri, strings.NewReader(leaderStatusJson))
		req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec1 := httptest.NewRecorder()
		c1 := e.Echo.NewContext(req1, rec1)
		nc1 := &NexusContext{
			NexusURI:  statusUriPath,
			Context:   c1,
			CrdType:   "leaders.management.vmware.org",
			GroupName: "management.vmware.org",
			Resource:  "leaders",
		}
		err = putHandler(nc1)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec1.Code).To(Equal(200))

		// ============ status GET
		req2 := httptest.NewRequest(http.MethodPost, targetUri, nil)
		req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec2 := httptest.NewRecorder()
		c2 := e.Echo.NewContext(req2, rec2)
		nc2 := &NexusContext{
			NexusURI: statusUriPath,
			//Codes: nexus.DefaultHTTPMethodsResponses,
			Context:   c2,
			CrdType:   "leaders.management.vmware.org",
			GroupName: "management.vmware.org",
			Resource:  "leaders",
		}

		err = getHandler(nc2)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec2.Code).To(Equal(200))
		Expect(rec2.Body.String()).Should(Equal("{\"status\":{\"DaysLeftToEndOfVacations\":139,\"IsOnVacations\":true}}\n"))

		// ============ GET Manager with status subresource
		nc3, rec3 := initNode(e, "leaders.management.vmware.org", "management.vmware.org",
			"leaders", "management.Leader", http.MethodGet, "",
			"/root/:orgchart.Root/leader/:management.Leader", restUri)
		n := constructTestAnnotation()
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
			[]string{}, n.Children, n.Links, true, "some description")
		err = getHandler(nc3)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec3.Code).To(Equal(200))
		Expect(rec3.Body.String()).Should(Equal("{\"spec\":{\"designation\":\"abc\",\"employeeID\":100,\"name\":\"xyz\"},\"status\":{\"status\":{\"DaysLeftToEndOfVacations\":139,\"IsOnVacations\":true}}}\n"))
	})

	It("shouldn't handle PUT status subresource if nexus native status subresource is presnet in status subresource payload", func() {
		// Create `Leader` object
		leaderJson := `{
					"designation": "abc",
					"employeeID": 100,
					"name": "xyz"
				  }`
		restUri := nexus.RestURIs{
			Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}

		nc, rec := initNode(e, "leaders.management.vmware.org", "management.vmware.org",
			"leaders", "management.Leader", http.MethodPost, leaderJson,
			"/root/:orgchart.Root/leader/:management.Leader", restUri)

		err := putHandler(nc)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))

		// =========== Status subresource
		leaderStatusJson := `{
					"nexus": {
						"sourceGeneration": 101,
						"remoteGeneration": 100
					},
					"status": {
					  "DaysLeftToEndOfVacations": 139,
					  "IsOnVacations": true
					}
				  }`
		statusUriPath := "/root/{orgchart.Root}/leader/{management.Leader}/status"
		targetUri := "/root/:orgchart.Root/leader/:management.Leader/status"
		model.UriToUriInfo[statusUriPath] = model.RestUriInfo{TypeOfURI: model.StatusURI}
		restUriForStatus := nexus.RestURIs{
			Uri: statusUriPath,
			// Methods: nexus.DefaultHTTPMethodsResponses,
			Methods: nexus.HTTPMethodsResponses{
				http.MethodPut: nexus.DefaultHTTPPUTResponses,
			},
		}
		urisMap := map[string]model.RestUriInfo{
			statusUriPath: {
				TypeOfURI: model.StatusURI,
			},
		}
		model.ConstructMapUriToUriInfo(model.Upsert, urisMap)

		// =========== status PUT
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.management.vmware.org", "management.Leader",
			[]string{}, nil, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.management.vmware.org", []nexus.RestURIs{restUriForStatus})

		req1 := httptest.NewRequest(http.MethodPost, targetUri, strings.NewReader(leaderStatusJson))
		req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec1 := httptest.NewRecorder()
		c1 := e.Echo.NewContext(req1, rec1)
		nc1 := &NexusContext{
			NexusURI:  statusUriPath,
			Context:   c1,
			CrdType:   "leaders.management.vmware.org",
			GroupName: "management.vmware.org",
			Resource:  "leaders",
		}
		err = putHandler(nc1)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec1.Code).ToNot(Equal(200))

	})

	Context("should handle PATCH API", func() {
		BeforeEach(func() {
			// Create `Leader` object with below fields
			leaderJson := `{
          "designation": "NexusLead",
          "employeeID": 1,
          "description": "Hello World!"
        } `

			restUri := nexus.RestURIs{
				Uri:     "/leader",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}
			e.RegisterRouter(restUri)
			model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
				[]string{}, nil, nil, false, "some description")
			model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

			req := httptest.NewRequest(http.MethodPost, "/:orgchart.Leader", strings.NewReader(leaderJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.Echo.NewContext(req, rec)
			c.SetParamNames("orgchart.Leader")
			c.SetParamValues("leader1")
			nc := &NexusContext{
				NexusURI:  "/leader",
				Context:   c,
				CrdType:   "leaders.orgchart.vmware.org",
				GroupName: "orgchart.vmware.org",
				Resource:  "leaders",
			}
			err := putHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))
		})

		It("should handle PATCH request", func() {
			// Modify `designation` value from `NexusLead` to `Manager`
			patchJson := `{
          "designation": "Manager",
          "new-field": "new-value"
        } `
			req := httptest.NewRequest(http.MethodPatch, "/:orgchart.Leader", strings.NewReader(patchJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.Echo.NewContext(req, rec)
			c.SetParamNames("orgchart.Leader")
			c.SetParamValues("leader1")
			nc := &NexusContext{
				NexusURI:  "/leader",
				Context:   c,
				CrdType:   "leaders.orgchart.vmware.org",
				GroupName: "orgchart.vmware.org",
				Resource:  "leaders",
			}

			// patch should be applied successfully
			err := patchHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))
			Expect(rec.Body.String()).To(Equal("{\"message\":\"Patch applied successfully\"}\n"))

			req = httptest.NewRequest(http.MethodGet, "/:orgchart.Leader", strings.NewReader(patchJson))
			rec = httptest.NewRecorder()
			c = e.Echo.NewContext(req, rec)
			c.SetParamNames("orgchart.Leader")
			c.SetParamValues("leader1")
			nc = &NexusContext{
				NexusURI:  "/leader",
				Context:   c,
				CrdType:   "leaders.orgchart.vmware.org",
				GroupName: "orgchart.vmware.org",
				Resource:  "leaders",
			}

			// `designation` field and new-field should only be modified,
			err = getHandler(nc)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Body.String()).To(Equal("{\"spec\":{\"description\":\"Hello World!\",\"designation\":\"Manager\"," +
				"\"employeeID\":1,\"new-field\":\"new-value\"},\"status\":{}}\n"))
		})

		It("should fail PATCH request when patch format is wrong", func() {
			patchJson := `[{
          "designation": "Manager",
          "new-field": "new-value"
        }]`
			req := httptest.NewRequest(http.MethodPatch, "/:orgchart.Leader", strings.NewReader(patchJson))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.Echo.NewContext(req, rec)
			c.SetParamNames("orgchart.Leader")
			c.SetParamValues("leader1")
			nc := &NexusContext{
				NexusURI:  "/leader",
				Context:   c,
				CrdType:   "leaders.orgchart.vmware.org",
				GroupName: "orgchart.vmware.org",
				Resource:  "leaders",
				Codes: map[nexus.ResponseCode]nexus.HTTPResponse{
					http.StatusBadRequest: {Description: http.StatusText(http.StatusBadRequest)},
				},
			}

			// patch should be failed when wrong format provided
			err := patchHandler(nc)
			Expect(err).To(HaveOccurred())
		})
	})
})

func createTestNode(apiVersion, kind, name string) string {
	return fmt.Sprintf(`{
		"apiVersion": "%s",
		"kind": "%s",
	   "metadata": {
	       "name": "%s"
	   },
	   "spec": {
	       "foo": "bar2"
	   }
	}`, apiVersion, kind, name)
}

func expectedEngManagersRestURIs() []nexus.RestURIs {
	return []nexus.RestURIs{
		{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}/status",
			Methods: map[nexus.HTTPMethod]nexus.HTTPCodesResponse{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
				http.MethodPut: nexus.DefaultHTTPPUTResponses,
			},
		},
		{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}/EngManagers",
			Methods: map[nexus.HTTPMethod]nexus.HTTPCodesResponse{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
	}
}

func expectedHRRestURIs() []nexus.RestURIs {
	return []nexus.RestURIs{
		{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}/status",
			Methods: map[nexus.HTTPMethod]nexus.HTTPCodesResponse{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
				http.MethodPut: nexus.DefaultHTTPPUTResponses,
			},
		},
		{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}/HR",
			Methods: map[nexus.HTTPMethod]nexus.HTTPCodesResponse{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
	}
}

func constructTestLinkAnnotation() model.NexusAnnotation {
	return model.NexusAnnotation{
		Name: "management.Leader",
		Links: map[string]model.NodeHelperChild{
			"mgrs.management.vmware.org": {
				FieldName:    "EngManagers",
				FieldNameGvk: "engManagersGvk",
				IsNamed:      true,
			},
		},
		NexusRestAPIGen: nexus.RestAPISpec{
			Uris: []nexus.RestURIs{{
				Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}},
		},
		IsSingleton: false,
	}
}

func constructTestChildrenAnnotation() model.NexusAnnotation {
	return model.NexusAnnotation{
		Name: "management.Leader",
		Children: map[string]model.NodeHelperChild{
			"humanresourceses.hr.vmware.org": {
				FieldName:    "HR",
				FieldNameGvk: "hRGvk",
				IsNamed:      false,
			},
		},
		NexusRestAPIGen: nexus.RestAPISpec{
			Uris: []nexus.RestURIs{{
				Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}},
		},
		IsSingleton: false,
	}
}

func constructTestAnnotation() model.NexusAnnotation {
	return model.NexusAnnotation{
		Name: "management.Leader",
		Children: map[string]model.NodeHelperChild{
			"humanresourceses.hr.vmware.org": {
				FieldName:    "HR",
				FieldNameGvk: "hRGvk",
				IsNamed:      false,
			},
		},
		Links: map[string]model.NodeHelperChild{
			"mgrs.management.vmware.org": {
				FieldName:    "EngManagers",
				FieldNameGvk: "engManagersGvk",
				IsNamed:      true,
			},
			"roles.management.vmware.org": {
				FieldName:    "roles",
				FieldNameGvk: "roleGvk",
				IsNamed:      true,
			},
		},
		NexusRestAPIGen: nexus.RestAPISpec{
			Uris: []nexus.RestURIs{{
				Uri:     "/root/{orgchart.Root}/leader/{management.Leader}",
				Methods: nexus.DefaultHTTPMethodsResponses,
			}},
		},
		IsSingleton: false,
	}
}

func initNode(e *EchoServer, crdType, groupName, resourceName, name, method, body, targetURI string, restUri nexus.RestURIs) (*NexusContext, *httptest.ResponseRecorder) {
	e.RegisterRouter(restUri)
	model.ConstructMapCRDTypeToNode(model.Upsert, crdType, name,
		[]string{}, nil, nil, true, "some description")
	model.ConstructMapURIToCRDType(model.Upsert, crdType, []nexus.RestURIs{restUri})

	return createTestNexusContext(e, crdType, groupName, resourceName, name, method, body, targetURI, restUri)
}

func createTestNexusContext(e *EchoServer, crdType, groupName, resourceName, name, method, body, targetURI string, restUri nexus.RestURIs) (*NexusContext, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, targetURI, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.Echo.NewContext(req, rec)
	c.SetParamNames(name)
	c.SetParamValues("default")
	return &NexusContext{
		NexusURI:  restUri.Uri,
		Context:   c,
		CrdType:   crdType,
		GroupName: groupName,
		Resource:  resourceName,
	}, rec
}

func createSampleEngManagerRequest(e *EchoServer) (*NexusContext, *httptest.ResponseRecorder) {
	return createTestNexusContext(e, "leaders.management.vmware.org", "management.vmware.org",
		"leaders", "management.Leader", http.MethodGet, "",
		"/root/:orgchart.Root/leader/:management.Leader/EngManagers", nexus.RestURIs{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}/EngManagers",
		})
}

func createSampleHRRequest(e *EchoServer) (*NexusContext, *httptest.ResponseRecorder) {
	return createTestNexusContext(e, "leaders.management.vmware.org", "management.vmware.org",
		"leaders", "management.Leader", http.MethodGet, "",
		"/root/:orgchart.Root/leader/:management.Leader/HR", nexus.RestURIs{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}/HR",
		})
}

func createSampleRoleRequest(e *EchoServer) (*NexusContext, *httptest.ResponseRecorder) {
	return createTestNexusContext(e, "leaders.management.vmware.org", "management.vmware.org",
		"leaders", "management.Leader", http.MethodGet, "",
		"/root/:orgchart.Root/leader/:management.Leader/Role", nexus.RestURIs{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}/Role",
		})
}

func createSampleLeaderRequest(e *EchoServer) (*NexusContext, *httptest.ResponseRecorder) {
	return createTestNexusContext(e, "leaders.management.vmware.org", "management.vmware.org",
		"leaders", "management.Leader", http.MethodGet, "",
		"/root/:orgchart.Root/leader/:management.Leader", nexus.RestURIs{
			Uri: "/root/{orgchart.Root}/leader/{management.Leader}",
		})
}

func getLeaderChildrenJson(key, val string) string {
	var additionalKey string
	if len(key) > 0 {
		additionalKey = fmt.Sprintf(`,"%s": "%s"`, key, val)
	}
	return fmt.Sprintf(`{
          "designation": "NexusLead",
          "employeeID": 1,
          "engManagersGvk": {
            "default": {
              "group": "management.vmware.org",
              "kind": "Mgr",
              "name": "eac9763b09291c96b4973c41036f841ba46aa502"
          },
            "foo": {
              "group": "management.vmware.org",
              "kind": "Mgr",
              "name": "545edb30e5d0b62628be5e5455843908f1d76b34"
          }
        },
          "hRGvk": {
            "group": "hr.vmware.org",
            "kind": "HumanResources",
            "name": "71d2f43510c62c8a4cc08ed4fffa58839d722608"
        } ` + additionalKey + `
	}`)
}
