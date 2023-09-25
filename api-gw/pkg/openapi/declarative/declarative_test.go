package declarative_test

import (
	"api-gw/pkg/config"
	"api-gw/pkg/openapi/declarative"
	"api-gw/pkg/server/echo_server"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	nexus_client "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/nexus-client"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("OpenAPI tests", func() {
	It("should setup and load openapi file", func() {
		openApiSpecFile := "testFile"
		f, err := os.Create(openApiSpecFile)
		defer os.RemoveAll(openApiSpecFile)
		Expect(err).To(BeNil())
		f.Sync()
		defer f.Close()
		bytesWritten, err := f.Write(spec)
		Expect(err).To(BeNil())
		Expect(bytesWritten).ToNot(Equal(0))
		f.Sync()
		err = declarative.Setup(openApiSpecFile)
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
		e := echo_server.NewEchoServer(config.Cfg, &kubernetes.Clientset{}, &nexus_client.Clientset{})
		e.RegisterDeclarativeRouter()

		c := e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/gns.vmware.org/v1/globalnamespaces", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces"))

		// short name
		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/v1/gns", c)
		Expect(c.Path()).To(Equal("/apis/v1/gns"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/gns.vmware.org/v1/globalnamespaces/:name", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces/:name"))

		// short name
		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/v1/gns/:name", c)
		Expect(c.Path()).To(Equal("/apis/v1/gns/:name"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodPut, "/apis/gns.vmware.org/v1/globalnamespaces", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces"))

		// short name
		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/v1/gns", c)
		Expect(c.Path()).To(Equal("/apis/v1/gns"))

		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodDelete, "/apis/gns.vmware.org/v1/globalnamespaces/:name", c)
		Expect(c.Path()).To(Equal("/apis/gns.vmware.org/v1/globalnamespaces/:name"))

		// short name
		c = e.Echo.NewContext(nil, nil)
		e.Echo.Router().Find(http.MethodGet, "/apis/v1/gns/:name", c)
		Expect(c.Path()).To(Equal("/apis/v1/gns/:name"))
	})

	It("should parse schema for GlobalNamespace", func() {
		ec := declarative.SetupContext(Uri, http.MethodGet, declarative.Paths[Uri].Get)
		declarative.AddApisEndpoint(ec)

		Expect(declarative.ApisList).To(HaveKey("/apis/gns.vmware.org/v1/globalnamespaces"))
		Expect(declarative.ApisList["/apis/gns.vmware.org/v1/globalnamespaces"]).To(HaveKey("yaml"))

		expectedYaml := `apiVersion: gns.vmware.org/v1
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
		Expect(declarative.ApisList["/apis/gns.vmware.org/v1/globalnamespaces"]["yaml"]).To(Equal(expectedYaml))
	})

	//It("should create a list of short names", func() {
	//	apisList := map[string]map[string]interface{}{
	//		"/one": {
	//			"POST": map[string]interface{}{
	//				"group": "vmware.org",
	//				"kind":  "one",
	//			},
	//			"GET": map[string]interface{}{
	//				"group": "vmware.org",
	//				"kind":  "one",
	//			},
	//		},
	//		"/two": {
	//			"POST": map[string]interface{}{
	//				"group": "vmware.org",
	//				"kind":  "two",
	//			},
	//			"GET": map[string]interface{}{
	//				"group": "vmware.org",
	//				"kind":  "two",
	//			},
	//			"PUT": map[string]interface{}{
	//				"group": "different.vmware.org",
	//				"kind":  "different",
	//			},
	//			"DELETE": map[string]interface{}{
	//				"group": "different.vmware.org",
	//				"kind":  "different-kind",
	//			},
	//		},
	//		"/oneone": {
	//			"POST": map[string]interface{}{
	//				"group": "vmware.org",
	//				"kind":  "one",
	//			},
	//			"GET": map[string]interface{}{
	//				"group": "vmware.org",
	//				"kind":  "one",
	//			},
	//		},
	//		"/three": {
	//			"POST": map[string]interface{}{
	//				"group": "vmware.org",
	//				"kind":  "three",
	//			},
	//			"PUT": map[string]interface{}{
	//				"group": "different.vmware.org",
	//				"kind":  "different",
	//			},
	//		},
	//	}
	//
	//	expectedShortNames := map[string]string{
	//		"threes":          "threes.vmware.org",
	//		"twos":            "twos.vmware.org",
	//		"different-kinds": "different-kinds.different.vmware.org",
	//	}
	//	shortNames := declarative.ShortNames(apisList)
	//	Expect(shortNames).To(BeEquivalentTo(expectedShortNames))
	//
	//})
})
