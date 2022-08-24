package controllers

import (
	"api-gw/pkg/model"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Datamodel controller", func() {
	It("should create datamodel crd", func() {
		gvr := schema.GroupVersionResource{
			Group:    "nexus.org",
			Version:  "v1",
			Resource: "datamodels",
		}

		unstructuredObject := unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "nexus.org/v1",
				"kind":       "Datamodel",
				"metadata": map[string]interface{}{
					"name": "nexus.org",
				},
				"spec": map[string]interface{}{
					"name":  "nexus.org",
					"title": "Example title",
				},
			},
		}
		_, err := dynamicClient.Resource(gvr).Create(context.TODO(), &unstructuredObject, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			if _, ok := model.DatamodelToDatamodelInfo["nexus.org"]; ok {
				return true
			}
			return false
		}).Should(BeTrue())
	})
})
