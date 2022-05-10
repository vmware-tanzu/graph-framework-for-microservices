package validate

import (
	"context"
	"encoding/json"
	"errors"
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

func CrdType(client *kubernetes.Clientset, r admissionv1.AdmissionReview) *admissionv1.AdmissionReview {
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
		err = ProcessCRDType(crd)
		if err != nil {
			admRes.Response.Allowed = false
			admRes.Response.Result.Message = err.Error()
			return admRes
		}
		UpdateValidationWebhook(client)
	}

	if r.Request.Operation == admissionv1.Update {
		oldCrd := v1.CustomResourceDefinition{}
		_, _, err := deserializer.Decode(r.Request.OldObject.Raw, nil, &oldCrd)
		if err != nil {
			admRes.Response.Allowed = false
			admRes.Response.Result.Message = "could not unmarshal old crd type"
			return admRes
		}

		if newNexus, ok := crd.Annotations["nexus"]; ok {
			if oldNexus, oldOk := oldCrd.Annotations["nexus"]; oldOk {
				if strings.Compare(newNexus, oldNexus) != 0 {
					admRes.Response.Allowed = false
					admRes.Response.Result.Message = "You are not allowed to change nexus annotation for this CRD"
					return admRes
				}
			}
		}
	}

	return admRes
}

func ProcessCRDType(crd v1.CustomResourceDefinition) error {
	nexusStr, ok := crd.Annotations["nexus"]
	if !ok {
		return nil
	}

	nexus := &NexusAnnotation{}
	err := json.Unmarshal([]byte(nexusStr), &nexus)
	if err != nil {
		log.Errorf("could not unmarshal nexus annotation: %v", err)
		return errors.New("could not unmarshal nexus annotation")
	}

	CrdParentsMap[crd.Name] = nexus.Hierarchy
	log.Infof("Added %s to parents map (%v)", crd.Name, nexus.Hierarchy)

	return nil
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
	parents := CrdParentsMap[crdName]
	labels := obj.ObjectMeta.GetLabels()

	for _, parent := range parents {
		parts := strings.Split(parent, ".")
		gvr := schema.GroupVersionResource{
			Group:    strings.Join(parts[1:], "."),
			Version:  "v1",
			Resource: parts[0],
		}
		parentParents := CrdParentsMap[parent]

		isNameHashed := false
		if val, ok := labels["nexus/is_name_hashed"]; ok {
			if val == "true" {
				isNameHashed = true
			}
		}

		var name string
		if label, ok := labels[parent]; ok {
			log.Infof("label %s found, val: %s", parent, label)
			if isNameHashed {
				name = nexus.GetHashedName(parent, parentParents, labels, label)
			}
			if getCrdObject(client, gvr, name) == nil {
				message := fmt.Sprintf("required parent %s with name %s not found", parent, name)
				log.Warn(message)

				admRes.Response.Allowed = false
				admRes.Response.Result.Message = message
				return admRes, nil
			}
		} else {
			if isNameHashed {
				name = nexus.GetHashedName(parent, parentParents, labels, "default")
			}

			log.Warnf("label %s not found", parent)
			if getCrdObject(client, gvr, name) == nil {
				message := fmt.Sprintf("required parent %s with name default not found", parent)
				log.Warn(message)

				admRes.Response.Allowed = false
				admRes.Response.Result.Message = message
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

		err = ProcessCRDType(crd)
		if err != nil {
			panic(err)
		}
	}

	log.Infof("parents map: %v", CrdParentsMap)
}
