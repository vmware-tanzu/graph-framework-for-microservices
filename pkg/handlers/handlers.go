package handlers

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"connector/pkg/utils"
)

type RemoteHandler struct {
	Gvr         schema.GroupVersionResource
	CrdType     string
	LocalClient dynamic.Interface
}

func NewRemoteHandler(gvr schema.GroupVersionResource, crdType string, localClient dynamic.Interface) *RemoteHandler {
	return &RemoteHandler{
		Gvr:         gvr,
		CrdType:     crdType,
		LocalClient: localClient,
	}
}

func (h *RemoteHandler) Create(obj interface{}) error {
	if err := Replicator(obj, h, utils.Create); err != nil {
		return err
	}
	return nil
}

func (h *RemoteHandler) Update(obj interface{}, oldObj interface{}) error {
	if err := Replicator(obj, h, utils.Update); err != nil {
		return err
	}
	return nil
}

// TODO: Need to be tracked in a separate JIRA
func (h *RemoteHandler) Delete(obj interface{}) error {
	return nil
}
