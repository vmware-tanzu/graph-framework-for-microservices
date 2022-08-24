package utils_test

import (
	"api-gw/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common tests", func() {

	It("should get correct datamodel name from crd", func() {
		datamodelName := utils.GetDatamodelName("datamodels.nexus.org")
		Expect(datamodelName).To(Equal("nexus.org"))
	})
})
