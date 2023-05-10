package handlers

import (
	"context"
	"fmt"
	"reflect"

	"encoding/json"

	log "github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"

	"connector/pkg/utils"
)

func getStatus(obj map[string]interface{}) (map[string]interface{}, error) {
	status, _, err := unstructured.NestedMap(obj, "status")
	if err != nil {
		return nil, fmt.Errorf("error occurred in obtaining the status object: %v", err)
	}

	return status, nil
}

func extractResourceInfo(annotations map[string]string) (*utils.ResourceAnnotation, bool, error) {
	resourceInfo, ok := annotations[utils.NexusReplicationResource]
	if !ok {
		return nil, false, fmt.Errorf("CR annotation doesn't contain `NexusReplicationResource[GVR]` : %v", annotations)
	}

	// unmarshall to nexus sourceCR struct
	r := utils.ResourceAnnotation{}
	err := json.Unmarshal([]byte(resourceInfo), &r)
	if err != nil {
		return nil, false, fmt.Errorf("error unmarshalling resource info from CR annotation: %v", annotations)
	}

	return &r, true, nil
}

func nexusReplicationManaged(annotations map[string]string) bool {
	_, ok := annotations[utils.NexusReplicationManager]
	return ok
}

func ProcessStatus(res, oldObj *unstructured.Unstructured, remoteClient dynamic.Interface) error {
	annotations := res.GetAnnotations()
	if !nexusReplicationManaged(annotations) {
		log.Debugf("CR %q not replicated by nexus connector, skipping: %v", res.GetName(), annotations)
		return nil
	}

	r, found, err := extractResourceInfo(annotations)
	if err != nil || !found || r == nil {
		return err
	}

	// Get status from the old and current object
	oldStatus := map[string]interface{}{}
	if oldObj != nil {
		if oldStatus, err = getStatus(oldObj.UnstructuredContent()); err != nil {
			return err
		}
	}

	currentStatus, err := getStatus(res.UnstructuredContent())
	if err != nil {
		return err
	}

	log.Debugf("Old status %v and Current status %v of CR %q", oldStatus, currentStatus, res.GetName())

	if reflect.DeepEqual(oldStatus, currentStatus) {
		log.Debugf("No status changes %v found for CR %q, skip handling", currentStatus, res.GetName())
		return nil
	}

	patchBytes, err := CreateStatusPatch(oldStatus, currentStatus)
	if err != nil {
		return fmt.Errorf("could not create patch for the CR(%q): %v", res.GetName(), err)
	}
	log.Debugf("patchBytes %+v for CR %q", string(patchBytes), res.GetName())

	_, err = remoteClient.Resource(r.GVR).Patch(context.TODO(), r.Name, types.JSONPatchType, patchBytes, metav1.PatchOptions{}, "status")
	if err != nil {
		log.Errorf("Resource %s patching failed with an error: %v", r.Name, err)
		return err
	}

	return nil
}

func removeFromStatus(patch *Patch, oldStatus, currentStatus map[string]interface{}) {
	for k1 := range oldStatus {
		found := false
		if currentStatus != nil {
			if _, ok := currentStatus[k1]; ok {
				found = true
				break
			}
		}

		if !found {
			patchOp := PatchOp{
				Op:   "remove",
				Path: "/status/" + k1,
			}
			*patch = append(*patch, patchOp)
		}
	}
}

func addToStatus(patch *Patch, currentStatus map[string]interface{}) {
	for k1, v1 := range currentStatus {
		patchOp := PatchOp{
			Op:    "replace",
			Path:  "/status/" + k1,
			Value: v1,
		}
		*patch = append(*patch, patchOp)
	}
}

func CreateStatusPatch(oldStatus, currentStatus map[string]interface{}) ([]byte, error) {
	patch := &Patch{}

	removeFromStatus(patch, oldStatus, currentStatus)
	addToStatus(patch, currentStatus)

	return patch.Marshal()
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
