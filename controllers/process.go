package controllers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"connector/pkg/utils"
)

func (r *CustomResourceDefinitionReconciler) ProcessAnnotation(crdType string,
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

	// Store Children and Parent information for a given CRD Type.
	utils.ConstructMapCRDTypeToParentHierarchy(eventType, crdType, n.Hierarchy)
	utils.ConstructMapCRDTypeToChildren(eventType, crdType, children)
	return nil
}
