package controllers

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	handler "gitlab.eng.vmware.com/nsx-allspark_users/m7/handler.git"

	"connector/controllers"
	"connector/pkg/config"
	"connector/pkg/handlers"
	"connector/pkg/utils"
)

func Start() {
	conf, err := config.LoadConfig("/config/connector-config")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	for {
		stopCh := make(chan struct{})
		go StartControllers(conf, stopCh)
		go StartReplicationConfigController(conf, stopCh)
		<-stopCh
		log.Info("Stop signal received, restarting crd watcher")
	}
}

func StartControllers(conf *config.Config, stopCh chan struct{}) {
	log.Infoln("Starting Controllers...")
	localAPI, err := utils.SetUpLocalAPI()
	if err != nil {
		log.Fatalf("Error creating local API: %v", err)
	}

	localDynamicAPI, err := utils.SetUpDynamicLocalAPI()
	if err != nil {
		log.Fatalf("Error creating local API: %v", err)
		return
	}

	for {
		select {
		case <-stopCh:
			log.Errorln("StartController() terminated..")
			return
		case crd := <-controllers.CrdCh:
			parts := strings.Split(crd.CrdType, ".")
			gvr := schema.GroupVersionResource{
				Group:    strings.Join(parts[1:], "."),
				Version:  "v1",
				Resource: parts[0],
			}
			// Skipping watcher creation for nexus datamodel CRDs.
			if gvr.Group == "connect.nexus.org" || gvr.Group == "apis.nexus.org" ||
				gvr.Group == "config.nexus.org" {
				continue
			}
			informer := GetInformer(gvr)
			handler := handlers.NewRemoteHandler(gvr, crd.CrdType, localDynamicAPI)
			c := newController(
				fmt.Sprintf("controller-%s", crd.CrdType),
				localAPI,
				informer,
				handler,
				conf)

			crdInfo := crd.CrdCache.Get(crd.CrdType)
			if crdInfo.Controller == nil {
				go c.Run(stopCh)
			}
			crd.CrdCache.UpsertController(crd.CrdType, c)
		}
	}
}

func newController(name string, client kubernetes.Interface,
	informer cache.SharedIndexInformer, handl handler.Handler,
	conf *config.Config) *handler.Controller {
	return handler.NewParallelController(
		name,
		client,
		informer,
		handl,
		conf.Dispatcher.WorkerTTL,
		conf.Dispatcher.MaxWorkerCount,
		conf.Dispatcher.CloseRequestsQueueSize,
		conf.Dispatcher.EventProcessedQueueSize,
		conf.IgnoredNamespaces.MatchNames,
		nil)
}

func GetInformer(gvr schema.GroupVersionResource) cache.SharedIndexInformer {
	cfg, err := utils.GetRestConfig()
	if err != nil {
		log.WithError(err).Fatal("could not get config")
	}
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.WithError(err).Fatal("could not generate dynamic client for config")
	}
	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, v1.NamespaceAll, nil)
	return f.ForResource(gvr).Informer()
}
