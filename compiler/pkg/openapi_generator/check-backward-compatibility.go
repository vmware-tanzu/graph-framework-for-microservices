package openapi_generator

import (
	"bytes"
	"fmt"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	nexus_compare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func CheckBackwardCompatibility(inCompatibleCRDs []*bytes.Buffer, crd extensionsv1.CustomResourceDefinition, oldCRDContent []byte) ([]*bytes.Buffer, error) {
	oldCRDParts := splitCRDs(oldCRDContent)
	for _, oldPart := range oldCRDParts {
		if oldPart == "" {
			continue
		}
		oldCRD := &extensionsv1.CustomResourceDefinition{}
		err := yaml.Unmarshal([]byte(oldPart), oldCRD)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling existing crd: %v", err)
		}

		newCRDPart, err := yaml.Marshal(&crd)
		if err != nil {
			return nil, fmt.Errorf("error marshaling new crd: %v", err)
		}

		if oldCRD.Name != crd.Name {
			continue
		}

		// When there is a backward incompatibility, we fail the build if we don't force an upgrade.
		isInCompatible, message, err := nexus_compare.CompareFiles([]byte(oldPart), newCRDPart)
		if err != nil {
			log.Errorf("Error occurred while checking CRD's %q backward compatibility: %v", crd.Name, err)
		}
		if isInCompatible {
			log.Warnf("CRD %q is incompatible with the previous version", crd.Name)
			inCompatibleCRDs = append(inCompatibleCRDs, message)
		}
	}
	return inCompatibleCRDs, nil
}
