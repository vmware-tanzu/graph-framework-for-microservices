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
