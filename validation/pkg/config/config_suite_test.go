package config_test

import (
	"bytes"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/config"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = Describe("Config", func() {
	var logBuffer bytes.Buffer

	It("Should read config based on config-file existence.", func() {
		log.SetOutput(&logBuffer)
		log.SetLevel(log.TraceLevel)

		// When config-file is found, read from config-file.
		_, err := config.GetConfig("./config_test.yaml")
		Expect(logBuffer.String()).To(ContainSubstring("Setting context to K8s local-api-server"))

		Expect(err).To(HaveOccurred())

		// When config-file not found, read default K8s config.
		_, err = config.GetConfig("./incorrect_config_test.yaml")
		Expect(logBuffer.String()).To(ContainSubstring("Setting context to K8s base-api-server"))

		Expect(err).To(HaveOccurred())
	})
})
