package dir

import (
	"os"

	nexus_compare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
)

var _ = Describe("Generator", func() {

	ff := fake.NewSimpleClientset()
	c := kubewrapper.Client{Clientset: ff}
	c.Clientset.ApiextensionsV1()

	It("apply crds to empty cluster", func() {
		err := ApplyDir("./test_dir", false, &c, getNoDiff)
		Expect(err).To(Not(HaveOccurred()))
		err = c.FetchCrds()
		Expect(err).To(Not(HaveOccurred()))
		crd := c.GetCrd("my-crds.com.example")
		Expect(crd).ToNot(BeNil())
	})
	It("return err if changes and force=false", func() {
		patt1f, err := os.ReadFile("./test_dir/patt1.yaml")
		Expect(err).To(Not(HaveOccurred()))
		var patt1 v1.CustomResourceDefinition
		err = yaml.Unmarshal(patt1f, &patt1)
		Expect(err).To(Not(HaveOccurred()))
		patt1.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["propertyOne"] = v1.JSONSchemaProps{Type: "string"}
		err = c.ApplyCrd(patt1)
		Expect(err).To(Not(HaveOccurred()))

		err = ApplyDir("./test_dir", false, &c, nexus_compare.CompareFiles)
		Expect(err).To(HaveOccurred())
	})

})
