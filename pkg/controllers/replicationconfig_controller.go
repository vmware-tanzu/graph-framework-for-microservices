package controllers

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"connector/pkg/config"
	"connector/pkg/handlers"
	"connector/pkg/utils"
)

func StartReplicationConfigController(conf *config.Config, stopCh chan struct{}) {
	log.Info("Starting ReplicationConfig Controller...")
	localAPI, localDynamicAPI, err := utils.SetUpAPIs()
	if err != nil {
		log.Fatalf("Error creating APIs: %v", err)
		return
	}

	gvr := utils.GetGVRFromCrdType(utils.ReplicationConfigCRD)
	informer := GetInformer(gvr)
	handler := handlers.NewReplicationConfigHandler(gvr, conf, localDynamicAPI)
	c := newController(
		fmt.Sprintf("controller-%s", utils.ReplicationConfigCRD),
		localAPI,
		informer,
		handler,
		conf)
	go c.Run(stopCh)
}
