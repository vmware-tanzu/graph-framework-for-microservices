package test_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"

	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"authz-controller/controllers"
	"authz-controller/pkg/utils"
)

var _ = Describe("CustomresourcedefinitionController", func() {
	var (
		fakeClient          client.WithWatch
		r                   controllers.CustomResourceDefinitionReconciler
		reconcileReq        reconcile.Request
		objectName, crdType string
	)

	BeforeEach(func() {
		crdType = "apicollaborationspaces.config.mazinger.com"
		nexusAnnotations := apiCollaborationNexusAnnotation()

		crd := exampleCRD(crdType, nexusAnnotations)
		scheme.AddKnownTypes(apiextensionsv1.SchemeGroupVersion, crd)
		// Register the object in the fake client.
		objs := []runtime.Object{crd}
		objectName = crd.Name
		fakeClient = fakeclient.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(objs...).Build()
		//Create the context for the controller
		r = controllers.CustomResourceDefinitionReconciler{
			Client: fakeClient,
			Scheme: scheme,
		}

		reconcileReq = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: objectName,
			},
		}
	})

	AfterEach(func() {
		clearNexusAPI()
	})

	Context("CustomResourceDefinition", func() {
		It("should create successfully", func() {
			// Process the request
			_, err := r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			children := utils.GetChildrenByCRDType(crdType)
			m := gstruct.MatchAllKeys(gstruct.Keys{
				"apidevspaces.config.mazinger.com": Equal(utils.NodeHelperChild{
					FieldName:    "ApiDevSpaces",
					FieldNameGvk: "apiDevSpacesGvk",
					IsNamed:      true,
				}),
				"apilos.config.mazinger.com": Equal(utils.NodeHelperChild{
					FieldName:    "Apilos",
					FieldNameGvk: "apilosGvk",
					IsNamed:      true,
				}),
			})
			Expect(children).Should(m, "should match all keys")
		})

		It("should delete successfully", func() {
			key := types.NamespacedName{Name: objectName}
			By("Expecting to delete CRD successfully")
			Eventually(func() error {
				f := &apiextensionsv1.CustomResourceDefinition{}
				fakeClient.Get(ctx, key, f)
				return fakeClient.Delete(ctx, f)
			}).Should(Succeed())

			res, err := r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))

			f := &apiextensionsv1.CustomResourceDefinition{}
			err = fakeClient.Get(ctx, key, f)
			Expect(err).To(HaveOccurred())
		})

		It("should update the crd type to clusterrole successfully", func() {
			apidevCRD := exampleCRD("apidevspaces.config.mazinger.com", apidevSpaceNexusAnnotations())
			fakeClient := fakeclient.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects([]runtime.Object{apidevCRD}...).Build()
			r = controllers.CustomResourceDefinitionReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			clusterRole := &rbacv1.ClusterRole{
				ObjectMeta: v1.ObjectMeta{
					Name: "apicbRole",
				},
				Rules: []rbacv1.PolicyRule{{
					Verbs:     []string{"get", "put"},
					APIGroups: []string{"config.mazinger.com"},
					Resources: []string{"apicollaborationspaces", "apilos"},
				}},
			}
			err := fakeClient.Create(ctx, clusterRole)
			Expect(clusterRole.Rules[0].Resources).Should(HaveLen(2))

			utils.SetParentCRDTypeToChildren("apicbRole", crdType, []string{
				"apilos.config.mazinger.com"})

			// Process the request
			reconcileReq.Name = apidevCRD.Name
			_, err = r.Reconcile(ctx, reconcileReq)
			Expect(err).NotTo(HaveOccurred())

			err = fakeClient.Get(ctx, types.NamespacedName{Name: clusterRole.Name}, clusterRole)
			Expect(err).NotTo(HaveOccurred())
			// New CRD Type added to cluster role resources when parent set to hierarchical
			Expect(clusterRole.Rules[0].Resources).Should(HaveLen(3))
		})
	})
})

func exampleCRD(crdType, nexusAnnotations string) *apiextensionsv1.CustomResourceDefinition {
	return &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: crdType,
			Annotations: map[string]string{
				"nexus": nexusAnnotations,
			},
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group:                 "config.mazinger.com",
			Names:                 apiextensionsv1.CustomResourceDefinitionNames{},
			Scope:                 "",
			Versions:              nil,
			Conversion:            nil,
			PreserveUnknownFields: false,
		},
	}
}

func apiCollaborationNexusAnnotation() string {
	return `{
"name":"ApiCollaborationSpace.config","hierarchy":["roots.apix.mazinger.com","projects.apix.mazinger.com","configs.apix.mazinger.com"],
"children":{"apidevspaces.config.mazinger.com":{"fieldName":"ApiDevSpaces","fieldNameGvk":"apiDevSpacesGvk","isNamed":true},
"apilos.config.mazinger.com":{"fieldName":"Apilos","fieldNameGvk":"apilosGvk","isNamed":true}}}`
}

func apidevSpaceNexusAnnotations() string {
	return `{"name":"ApiDevSpace.config",
"hierarchy":["roots.apix.mazinger.com","projects.apix.mazinger.com","configs.apix.mazinger.com","apicollaborationspaces.config.mazinger.com"]}`
}
