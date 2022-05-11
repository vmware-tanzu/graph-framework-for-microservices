package client

import (
	"api-gw/pkg/model"
	"context"
	"encoding/json"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"strings"
)

var Client dynamic.Interface

func New(config *rest.Config) (err error) {
	Client, err = dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func CreateObject(gvr schema.GroupVersionResource, kind, hashedName string, labels map[string]string, body map[string]interface{}) error {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvr.GroupVersion().String(),
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":   hashedName,
				"labels": labels,
			},
			"spec": body,
		},
	}

	// Create resource
	_, err := Client.Resource(gvr).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

type PatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type Patch []PatchOp

func (p Patch) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func UpdateParentGvk(parentCrdName string, parentCrd model.NodeInfo, labels map[string]string, childCrdName string, childName string) error {
	var (
		patchType types.PatchType
		marshaled []byte
	)

	parentParts := strings.Split(parentCrdName, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parentParts[1:], "."),
		Version:  "v1",
		Resource: parentParts[0],
	}

	parentName := labels[parentCrdName]
	hashedParentName := nexus.GetHashedName(parentCrdName, parentCrd.ParentHierarchy, labels, parentName)
	childGvk := parentCrd.Children[childCrdName]

	if childGvk.IsNamed {
		payload := "{\"spec\": {\"" + childGvk.FieldNameGvk + "\": {\"" + childName + "\": {\"name\": \"" + childName + "\"}}}}"
		patchType = types.MergePatchType
		marshaled = []byte(payload)
	} else {
		var patch Patch
		patchOp := PatchOp{
			Op:   "replace",
			Path: "/spec/" + childGvk.FieldNameGvk,
			Value: map[string]interface{}{
				"name": childName,
			},
		}
		patch = append(patch, patchOp)
		patchBytes, err := patch.Marshal()
		if err != nil {
			return err
		}
		marshaled = patchBytes
		patchType = types.JSONPatchType
	}

	_, err := Client.Resource(gvr).Patch(context.TODO(), hashedParentName, patchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}
