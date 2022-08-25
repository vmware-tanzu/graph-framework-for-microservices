package controllers_test

import (
	"bytes"
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	handler "gitlab.eng.vmware.com/nsx-allspark_users/m7/handler.git"

	log "github.com/sirupsen/logrus"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"connector/controllers"
	"connector/pkg/utils"
)

const (
	apicollaborationspace = "apicollaborationspaces.config.mazinger.com"
	nexuscrd              = "nexuses.api.nexus.org"
)

var (
	ctx        context.Context
	nClient    client.Client
	reconciler *controllers.CustomResourceDefinitionReconciler
	crd        *apiextensionsv1.CustomResourceDefinition
	logBuffer  bytes.Buffer
)

var _ = Describe("Nexus Connector Reconciler Tests", func() {
	When("Crd Type is created", func() {
		BeforeEach(func() {
			log.SetOutput(&logBuffer)
			log.SetLevel(log.DebugLevel)

			crd = getCRD()
			nClient = getHelperClient(ctx, crd)
			reconciler = &controllers.CustomResourceDefinitionReconciler{
				Client: nClient,
				Scheme: runtime.NewScheme(),
				Cache:  controllers.NewCrdCache(),
			}
		})

		It("should construct children and parent map by processing the annotation.", func() {
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(utils.CRDTypeToChildren).Should(HaveKey(apicollaborationspace))
			Expect(utils.CRDTypeToParentHierarchy).Should(HaveKey(apicollaborationspace))
			Expect(reconciler.Cache.CrdMap).To(HaveLen(1))
			Expect(controllers.CrdCh).To(HaveLen(1))
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
			crd.Name = nexuscrd
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(utils.CRDTypeToChildren).Should(HaveLen(0))
			Expect(utils.CRDTypeToParentHierarchy).Should(HaveLen(0))
			Expect(reconciler.Cache.CrdMap).To(HaveLen(0))
		})

		It("should return without error if annotation not present.", func() {
			crd.SetAnnotations(nil)
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(utils.CRDTypeToChildren).Should(HaveLen(0))
			Expect(utils.CRDTypeToParentHierarchy).Should(HaveLen(0))
			Expect(reconciler.Cache.CrdMap).To(HaveLen(1))
		})

		AfterEach(func() {
			delete(utils.CRDTypeToChildren, apicollaborationspace)
			delete(utils.CRDTypeToParentHierarchy, apicollaborationspace)
		})
	})

	When("Crd is deleted", func() {
		BeforeEach(func() {
			nClient = getHelperClientForDeleteEvent(ctx)
			reconciler = &controllers.CustomResourceDefinitionReconciler{
				Client: nClient,
				Scheme: runtime.NewScheme(),
				Cache:  controllers.NewCrdCache(),
			}
		})

		It("should process annotation and populate parent and children map", func() {
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(reconciler.Cache.CrdMap).To(HaveLen(0))
		})
	})

	Context("Cache tests", func() {
		BeforeEach(func() {
			crd = getCRD()
			nClient = getHelperClient(ctx, crd)
			reconciler = &controllers.CustomResourceDefinitionReconciler{
				Client: nClient,
				Scheme: runtime.NewScheme(),
				Cache:  controllers.NewCrdCache(),
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
			reconciler.Cache.UpsertController(crd.Name, c)
			Expect(reconciler.Cache.CrdMap).To(HaveLen(1))
			Expect(reconciler.Cache.CrdMap[apicollaborationspace].Controller).To(Not(BeNil()))
		})

		It("Get", func() {
			_, err := reconciler.Reconcile(ctx, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name: crd.Name,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			controller := &handler.Controller{}
			reconciler.Cache.UpsertController(crd.Name, controller)

			c := reconciler.Cache.Get(crd.Name)
			Expect(c.Spec).To(Not(BeNil()))
			Expect(c.Controller).To(Not(BeNil()))
		})
	})
})
