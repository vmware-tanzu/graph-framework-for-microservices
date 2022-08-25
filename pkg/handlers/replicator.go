package handlers

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"

	"connector/pkg/utils"
)

func createObject(gvr schema.GroupVersionResource, res *unstructured.Unstructured, children utils.Children, hierarchy string,
	client, localClient dynamic.Interface, source utils.ReplicationSource, destination utils.ReplicationDestination) error {

	objData := res.UnstructuredContent()
	labels := make(map[string]string)
	for key, val := range res.GetLabels() {
		labels[key] = val
	}
	spec := objData["spec"]

	if source.Kind == utils.Object && !destination.Hierarchical && children != nil {
		// Ignore relationships.
		utils.DeleteChildGvkFields(spec.(map[string]interface{}), children)
		res.UnstructuredContent()["spec"] = spec
	}

	if destination.Hierarchical {
		for _, label := range destination.Hierarchy.Labels {
			labels[label.Key] = label.Value
		}
	}

	annotations := map[string]string{}
	annotations = utils.GenerateAnnotations(annotations, gvr, res.GetName())

	// If destination object type is specified, then replicate the object to that type.
	destGvr, destKind := utils.GetDestinationGvrAndKind(destination, gvr, res.GetKind())

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": destGvr.GroupVersion().String(),
			"kind":       destKind,
			"metadata": map[string]interface{}{
				"name":        res.GetName(),
				"labels":      labels,
				"namespace":   destination.Namespace,
				"annotations": annotations,
			},
			"spec": res.UnstructuredContent()["spec"],
		},
	}

	if destObj, err := client.Resource(destGvr).Namespace(destination.Namespace).Get(context.TODO(), obj.GetName(),
		metav1.GetOptions{}); destObj != nil && err == nil {
		return nil
	}

	log.Infof("Replication is enabled for the resource: %v, creating...", hierarchy)
	// If the object was successfully replicated, we need to patch the source and remote generation ID.
	destObj, err := client.Resource(destGvr).Namespace(destination.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	if err == nil && destObj != nil {
		if err = patchStatusObject(localClient, gvr, res.GetName(), res.GetGeneration(), destObj.GetGeneration()); err != nil {
			log.Errorf("CR %q status patch failed with an error: %v", res.GetName(), err)
			return err
		}
	}
	return err
}

func updateObject(gvr schema.GroupVersionResource, res *unstructured.Unstructured, children utils.Children, hierarchy string,
	client, localClient dynamic.Interface, source utils.ReplicationSource, destination utils.ReplicationDestination) error {

	objData := res.UnstructuredContent()
	spec := objData["spec"]

	// If destination object type is specified, then replicate the object to that type.
	destGvr, _ := utils.GetDestinationGvrAndKind(destination, gvr, "")

	oldObject, err := client.Resource(destGvr).Namespace(destination.Namespace).Get(context.TODO(), res.GetName(), metav1.GetOptions{})
	if err != nil && strings.Contains(err.Error(), "not found") {
		log.Infof("Resource %s not found, creating instead", res.GetName())
		return createObject(gvr, res, children, hierarchy, client, localClient, source, destination)
	}

	if source.Object.Hierarchical && !destination.Hierarchical {
		// Ignore relationships.
		utils.DeleteChildGvkFields(spec.(map[string]interface{}), children)
	}
	oldObject.UnstructuredContent()["spec"] = spec

	annotations := oldObject.GetAnnotations()
	annotations = utils.GenerateAnnotations(annotations, gvr, res.GetName())
	oldObject.SetAnnotations(annotations)

	log.Infof("Replication is enabled for the resource: %v, updating...", hierarchy)
	// If the object was successfully replicated, we need to patch the source and remote generation ID.
	destObj, err := client.Resource(destGvr).Namespace(destination.Namespace).Update(context.TODO(), oldObject, metav1.UpdateOptions{})
	if err == nil && destObj != nil {
		if err = patchStatusObject(localClient, gvr, res.GetName(), res.GetGeneration(), destObj.GetGeneration()); err != nil {
			log.Errorf("CR %q status patch failed with an error: %v", res.GetName(), err)
			return err
		}
	}
	return err
}

func patchStatusObject(localClient dynamic.Interface, gvr schema.GroupVersionResource,
	repObjName string, srcGeneration, destGeneration int64) error {
	patchBytes, err := utils.CreatePatch(srcGeneration, destGeneration)
	if err != nil {
		log.Errorf("Could not create patch for the CR(%q) %v", repObjName, err)
		return err
	}

	log.Debugf("Patching status of CR %q: %v", repObjName, string(patchBytes))

	_, err = localClient.Resource(gvr).Patch(context.TODO(), repObjName, types.MergePatchType, patchBytes, metav1.PatchOptions{}, "status")
	if err != nil {
		log.Errorf("Resource %s patching failed with an error: %v", repObjName, err)
		return err
	}

	return nil
}

func processEvents(h *RemoteHandler, eventType, hierarchy string, res *unstructured.Unstructured, children utils.Children,
	repInfo utils.ReplicationConfigSpec, replicationEnabledNode bool) error {

	switch eventType {
	case utils.Create:
		err := createObject(h.Gvr, res, children, hierarchy, repInfo.Client, h.LocalClient, repInfo.Source, repInfo.Destination)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Errorf("Resource %v create failed with an error: %v", hierarchy, err)
			return err
		}
		if repInfo.Source.Object.Hierarchical && replicationEnabledNode {
			// Replicate children only if the obj exactly matches the source object.
			newRepDestination := repInfo.Destination
			newRepDestination.IsChild = true
			if err := ReplicateChildren(h.CrdType, res, h.LocalClient, repInfo.Client, repInfo.Source, newRepDestination); err != nil {
				log.Errorf("Children replication failed for the resource %v: %v", hierarchy, err)
				return err
			}
		}
	case utils.Update:
		err := updateObject(h.Gvr, res, children, hierarchy, repInfo.Client, h.LocalClient, repInfo.Source, repInfo.Destination)
		if err != nil {
			log.Errorf("Resource %v update failed with an error: %v", hierarchy, err)
			return err
		}
	}
	return nil
}

// For every object notification received, replicator replicates in the following order.
// 1. If the CRDType is of interest, simply replicate the object.
// 2. If only the object is of interest and if it is a part of a graph, then replicate the object and its immediate children.
// 3. If the oject is not a part of a graph, then replicate only that object.
func Replicator(obj interface{}, h *RemoteHandler, eventType string) error {

	res := obj.(*unstructured.Unstructured)

	var (
		repInfo   utils.ReplicationConfigSpec
		repObj    utils.ReplicationObject
		hierarchy string
	)

	// Verify if the type is of interest.
	// If the CRDType is enabled for replication, simply replicate the object.
	if repInfo, replicationEnabledType := utils.ReplicationEnabledCRDType[h.CrdType]; replicationEnabledType {
		if err := processEvents(h, eventType, res.GetName(), res, nil, repInfo, false); err != nil {
			return err
		}
		return nil
	}

	labels := res.GetLabels()
	parents := utils.CRDTypeToParentHierarchy[h.CrdType]
	children := utils.CRDTypeToChildren[h.CrdType]

	name := res.GetName()
	if objName, ok := labels[utils.DisplayNameKey]; ok {
		name = objName
	}
	repObj = utils.GetReplicationObject(res.GroupVersionKind().Group, res.GetKind(), name)

	// Verify if the obj is of interest.
	repConf, replicationEnabledNode := utils.ReplicationEnabledNode[repObj]
	if replicationEnabledNode {
		hierarchy = name
		if repConf.Source.Object.Hierarchical {
			hierarchy = utils.GetNodeHierarchy(parents, labels, h.CrdType)
		}
		if err := processEvents(h, eventType, hierarchy, res, children, repConf, true); err != nil {
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

	replicationEnabledNode = false
	for _, item := range c.Items {
		name := item.GetName()
		if objName, ok := item.GetLabels()[utils.DisplayNameKey]; ok {
			name = objName
		}
		repObj := utils.GetReplicationObject(item.GroupVersionKind().Group, item.GetKind(), name)
		if repInfo, replicationEnabledNode = utils.ReplicationEnabledNode[repObj]; replicationEnabledNode {
			hierarchy = name
			if repInfo.Source.Object.Hierarchical {
				hierarchy = utils.GetNodeHierarchy(parents, labels, h.CrdType)
			}
			replicate = true
			break
		}
	}
	if replicate {
		newRepInfo := repInfo
		newRepInfo.Destination.IsChild = true
		if err := processEvents(h, eventType, hierarchy, res, children, newRepInfo, false); err != nil {
			return err
		}
	}
	return nil
}

// When ReplicationConfig Create events occurs, ReplicateNode() replicates the replication node if it exists.
// If not, simply returns.
// Replication occurs based on Source and Destination Kind.
func ReplicateNode(source utils.ReplicationSource, destination utils.ReplicationDestination,
	localClient, remoteClient dynamic.Interface) error {

	var (
		children      utils.Children
		labels        []string
		labelSelector string
	)

	repConfigSpec := utils.ReplicationConfigSpec{Source: source, Destination: destination, Client: remoteClient}

	// If the source kind is "Type", replicate all the objects of that type.
	if source.Kind == utils.Type {
		crdType := utils.GetCrdType(source.Type.Kind, source.Type.Group)

		// Add the entry to the ReplicationEnabledCRDType map.
		utils.ConstructMapReplicationEnabledCRDType(crdType, repConfigSpec)

		gvr := utils.GetGVRFromCrdType(crdType)
		list, err := localClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
		if err != nil && !errors.IsNotFound(err) {
			log.Errorf("Error getting objects of the crd type %v: %v", crdType, err)
			return err
		}

		for _, item := range list.Items {
			err = createObject(gvr, &item, children, item.GetName(), remoteClient, localClient, source, destination)
			if err != nil && !errors.IsAlreadyExists(err) {
				log.Errorf("Resource %v creation failed with an error: %v", item.GetName(), err)
				return err
			}
		}
		return nil
	}

	// If the source kind is "Object", then:
	// 1. If the source is hierarchical, replicate the object and its immediate children.
	// 2. If not, replicate only the object.
	crdType := utils.GetCrdType(source.Object.Kind, source.Object.Group)
	gvr := utils.GetGVRFromCrdType(crdType)
	children = utils.CRDTypeToChildren[crdType]

	// Add the entry to ReplicationEnabledNode Map.
	repObject := utils.GetReplicationObject(source.Object.Group, source.Object.Kind, source.Object.Name)
	utils.ConstructMapReplicationEnabledNode(repObject, repConfigSpec)

	if source.Object.Hierarchical {
		for _, label := range source.Object.Hierarchy.Labels {
			labels = append(labels, label.Key+"="+label.Value)
		}
		labelSelector = strings.Join(labels, ",")
	}

	list, err := localClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil && !errors.IsNotFound(err) {
		log.Errorf("Error getting objects of the crd type %v: %v", crdType, err)
		return err
	}

	for _, item := range list.Items {
		err = createObject(gvr, &item, children, item.GetName(), remoteClient, localClient, source, destination)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Errorf("Resource %v creation failed with an error: %v", item.GetName(), err)
			return err
		}
		if source.Object.Hierarchical {
			if err := ReplicateChildren(crdType, &item, localClient, remoteClient, source, destination); err != nil {
				log.Errorf("Children replication failed for replication object %v: %v", &item, err)
				return err
			}
		}
	}
	return nil
}

// ReplicationChildren replicates the immediate children of desired node.
func ReplicateChildren(crdType string, repEnabledNode *unstructured.Unstructured, localClient, remoteClient dynamic.Interface,
	source utils.ReplicationSource, destination utils.ReplicationDestination) error {

	children := utils.CRDTypeToChildren[crdType]
	parents := utils.CRDTypeToParentHierarchy[crdType]
	labels := repEnabledNode.GetLabels()

	opts := utils.GetNodeLabels(parents, labels, crdType)
	for child := range children {
		parents = utils.CRDTypeToParentHierarchy[child]
		gvr := utils.GetGVRFromCrdType(child)
		c, err := localClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{LabelSelector: opts})
		if err != nil && !errors.IsNotFound(err) {
			log.Errorf("Error getting child objects of %v: %v", repEnabledNode.GetName(), err)
			return err
		}

		for _, item := range c.Items {
			hierarchy := utils.GetNodeHierarchy(parents, item.GetLabels(), child)
			destination.IsChild = true
			err := createObject(gvr, &item, children, hierarchy, remoteClient, localClient, source, destination)
			if err != nil && !errors.IsAlreadyExists(err) {
				log.Errorf("Error creating resource, skipping %v: %v", hierarchy, err)
				return err
			}
		}
	}
	return nil
}
