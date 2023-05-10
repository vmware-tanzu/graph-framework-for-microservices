package config_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"connector/pkg/config"
	"connector/pkg/utils"
)

var _ = Describe("Config tests", func() {
	BeforeEach(func() {
		err := os.Setenv(utils.RemoteEndpointHost, "http://a1eb8ab4a5a2d4a0b9c898200d636cbd-1190706442.us-east-2.elb.amazonaws.com")
		Expect(err).NotTo(HaveOccurred())

		err = os.Setenv(utils.RemoteEndpointPort, "80")
		Expect(err).NotTo(HaveOccurred())
	})

	It("Should read config successfully", func() {
		c, err := config.LoadConfig("./test_utils/correct.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(c).NotTo(BeNil())
		Expect(c.Dispatcher.WorkerTTL).To(Equal(time.Duration(time.Second * 15)))
		Expect(c.Dispatcher.MaxWorkerCount).To(Equal(uint(100)))
		Expect(c.Dispatcher.CloseRequestsQueueSize).To(Equal(uint(15)))
		Expect(c.Dispatcher.EventProcessedQueueSize).To(Equal(uint(100)))
		Expect(c.RemoteEndpointHost).To(Equal("http://a1eb8ab4a5a2d4a0b9c898200d636cbd-1190706442.us-east-2.elb.amazonaws.com"))
		Expect(c.RemoteEndpointPort).To(Equal("80"))
	})

	It("Should return opening error", func() {
		c, err := config.LoadConfig("./not_existing.yaml")
		Expect(err).To(
			MatchError(ContainSubstring("failed to open config file: ")))
		Expect(c).To(BeNil())
	})

	It("Should return read error", func() {
		c, err := config.LoadConfig("test_utils/")
		Expect(err).To(
			MatchError(ContainSubstring("failed to read config file: ")))
		Expect(c).To(BeNil())
	})

	It("Should return parse error", func() {
		c, err := config.LoadConfig("test_utils/incorrect.yaml")
		Expect(err).To(
			MatchError(ContainSubstring("failed to parse config file: ")))
		Expect(c).To(BeNil())
	})
})
