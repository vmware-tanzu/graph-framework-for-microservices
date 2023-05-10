package controllers_test

import (
	"bytes"
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	handler "gitlab.eng.vmware.com/nsx-allspark_users/m7/handler.git"

	"connector/controllers"
	"connector/pkg/utils"
)

const nexuscrd = "nexuses.api.nexus.vmware.com"

var (
	ctx        context.Context
	nClient    client.Client
	reconciler *controllers.CustomResourceDefinitionReconciler
	crd        *apiextensionsv1.CustomResourceDefinition
	gvr        schema.GroupVersionResource
	logBuffer  bytes.Buffer
)

var _ = Describe("Nexus Connector Reconciler Tests", func() {
	When("Crd Type is created", func() {
		BeforeEach(func() {
			log.SetOutput(&logBuffer)
			log.SetLevel(log.DebugLevel)

			crd = getCRD()
			gvr = schema.GroupVersionResource{Group: "config.mazinger.com", Version: "v1", Resource: "apicollaborationspaces"}
			nClient = getHelperClient(ctx, crd)
			reconciler = &controllers.CustomResourceDefinitionReconciler{
				Client: nClient,
				Scheme: runtime.NewScheme(),
				Cache:  controllers.NewGvrCache(),
			}
		})

		It("should construct children and parent map by processing the annotation.", func() {
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(utils.GVRToChildren).Should(HaveKey(gvr))
			Expect(utils.GVRToParentHierarchy).Should(HaveKey(gvr))
			Expect(controllers.GvrCh).To(HaveLen(1))
		})

		It("should fail if nexus annotation not present in the proper format.", func() {
			crd.SetAnnotations(map[string]string{"nexus": `{"name":"ApiCollaborationSpace.config","hierarchy":["roots.apix.mazinger.com","projects.apix.mazinger.com","configs.apix.mazinger.com"],"children":{"apidevspaces.config.mazinger.com":{"fieldName":"ApiDevSpaces","fieldNameGvk":INVALID,"isNamed":true}}}`})
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(logBuffer.String()).To(ContainSubstring("Error unmarshaling Nexus annotation"))
		})

		It("should not process nexus datamodel crd.", func() {
			newCrd := getCRD()
			newCrd.Name = nexuscrd
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: newCrd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(utils.GVRToChildren).Should(HaveLen(1))
			Expect(utils.GVRToParentHierarchy).Should(HaveLen(1))
		})

		It("should return without error if annotation not present.", func() {
			crd.SetAnnotations(nil)
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(utils.GVRToChildren).Should(HaveLen(0))
			Expect(utils.GVRToParentHierarchy).Should(HaveLen(0))
		})

		AfterEach(func() {
			delete(utils.GVRToChildren, gvr)
			delete(utils.GVRToParentHierarchy, gvr)
		})
	})

	When("Crd is deleted", func() {
		BeforeEach(func() {
			nClient = getHelperClientForDeleteEvent(ctx)
			reconciler = &controllers.CustomResourceDefinitionReconciler{
				Client: nClient,
				Scheme: runtime.NewScheme(),
				Cache:  controllers.NewGvrCache(),
			}
		})

		It("should process annotation and populate parent and children map", func() {
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Cache tests", func() {
		BeforeEach(func() {
			crd = getCRD()
			nClient = getHelperClient(ctx, crd)
			reconciler = &controllers.CustomResourceDefinitionReconciler{
				Client: nClient,
				Scheme: runtime.NewScheme(),
				Cache:  controllers.NewGvrCache(),
			}
		})

		It("Upsert Controller", func() {
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			c := &handler.Controller{}
			reconciler.Cache.UpsertController(gvr, c)
			Expect(reconciler.Cache.GvrMap[gvr].Controller).To(Not(BeNil()))
		})

		It("Get", func() {
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			controller := &handler.Controller{}
			reconciler.Cache.UpsertController(gvr, controller)

			c := reconciler.Cache.Get(gvr)
			Expect(c.Controller).To(Not(BeNil()))
		})
	})
})
