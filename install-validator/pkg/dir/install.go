package dir

import (
	"os"
	"path/filepath"
	"strings"

	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/kube-wrapper"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func InstallDir(dir string, c kubewrapper.ClientInt) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var crd v1.CustomResourceDefinition
		if err := yaml.Unmarshal(data, &crd); err != nil {
			return err
		}

		err = c.ApplyCrd(crd)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
