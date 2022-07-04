package controllers

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"connector/pkg/config"
	"connector/pkg/handlers"
	"connector/pkg/utils"
)

func StartReplicationConfigController(conf *config.Config, stopCh chan struct{}) {
	log.Info("Starting ReplicationConfig Controller...")
	localAPI, err := utils.SetUpLocalAPI()
	if err != nil {
		log.Fatalf("Error creating remote API: %v", err)
		return
	}
	localDynamicAPI, err := utils.SetUpDynamicLocalAPI()
	if err != nil {
		log.Fatalf("Error creating local API: %v", err)
		return
	}

	parts := strings.Split(utils.ReplicationConfigCRD, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}
	informer := GetInformer(gvr)
	handler := handlers.NewReplicationConfigHandler(gvr, conf.RemoteEndpoint, localDynamicAPI)
	c := newController(
		fmt.Sprintf("controller-%s", utils.ReplicationConfigCRD),
		localAPI,
		informer,
		handler,
		conf)
	go c.Run(stopCh)
}
