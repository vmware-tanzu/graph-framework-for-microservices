package handlers_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
		repObj                    utils.ReplicationObject
	)

	BeforeEach(func() {
		localClient = fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(),
			GetObject("A", AcKind, "example"), GetObject("B", AcKind, "example"), GetObject("update", AcKind, "example"),
			GetObject("C", AcKind, "example"), GetObject("D", AcKind, "example"), GetParentObject("foo", AcKind), GetChildObject("bar", AdKind),
			GetObject("Root", RootKind, "example"), GetObject("Config", ConfigKind, "example"), GetObject("Project", ProjectKind, "example"),
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
			source := GetTypeConfig()
			destination := GetNonHierarchicalDestConfig()
			remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(ApiCollaborationSpace), ApiCollaborationSpace, localClient, nil, conf)
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
			utils.ReplicationEnabledCRDType[ApiCollaborationSpace] = replicationConfigSpec

			// Create objA
			err = remoteHandler.Create(GetObject("A", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objB
			err = remoteHandler.Create(GetObject("B", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledCRDType, ApiCollaborationSpace)
		})
	})

	When("Replication is configured for CRD Type to be replicated from one type to another", func() {
		BeforeEach(func() {
			source := GetTypeConfig()
			destination := GetDifferentTypeDestConfig()
			remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(ApiCollaborationSpace), ApiCollaborationSpace, localClient, nil, conf)
			replicationConfigSpec = utils.ReplicationConfigSpec{Client: remoteClient, Source: source, Destination: destination}
		})

		It("Should replicate all the objects to the type configured in the replication config", func() {
			// server receives both objA and objB of type ApiDevSpace
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/A"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{},\"name\":\"A\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/B"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{},\"name\":\"B\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
			)

			// Enable replication for type apicollaborationspaces.config.mazinger.com
			utils.ReplicationEnabledCRDType[ApiCollaborationSpace] = replicationConfigSpec

			// Create objA of type ApiCollaborationSpace
			err = remoteHandler.Create(GetObject("A", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objB of type ApiCollaborationSpace
			err = remoteHandler.Create(GetObject("B", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledCRDType, ApiCollaborationSpace)
		})
	})

	When("Replication is enabled for an individual object and if source is non-hierarchical", func() {
		BeforeEach(func() {
			source := GetNonHierarchicalSourceConfig()
			destination := GetNonHierarchicalDestConfig()
			remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(ApiCollaborationSpace), ApiCollaborationSpace, localClient, nil, conf)
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
			repObj := utils.GetReplicationObject(Group, AcKind, "C")
			utils.ReplicationEnabledNode[repObj] = replicationConfigSpec

			// Create objC.
			err = remoteHandler.Create(GetObject("C", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objD but the destination endpoint will not receive the request.
			err = remoteHandler.Create(GetObject("D", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledNode, repObj)
		})
	})

	Context("Hierarchical source and non-hierarchical destination", func() {
		BeforeEach(func() {
			source := GetHierarchicalSourceConfig()
			destination := GetNonHierarchicalDestConfig()
			replicationConfigSpec = utils.ReplicationConfigSpec{Client: remoteClient, Source: source, Destination: destination}

			// Enable replication for object foo of type apicollaborationspaces.config.mazinger.com
			repObj = utils.GetReplicationObject(Group, AcKind, "foo")
			utils.ReplicationEnabledNode[repObj] = replicationConfigSpec
		})
		When("Replication is enabled for the object's parent", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(ApiDevSpace), ApiDevSpace, localClient, nil, conf)
				utils.CRDTypeToParentHierarchy[ApiDevSpace] = []string{Root, Project,
					Config, ApiCollaborationSpace}
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
				err = remoteHandler.Create(GetChildObject("bar", AdKind))
				Expect(err).NotTo(HaveOccurred())

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
		When("Replication is enabled for an object", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(ApiCollaborationSpace), ApiCollaborationSpace, localClient, nil, conf)
				utils.CRDTypeToChildren[ApiCollaborationSpace] = utils.Children{
					ApiDevSpace: utils.NodeHelperChild{
						FieldNameGvk: "apiDevSpaceGvk",
					},
				}
				utils.CRDTypeToParentHierarchy[ApiCollaborationSpace] = []string{Root, Project, Config}
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
				err = remoteHandler.Create(GetParentObject("foo", AcKind))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("Non-hierarchical source and Hierarchical destination", func() {
		BeforeEach(func() {
			source := GetNonHierarchicalSourceConfig()
			destination := GetHierarchicalDestConfig()
			replicationConfigSpec = utils.ReplicationConfigSpec{Client: remoteClient, Source: source, Destination: destination}

			// Enable replication for object bar of type apidevspaces.config.mazinger.com
			repObj = utils.GetReplicationObject(Group, AdKind, "bar")
			utils.ReplicationEnabledNode[repObj] = replicationConfigSpec
		})
		When("Replication is enabled for an object of individual type", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(ApiDevSpace), ApiDevSpace, localClient, nil, conf)
				utils.CRDTypeToParentHierarchy[ApiDevSpace] = []string{Root, Project,
					Config, ApiCollaborationSpace}
			})
			It("Should replicate that object with parent info as labels", func() {
				// server receives obj "bar" with parent labels set.
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/bar"),
						ghttp.RespondWith(404, "not found"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"bar\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"foo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)

				// Create obj bar without parent labels.
				err = remoteHandler.Create(GetObject("bar", AdKind, "example"))
				Expect(err).NotTo(HaveOccurred())

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
	})

	Context("Update", func() {
		When("Update event occurs for a replication enabled object", func() {
			BeforeEach(func() {
				source := GetNonHierarchicalSourceConfig()
				destination := GetNonHierarchicalDestConfig()
				remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(ApiCollaborationSpace), ApiCollaborationSpace, localClient, nil, conf)
				remoteClient = fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(), GetObject("update", AcKind, "example"))
				replicationConfigSpec = utils.ReplicationConfigSpec{Client: remoteClient, Source: source, Destination: destination}
			})

			It("Should update the spec of that object to the destination endpoint", func() {
				// Enable replication for object New of type apicollaborationspaces.config.mazinger.com
				repObj := utils.GetReplicationObject(Group, AcKind, "update")
				utils.ReplicationEnabledNode[repObj] = replicationConfigSpec

				// Update objNew.
				expectedObj := GetObject("update", AcKind, "NEW_VALUE")
				err = remoteHandler.Update(expectedObj, GetObject("update", AcKind, "example"))
				Expect(err).NotTo(HaveOccurred())

				gvr := schema.GroupVersionResource{
					Group:    "config.mazinger.com",
					Version:  "v1",
					Resource: "apicollaborationspaces",
				}
				newObj, err := remoteClient.Resource(gvr).Get(context.TODO(), "update", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				Expect(newObj.UnstructuredContent()["spec"]).To(Equal(expectedObj.UnstructuredContent()["spec"]))

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
	})
})
