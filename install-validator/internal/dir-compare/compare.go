package dir_compare

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/rest"

	nexuscompare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	ext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func CheckDir(dir string) (bool, *bytes.Buffer, error) {
	var changes []*bytes.Buffer
	incompatible := false
	cs, err := getKubeClientSet()
	if err != nil {
		return true, nil, err
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
		actData, err := getCrdConfig(name, *cs)
		if err != nil || actData == "" {
			return err
		}
		inc, text, err := nexuscompare.CompareFiles([]byte(actData), newData)
		if err != nil {
			return err
		}
		incompatible = incompatible || inc // incompatible var is true if any of datamodel is not compatible
		if inc {                           // if this is true, then there are some incompatible changes
			changes = append(changes, text) // sums up the buffers with what was changed
		}
		return nil
	})

	change := new(bytes.Buffer)
	for _, c := range changes {
		change.Write(c.Bytes())
	}
	return incompatible, change, err
}

func getKubeClientSet() (*ext.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = getFlagConfig()
	}
	if err != nil {
		return nil, err
	}

	clientset, err := ext.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func getFlagConfig() (*rest.Config, error) {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}
	configPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getCrdConfig(name string, clientset ext.Clientset) (string, error) {
	res, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
	switch {
	case errors.IsNotFound(err):
		return "", nil
	case err != nil:
		return "", err
	}
	crd := res.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]

	return crd, nil
}
