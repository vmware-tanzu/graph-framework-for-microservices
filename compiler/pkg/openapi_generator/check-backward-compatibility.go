package openapi_generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	nexus_compare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func compareCRDs(inCompatibleCRDs []*bytes.Buffer, existingCRDContent string, newCRDContent []byte) ([]*bytes.Buffer, error) {
	newCRDParts := splitCRDs(newCRDContent)
	for _, newCRDPart := range newCRDParts {
		if newCRDPart == "" {
			continue
		}

		newCRD := &extensionsv1.CustomResourceDefinition{}
		err := yaml.Unmarshal([]byte(newCRDPart), newCRD)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling new CRD: %v", err)
		}

		existingCRD := &extensionsv1.CustomResourceDefinition{}
		err = yaml.Unmarshal([]byte(existingCRDContent), existingCRD)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling existing CRD: %v", err)
		}

		if newCRD.Name != existingCRD.Name {
			continue
		}

		// When there is a backward incompatibility, we fail the build if we don't force an upgrade.
		isInCompatible, message, err := nexus_compare.CompareFiles([]byte(existingCRDContent), []byte(newCRDPart))
		if err != nil {
			panic(fmt.Sprintf("Error occurred while checking CRD's %q backward compatibility: %v", existingCRD.Name, err))
		}
		if isInCompatible {
			log.Warnf("CRD %q is incompatible with the previous version", existingCRD.Name)
			inCompatibleCRDs = append(inCompatibleCRDs, message)
		}
	}
	return inCompatibleCRDs, nil
}

func CheckBackwardCompatibility(existingCRDsPath, yamlsPath string, force bool) error {
	var inCompatibleCRDs []*bytes.Buffer
	if err := filepath.Walk(existingCRDsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking existing CRD files: %v", err)
		}

		if !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}

		if info.IsDir() {
			fmt.Printf("Skipping dir %q\n", path)
			return nil
		}

		existingCRDContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading the existing CRD file on path %q: %v", path, err)
		}

		/*
			Find the new CRD for the file with the same name as the existing CRD file and compare it to the new CRD spec provided.
				existingCRDsPath - indicates the directory path of the existing crd yamls
					Ex: the path will be `example/output/generated/crds`

				path - indicates the file path of node
					Ex: For the node `Config` the path will be `example/output/generated/crds/config_config.yaml`

				yamlsPath - indicates the directory path of new crd yamls
					Ex: the path will be `_generated/crds`
		*/

		newFilePath := yamlsPath + strings.TrimPrefix(path, existingCRDsPath)
		newCRDContent, err := os.ReadFile(newFilePath)
		if err != nil || len(newCRDContent) == 0 {
			return fmt.Errorf("error reading the crd file on the path, appears CRD is removed %q: %v", newFilePath, err)
		}

		existingCRDParts := splitCRDs(existingCRDContent)
		for _, existingCRDPart := range existingCRDParts {
			if existingCRDPart == "" {
				continue
			}

			if inCompatibleCRDs, err = compareCRDs(inCompatibleCRDs, existingCRDPart, newCRDContent); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if inCompatibleCRDs != nil {
		if !force {
			// If the CRD are incompatible with the previous version, this will fail the build.
			return fmt.Errorf("datamodel upgrade failed due to incompatible datamodel changes:\n %v", inCompatibleCRDs)
		}
		log.Warnf("Upgrading the data model that is incompatible with the previous version: %v", inCompatibleCRDs)
	}
	return nil
}
