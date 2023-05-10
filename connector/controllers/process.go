package controllers

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"connector/pkg/utils"
)

func (r *CustomResourceDefinitionReconciler) ProcessAnnotation(crdType string, gvr schema.GroupVersionResource,
	annotations map[string]string, eventType utils.EventType) error {
	n := utils.NexusAnnotation{}

	if eventType != utils.Delete {
		apiInfo, ok := annotations["nexus"]
		if !ok {
			return nil
		}

		// unmarshall to nexus annotation struct
		err := json.Unmarshal([]byte(apiInfo), &n)
		if err != nil {
			log.Errorf("Error unmarshaling Nexus annotation %v\n", err)
			return err
		}
	}

	children := make(map[string]utils.NodeHelperChild)
	if n.Children != nil {
		children = n.Children
	}

	// Store Children and Parent information for a given GVR.
	utils.ConstructMapGVRToParentHierarchy(eventType, gvr, n.Hierarchy)
	utils.ConstructMapGVRToChildren(eventType, gvr, children)

	// Store CRD version for a given CRD Type.
	utils.ConstructCRDTypeToCrdVersion(eventType, crdType, gvr.Version)
	return nil
}

/* ConstructGVR constructs group, version, resource for a CRD Type.
Eg: For a given CRD type: roots.vmware.org and ApiVersion: vmware.org/v1,
      group => vmware.org
	  resource => roots
	  version => v1
*/
func (r *CustomResourceDefinitionReconciler) ConstructGVR(crdType, crdVersion string) schema.GroupVersionResource {
	parts := strings.Split(crdType, ".")
	return schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  crdVersion,
		Resource: parts[0],
	}
}
