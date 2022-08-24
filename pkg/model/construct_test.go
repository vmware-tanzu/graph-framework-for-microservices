package model_test

import (
	"api-gw/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Construct tests", func() {

	It("should construct new datamodel (vmware-test.org)", func() {
		unstructuredObj := unstructured.Unstructured{
			Object: map[string]interface{}{
				"spec": map[string]interface{}{
					"title": "VMWare Datamodel",
				},
			},
		}

		model.ConstructDatamodel(model.Upsert, "vmware-test.org", &unstructuredObj)
		Expect(model.DatamodelToDatamodelInfo).To(HaveKey("vmware-test.org"))
	})

	It("should delete datamodel vmware-test.org", func() {
		unstructuredObj := unstructured.Unstructured{
			Object: map[string]interface{}{
				"spec": map[string]interface{}{
					"title": "VMWare Datamodel",
				},
			},
		}

		model.ConstructDatamodel(model.Delete, "vmware-test.org", &unstructuredObj)
		Expect(model.DatamodelToDatamodelInfo).ToNot(HaveKey("vmware-test.org"))
	})
})
