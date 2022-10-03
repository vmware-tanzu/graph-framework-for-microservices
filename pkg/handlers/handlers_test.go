package handlers_test

import (
	"bytes"
	"context"
	"fmt"

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
		server                *ghttp.Server
		err                   error
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
		}
		server = ghttp.NewServer()
		remoteClient, err = utils.SetUpDynamicRemoteAPI(fmt.Sprintf("http://%s", server.Addr()), "", "")
		Expect(err).NotTo(HaveOccurred())

		utils.CRDTypeToCrdVersion[ApiCollaborationSpace] = utils.V1Version
		utils.CRDTypeToCrdVersion[ApiDevSpace] = utils.V1Version
	})

	When("Replication is enabled for CRD Type", func() {
		BeforeEach(func() {
			source := getTypeConfig(Group, AcKind)
			destination := getNonHierarchicalDestConfig()
			remoteHandler = h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
			replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: source, Destination: destination}
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

			// Enable replication for gvr: config.mazinger.com/v1, apicollaborationspaces.
			utils.ReplicationEnabledGVR[apicollaborationspace] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledGVR[apicollaborationspace][repConfName] = replicationConfigSpec

			// Create objA
			err = remoteHandler.Create(getObject("A", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objB
			err = remoteHandler.Create(getObject("B", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledGVR, apicollaborationspace)
		})

		It("Should fail object creation when destination is down", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apicollaborationspaces/A"),
					ghttp.RespondWith(404, "not found"),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apicollaborationspaces"),
					ghttp.RespondWith(400, "NotOK"),
				),
			)

			// Enable replication for gvr: config.mazinger.com/v1, apicollaborationspaces.
			utils.ReplicationEnabledGVR[apicollaborationspace] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledGVR[apicollaborationspace][repConfName] = replicationConfigSpec

			// Create objA
			err = remoteHandler.Create(getObject("A", AcKind, "example"))
			Expect(logBuffer.String()).To(ContainSubstring("Resource A create failed with an error"))

			delete(utils.ReplicationEnabledGVR, apicollaborationspace)
		})
	})

	When("Replication is configured for CRD Type to be replicated from one type to another", func() {
		BeforeEach(func() {
			source := getTypeConfig(Group, AcKind)
			destination := getDifferentTypeDestConfig()
			remoteHandler = h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
			replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: source, Destination: destination}
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

			// Enable replication for gvr: config.mazinger.com/v1, apicollaborationspaces.
			utils.ReplicationEnabledGVR[apicollaborationspace] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledGVR[apicollaborationspace][repConfName] = replicationConfigSpec

			// Create objA of type ApiCollaborationSpace
			err = remoteHandler.Create(getObject("A", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objB of type ApiCollaborationSpace
			err = remoteHandler.Create(getObject("B", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledGVR, apicollaborationspace)
		})
	})

	When("Replication is enabled for an individual object and if source is non-hierarchical", func() {
		BeforeEach(func() {
			source := getNonHierarchicalSourceConfig()
			destination := getNonHierarchicalDestConfig()
			remoteHandler = h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
			replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: source, Destination: destination}
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

			utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

			// Create objC.
			err = remoteHandler.Create(getObject("C", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			// Create objD but the destination endpoint will not receive the request.
			err = remoteHandler.Create(getObject("D", AcKind, "example"))
			Expect(err).NotTo(HaveOccurred())

			delete(utils.ReplicationEnabledNode, repObj)
		})
	})

	Context("Hierarchical source and non-hierarchical destination", func() {
		BeforeEach(func() {
			source := getHierarchicalSourceConfig()
			destination := getNonHierarchicalDestConfig()
			replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: source, Destination: destination}

			// Enable replication for object foo of type apicollaborationspaces.config.mazinger.com
			repObj = utils.GetReplicationObject(Group, AcKind, "foo")
			utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
			utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec
		})
		When("Replication is enabled for the object's parent", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(apidevspace, localClient, nil, conf)
				utils.GVRToParentHierarchy[apidevspace] = []string{Root, Project,
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
				err = remoteHandler.Create(getChildObject("bar", AdKind))
				Expect(err).NotTo(HaveOccurred())

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
		When("Replication is enabled for an object", func() {
			BeforeEach(func() {
				remoteHandler = h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)
				utils.GVRToChildren[apicollaborationspace] = utils.Children{
					ApiDevSpace: utils.NodeHelperChild{
						FieldNameGvk: "apiDevSpaceGvk",
					},
				}
				utils.GVRToParentHierarchy[apicollaborationspace] = []string{Root, Project, Config}
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
				err = remoteHandler.Create(getParentObject("foo", AcKind))
				Expect(err).NotTo(HaveOccurred())
			})

			It("Should fail when children creation fails", func() {
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
						ghttp.RespondWith(400, "NotOk"),
					),
				)

				// Create obj foo.
				err = remoteHandler.Create(getParentObject("foo", AcKind))
				Expect(logBuffer.String()).To(ContainSubstring("Children replication failed for the resource" +
					" roots.config.mazinger.com/root/projects.config.mazinger.com/project/configs.config.mazinger.com/config/apicollaborationspaces.config.mazinger.com/foo"))
			})
		})
	})

	Context("Non-hierarchical source and Hierarchical destination", func() {
		BeforeEach(func() {
			source := getNonHierarchicalSourceConfig()
			destination := getHierarchicalDestConfig()
			replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient,
				Source: source, Destination: destination, StatusEndpoint: utils.Source}

			// Enable replication for object bar of type apidevspaces.config.mazinger.com
			repObj = utils.GetReplicationObject(Group, AdKind, "bar")
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
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"bar\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"foo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)

				// Create obj bar without parent labels.
				err = remoteHandler.Create(getObject("bar", AdKind, "example"))
				Expect(err).NotTo(HaveOccurred())

				delete(utils.ReplicationEnabledNode, repObj)
			})

			It("Should fail creation when destination is down", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/bar"),
						ghttp.RespondWith(404, "not found"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
						ghttp.RespondWith(400, "NotOK"),
					),
				)

				err = remoteHandler.Create(getObject("bar", AdKind, "example"))
				Expect(logBuffer.String()).To(ContainSubstring("create failed with an error"))

				delete(utils.ReplicationEnabledNode, repObj)
			})

			It("Should fail status patch with invalid object", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/apis/config.mazinger.com/v1/apidevspaces/status_fail"),
						ghttp.RespondWith(404, "not found"),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/apis/config.mazinger.com/v1/apidevspaces"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"config.mazinger.com/v1\",\"kind\":\"ApiDevSpace\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"status_fail\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"status_fail\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)
				repObj = utils.GetReplicationObject(Group, AdKind, "status_fail")
				utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
				utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

				obj := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": ApiVersion,
						"kind":       AdKind,
						"metadata": map[string]interface{}{
							"name": "status_fail",
						},
					},
				}

				err = remoteHandler.Create(obj)
				Expect(logBuffer.String()).To(ContainSubstring("Resource status_fail patching failed with an error"))

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
	})

	Context("Default K8s Resource Types", func() {
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
						ghttp.VerifyRequest("POST", "/apis/apps/v1/deployments"),
						ghttp.RespondWith(200, "{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"labels\":{\"configs.apix.mazinger.com\":\"config\",\"nexus/display_name\":\"zoo\",\"projects.apix.mazinger.com\":\"project\",\"roots.apix.mazinger.com\":\"root\"},\"name\":\"zoo\",\"namespace\":\"\"},\"spec\":{\"example\":\"example\"}}"),
					),
				)

				err = remoteHandler.Create(getDefaultResourceObj())
				Expect(err).NotTo(HaveOccurred())

				delete(utils.ReplicationEnabledNode, repObj)
			})
		})
	})

	Context("Update", func() {
		When("Update event occurs for a replication enabled object", func() {
			BeforeEach(func() {
				source := getHierarchicalSourceConfig()
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
				err = remoteHandler.Update(expectedObj, getObject("update", AcKind, "example"))
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
		err = remoteHandler.Create(0)
		Expect(err.Error()).To(ContainSubstring("unstructured client did not understand object during create event"))
	})

	It("Should fail if Update Event for object of wrong type is received", func() {
		err = remoteHandler.Update(getObject("update", AcKind, "NEW_VALUE"), "wrongtype")
		Expect(err.Error()).To(ContainSubstring("unstructured client did not understand object during update event"))
	})

	It("Should skip creation if already exists", func() {
		source := getNonHierarchicalSourceConfig()
		destination := getNonHierarchicalDestConfig()
		remoteHandler = h.NewRemoteHandler(apicollaborationspace, localClient, nil, conf)

		remoteClient := fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(), getObject("C", AcKind, "example"))
		replicationConfigSpec = utils.ReplicationConfigSpec{LocalClient: localClient, RemoteClient: remoteClient, Source: source, Destination: destination}

		repObj := utils.GetReplicationObject(Group, AcKind, "C")

		utils.ReplicationEnabledNode[repObj] = make(map[string]utils.ReplicationConfigSpec)
		utils.ReplicationEnabledNode[repObj][repConfName] = replicationConfigSpec

		err = remoteHandler.Create(getObject("C", AcKind, "example"))
		Expect(err).NotTo(HaveOccurred())
	})
})
