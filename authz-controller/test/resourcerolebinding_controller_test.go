package test_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"authz-controller/controllers"
	auth_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authorization.nexus.org/v1"
)

var _ = Describe("ResourcerolebindingController", func() {
	var (
		fakeClient   client.WithWatch
		r            controllers.ResourceRoleBindingReconciler
		reconcileReq reconcile.Request
		objectName   string
	)

	BeforeEach(func() {
		// Create nexus Resource RoleBinding
		authNode = createParentNodes()
		objectToCreate := exampleResourceRoleBinding()

		nexusRoleBinding, err := authNode.AddResourceRoleBinding(ctx, objectToCreate)
		Expect(err).NotTo(HaveOccurred())
		Expect(nexusRoleBinding).NotTo(BeNil())

		nexusRoleBinding.Spec = addRoleBindingSpec()

		objectName = nexusRoleBinding.ResourceRoleBinding.Name
		scheme.AddKnownTypes(auth_nexus_org.SchemeGroupVersion, nexusRoleBinding.ResourceRoleBinding)

		// Register the object in the fake client.
		objs := []runtime.Object{nexusRoleBinding.ResourceRoleBinding}

		fakeClient = fakeclient.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(objs...).Build()
		//Create the context for the controller
		r = controllers.ResourceRoleBindingReconciler{
			Client: fakeClient,
			Scheme: scheme,
		}

		//Mock the event processing
		reconcileReq = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: objectName,
			},
		}
	})

	AfterEach(func() {
		clearNexusAPI()
	})

	Context("Resource RoleBinding", func() {
		It("should create successfully", func() {
			// Process the request
			_, err := r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			err = r.Client.Get(ctx, reconcileReq.NamespacedName, clusterRoleBinding)
			Expect(err).NotTo(HaveOccurred())
			Expect(clusterRoleBinding).ToNot(BeNil())
		})

		It("should update successfully", func() {
			key := types.NamespacedName{
				Name: objectName,
			}
			By("Expecting to update resource rolebinding spec successfully")

			f := &auth_nexus_org.ResourceRoleBinding{}
			err := fakeClient.Get(ctx, key, f)
			Expect(err).NotTo(HaveOccurred())
			Expect(f.Spec.RoleGvk.Name).To(Equal("1234"))

			f.Annotations["some-key"] = "some-value"
			f.Spec.RoleGvk.Name = "doe"

			err = fakeClient.Update(ctx, f)
			Expect(err).NotTo(HaveOccurred())
			Expect(f.Annotations["some-key"]).To(Equal("some-value"))
			Expect(f.Spec.RoleGvk.Name).To(Equal("doe"))

			_, err = r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			err = fakeClient.Get(ctx, key, clusterRoleBinding)
			Expect(err).NotTo(HaveOccurred())
			Expect(clusterRoleBinding.Annotations["some-key"]).Should(Equal("some-value"))
			Expect(f.Spec.RoleGvk.Name).To(Equal("doe"))
		})

		It("should delete successfully", func() {
			key := types.NamespacedName{Name: objectName}
			By("Expecting to delete resource rolebinding successfully")
			Eventually(func() error {
				f := &auth_nexus_org.ResourceRoleBinding{}
				fakeClient.Get(ctx, key, f)
				return fakeClient.Delete(ctx, f)
			}).Should(Succeed())

			_, err := r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() error {
				clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
				return fakeClient.Get(ctx, key, clusterRoleBinding)
			}).ShouldNot(Succeed())
		})
	})
})

func exampleResourceRoleBinding() *auth_nexus_org.ResourceRoleBinding {
	return &auth_nexus_org.ResourceRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				"key1": "value1",
			},
		},
		Spec: auth_nexus_org.ResourceRoleBindingSpec{},
	}
}

func addRoleBindingSpec() auth_nexus_org.ResourceRoleBindingSpec {
	return auth_nexus_org.ResourceRoleBindingSpec{
		RoleGvk: &auth_nexus_org.Link{
			Group: "authorization.test.org",
			Kind:  "ResourceRole",
			Name:  "1234",
		},
		UsersGvk: map[string]auth_nexus_org.Link{
			"bob": {
				Group: "authorization.test.org",
				Kind:  "User",
				Name:  "bob",
			},
		},
	}
}
