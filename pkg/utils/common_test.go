package utils_test

import (
	"api-gw/pkg/config"
	"api-gw/pkg/utils"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common tests", func() {

	It("should get correct datamodel name from crd", func() {
		datamodelName := utils.GetDatamodelName("datamodels.nexus.org")
		Expect(datamodelName).To(Equal("nexus.org"))
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
})
