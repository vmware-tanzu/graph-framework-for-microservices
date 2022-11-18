package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
		localClient           *fake_dynamic.FakeDynamicClient
		remoteClient          dynamic.Interface
		remoteHandler         *h.RemoteHandler
		replicationConfigSpec utils.ReplicationConfigSpec
		conf                  *config.Config
		repObj                utils.ReplicationObject
		logBuffer             bytes.Buffer
		repConfName           string
	)

	BeforeEach(func() {
		log.SetOutput(&logBuffer)
		log.SetLevel(log.DebugLevel)

		localClient = fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(),
			getObject("A", AcKind, "example"), getObject("B", AcKind, "example"), getObject("update", AcKind, "example"), getObject("create", AcKind, "example"),
			getObject("C", AcKind, "example"), getObject("D", AcKind, "example"), getParentObject("foo", AcKind), getChildObject("bar", AdKind),
			getObject("Root", RootKind, "example"), getObject("Config", ConfigKind, "example"), getObject("Project", ProjectKind, "example"), getDefaultResourceObj(),
		)

		repConfName = "one"
		conf = &config.Config{
			StatusReplicationEnabled: false,
			PeriodicSyncInterval:     3 * time.Second,
		}

		utils.CRDTypeToCrdVersion[ApiCollaborationSpace] = utils.V1Version
		utils.CRDTypeToCrdVersion[ApiDevSpace] = utils.V1Version
	})

	When("Replication is enabled for CRD Type", func() {
		It("Should replicate all the objects of that type to the destination endpoint", func() {
			server := ghttp.NewServer()
			remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
			Expect(err).NotTo(HaveOccurred())

			remoteHandler := h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
			replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
				Source: getTypeConfig(Group, AcKind), Destination: getNonHierarchicalDestConfig()}

			// server receives both objA and objB.
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

			// Enable replication for gvr: config.mazinger.com/v1, apicollaborationspaces.
			utils.ReplicationEnabledGVR[apicollaborationspace] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledGVR[apicollaborationspace][repConfName] = replicationConfigSpec

			// Create objA of type ApiCollaborationSpace.
			err = remoteHandler.Create(getObject("A", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objB of type ApiCollaborationSpace.
			err = remoteHandler.Create(getObject("B", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() []*http.Request { return server.ReceivedRequests() }, 3*time.Second).Should(HaveLen(4))

			delete(utils.ReplicationEnabledGVR, apicollaborationspace)
		})
	})

	When("Replication is configured for CRD Type to be replicated from one type to another", func() {
		It("Should replicate all the objects to the type configured in the replication config", func() {
			server := ghttp.NewServer()
			remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
			Expect(err).NotTo(HaveOccurred())

			remoteHandler := h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
			replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
				Source: getTypeConfig(Group, AcKind), Destination: getDifferentTypeDestConfig()}

			// server receives both objA and objB of Kind Config.
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/configs/A"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/configs"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"Config\",\"metadata\":{\"labels\":{},\"name\":\"A\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/configs/B"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/configs"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"Config\",\"metadata\":{\"labels\":{},\"name\":\"B\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
			)

			// Enable replication for gvr: config.mazinger.com/v1, apicollaborationspaces.
			utils.ReplicationEnabledGVR[apicollaborationspace] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledGVR[apicollaborationspace][repConfName] = replicationConfigSpec

			// Create objA of type ApiCollaborationSpace.
			err = remoteHandler.Create(getObject("A", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objB of type ApiCollaborationSpace.
			err = remoteHandler.Create(getObject("B", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() []*http.Request { return server.ReceivedRequests() }, 3*time.Second).Should(HaveLen(4))

			delete(utils.ReplicationEnabledGVR, apicollaborationspace)
		})
	})

	When("Replication is enabled for an individual object and if source is non-hierarchical", func() {
		It("Should replicate only that object to the destination endpoint", func() {
			server := ghttp.NewServer()
			remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
			Expect(err).NotTo(HaveOccurred())

			remoteHandler := h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
			replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
				Source: getNonHierarchicalSourceConfig(), Destination: getNonHierarchicalDestConfig()}

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

			utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

			// Create objC of type ApiCollaborationSpace.
			err = remoteHandler.Create(getObject("C", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objD but the destination endpoint will not receive the request.
			err = remoteHandler.Create(getObject("D", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() []*http.Request { return server.ReceivedRequests() }, 6*time.Second).Should(HaveLen(2))

			delete(utils.ReplicationEnabledNode, repObj)
		})
	})

	Context("Hierarchical source and non-hierarchical destination", func() {

		When("Replication is enabled for the object's parent", func() {
			It("Should replicate that object to the destination endpoint", func() {
				server := ghttp.NewServer()
				remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
				Expect(err).NotTo(HaveOccurred())

				replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
					Source: getHierarchicalSourceConfig("foo"), Destination: getNonHierarchicalDestConfig()}

				remoteHandler := h.NewRemoteHandler(apidevspace, localClient, nil, conf)
				utils.GVRToParentHierarchy[apidevspace] = []string{Root, Project,
					Config, ApiCollaborationSpace}

				// Enable replication for object foo of type apicollaborationspaces.config.mazinger.com
				repObj := utils.GetReplicationObject(Group, AcKind, "foo")
				utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
				utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

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

				// Create obj bar of type ApiDevSpace.
				err = remoteHandler.Create(getChildObject("bar", AdKind))
				Expect(err).NotTo(HaveOccurred())
				Eventually(func() []*http.Request { return server.ReceivedRequests() }, 6*time.Second).Should(HaveLen(2))

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})

		When("Replication is enabled for an object", func() {
			It("Should replicate that object and its immediate children to the destination endpoint", func() {
				server := ghttp.NewServer()
				remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
				Expect(err).NotTo(HaveOccurred())

				replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
					Source: getHierarchicalSourceConfig("foo"), Destination: getNonHierarchicalDestConfig()}

				remoteHandler := h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
				utils.GVRToChildren[apicollaborationspace] = utils.Children{
					ApiDevSpace: utils.NodeHelperChild{
						FieldNameGvk: "apiDevSpaceGvk",
					},
				}
				utils.GVRToParentHierarchy[apicollaborationspace] = []string{Root, Project, Config}

				// Enable replication for object foo of type apicollaborationspaces.config.mazinger.com
				repObj := utils.GetReplicationObject(Group, AcKind, "foo")
				utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
				utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

				// server receives "foo" and its immediate child "bar".
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

				// Create obj foo of type ApiCollaborationSpace.
				err = remoteHandler.Create(getParentObject("foo", AcKind))
				Expect(err).NotTo(HaveOccurred())
				Eventually(func() []*http.Request { return server.ReceivedRequests() }, 3*time.Second).Should(HaveLen(4))

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})

		When("There are two objects with same name under different hierarchy and replication is enabled for only one object", func() {
			It("Should not replicate the other object to the destination endpoint", func() {
				server := ghttp.NewServer()
				remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
				Expect(err).NotTo(HaveOccurred())

				source := getHierarchicalSourceConfig("foo")
				source.Object.Hierarchy = utils.Hierarchy{Labels: []utils.KVP{{Key: Root, Value: "root"}, {Key: Project, Value: "project"}, {Key: Config, Value: "different"}}}

				replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
					Source: source, Destination: getNonHierarchicalDestConfig()}

				remoteHandler := h.NewRemoteHandler(apidevspace, localClient, nil, conf)
				utils.GVRToParentHierarchy[apidevspace] = []string{Root, Project,
					Config, ApiCollaborationSpace}

				// Enable replication for object "Root/root/Project/default/Config/different/ApiCollaborationSpace/foo"
				repObj := utils.GetReplicationObject(Group, AcKind, "foo")
				utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
				utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

				// Create object "Root/root/Project/default/Config/config/ApiCollaborationSpace/foo"
				err = remoteHandler.Create(getObject("foo", AcKind, "example"))
				Expect(err).NotTo(HaveOccurred())

				// Destination didn't receive it.
				Eventually(func() []*http.Request { return server.ReceivedRequests() }).Should(HaveLen(0))

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
	})

	Context("Non-hierarchical source and Hierarchical destination", func() {
		server := ghttp.NewServer()
		remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
		Expect(err).NotTo(HaveOccurred())

		BeforeEach(func() {
			source := getNonHierarchicalSourceConfig()
			destination := getHierarchicalDestConfig()
			replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
				Source: source, Destination: destination, StatusEndpoint: utils.Source}

			// Enable replication for object bar of type apidevspaces.config.mazinger.com
			repObj := utils.GetReplicationObject(Group, AdKind, "bar")
			utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec
		})
		When("Replication is enabled for an object of individual type", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(apidevspace, localClient, nil, conf)
				utils.GVRToParentHierarchy[apidevspace] = []string{Root, Project,
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
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/configs"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"v1\",\"kind\":\"List\",\"items\":[{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"name\":\"status_fail\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}},{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"name\":\"status_fail\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}]}"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"bar\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"foo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)

				// Create obj bar without parent labels.
				err := remoteHandler.Create(getObject("bar", AdKind, "example"))
				Expect(err).NotTo(HaveOccurred())
				Eventually(func() []*http.Request { return server.ReceivedRequests() }, 3*time.Second).Should(HaveLen(3))

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
	})

	Context("Default K8s Resource Types", func() {
		server := ghttp.NewServer()
		remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
		Expect(err).NotTo(HaveOccurred())

		BeforeEach(func() {
			source := getTypeConfig("apps", "Deployment")
			destination := getHierarchicalDestConfig()
			replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
				Source: source, Destination: destination, StatusEndpoint: utils.Source}

			// Enable replication for Deployment.
			repObj = utils.GetReplicationObject("apps", "Deployment", "zoo")
			utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec
		})
		When("Replication is enabled for default K8s resource types", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(deployment, localClient, nil, conf)
			})
			It("Should replicate deployment objects to the desired destination hierarchy", func() {
				// server receives obj "bar" with parent labels set.
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/apps/v1/deployments/zoo"),
						ghttp.RespondWith(404, "not found"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/configs"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"v1\",\"kind\":\"List\",\"items\":[{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"name\":\"status_fail\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}},{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"name\":\"status_fail\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}]}"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/apps/v1/deployments"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"Deployment\",\"metadata\":{\"name\":\"zoo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)

				err = remoteHandler.Create(getDefaultResourceObj())
				Expect(err).NotTo(HaveOccurred())
				Eventually(func() []*http.Request { return server.ReceivedRequests() }, 3*time.Second).Should(HaveLen(3))

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
	})

	Context("Update", func() {
		When("Update event occurs for a replication enabled object", func() {
			BeforeEach(func() {
				source := getHierarchicalSourceConfig("update")
				destination := getNonHierarchicalDestConfig()
				remoteHandler = h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
				remoteClient = fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(), getObject("update", AcKind, "example"))
				replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: source, Destination: destination}

				utils.GVRToChildren[apicollaborationspace] = utils.Children{
					ApiDevSpace: utils.NodeHelperChild{
						FieldNameGvk: "apiDevSpaceGvk",
					},
				}
			})

			It("Should update the spec of that object to the destination endpoint", func() {
				// Enable replication for object New of type apicollaborationspaces.config.mazinger.com
				repObj := utils.GetReplicationObject(Group, AcKind, "update")
				utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
				utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

				// Update objNew.
				expectedObj := getObject("update", AcKind, "NEW_VALUE")
				err := remoteHandler.Update(expectedObj, getObject("update", AcKind, "example"))
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

	It("Should fail if Create Event for object of wrong type is received", func() {
		err := remoteHandler.Create(0)
		Expect(err.Error()).To(ContainSubstring("unstructured client did not understand object during create event"))
	})

	It("Should fail if Update Event for object of wrong type is received", func() {
		err := remoteHandler.Update(getObject("update", AcKind, "NEW_VALUE"), "wrongtype")
		Expect(err.Error()).To(ContainSubstring("unstructured client did not understand object during update event"))
	})

	It("Should skip creation if already exists", func() {
		source := getNonHierarchicalSourceConfig()
		destination := getNonHierarchicalDestConfig()
		remoteHandler = h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)

		remoteClient := fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(), getObject("C", AcKind, "example"))
		replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: source, Destination: destination}

		repObj := utils.GetReplicationObject(Group, AcKind, "C")

		utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
		utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

		err := remoteHandler.Create(getObject("C", AcKind, "example"))
		Expect(err).NotTo(HaveOccurred())
	})

	It("Should fail when source object not found during resync.", func() {
		server := ghttp.NewServer()
		remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
		Expect(err).NotTo(HaveOccurred())

		destination := getHierarchicalDestConfig()
		destination.Hierarchy.Labels = append(destination.Hierarchy.Labels,
			utils.KVP{Key: "invalid.config.mazinger.com", Value: "default"})

		replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: getNonHierarchicalSourceConfig(), Destination: destination}

		repObj := utils.GetReplicationObject(Group, AdKind, "new")
		utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
		utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

		remoteHandler := h.NewRemoteHandler(apidevspace, localClient, nil, conf)
		utils.GVRToParentHierarchy[apidevspace] = []string{Root, Project,
			Config, ApiCollaborationSpace}

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/new"),
				ghttp.RespondWith(404, "not found"),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/invalid"),
				ghttp.RespondWith(404, "not found"),
			),
		)

		obj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": ApiVersion,
				"kind":       AdKind,
				"metadata": map[string]interface{}{
					"name": "new",
				},
			},
		}

		err = remoteHandler.Create(obj)
		Eventually(func() string { return logBuffer.String() }, time.Second*35).
			Should(ContainSubstring("Source object new not found: apidevspaces.config.mazinger.com"))
		Eventually(func() []*http.Request { return server.ReceivedRequests() }).Should(HaveLen(2))

		delete(utils.ReplicationEnabledNode, repObj)
	})

	It("Should retry when sync fails the first time.", func() {
		server := ghttp.NewServer()
		remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
		Expect(err).NotTo(HaveOccurred())

		replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
			Source: getHierarchicalSourceConfig("foo"), Destination: getNonHierarchicalDestConfig()}

		repObj := utils.GetReplicationObject(Group, AcKind, "foo")
		utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
		utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

		remoteHandler := h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
		utils.GVRToParentHierarchy[apicollaborationspace] = []string{Root, Project,
			Config, ApiCollaborationSpace}

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/foo"),
				ghttp.RespondWith(404, "not found"),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
				ghttp.RespondWith(404, "Error"),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
				ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"foo\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"foo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
			),
		)

		err = remoteHandler.Create(getObject("foo", AcKind, "example"))
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() []*http.Request { return server.ReceivedRequests() }).Should(HaveLen(3))

		delete(utils.ReplicationEnabledNode, repObj)
	})

	When("Replication is filtered based on namespace", func() {
		It("Should sync objects only from the namespace of interest.", func() {
			server := ghttp.NewServer()
			remoteClient, err := utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "", nil)
			Expect(err).NotTo(HaveOccurred())

			source := getTypeConfig(Group, AcKind)
			source.Filters.Namespace = "required_ns"
			replicationConfigSpec := utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
				Source: source, Destination: getNonHierarchicalDestConfig()}

			utils.ReplicationEnabledGVR[apicollaborationspace] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledGVR[apicollaborationspace][repConfName] = replicationConfigSpec

			remoteHandler := h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
			utils.GVRToParentHierarchy[apicollaborationspace] = []string{Root, Project,
				Config, ApiCollaborationSpace}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/foo"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiCollaborationSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"foo\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"foo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
				),
			)

			// Object from non-desired namespace should not be synced.
			obj1 := getObject("foo", AcKind, "example")
			obj1.UnstructuredContent()["metadata"].(map[string]interface{})["namespace"] = "NOT_required_ns"

			err = remoteHandler.Create(obj1)
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() []*http.Request { return server.ReceivedRequests() }).Should(HaveLen(0))

			// Object from desired namespace should be synced.
			obj2 := getObject("foo", AcKind, "example")
			obj2.UnstructuredContent()["metadata"].(map[string]interface{})["namespace"] = "required_ns"

			err = remoteHandler.Create(obj2)
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() []*http.Request { return server.ReceivedRequests() }).Should(HaveLen(2))

			delete(utils.ReplicationEnabledNode, repObj)
		})
	})
})
