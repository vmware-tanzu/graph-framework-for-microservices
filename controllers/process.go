package controllers

import (
	"fmt"
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
			fmt.Printf("Error unmarshaling Nexus annotation %v\n", err)
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
	/*
	 populateEndpointCache populates the cache with CRD Type to restURIs attribute ( URL and method [GET, DELETE...]).
	 if any of the attribute removed in the new event notification, that should be removed from the cache and
	 triggers the server restart to remove the routes.
	 If any of the attribute added newly, notify that to `GlobalRestURIChan`.
	*/
	removed, added := model.PopulateEndpointCache(eventType, crdType, n.NexusRestAPIGen.Uris)
	if removed > 0 {
		r.StopCh <- struct{}{}
		model.GlobalRestURIChan <- model.GetGlobalEndpointCache()

		for k, _ := range model.GlobalCRDTypeToNodes {
			model.GlobalCRDChan <- k
		}
		return nil
	}

	if len(added) > 0 {
		model.GlobalRestURIChan <- added
	}
	return nil
}

func (r *CustomResourceDefinitionReconciler) ProcessCrdSpec(crdType string,
	spec apiextensionsv1.CustomResourceDefinitionSpec, eventType model.EventType) error {
	// It has stored the CRD type with the CRD spec
	model.ConstructMapCRDTypeToSpec(eventType, crdType, spec)
	return nil
}
