package controllers

import (
	"api-gw/pkg/model"
	"context"
	"encoding/json"

	yamlv1 "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	authenticationnexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/authentication.nexus.vmware.com/v1"
)

var _ = Describe("OidcConfig controller", func() {
	It("should process oidc config", func() {
		crdJson, err := yamlv1.YAMLToJSON([]byte(oidcCrdObjectExample))
		Expect(err).NotTo(HaveOccurred())

		var obj authenticationnexusv1.OIDC
		err = json.Unmarshal(crdJson, &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())

		event := <-model.OidcChan
		Expect(event.Type).To(Equal(model.Upsert))
		Expect(event.Oidc.Name).To(Equal("okta"))
	})
})
