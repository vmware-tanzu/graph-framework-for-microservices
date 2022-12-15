package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	ext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"

	nexuscompare "github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus-compare"
	"github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/dir"
	kubewrapper "github.com/vmware-tanzu/graph-framework-for-microservices/install-validator/pkg/k8s-utils"
)

const (
	dirEnv        = "CRD_SPEC_DIR"
	forceEnv      = "CRD_FORCE"
	groupNamePath = "/NAME"
)

func main() {
	directory := ""
	if directory = os.Getenv(dirEnv); directory == "" {
		directory = "/crds"
	}
	force := false
	if crdForceStr := os.Getenv(forceEnv); strings.ToLower(crdForceStr) == "true" {
		force = true
	}

	// setup k8s client
	config, err := kubewrapper.GetRestConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	clientset, err := ext.NewForConfig(config)
	if err != nil {
		logrus.Fatal(err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logrus.Fatal(err)
	}
	c := kubewrapper.Client{Clientset: clientset, DynamicClient: dynamicClient}
	err = c.FetchGroup(groupNamePath)
	if err != nil {
		logrus.Fatal(err)
	}

	err = dir.ApplyDir(directory, force, &c, nexuscompare.CompareFiles)
	if err != nil {
		logrus.Fatal(err)
	}

}
