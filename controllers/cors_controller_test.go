package controllers

import (
	"api-gw/pkg/model"
	"context"
	"encoding/json"

	yamlv1 "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	domain_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/domain.nexus.vmware.com/v1"
)

var _ = Describe("OidcConfig controller", func() {
	It("should process oidc config", func() {
		crdJson, err := yamlv1.YAMLToJSON([]byte(corsConfigExample))
		Expect(err).NotTo(HaveOccurred())

		var obj domain_nexus_org.CORSConfig
		err = json.Unmarshal(crdJson, &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())

		event := <-model.CorsChan
		Expect(event.Type).To(Equal(model.Upsert))
		Expect(event.Cors.Name).To(Equal("default"))

		err = k8sClient.Delete(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())

		event = <-model.CorsChan
		Expect(event.Type).To(Equal(model.Delete))
	})
})
