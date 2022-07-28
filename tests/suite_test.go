//go:build unit
// +build unit

package tests

import (
	"context"
	"testing"

	nxcontroller "gitlab.eng.vmware.com/nexus/controller/controllers"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	nexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

// TODO
// var cfg *rest.Config
// var k8sClient client.Client

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"Nexus Controller Tests",
		[]Reporter{printer.NewlineReporter{}})
}

var (
	sn           *SingletonNodes
	ctx          context.Context
	fakeClient   *nexus_client.Clientset
	k8sClientSet kubernetes.Interface
)

var _ = BeforeSuite(func() {
	ctx = context.Background()
	fakeClient = nexus_client.NewFakeClient()
	cleanupEnv()
	k8sClientSet = testclient.NewSimpleClientset()
	sn = initDatamodel(ctx, fakeClient)
	intiEnvVars()
}, 60)

var _ = Describe("Nexus Controller Tests", func() {
	name := "endpoint1"
	var endPoint *nexus_client.ConnectNexusEndpoint
	When("an endpoint is created", func() {
		var nClient client.Client
		var err error
		var r *nxcontroller.NexusConnectorReconciler
		BeforeEach(func() {
			endPoint, err = sn.Connect.AddEndpoints(ctx, &nexusv1.NexusEndpoint{
				ObjectMeta: v1.ObjectMeta{
					Name: name,
				},
				Spec: nexusv1.NexusEndpointSpec{
					Host: "https://localhost",
					Port: "8080",
					Cert: "NA",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(endPoint).NotTo(BeNil())
			nClient = getHelperClient(ctx, endPoint)
			r = &nxcontroller.NexusConnectorReconciler{
				K8sClient: k8sClientSet,
				Client:    nClient,
				Scheme:    newScheme(),
			}
			_, err := r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create nexus connector deployment, service and config map", func() {
			deploymentList, err := r.K8sClient.AppsV1().Deployments("default").List(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentList.Items).To(HaveLen(1))
			Expect(deploymentList.Items[0].Name).To(Equal("nexus-connector-9876611c09489e8c75cc3691066480420a010434"))

			serviceList, err := r.K8sClient.CoreV1().Services("default").List(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(serviceList.Items).To(HaveLen(1))
			Expect(serviceList.Items[0].Name).To(Equal("nexus-connector-9876611c09489e8c75cc3691066480420a010434"))

			configMapList, err := r.K8sClient.CoreV1().ConfigMaps("default").List(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(configMapList.Items).To(HaveLen(1))
			Expect(configMapList.Items[0].Name).To(Equal("connector-kubeconfig-local"))

		})
	})

	When("an endpoint is deleted", func() {
		var nClient client.Client
		// var endPoint *nexus_client.ConnectNexusEndpoint
		var err error
		var r *nxcontroller.NexusConnectorReconciler
		BeforeEach(func() {
			err = sn.Connect.DeleteEndpoints(ctx, name)
			Expect(err).NotTo(HaveOccurred())

			nClient = getHelperClientForDeleteEvent(ctx)
			r = &nxcontroller.NexusConnectorReconciler{
				K8sClient: k8sClientSet,
				Client:    nClient,
				Scheme:    newScheme(),
			}
			_, err := r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
		It("should delete nexus connector deployment", func() {
			deploymentList, err := r.K8sClient.AppsV1().Deployments("default").List(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentList.Items).To(HaveLen(0))
		})
	})
})

var _ = AfterSuite(func() {
	cleanupEnv()
})
