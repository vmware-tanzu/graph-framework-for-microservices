package handlers

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"connector/pkg/utils"
)

func createObject(gvr schema.GroupVersionResource, res *unstructured.Unstructured,
	children utils.Children, client dynamic.Interface) error {
	objData := res.UnstructuredContent()
	spec := objData["spec"]

	// Ignore relationships.
	utils.DeleteChildGvkFields(spec.(map[string]interface{}), children)
	res.UnstructuredContent()["spec"] = spec

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvr.GroupVersion().String(),
			"kind":       res.GetKind(),
			"metadata": map[string]interface{}{
				"name":   res.GetName(),
				"labels": res.GetLabels(),
			},
			"spec": res.UnstructuredContent()["spec"],
		},
	}

	_, err := client.Resource(gvr).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func updateObject(gvr schema.GroupVersionResource, res *unstructured.Unstructured,
	children utils.Children, client dynamic.Interface) error {
	objData := res.UnstructuredContent()
	spec := objData["spec"]

	oldObject, err := client.Resource(gvr).Get(context.TODO(), res.GetName(), metav1.GetOptions{})
	if err != nil && strings.Contains(err.Error(), "not found") {
		log.Infof("Resource %s not found, creating instead", res.GetName())
		err := createObject(gvr, res, children, client)
		if err != nil {
			log.Errorf("Resource %v creation failed with an error: %v", res.GetName(), err)
			return err
		}
		return nil
	}

	// Ignore relationships.
	utils.DeleteChildGvkFields(spec.(map[string]interface{}), children)
	oldObject.UnstructuredContent()["spec"] = spec

	_, err = client.Resource(gvr).Update(context.TODO(), oldObject, metav1.UpdateOptions{})
	return err
}

func processEvents(h *RemoteHandler, eventType, hierarchy string, res *unstructured.Unstructured,
	children utils.Children, client dynamic.Interface, replicationEnabledNode bool) error {

	switch eventType {
	case utils.Create:
		log.Infof("Replication is enabled for the resource: %v, creating...", hierarchy)
		err := createObject(h.Gvr, res, children, client)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Errorf("Resource %v create failed with an error: %v", hierarchy, err)
			return err
		}
		if replicationEnabledNode {
			// Replicate children only if the obj exactly matches the source object.
			if err := ReplicateChildren(h.CrdType, res, h.LocalClient, client); err != nil {
				log.Errorf("Children replication failed for the resource %v: %v", hierarchy, err)
				return err
			}
		}
	case utils.Update:
		log.Infof("Replication is enabled for the resource: %v, updating...", hierarchy)
		err := updateObject(h.Gvr, res, children, client)
		if err != nil {
			log.Errorf("Resource %v update failed with an error: %v", hierarchy, err)
			return err
		}
	}
	return nil
}

// For every object notification received, replicator inspects if the object is of interest.
func Replicator(obj interface{}, h *RemoteHandler, eventType string) error {
	res := obj.(*unstructured.Unstructured)
	parents := utils.CRDTypeToParentHierarchy[h.CrdType]
	children := utils.CRDTypeToChildren[h.CrdType]
	labels := res.GetLabels()
	hierarchy := utils.GetNodeHierarchy(parents, labels, h.CrdType)

	repObj := &utils.ReplicationObject{
		Group: res.GroupVersionKind().Group,
		Kind:  res.GetKind(),
		Name:  res.GetName(),
	}

	var (
		client                 dynamic.Interface
		replicationEnabledNode bool
	)

	// Verify if the obj is of interest.
	// If obj exactly matches replication object source, replicate obj and its immediate children.
	if client, replicationEnabledNode = utils.ReplicationEnabledNode[*repObj]; replicationEnabledNode {
		if err := processEvents(h, eventType, hierarchy, res, children, client, true); err != nil {
			return err
		}
		return nil
	}

	// Verify if obj's immediate parent matches replication object source.
	// If yes, replicate obj.
	if len(parents) <= 0 {
		return nil
	}
	replicate := false
	immediateParent := parents[len(parents)-1]
	gvr := utils.GetGVRFromCrdType(immediateParent)

	// Get parent information from the object's labels to verify if the object's immediate parent is of interest.
	opts := metav1.ListOptions{LabelSelector: utils.GetParentLabels(parents, labels)}
	c, err := h.LocalClient.Resource(gvr).List(context.TODO(), opts)
	if err != nil {
		log.Errorf("error getting child objects of %v: %v", repObj.Name, err)
		return err
	}

	for _, item := range c.Items {
		repEnabledObj := &utils.ReplicationObject{
			Group: item.GroupVersionKind().Group,
			Kind:  item.GetKind(),
			Name:  item.GetName(),
		}

		if client, replicationEnabledNode = utils.ReplicationEnabledNode[*repEnabledObj]; replicationEnabledNode {
			replicate = true
			break
		}
	}

	if replicate {
		if err := processEvents(h, eventType, hierarchy, res, children, client, false); err != nil {
			return err
		}
	}
	return nil
}

// ReplicateNode replicates the exact node that is configured in the replicationConfig.
func ReplicateNode(repObj utils.ReplicationObject, localClient, remoteClient dynamic.Interface) error {
	utils.ConstructMapReplicationEnabledNode(repObj, remoteClient)
	crdType := utils.GetCrdType(repObj.Kind, repObj.Group)
	gvr := utils.GetGVRFromCrdType(crdType)
	children := utils.CRDTypeToChildren[crdType]

	repEnabledNode, err := localClient.Resource(gvr).Get(context.TODO(), repObj.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Errorf("error getting replication enabled object %v: %v", repObj.Name, err)
		return err
	}

	// If replication object is not found, simply return.
	if repEnabledNode == nil {
		return nil
	}

	log.Infof("Replication is enabled for the resource: %v, creating...", repObj.Name)
	err = createObject(gvr, repEnabledNode, children, remoteClient)
	if err != nil && !errors.IsAlreadyExists(err) {
		log.Errorf("Resource %v creation failed with an error: %v", repEnabledNode, err)
		return err
	}
	if err := ReplicateChildren(crdType, repEnabledNode, localClient, remoteClient); err != nil {
		log.Errorf("Children replication failed for replication object %v: %v", repEnabledNode, err)
		return err
	}

	return nil
}

// ReplicationChildren replicates the immediate children of desired node.
func ReplicateChildren(crdType string, repEnabledNode *unstructured.Unstructured, localClient, remoteClient dynamic.Interface) error {
	children := utils.CRDTypeToChildren[crdType]
	parents := utils.CRDTypeToParentHierarchy[crdType]
	labels := repEnabledNode.GetLabels()

	opts := utils.GetNodeLabels(parents, labels, crdType)
	for child := range children {
		parents = utils.CRDTypeToParentHierarchy[child]
		gvr := utils.GetGVRFromCrdType(child)
		c, err := localClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{LabelSelector: opts})
		if err != nil && !errors.IsNotFound(err) {
			log.Errorf("error getting child objects of %v: %v", repEnabledNode.GetName(), err)
			return err
		}

		for _, item := range c.Items {
			hierarchy := utils.GetNodeHierarchy(parents, item.GetLabels(), child)
			log.Infof("Replication is enabled for the resource: %v, creating...", hierarchy)

			err := createObject(gvr, &item, children, remoteClient)
			if err != nil && !errors.IsAlreadyExists(err) {
				log.Errorf("error creating resource %v: %v", hierarchy, err)
				return err
			}
		}
	}
	return nil
}
