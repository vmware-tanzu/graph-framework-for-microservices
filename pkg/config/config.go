package config

import (
	"os"

	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

func GetConfig(filepath string) (*rest.Config, error) {
	_, err := os.Stat(filepath)

	var config *rest.Config
	if os.IsNotExist(err) {
		log.Tracef("Setting context to K8s base-api-server...")
		config, err = rest.InClusterConfig()
	} else {
		log.Tracef("Setting context to K8s local-api-server...")
		config, err = clientcmd.BuildConfigFromFlags("", filepath)
	}
	if err != nil {
		return nil, err
	}
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(200, 300)
	return config, nil
}
