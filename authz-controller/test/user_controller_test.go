package test_test

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"authz-controller/controllers"
	auth_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authorization.nexus.org/v1"
)

var _ = Describe("UserController", func() {
	var (
		fakeClient   client.WithWatch
		r            controllers.UserReconciler
		reconcileReq reconcile.Request
		objectName   string
	)

	BeforeEach(func() {
		// Create nexus User
		authNode = createParentNodes()
		objectToCreate := exampleNexusUser()
		nexusUser, err := authNode.AddUsers(ctx, objectToCreate)
		Expect(err).NotTo(HaveOccurred())
		Expect(nexusUser).NotTo(BeNil())

		objectName = nexusUser.User.Name
		scheme.AddKnownTypes(auth_nexus_org.SchemeGroupVersion, nexusUser.User)

		// Register the object in the fake client.
		objs := []runtime.Object{nexusUser.User}

		fakeClient = fakeclient.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(objs...).Build()
		//Create the context for the controller
		r = controllers.UserReconciler{
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

	Context("Nexus User", func() {
		It("should create successfully", func() {
			// should throw error due to undefined cert path
			_, err := r.Reconcile(ctx, reconcileReq)
			Expect(err).To(HaveOccurred())

			By("Expecting to create user signed certificate successfully")
			certPrivateKey, cert, err := r.CreateUserCert(objectName, "../test/sample_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(certPrivateKey).NotTo(BeNil())
			Expect(cert).NotTo(BeNil())

			By("Expecting to store private cert in user certificate CR successfully")
			err = r.CreateUserCertificate(ctx, certPrivateKey.String(), cert.String(), objectName, *exampleNexusUser())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete successfully", func() {
			key := types.NamespacedName{Name: objectName}
			By("Expecting to delete nexus user successfully")
			Eventually(func() error {
				f := &auth_nexus_org.User{}
				fakeClient.Get(ctx, key, f)
				return fakeClient.Delete(ctx, f)
			}).Should(Succeed())

			userCertificateCRDName := "usercertificates.authorization.nexus.org"
			certObjName := fmt.Sprintf("%s:%s", userCertificateCRDName, objectName)
			h := sha1.New()
			h.Write([]byte(certObjName))

			userCertificateName := hex.EncodeToString(h.Sum(nil))

			// nexus user deletion will delete the user certificate too
			uCert := &auth_nexus_org.UserCertificate{}
			err := fakeClient.Get(ctx, types.NamespacedName{Name: userCertificateName}, uCert)
			Expect(err).To(HaveOccurred())
			Expect(err).Should(MatchError("usercertificates.authorization.nexus.org \"21477d3af6e858b7717eb6578951311835653ee2\" not found"))
		})
	})
})

func exampleNexusUser() *auth_nexus_org.User {
	return &auth_nexus_org.User{
		ObjectMeta: v1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				"key1": "value1",
			},
		},
	}
}
