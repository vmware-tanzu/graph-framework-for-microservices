package openapi_generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	nexus_compare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func splitCRDs(content []byte) []string {
	return strings.Split(string(content), "---")
}

func CheckBackwardCompatibility(existingCRDsPath, yamlsPath string, force bool) error {
	var (
		removedCRDs      []string
		inCompatibleCRDs []*bytes.Buffer
	)

	if err := filepath.Walk(existingCRDsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("No files exists on the path %q: %v", path, err)
				return nil
			}
			return fmt.Errorf("walking existing CRD's failed with error: %v", err)
		}

		if info.IsDir() {
			log.Debugf("Skipping dir %q\n", path)
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".yaml") {
			log.Debugf("Expected filename with suffix %v but got %v", ".yaml", info.Name())
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
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("error reading the crd file on the path %q: %v", newFilePath, err)
		}

		for _, existingCRDPart := range splitCRDs(existingCRDContent) {
			if existingCRDPart == "" {
				continue
			}

			existingCRD := &extensionsv1.CustomResourceDefinition{}
			err := yaml.Unmarshal([]byte(existingCRDPart), existingCRD)
			if err != nil {
				return fmt.Errorf("error unmarshaling existing CRD: %v", err)
			}

			found := false
			for _, newCRDPart := range splitCRDs(newCRDContent) {
				if newCRDPart == "" {
					continue
				}

				newCRD := &extensionsv1.CustomResourceDefinition{}
				err := yaml.Unmarshal([]byte(newCRDPart), newCRD)
				if err != nil {
					return fmt.Errorf("error unmarshaling new CRD: %v", err)
				}

				if newCRD.Name != existingCRD.Name {
					continue
				}

				found = true
				isInCompatible, message, err := nexus_compare.CompareFiles([]byte(existingCRDPart), []byte(newCRDPart))
				if err != nil {
					return err
				}
				if isInCompatible {
					log.Warnf("CRD %q is incompatible with the previous version", existingCRD.Name)
					inCompatibleCRDs = append(inCompatibleCRDs, message)
				}
			}

			// Appears node is removed in the latest version
			if !found {
				removedCRDs = append(removedCRDs, existingCRD.Name)
				continue
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if len(inCompatibleCRDs) > 0 || len(removedCRDs) > 0 {
		inCompatibleCRDsChanges := &bytes.Buffer{}
		for _, crd := range inCompatibleCRDs {
			inCompatibleCRDsChanges.Write(crd.Bytes())
		}
		for _, crd := range removedCRDs {
			inCompatibleCRDsChanges.WriteString(fmt.Sprintf("%q is deleted\n", crd))
		}
		// If the CRD are incompatible with the previous version, this will fail the build.
		if !force {
			return fmt.Errorf("datamodel upgrade failed due to incompatible datamodel changes: \n %v", inCompatibleCRDsChanges)
		}
		log.Warnf("Upgrading the data model that is incompatible with the previous version: \n %v", inCompatibleCRDsChanges)
	}

	return nil
}
