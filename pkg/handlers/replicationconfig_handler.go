package handlers

import (
	"context"
	"encoding/json"
	"fmt"

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

	spec := res.UnstructuredContent()["spec"]
	specInByte, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal replicationconfig spec %v: %v", res.GetName(), err)
	}
	repConf := &utils.ReplicationConfig{}
	if err := json.Unmarshal(specInByte, &repConf); err != nil {
		return fmt.Errorf("failed to unmarshal replicationconfig spec %v: %v", res.GetName(), err)
	}

	eObj, err := h.GetEndpointObject(repConf.RemoteEndpoint.Name)
	if err != nil {
		return fmt.Errorf("failed to get endpoint object %v: %v", repConf.RemoteEndpoint.Name, err)
	}

	if !h.isValidReplicationConfig(eObj) {
		return nil
	}

	host := utils.ConstructURL(h.Config.RemoteEndpointHost, h.Config.RemoteEndpointPort, h.Config.RemoteEndpointPath)

	log.Infof("Connecting to the destination host: %v", host)
	remoteClient, err := utils.SetUpDynamicRemoteAPI(host, repConf.AccessToken, eObj.Cert)
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
	gvr := utils.GetGVRFromCrdType(utils.NexusEndpointCRD, utils.V1Version)

	// Get the desired endpoint object to get the endpoint information.
	endpointObj, err := h.LocalClient.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	spec := endpointObj.UnstructuredContent()["spec"]
	specInByte, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal endpoint spec of %v: %v", endpointObj.GetName(), err)
	}
	eObj := &utils.NexusEndpoint{}
	if err := json.Unmarshal(specInByte, &eObj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal endpoint spec of %v: %v", endpointObj.GetName(), err)
	}

	if eObj.Host == "" {
		return nil, fmt.Errorf("host can't be empty in endpoint spec: %v", endpointObj.GetName())
	}
	return eObj, nil
}

// isValidReplicationConfig checks if the replication-config is of interest to the desired destination.
// If not, skips this replication-config.
func (h *ReplicationConfigHandler) isValidReplicationConfig(eObj *utils.NexusEndpoint) bool {
	if eObj.Host != h.Config.RemoteEndpointHost ||
		eObj.Port != h.Config.RemoteEndpointPort ||
		eObj.Path != h.Config.RemoteEndpointPath {
		return false
	}
	return true
}
