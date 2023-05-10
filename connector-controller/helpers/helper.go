package helpers

import (
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNexusConfig() (*rest.Config, error) {
	filePath := os.Getenv("INNER_KUBECONFIG")
	if filePath == "" {
		filePath = "localapiserver.config"
	}

	config, err := clientcmd.BuildConfigFromFlags("", filePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GetNexusRestConfig() *rest.Config {
	return &rest.Config{
		Host: "http://nexus-proxy-container:80",
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
}
