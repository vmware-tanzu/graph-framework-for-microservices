package utils_test

import (
	"api-gw/pkg/client"
	"api-gw/pkg/config"
	"api-gw/pkg/utils"
	"os"

	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common tests", func() {

	It("should get correct datamodel name from crd", func() {
		datamodelName := utils.GetDatamodelName("route.route.nexus.vmware.com")
		Expect(datamodelName).To(Equal("nexus.vmware.com"))
	})

	It("should check if file exist", func() {
		file, err := os.Create("test-file.txt")
		Expect(err).ToNot(HaveOccurred())

		check := utils.IsFileExists(file.Name())
		Expect(check).To(BeTrue())

		err = os.Remove("test-file.txt")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should check if file not exist", func() {
		check := utils.IsFileExists("non-existent-file")
		Expect(check).To(BeFalse())
	})

	It("should check if server config is valid", func() {
		isValid := utils.IsServerConfigValid(&config.Config{
			Server: config.ServerConfig{
				Address:  "address",
				CertPath: "cert_path",
				KeyPath:  "key_path",
			},
		})
		Expect(isValid).To(BeTrue())
	})

	It("should check if server config is not valid", func() {
		isValid := utils.IsServerConfigValid(&config.Config{})
		Expect(isValid).To(BeFalse())
	})

	It("should get crd type", func() {
		crdType := utils.GetCrdType("Test", "vmware.org")
		Expect(crdType).To(Equal("tests.vmware.org"))
	})

	It("should get resource name", func() {
		resource := utils.GetGroupResourceName("Test")
		Expect(resource).To(Equal("tests"))
	})

	It("should GetEnvoyInitParams without error", func() {
		client.NexusClient = nexus_client.NewFakeClient()
		//TODO: more detailed test with mock server
		//client.NexusClient.Authentication().CreateOIDCByName(context.TODO(), &v1.OIDC{
		//	ObjectMeta: metav1.ObjectMeta{
		//		Name: "okta",
		//	},
		//	Spec: v1.OIDCSpec{
		//		Config: v1.IDPConfig{
		//			ClientId:       "XXX",
		//			ClientSecret:   "XXX",
		//			OAuthIssuerUrl: "https://dev-XXX.okta.com/oauth2/default",
		//			Scopes: []string{
		//				"openid",
		//				"profile",
		//				"offline_access",
		//			},
		//			OAuthRedirectUrl: "http://<API-GW-DNS/IP>:<PORT>/<CALLBACK_PATH>",
		//		},
		//		ValidationProps: v1.ValidationProperties{
		//			InsecureIssuerURLContext: false,
		//			SkipIssuerValidation:     false,
		//			SkipClientIdValidation:   false,
		//			SkipClientAudValidation:  false,
		//		},
		//	},
		//})
		_, _, _, err := utils.GetEnvoyInitParams()
		Expect(err).ToNot(HaveOccurred())
	})
})
