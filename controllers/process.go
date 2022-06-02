package controllers

import (
	log "github.com/sirupsen/logrus"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/json"

	"api-gw/pkg/model"
)

func (r *CustomResourceDefinitionReconciler) ProcessAnnotation(crdType string,
	annotations map[string]string, eventType model.EventType) error {
	n := model.NexusAnnotation{}

	if eventType != model.Delete {
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

	children := make(map[string]model.NodeHelperChild)
	if n.Children != nil {
		children = n.Children
	}

	// It has stored the URI with the CRD type and CRD type with the Node Info.
	model.ConstructMapURIToCRDType(eventType, crdType, n.NexusRestAPIGen.Uris)
	model.ConstructMapCRDTypeToNode(eventType, crdType, n.Name, n.Hierarchy, children)
	model.ConstructMapCRDTypeToRestUris(eventType, crdType, n.NexusRestAPIGen)

	// Restart echo server
	log.Debugln("Restarting echo server...")
	r.StopCh <- struct{}{}

	for cType, uris := range model.CrdTypeToRestUris {
		model.RestURIChan <- uris
		model.CrdTypeChan <- cType
	}
	return nil
}

func (r *CustomResourceDefinitionReconciler) ProcessCrdSpec(crdType string,
	spec apiextensionsv1.CustomResourceDefinitionSpec, eventType model.EventType) error {
	// It has stored the CRD type with the CRD spec
	model.ConstructMapCRDTypeToSpec(eventType, crdType, spec)
	return nil
}
