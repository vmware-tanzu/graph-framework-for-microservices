package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"connector/pkg/config"
	"connector/pkg/utils"
)

type ReplicationConfigHandler struct {
	Gvr         schema.GroupVersionResource
	Config      *config.Config
	LocalClient dynamic.Interface
}

func NewReplicationConfigHandler(gvr schema.GroupVersionResource,
	conf *config.Config, localClient dynamic.Interface) *ReplicationConfigHandler {
	return &ReplicationConfigHandler{
		Gvr:         gvr,
		Config:      conf,
		LocalClient: localClient,
	}
}

func (h *ReplicationConfigHandler) Create(obj interface{}) error {
	res := obj.(*unstructured.Unstructured)
	log.Infof("Received Create event for Replication Config %s", res.GetName())

	m := res.UnstructuredContent()
	spec := m["spec"]
	specInByte, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal replicationconfig spec %v: %v", res.GetName(), err)
	}
	repConf := &utils.ReplicationConfig{}
	if err := json.Unmarshal(specInByte, &repConf); err != nil {
		return fmt.Errorf("failed to unmarshal replicationconfig spec %v: %v", res.GetName(), err)
	}
	endpoint, err := h.GetEndpointObject(repConf.RemoteEndpoint.Name)
	if err != nil {
		return fmt.Errorf("failed to get endpoint object %v: %v", repConf.RemoteEndpoint.Name, err)
	}

	if endpoint.Host != h.Config.RemoteEndpointHost || endpoint.Port != h.Config.RemoteEndpointPort {
		return nil
	}
	host := fmt.Sprintf("%s:%s", h.Config.RemoteEndpointHost, h.Config.RemoteEndpointPort)
	remoteClient, err := utils.SetUpDynamicRemoteAPI(host, repConf.AccessToken)
	if err != nil {
		return fmt.Errorf("error creating dynamic remote API: %v", err)
	}

	rc := utils.ReplicationConfigSpec{LocalClient: h.LocalClient, Source: repConf.Source, Destination: repConf.Destination,
		RemoteClient: remoteClient, StatusEndpoint: repConf.StatusEndpoint}

	if err := ReplicateNode(res.GetName(), rc); err != nil {
		return fmt.Errorf("error replicating desired nodes: %v", err)
	}
	return nil
}

// ReplicationConfig will not be allowed to modify.
func (h *ReplicationConfigHandler) Update(obj interface{}, oldObj interface{}) error {
	return nil
}

// TODO: https://jira.eng.vmware.com/browse/NPT-378
func (h *ReplicationConfigHandler) Delete(obj interface{}) error {
	return nil
}

func (h *ReplicationConfigHandler) GetEndpointObject(name string) (*utils.NexusEndpoint, error) {
	parts := strings.Split(utils.NexusEndpointCRD, ".")
	gvr := schema.GroupVersionResource{
		Group:    strings.Join(parts[1:], "."),
		Version:  "v1",
		Resource: parts[0],
	}

	// Get the desired endpoint object to get the endpoint information.
	endpointObj, err := h.LocalClient.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	r := endpointObj.UnstructuredContent()
	spec := r["spec"]
	specInByte, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal endpoint spec of %v: %v", endpointObj.GetName(), err)
	}
	eObj := &utils.NexusEndpoint{}
	if err := json.Unmarshal(specInByte, &eObj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal endpoint spec of %v: %v", endpointObj.GetName(), err)
	}
	return &utils.NexusEndpoint{
		Host: eObj.Host,
		Cert: eObj.Cert,
		Port: eObj.Port,
	}, nil
}
