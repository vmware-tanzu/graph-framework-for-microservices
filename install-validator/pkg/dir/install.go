package dir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func InstallDir(dir string, c kubewrapper.ClientInt) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Debugf("unexpected error while walking through dir: %s ", dir)
			return err
		}
		if info.IsDir() {
			logrus.Debugf("%s is a directory, skipping", info.Name())
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".yaml") {
			logrus.Debugf("%s is not a yaml file, skipping", info.Name())
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			logrus.Debugf("error while reading file %s", info.Name())
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
