package controllers_test

import (
	"context"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	nxcontroller "gitlab.eng.vmware.com/nexus/controller/controllers"
	nexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	sn           *SingletonNodes
	ctx          context.Context
	fakeClient   *nexus_client.Clientset
	k8sClientSet *testclient.Clientset
)

var _ = BeforeSuite(func() {
	ctx = context.Background()
	fakeClient = nexus_client.NewFakeClient()
	cleanupEnv()
	k8sClientSet = testclient.NewSimpleClientset()
	sn = initDatamodel(ctx, fakeClient)
}, 60)

var _ = Describe("Nexus Controller Tests", func() {
	name := "endpoint1"
	var (
		endPoint *nexus_client.ConnectNexusEndpoint
		nClient  client.Client
		err      error
		crd      *nexusv1.NexusEndpoint
		r        *nxcontroller.NexusConnectorReconciler
	)

	BeforeEach(func() {
		initEnvVars()

		crd = &nexusv1.NexusEndpoint{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: nexusv1.NexusEndpointSpec{
				Host: "https://localhost",
				Port: "8080",
				Cert: "NA",
			},
		}

		endPoint, err = sn.Connect.AddEndpoints(ctx, crd)
		Expect(err).NotTo(HaveOccurred())
		Expect(endPoint).NotTo(BeNil())

		nClient = getHelperClient(ctx, endPoint)
		r = &nxcontroller.NexusConnectorReconciler{
			K8sClient: k8sClientSet,
			Client:    nClient,
			Scheme:    newScheme(),
		}
		Expect(err).NotTo(HaveOccurred())
	})

	When("an endpoint is created", func() {
		It("should create nexus connector deployment, service and config map", func() {
			_, err := r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			deploymentList, err := r.K8sClient.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentList.Items).To(HaveLen(1))
			Expect(deploymentList.Items[0].Name).To(Equal("nexus-connector-9876611c09489e8c75cc3691066480420a010434"))

			serviceList, err := r.K8sClient.CoreV1().Services("default").List(ctx, metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(serviceList.Items).To(HaveLen(1))
			Expect(serviceList.Items[0].Name).To(Equal("nexus-connector-9876611c09489e8c75cc3691066480420a010434"))

			configMapList, err := r.K8sClient.CoreV1().ConfigMaps("default").List(ctx, metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(configMapList.Items).To(HaveLen(1))
			Expect(configMapList.Items[0].Name).To(Equal("connector-kubeconfig-local"))

		})

		It("should update nexus connector configurations", func() {
			endPoint.Spec.Port = "8081"
			_, err := r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			deploymentList, err := r.K8sClient.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentList.Items).To(HaveLen(1))
			Expect(deploymentList.Items[0].Name).To(Equal("nexus-connector-9876611c09489e8c75cc3691066480420a010434"))

			env := getEnv(deploymentList)
			Expect(env).To(Equal("8081"))
		})

		It("should fail if endpoint port or host is empty", func() {
			endPoint.Spec.Port = ""
			_, err := r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("endpoint Host/Port is empty"))
		})

		It("should fail if nexus connector version image is empty", func() {
			err := os.Setenv("NEXUS_CONNECTOR_VERSION", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("env var NEXUS_CONNECTOR_VERSION is missing"))
		})

		It("Should fail services update when error occurs in updation", func() {
			ks := testclient.NewSimpleClientset()
			ks.Fake.PrependReactor("update", "services",
				func(action testing.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("nope")
				})

			r := &nxcontroller.NexusConnectorReconciler{
				K8sClient: ks,
				Client:    nClient,
				Scheme:    newScheme(),
			}

			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			endPoint.Spec.Port = "8083"
			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).To(HaveOccurred())
		})

		It("Should fail configmaps update when error occurs in updation", func() {
			ks := testclient.NewSimpleClientset()
			ks.Fake.PrependReactor("update", "configmaps",
				func(action testing.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("nope")
				})

			r := &nxcontroller.NexusConnectorReconciler{
				K8sClient: ks,
				Client:    nClient,
				Scheme:    newScheme(),
			}

			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			endPoint.Spec.Port = "8082"
			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).To(HaveOccurred())
		})

		It("Should fail deployment update when error occurs in updation", func() {
			ks := testclient.NewSimpleClientset()
			ks.Fake.PrependReactor("update", "deployments",
				func(action testing.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("nope")
				})

			r := &nxcontroller.NexusConnectorReconciler{
				K8sClient: ks,
				Client:    nClient,
				Scheme:    newScheme(),
			}

			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			endPoint.Spec.Port = "8083"
			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).To(HaveOccurred())
		})

		AfterEach(func() {
			err = sn.Connect.DeleteEndpoints(ctx, name)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("an endpoint is deleted", func() {
		BeforeEach(func() {
			err := sn.Connect.DeleteEndpoints(ctx, name)
			Expect(err).NotTo(HaveOccurred())

			nClient = getHelperClientForDeleteEvent(ctx)
		})
		It("should delete nexus connector deployment", func() {
			r := &nxcontroller.NexusConnectorReconciler{
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

			deploymentList, err := r.K8sClient.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentList.Items).To(HaveLen(0))
		})

		It("should fail when error occurs in deployment delete", func() {
			k8sClientSet.Fake.PrependReactor("delete", "deployments",
				func(action testing.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("nope")
				})

			r := &nxcontroller.NexusConnectorReconciler{
				K8sClient: k8sClientSet,
				Client:    nClient,
				Scheme:    newScheme(),
			}

			_, err = r.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: endPoint.Name,
				},
			})
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = AfterSuite(func() {
	cleanupEnv()
})
