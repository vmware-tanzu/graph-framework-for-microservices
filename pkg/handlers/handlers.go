package handlers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"connector/pkg/config"
	"connector/pkg/utils"
)

type RemoteHandler struct {
	Gvr                       schema.GroupVersionResource
	LocalClient, RemoteClient dynamic.Interface
	Config                    *config.Config
}

func NewRemoteHandler(gvr schema.GroupVersionResource,
	localClient, remoteDynamicClient dynamic.Interface, conf *config.Config) *RemoteHandler {
	return &RemoteHandler{
		Gvr:          gvr,
		LocalClient:  localClient,
		RemoteClient: remoteDynamicClient,
		Config:       conf,
	}
}

func (h *RemoteHandler) handleEvent(obj, oldObj *unstructured.Unstructured, eventType string) error {
	if h.Config.StatusReplicationEnabled {
		if err := ProcessStatus(obj, oldObj, h.RemoteClient); err != nil {
			log.Errorf("Error processing CR custom status during create %q: %v", obj.GetName(), err)
			// TODO: Need to continue or return on error?
			return err
		}
	}

	if err := Replicator(obj, h, eventType); err != nil {
		return err
	}

	return nil
}

func (h *RemoteHandler) Create(obj interface{}) error {
	currObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object during create event: %T", obj)
	}

	if err := h.handleEvent(currObj, nil, utils.Create); err != nil {
		return err
	}

	return nil
}

func (h *RemoteHandler) Update(obj interface{}, oldObj interface{}) error {
	currObj, ok := obj.(*unstructured.Unstructured)
	if ok {
		oldObject, ok := oldObj.(*unstructured.Unstructured)
		if !ok {
			return fmt.Errorf("unstructured client did not understand object during update event: %T", oldObj)
		}
		if err := h.handleEvent(currObj, oldObject, utils.Update); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Need to be tracked in a separate JIRA
func (h *RemoteHandler) Delete(obj interface{}) error {
	return nil
}
