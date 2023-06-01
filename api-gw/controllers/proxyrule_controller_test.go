package controllers

import (
	"api-gw/pkg/envoy"
	"context"
	"encoding/json"

	yamlv1 "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	adminv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/admin.nexus.vmware.com/v1"
)

var _ = Describe("ProxyRule controller", func() {

	It("should init envoy mulitple times without issues", func() {
		logLevel, _ := log.ParseLevel("debug")
		err := envoy.Init(nil, nil, nil, logLevel)
		Expect(err).NotTo(HaveOccurred())

		// calling the envoy Init multiple times , as the XDS listener should stop and restart each time
		logLevel, _ = log.ParseLevel("debug")
		err = envoy.Init(nil, nil, nil, logLevel)
		Expect(err).NotTo(HaveOccurred())

	})

	It("should create header based proxyrule", func() {
		crdJson, err := yamlv1.YAMLToJSON([]byte(proxyRuleHeaderExample))
		Expect(err).NotTo(HaveOccurred())

		var obj adminv1.ProxyRule
		err = json.Unmarshal(crdJson, &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Delete(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should create jwt based proxyrule", func() {
		crdJson, err := yamlv1.YAMLToJSON([]byte(proxyRuleJwtExample))
		Expect(err).NotTo(HaveOccurred())

		var obj adminv1.ProxyRule
		err = json.Unmarshal(crdJson, &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Delete(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())
	})
})
