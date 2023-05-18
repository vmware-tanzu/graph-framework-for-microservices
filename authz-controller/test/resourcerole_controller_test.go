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
	"authz-controller/pkg/utils"
	auth_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authorization.nexus.org/v1"
)

var _ = Describe("ResourceroleController", func() {
	var (
		fakeClient   client.WithWatch
		r            controllers.ResourceRoleReconciler
		reconcileReq reconcile.Request
		objectName   string
	)

	BeforeEach(func() {
		// Create nexus Resource Role
		authNode = createParentNodes()
		objectToCreate := exampleResourceRole()
		nexusRole, err := authNode.AddResourceRole(ctx, objectToCreate)
		Expect(err).NotTo(HaveOccurred())
		Expect(nexusRole).NotTo(BeNil())

		objectName = nexusRole.ResourceRole.Name
		scheme.AddKnownTypes(auth_nexus_org.SchemeGroupVersion, nexusRole.ResourceRole)

		// Register the object in the fake client.
		objs := []runtime.Object{nexusRole.ResourceRole}

		fakeClient = fakeclient.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(objs...).Build()
		//Create the context for the controller
		r = controllers.ResourceRoleReconciler{
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

	Context("Resource Role", func() {
		It("should create successfully", func() {
			// Process the request
			_, err := r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			clusterRole := &rbacv1.ClusterRole{}
			err = r.Client.Get(ctx, reconcileReq.NamespacedName, clusterRole)
			Expect(err).NotTo(HaveOccurred())
			Expect(clusterRole).ToNot(BeNil())
		})

		It("should update successfully", func() {
			key := types.NamespacedName{
				Name: objectName,
			}
			By("Expecting to update resource role rules successfully")

			f := &auth_nexus_org.ResourceRole{}
			err := fakeClient.Get(ctx, key, f)
			Expect(err).NotTo(HaveOccurred())
			Expect(f.Spec.Rules).Should(HaveLen(1))

			f.Annotations["some-key"] = "some-value"
			f.Spec.Rules = append(f.Spec.Rules, auth_nexus_org.ResourceRule{
				Resource: auth_nexus_org.ResourceType{
					Group: "hello.123.com",
					Kind:  "A",
				},
				Verbs:        []auth_nexus_org.Verb{"get", "list"},
				Hierarchical: true,
			})

			err = fakeClient.Update(ctx, f)
			Expect(err).NotTo(HaveOccurred())
			Expect(f.Spec.Rules).Should(HaveLen(2))
			Expect(f.Annotations["some-key"]).To(Equal("some-value"))

			_, err = r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			clusterRole := &rbacv1.ClusterRole{}
			err = fakeClient.Get(ctx, key, clusterRole)
			Expect(err).NotTo(HaveOccurred())
			Expect(clusterRole.Annotations["some-key"]).Should(Equal("some-value"))
			Expect(clusterRole.Rules).Should(HaveLen(2))
		})

		It("should delete successfully", func() {
			key := types.NamespacedName{Name: objectName}
			By("Expecting to delete resource role successfully")
			Eventually(func() error {
				f := &auth_nexus_org.ResourceRole{}
				fakeClient.Get(ctx, key, f)
				return fakeClient.Delete(ctx, f)
			}).Should(Succeed())

			_, err := r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() error {
				clusterRole := &rbacv1.ClusterRole{}
				return fakeClient.Get(ctx, key, clusterRole)
			}).ShouldNot(Succeed())
		})

		It("should rule contain all child hierarchical node successfully", func() {
			key := types.NamespacedName{Name: objectName}

			f := &auth_nexus_org.ResourceRole{}
			err := fakeClient.Get(ctx, key, f)
			Expect(err).NotTo(HaveOccurred())

			f.Spec.Rules = append(f.Spec.Rules, auth_nexus_org.ResourceRule{
				Resource: auth_nexus_org.ResourceType{
					Group: "config.mazinger.com",
					Kind:  "ApiCollaborationSpace",
				},
				Verbs:        []auth_nexus_org.Verb{"get", "list"},
				Hierarchical: true,
			})

			err = fakeClient.Update(ctx, f)
			Expect(err).NotTo(HaveOccurred())
			Expect(f.Spec.Rules).Should(HaveLen(2))

			c := map[string]utils.NodeHelperChild{
				"apidevspaces.config.mazinger.com": {
					FieldName:    "ApiDevSpaces",
					FieldNameGvk: "apiDevSpacesGvk",
					IsNamed:      true,
				},
			}
			utils.ConstructMapCRDTypeToChildren(utils.Upsert, "apicollaborationspaces.config.mazinger.com", c)

			_, err = r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			clusterRole := &rbacv1.ClusterRole{}
			err = fakeClient.Get(ctx, key, clusterRole)
			Expect(err).NotTo(HaveOccurred())

			for _, cr := range clusterRole.Rules {
				if utils.ContainsString(cr.Resources, "apicollaborationspaces") {
					Expect(cr.Resources).Should(ConsistOf([]string{"apicollaborationspaces", "apidevspaces"}))
				} else {
					Expect(cr.Resources).Should(ConsistOf([]string{"roots"}))
				}
			}
		})
	})
})

func exampleResourceRole() *auth_nexus_org.ResourceRole {
	return &auth_nexus_org.ResourceRole{
		ObjectMeta: v1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				"key1": "value1",
			},
		},
		Spec: auth_nexus_org.ResourceRoleSpec{
			Rules: []auth_nexus_org.ResourceRule{
				{
					Resource: auth_nexus_org.ResourceType{
						Group: "apix.mazinger.com",
						Kind:  "Root",
					},
					Verbs:        []auth_nexus_org.Verb{"put", "get"},
					Hierarchical: false,
				},
			},
		},
	}
}
