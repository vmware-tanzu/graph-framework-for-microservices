package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func CrdType(client kubernetes.Interface, r admissionv1.AdmissionReview) *admissionv1.AdmissionReview {
	admRes := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     r.Request.UID,
			Allowed: true,
			Result: &metav1.Status{
				Message: "",
			},
		},
	}

	scheme := runtime.NewScheme()
	codecFactory := serializer.NewCodecFactory(scheme)
	deserializer := codecFactory.UniversalDeserializer()

	crd := v1.CustomResourceDefinition{}
	_, _, err := deserializer.Decode(r.Request.Object.Raw, nil, &crd)
	if err != nil {
		admRes.Response.Allowed = false
		admRes.Response.Result.Message = "could not unmarshal crd type"
		return admRes
	}

	if r.Request.Operation == admissionv1.Create {
		err = CRDs.ProcessNewCRDType(crd)
		if err != nil {
			admRes.Response.Allowed = false
			admRes.Response.Result.Message = err.Error()
			return admRes
		}
		UpdateValidationWebhook(client)
	}

	if r.Request.Operation == admissionv1.Update {
		err = CRDs.ProcessNewCRDType(crd)
		if err != nil {
			admRes.Response.Allowed = false
			admRes.Response.Result.Message = err.Error()
			return admRes
		}
	}

	return admRes
}

func Crd(client dynamic.Interface, r admissionv1.AdmissionReview) (*admissionv1.AdmissionReview, error) {
	admRes := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     r.Request.UID,
			Allowed: true,
			Result: &metav1.Status{
				Message: "",
			},
		},
	}

	raw := r.Request.Object.Raw
	obj := struct {
		metav1.ObjectMeta `json:"metadata,omitempty"`
	}{}
	err := json.Unmarshal(raw, &obj)
	if err != nil {
		log.Warnf("could not unmarshal object meta: %v", err)
		return nil, fmt.Errorf("could not unmarshal object meta: %v", err)
	}

	crdName := fmt.Sprintf("%s.%s", r.Request.Resource.Resource, r.Request.Resource.Group)
	labels := obj.ObjectMeta.GetLabels()

	isCRDSingleton := CRDs.IsSingleton(crdName)
	if isCRDSingleton {
		valid := true
		if val, ok := labels[IS_NAME_HASHED_LABEL]; ok {
			if val == "true" {
				if labels[DISPLAY_NAME_LABEL] != DEFAULT_KEY {
					valid = false
				}
			}
		} else {
			if obj.GetName() != "default" {
				valid = false
			}
		}

		if !valid {
			message := "singleton object display name can only be 'default'"
			setResponseToNotAllowed(admRes, message)
			return admRes, nil
		}
	}

	parents, err := CRDs.GetParents(crdName, client)
	if err != nil {
		message := fmt.Sprintf("Couldn't determine parents info about for CRD %s, please make sure CRD definition is applied", crdName)
		setResponseToNotAllowed(admRes, message)
		return admRes, nil
	}

	for _, parent := range parents {
		parts := strings.Split(parent, ".")
		gvr := schema.GroupVersionResource{
			Group:    strings.Join(parts[1:], "."),
			Version:  "v1",
			Resource: parts[0],
		}
		parentParents, err := CRDs.GetParents(parent, client)
		if err != nil {
			message := fmt.Sprintf("Couldn't determine parent info %s for CRD %s, please make sure CRD definition is applied", parent, crdName)
			setResponseToNotAllowed(admRes, message)
			return admRes, nil
		}

		isNameHashed := false
		if val, ok := labels[IS_NAME_HASHED_LABEL]; ok {
			if val == "true" {
				isNameHashed = true
			}
		}

		var name string
		if label, ok := labels[parent]; ok {
			log.Infof("label %s found, val: %s", parent, label)
			displayName := label
			if isNameHashed {
				name = nexus.GetHashedName(parent, parentParents, labels, label)
			} else {
				name = label
			}
			if getCrdObject(client, gvr, name) == nil {
				message := fmt.Sprintf("required parent %s with display name %s not found", parent, displayName)
				setResponseToNotAllowed(admRes, message)
				return admRes, nil
			}
		} else {
			if isNameHashed {
				name = nexus.GetHashedName(parent, parentParents, labels, DEFAULT_KEY)
			} else {
				name = DEFAULT_KEY
			}

			log.Warnf("label %s not found", parent)
			if getCrdObject(client, gvr, name) == nil {
				message := fmt.Sprintf("required parent %s with name default not found", parent)
				setResponseToNotAllowed(admRes, message)
				return admRes, nil
			}
		}
	}

	return admRes, nil
}

func getCrdObject(client dynamic.Interface, gvr schema.GroupVersionResource, name string) interface{} {
	obj, err := client.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	return obj
}

func ProcessCRDs(client dynamic.Interface) {
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	list, err := client.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, unstructuredCrd := range list.Items {
		var crd v1.CustomResourceDefinition
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredCrd.UnstructuredContent(), &crd)
		if err != nil {
			panic(err)
		}

		err = CRDs.ProcessNewCRDType(crd)
		if err != nil {
			panic(err)
		}
	}
}

func setResponseToNotAllowed(admRes *admissionv1.AdmissionReview, message string) {
	log.Warn(message)
	admRes.Response.Allowed = false
	admRes.Response.Result.Message = message
}
