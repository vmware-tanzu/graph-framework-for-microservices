package utils

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func SetUpDynamicRemoteAPI(host, token string) (dynamic.Interface, error) {
	conf := &rest.Config{
		Host:        host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	client, err := dynamic.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("can not connect to Kubernetes API: %v", err)
	}
	return client, nil
}

func SetUpLocalAPI() (kubernetes.Interface, error) {
	cfg, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get config")
	}
	api, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("can not create kubernetes coreClient: %v", err)
	}
	return api, nil
}

func SetUpDynamicLocalAPI() (dynamic.Interface, error) {
	cfg, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get config")
	}
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not generate dynamic client for config")
	}
	return dc, nil
}

func GetRestConfig() (*rest.Config, error) {
	filePath := os.Getenv("KUBECONFIG")
	if filePath == "" {
		filePath = "localapiserver.config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", filePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func SetUpAPIs() (kubernetes.Interface, dynamic.Interface, error) {
	localAPI, err := SetUpLocalAPI()
	if err != nil {
		log.Errorf("Error creating local API: %v", err)
		return nil, nil, err
	}
	localDynamicAPI, err := SetUpDynamicLocalAPI()
	if err != nil {
		log.Errorf("Error creating local Dynamic API: %v", err)
		return nil, nil, err
	}
	return localAPI, localDynamicAPI, nil
}
