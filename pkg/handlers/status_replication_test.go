package handlers_test

import (
	"bytes"
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	fake_dynamic "k8s.io/client-go/dynamic/fake"

	"connector/pkg/config"
	h "connector/pkg/handlers"
	"connector/pkg/utils"
)

var _ = Describe("Status Replication", func() {
	var (
		remoteClient, localClient dynamic.Interface
		remoteHandler             *h.RemoteHandler
		err                       error
		conf                      *config.Config
		logBuffer                 bytes.Buffer
	)

	BeforeEach(func() {
		log.SetOutput(&logBuffer)
		log.SetLevel(log.DebugLevel)

		localClient = fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(), GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "ready"}),
			GetObject("Root", RootKind, "example"), GetObject("Config", ConfigKind, "example"), GetObject("Project", ProjectKind, "example"))
		conf = &config.Config{
			StatusReplicationEnabled: true,
		}
		remoteClient = fake_dynamic.NewSimpleDynamicClient(runtime.NewScheme(), GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "started"}),
			GetObject("Root", RootKind, "example"), GetObject("Config", ConfigKind, "example"), GetObject("Project", ProjectKind, "example"))
		Expect(err).NotTo(HaveOccurred())

		remoteHandler = h.NewRemoteHandler(utils.GetGVRFromCrdType(Config), Config, localClient, remoteClient, conf)
	})

	It("Should replicate the custom status back to source", func() {
		expectedObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "ready"})
		err := remoteHandler.Create(expectedObj)
		Expect(err).NotTo(HaveOccurred())

		gvr := schema.GroupVersionResource{
			Group:    "config.mazinger.com",
			Version:  "v1",
			Resource: "configs",
		}
		newObj, err := remoteClient.Resource(gvr).Get(context.TODO(), "New", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(newObj.UnstructuredContent()["status"]).To(Equal(expectedObj.UnstructuredContent()["status"]))
	})

	It("Should not add custom resource if not replicated by connector", func() {
		expectedObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "ready"})

		// Remove the nexus-replication-manager annotation.
		expectedObj.SetAnnotations(nil)
		err := remoteHandler.Create(expectedObj)
		Expect(err).NotTo(HaveOccurred())

		Expect(logBuffer.String()).To(ContainSubstring("not replicated by nexus connector, skipping"))
	})

	It("Should fail if nexus-replication-resource annotation not found", func() {
		expectedObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "ready"})

		// Remove the nexus-replication-resource annotation.
		expectedObj.SetAnnotations(map[string]string{utils.NexusReplicationManager: "connector"})
		err := remoteHandler.Create(expectedObj)
		Expect(err.Error()).To(ContainSubstring("CR annotation doesn't contain `NexusReplicationResource[GVR]`"))
	})

	It("Should fail if invalid annotation found", func() {
		expectedObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "ready"})

		// Add invalid annotation.
		expectedObj.SetAnnotations(map[string]string{utils.NexusReplicationManager: "connector",
			utils.NexusReplicationResource: `{"GVR":{"Group":"config.mazinger.com","Version":"v1","Resource":"configs"},"Name":INVALID}`})
		err := remoteHandler.Create(expectedObj)
		Expect(err.Error()).To(ContainSubstring("error unmarshalling resource info from CR annotation"))
	})

	It("Should skip patch if status is same", func() {
		oldObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "started"})
		expectedObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "started"})

		err := remoteHandler.Update(expectedObj, oldObj)
		Expect(err).NotTo(HaveOccurred())

		Expect(logBuffer.String()).To(ContainSubstring("No status changes map[state:started] found for CR"))
	})

	It("Should update patch if status is different", func() {
		oldObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "initializing"})
		expectedObj := GetReplicatedObject("New", ConfigKind, map[string]interface{}{"state": "started"})

		err := remoteHandler.Update(expectedObj, oldObj)
		Expect(err).NotTo(HaveOccurred())

		gvr := schema.GroupVersionResource{
			Group:    "config.mazinger.com",
			Version:  "v1",
			Resource: "configs",
		}
		newObj, err := remoteClient.Resource(gvr).Get(context.TODO(), "New", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(newObj.UnstructuredContent()["status"]).To(Equal(expectedObj.UnstructuredContent()["status"]))
	})
})
