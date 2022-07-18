package echo_server

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"api-gw/pkg/config"
	"api-gw/pkg/model"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
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

	It("should handle put query for singleton object with default as name", func() {
		leaderJson := `{
		"apiVersion": "orgchart.vmware.org/v1",
		"kind": "Leader",
	   "metadata": {
	       "name": "default"
	   },
	   "spec": {
	       "foo": "bar2"
	   }
	}`

		restUri := nexus.RestURIs{
			Uri:     "/leader",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "management.Leader",
			[]string{}, nil, true, "some description")
		model.ConstructMapURIToCRDType(model.Upsert, "leaders.orgchart.vmware.org", []nexus.RestURIs{restUri})

		req := httptest.NewRequest(http.MethodPost, "/leader", strings.NewReader(leaderJson))
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

	It("shouldn't handle put query for singleton object with not default as name", func() {
		leaderJson := `{
	"apiVersion": "orgchart.vmware.org/v1",
	"kind": "Leader",
    "metadata": {
        "name": "notdefault"
    },
    "spec": {
        "foo": "bar2"
    }
}`

		restUri := nexus.RestURIs{
			Uri:     "/leader/{orgchart.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		}
		e.RegisterRouter(restUri)
		model.ConstructMapCRDTypeToNode(model.Upsert, "leaders.orgchart.vmware.org", "orgchart.Leader",
			[]string{}, nil, true, "")
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
	})
})
