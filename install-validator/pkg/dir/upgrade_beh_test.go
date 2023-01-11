package dir

import (
	"context"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/runtime"

	nexus_compare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

var _ = Describe("Generator", func() {

	ff := fake.NewSimpleClientset()
	groupRes := schema.GroupVersionResource{
		Group:    "outdated.tsm.tanzu.vmware.com",
		Version:  "v1",
		Resource: "roots",
	}
	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), map[schema.GroupVersionResource]string{
		groupRes: "mycrdsList",
	})

	c := kubewrapper.Client{Clientset: ff, DynamicClient: dyn}
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
	It("return delete outdated model but leave the one with other group", func() {
		patt1f, err := os.ReadFile("./test_dir/patt1.yaml")
		Expect(err).To(Not(HaveOccurred()))
		var patt1 v1.CustomResourceDefinition
		err = yaml.Unmarshal(patt1f, &patt1)
		Expect(err).To(Not(HaveOccurred()))
		err = c.ApplyCrd(patt1)
		Expect(err).To(Not(HaveOccurred()))

		rootO, err := os.ReadFile("./test_dir/root_root_outdated.yaml")
		Expect(err).To(Not(HaveOccurred()))
		var root v1.CustomResourceDefinition
		err = yaml.Unmarshal(rootO, &root)
		Expect(err).To(Not(HaveOccurred()))
		err = c.ApplyCrd(root)
		Expect(err).To(Not(HaveOccurred()))

		err = c.FetchGroup("./test_dir/grpName")
		Expect(err).To(Not(HaveOccurred()))

		err = ApplyDir("./test_dir2", true, &c, nexus_compare.CompareFiles)
		Expect(err).To(Not(HaveOccurred()))

		err = c.FetchCrds()
		Expect(err).To(Not(HaveOccurred()))

		crd := c.GetCrd(patt1.Name)
		Expect(crd).ToNot(BeNil())

		crdRoot := c.GetCrd(root.Name)
		Expect(crdRoot).To(BeNil())

	})
	It("should fail if there exist any resource in crd that is outdated", func() {
		_, err := c.DynamicClient.Resource(groupRes).Create(context.TODO(), &unstructured.Unstructured{}, metav1.CreateOptions{})
		Expect(err).To(Not(HaveOccurred()))

		rootO, err := os.ReadFile("./test_dir/root_root_outdated.yaml")
		Expect(err).To(Not(HaveOccurred()))
		var root v1.CustomResourceDefinition
		err = yaml.Unmarshal(rootO, &root)
		Expect(err).To(Not(HaveOccurred()))
		err = c.ApplyCrd(root)
		Expect(err).To(Not(HaveOccurred()))

		err = c.FetchGroup("./test_dir/grpName")
		Expect(err).To(Not(HaveOccurred()))

		err = ApplyDir("./test_dir2", true, &c, nexus_compare.CompareFiles)
		Expect(err).To(HaveOccurred())
	})

})
