package handlers_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	fake_dynamic "k8s.io/client-go/dynamic/fake"

	"connector/pkg/config"
	h "connector/pkg/handlers"
	"connector/pkg/utils"
)

var _ = Describe("Create", func() {
	var (
		localClient, remoteClient dynamic.Interface
		remoteHandler             *h.RemoteHandler
		replicationConfigSpec     utils.ReplicationConfigSpec
		server                    *ghttp.Server
		err                       error
		conf                      *config.Config
	)

	BeforeEach(func() {
		localClient = fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(),
			h.GetObject("A", h.AcKind), h.GetObject("B", h.AcKind),
			h.GetObject("C", h.AcKind), h.GetObject("D", h.AcKind),
			h.GetParentObject("foo", h.AcKind), h.GetChildObject("bar", h.AdKind),
		)

		conf = &config.Config{
			StatusReplicationEnabled: false,
		}
		server = ghttp.NewServer()
		remoteClient, err = utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "")
		Expect(err).NotTo(HaveOccurred())
	})

	When("Replication is enabled for CRD Type", func() {
		BeforeEach(func() {
			source := h.GetTypeConfig()
			destination := h.GetNonHierarchicalDestConfig()
			remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(h.ApiCollaborationSpace), h.ApiCollaborationSpace, localClient, nil, conf)
			replicationConfigSpec = utils.ReplicationConfigSpec{Client: remoteClient, Source: source, Destination: destination}
		})

		It("Should replicate all the objects of that type to the destination endpoint", func() {
			// server receives both objA and objB
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/A"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{},\"name\":\"A\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/B"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{},\"name\":\"B\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
			)

			// Enable replication for type apicollaborationspaces.config.mazinger.com
			utils.ReplicationEnabledCRDType[h.ApiCollaborationSpace] = replicationConfigSpec

			// Create objA
			err = remoteHandler.Create(h.GetObject("A", h.AcKind))
			Expect(err).NotTo(HaveOccurred())

			// Create objB
			err = remoteHandler.Create(h.GetObject("B", h.AcKind))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledCRDType, h.ApiCollaborationSpace)
		})
	})

	When("Replication is enabled for an individual object and if source is non-hierarchical", func() {
		BeforeEach(func() {
			source := h.GetNonHierarchicalSourceConfig()
			destination := h.GetNonHierarchicalDestConfig()
			remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(h.ApiCollaborationSpace), h.ApiCollaborationSpace, localClient, nil, conf)
			replicationConfigSpec = utils.ReplicationConfigSpec{Client: remoteClient, Source: source, Destination: destination}
		})

		It("Should replicate only that object to the destination endpoint", func() {
			// server receives only objC.
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/C"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{},\"name\":\"C\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
			)

			// Enable replication for object C of type apicollaborationspaces.config.mazinger.com
			repObj := utils.GetReplicationObject(h.Group, h.AcKind, "C")
			utils.ReplicationEnabledNode[repObj] = replicationConfigSpec

			// Create objC.
			err = remoteHandler.Create(h.GetObject("C", h.AcKind))
			Expect(err).NotTo(HaveOccurred())

			// Create objD but the destination endpoint will not receive the request.
			err = remoteHandler.Create(h.GetObject("D", h.AcKind))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledNode, repObj)
		})
	})

	Context("Hierarchical source and non-hierarchical destination", func() {
		BeforeEach(func() {
			source := h.GetHierarchicalSourceConfig()
			destination := h.GetNonHierarchicalDestConfig()
			replicationConfigSpec = utils.ReplicationConfigSpec{Client: remoteClient, Source: source, Destination: destination}

			// Enable replication for object foo of type apicollaborationspaces.config.mazinger.com
			repObj := utils.GetReplicationObject(h.Group, h.AcKind, "foo")
			utils.ReplicationEnabledNode[repObj] = replicationConfigSpec
		})
		When("Replication is enabled for the object's parent", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(h.ApiDevSpace), h.ApiDevSpace, localClient, nil, conf)
				utils.CRDTypeToParentHierarchy[h.ApiDevSpace] = []string{h.Root, h.Project,
					h.Config, h.ApiCollaborationSpace}
			})
			It("Should replicate that object to the destination endpoint", func() {
				// server receives obj "bar".
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/bar"),
						ghttp.RespondWith(404, "not found"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{\"apicollaborationspaces.config.mazinger.com\":\"foo\",\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"bar\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"bar\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)

				// Create obj bar.
				err = remoteHandler.Create(h.GetChildObject("bar", h.AdKind))
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("Replication is enabled for an object", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(h.ApiCollaborationSpace), h.ApiCollaborationSpace, localClient, nil, conf)
				utils.CRDTypeToChildren[h.ApiCollaborationSpace] = utils.Children{
					h.ApiDevSpace: utils.NodeHelperChild{
						FieldNameGvk: "apiDevSpaceGvk",
					},
				}
				utils.CRDTypeToParentHierarchy[h.ApiCollaborationSpace] = []string{h.Root, h.Project, h.Config}
			})
			It("Should replicate that object and its immediate children to the destination endpoint", func() {
				// server receives obj "foo" and its immediate child "bar".
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/foo"),
						ghttp.RespondWith(404, "not found"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"foo\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"foo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
					//Child request
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/bar"),
						ghttp.RespondWith(404, "not found"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"bar\",\"apicollaborationspaces.config.mazinger.com\":\"foo\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"foo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)

				// Create obj foo.
				err = remoteHandler.Create(h.GetParentObject("foo", h.AcKind))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
