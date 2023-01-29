package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"

	nexus_client "../../example/test-utils/output-group-name-with-hyphen-datamodel/crd_generated/nexus-client"
	"../../example/test-utils/output-group-name-with-hyphen-datamodel/crd_generated/nexus-gql/graph/model"
)

var c = GrpcClients{
	mtx:     sync.Mutex{},
	Clients: map[string]GrpcClient{},
}
var nc *nexus_client.Clientset

func InitNexusClientSet() error {
	if nc == nil {
		k8sApiConfig := getK8sAPIEndpointConfig()
		nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
		if err != nil {
			return fmt.Errorf("failed to get k8s client config: %s", err)
		}
		nc = nexusClient
		nc.SubscribeAll()
		log.Debugf("Subscribed to all nodes in datamodel")
	}
	return nil
}

func getParentName(parentLabels map[string]interface{}, key string) string {
	if v, ok := parentLabels[key]; ok && v != nil {
		return v.(string)
	}
	return ""
}

// ////////////////////////////////////
// Nexus K8sAPIEndpointConfig
// ////////////////////////////////////
func getK8sAPIEndpointConfig() *rest.Config {
	var (
		config *rest.Config
		err    error
	)
	filePath := os.Getenv("KUBECONFIG")
	if filePath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", filePath)
		if err != nil {
			return nil
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil
		}
	}
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(200, 300)
	return config
}

// ////////////////////////////////////
// Singleton Resolver for Parent Node
// PKG: Config, NODE: Config
// ////////////////////////////////////
func getRootResolver() (*model.ConfigConfig, error) {
	if nc == nil {
		log.Debugf("nc is nil in getRootResolver calling initNexusClientSet")
		initNCErr := InitNexusClientSet()
		if initNCErr != nil {
			log.Errorf("[getRootResolver]Error initializing nexus client: %s", initNCErr)
			return nil, nil
		}
	}

	vConfig, err := nc.GetConfigConfig(context.TODO())
	if err != nil {
		log.Errorf("[getRootResolver]Error getting Config node %s", err)
		return nil, nil
	}
	if vConfig == nil {
		log.Errorf("[getRootResolver]Error getting Root node %s", err)
		return nil, nil
	}
	dn := vConfig.DisplayName()
	parentLabels := map[string]interface{}{"configs.config.tsm-tanzu.vmware.com": dn}
	vFieldX := string(vConfig.Spec.FieldX)
	vFieldY := int(vConfig.Spec.FieldY)
	MyStructField, _ := json.Marshal(vConfig.Spec.MyStructField)
	MyStructFieldData := string(MyStructField)

	ret := &model.ConfigConfig{
		Id:            &dn,
		ParentLabels:  parentLabels,
		FieldX:        &vFieldX,
		FieldY:        &vFieldY,
		MyStructField: &MyStructFieldData,
	}
	log.Debugf("[getRootResolver]Output Config object %+v", ret)
	return ret, nil
}
