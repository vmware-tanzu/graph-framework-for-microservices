package controllers

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Custom resource definition controller", func() {
	//FIXME: flaky test, crd is not available while executing this test
	//It("should create and process orgchart.Root crd", func() {
	//	Eventually(func() bool {
	//		res := &apiextensionsv1.CustomResourceDefinition{}
	//		err := k8sClient.Get(context.TODO(), client.ObjectKey{
	//			Name: "roots.orgchart.vmware.org",
	//		}, res)
	//		Expect(err).ToNot(HaveOccurred())
	//		if res != nil {
	//			return true
	//		}
	//
	//		return false
	//	}).Should(BeTrue())
	//
	//	Eventually(func() bool {
	//		if val, ok := model.UriToUriInfo["/root/{orgchart.Root}"]; ok {
	//			if val.TypeOfURI == model.DefaultURI {
	//				return true
	//			}
	//		}
	//		return false
	//	}).Should(BeTrue())
	//
	//	Eventually(func() bool {
	//		if val, ok := model.UriToCRDType["/root/{orgchart.Root}"]; ok {
	//			if val == "roots.orgchart.vmware.org" {
	//				return true
	//			}
	//		}
	//		return false
	//	}).Should(BeTrue())
	//
	//	Eventually(func() bool {
	//		if val, ok := model.CrdTypeToNodeInfo["roots.orgchart.vmware.org"]; ok {
	//			if val.Name == "orgchart.Root" {
	//				return true
	//			}
	//		}
	//		return false
	//	}).Should(BeTrue())
	//
	//	Eventually(func() bool {
	//		if val, ok := model.CrdTypeToRestUris["roots.orgchart.vmware.org"]; ok {
	//			if val[0].Uri == "/root/{orgchart.Root}" {
	//				return true
	//			}
	//		}
	//		return false
	//	}).Should(BeTrue())
	//})
})
