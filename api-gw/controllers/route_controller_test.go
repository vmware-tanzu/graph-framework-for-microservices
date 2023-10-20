package controllers

import (
	"context"
	"encoding/json"

	yamlv1 "github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/api.nexus.vmware.com/v1"
	configv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/config.nexus.vmware.com/v1"
	routev1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/route.nexus.vmware.com/v1"
)

var _ = Describe("Route controller", func() {
	It("should create nexus obj", func() {
		crdJson, err := yamlv1.YAMLToJSON([]byte(nexusExample))
		Expect(err).NotTo(HaveOccurred())

		var obj apiv1.Nexus
		err = json.Unmarshal(crdJson, &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())
	})
	It("should create config obj", func() {
		crdJson, err := yamlv1.YAMLToJSON([]byte(nexusExample))
		Expect(err).NotTo(HaveOccurred())

		var obj configv1.Config
		err = json.Unmarshal(crdJson, &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())
	})
	It("should create a route", func() {
		crdJson, err := yamlv1.YAMLToJSON([]byte(routeExample))
		Expect(err).NotTo(HaveOccurred())

		var obj routev1.Route
		err = json.Unmarshal(crdJson, &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Create(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.Delete(context.TODO(), &obj)
		Expect(err).ToNot(HaveOccurred())
	})
})
