package dir

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	nexuscompare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils"
	"sigs.k8s.io/yaml"
)

type compareFunc func([]byte, []byte) (bool, *bytes.Buffer, error)

func CheckDir(dir string, c kubewrapper.ClientInt, cFunc compareFunc) (map[string]*bytes.Buffer, error) {
	changes := make(map[string]*bytes.Buffer)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		newData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		name, err := nexuscompare.GetSpecName(newData)
		if err != nil {
			return err
		}
		crd := c.GetCrd(name)
		if crd == nil {
			return nil
		}

		if cFunc == nil {
			return errors.New("nil compare func passed")
		}
		actData, err := yaml.Marshal(crd)
		if err != nil {
			return err
		}

		inc, text, err := cFunc(actData, newData)
		if err != nil {
			return err
		}

		if inc { // if this is true, then there are some incompatible changes
			changes[name] = text
		}
		return nil
	})

	if err != nil {
		return changes, err
	}

	return changes, err
}

func GetOutdated(dir string, c kubewrapper.ClientInt) ([]string, error) {
	var toDelete []string
	var toInstall []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		newData, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		name, err := nexuscompare.GetSpecName(newData)
		if err != nil {
			return err
		}
		toInstall = append(toInstall, name)
		return nil
	})
	if err != nil {
		return toDelete, err
	}

	for _, crd := range c.GetCrds() {
		found := false
		for _, nameToInstall := range toInstall {
			if crd.Name == nameToInstall {
				found = true
				break
			}
		}
		if !found && strings.HasSuffix(crd.Spec.Group, c.GetGroup()) {
			toDelete = append(toDelete, crd.Name)
		}
	}
	return toDelete, nil
}
