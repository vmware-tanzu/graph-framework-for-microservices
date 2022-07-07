package validate_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	k8sFake "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/yaml"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/validate"
)

var _ = Describe("Validation tests", func() {
	var (
		dynamicClient dynamic.Interface
		fakeClient    kubernetes.Interface
	)
	BeforeEach(func() {
		dynamicClient = fake.NewSimpleDynamicClient(runtime.NewScheme())
		fakeClient = k8sFake.NewSimpleClientset()
		_, _ = fakeClient.AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(context.TODO(),
			&admissionregistrationv1.ValidatingWebhookConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nexus-validation.webhook.svc",
				},
			}, metav1.CreateOptions{})
	})
	It("should accept non singleton object with any name", func() {
		crdDefJson, err := yaml.YAMLToJSON([]byte(getCRDDef(false)))
		Expect(err).NotTo(HaveOccurred())
		admReq := admissionv1.AdmissionReview{
			Request: &admissionv1.AdmissionRequest{
				Operation: admissionv1.Create,
				Kind:      metav1.GroupVersionKind{Kind: "CustomResourceDefinition"},
				Object: runtime.RawExtension{
					Raw: crdDefJson,
				},
			},
		}
		admRes := validate.CrdType(fakeClient, admReq)
		Expect(admRes.Response.Allowed).To(BeTrue())

		crdObjJson, err := yaml.YAMLToJSON([]byte(getCRDObject("foo")))
		Expect(err).NotTo(HaveOccurred())
		admReqCRDObj := admissionv1.AdmissionReview{
			Request: &admissionv1.AdmissionRequest{
				Operation: admissionv1.Create,
				Kind:      metav1.GroupVersionKind{Kind: "Root"},
				Object: runtime.RawExtension{
					Raw: crdObjJson,
				},
				Resource: metav1.GroupVersionResource{
					Group:    "orgchart.vmware.org",
					Resource: "roots",
				},
			},
		}

		admResCrd, err := validate.Crd(dynamicClient, admReqCRDObj)
		Expect(err).NotTo(HaveOccurred())
		Expect(admResCrd.Response.Allowed).To(BeTrue())
	})

	It("should accept singleton object with default as a name", func() {
		crdDefJson, err := yaml.YAMLToJSON([]byte(getCRDDef(true)))
		Expect(err).NotTo(HaveOccurred())
		admReq := admissionv1.AdmissionReview{
			Request: &admissionv1.AdmissionRequest{
				Operation: admissionv1.Create,
				Kind:      metav1.GroupVersionKind{Kind: "CustomResourceDefinition"},
				Object: runtime.RawExtension{
					Raw: crdDefJson,
				},
			},
		}
		admRes := validate.CrdType(fakeClient, admReq)
		Expect(admRes.Response.Allowed).To(BeTrue())

		crdObjJson, err := yaml.YAMLToJSON([]byte(getCRDObject("default")))
		Expect(err).NotTo(HaveOccurred())
		admReqCRDObj := admissionv1.AdmissionReview{
			Request: &admissionv1.AdmissionRequest{
				Operation: admissionv1.Create,
				Kind:      metav1.GroupVersionKind{Kind: "Root"},
				Object: runtime.RawExtension{
					Raw: crdObjJson,
				},
				Resource: metav1.GroupVersionResource{
					Group:    "orgchart.vmware.org",
					Resource: "roots",
				},
			},
		}

		admResCrd, err := validate.Crd(dynamicClient, admReqCRDObj)
		Expect(err).NotTo(HaveOccurred())
		Expect(admResCrd.Response.Allowed).To(BeTrue())
	})

	It("should reject singleton object with not default as a name", func() {
		crdDefJson, err := yaml.YAMLToJSON([]byte(getCRDDef(true)))
		Expect(err).NotTo(HaveOccurred())
		admReq := admissionv1.AdmissionReview{
			Request: &admissionv1.AdmissionRequest{
				Operation: admissionv1.Create,
				Kind:      metav1.GroupVersionKind{Kind: "CustomResourceDefinition"},
				Object: runtime.RawExtension{
					Raw: crdDefJson,
				},
			},
		}
		admRes := validate.CrdType(fakeClient, admReq)
		Expect(admRes.Response.Allowed).To(BeTrue())

		crdObjJson, err := yaml.YAMLToJSON([]byte(getCRDObject("foo")))
		Expect(err).NotTo(HaveOccurred())
		admReqCRDObj := admissionv1.AdmissionReview{
			Request: &admissionv1.AdmissionRequest{
				Operation: admissionv1.Create,
				Kind:      metav1.GroupVersionKind{Kind: "Root"},
				Object: runtime.RawExtension{
					Raw: crdObjJson,
				},
				Resource: metav1.GroupVersionResource{
					Group:    "orgchart.vmware.org",
					Resource: "roots",
				},
			},
		}

		admResCrd, err := validate.Crd(dynamicClient, admReqCRDObj)
		Expect(err).NotTo(HaveOccurred())
		Expect(admResCrd.Response.Allowed).To(BeFalse())
	})
})
