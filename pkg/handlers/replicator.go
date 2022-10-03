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

	"connector/controllers"
	"connector/pkg/utils"
)

func createObject(gvr schema.GroupVersionResource, res *unstructured.Unstructured, children utils.Children,
	hierarchy string, rc utils.ReplicationConfigSpec) error {

	labels := make(map[string]string)
	for key, val := range res.GetLabels() {
		labels[key] = val
	}

	if rc.Source.Kind == utils.Object && !rc.Destination.Hierarchical && children != nil {
		// Ignore relationships.
		spec, ok := res.UnstructuredContent()["spec"]
		if ok {
			utils.DeleteChildGvkFields(spec.(map[string]interface{}), children)
			res.UnstructuredContent()["spec"] = spec
		}
	}

	if rc.Destination.Hierarchical {
		for _, label := range rc.Destination.Hierarchy.Labels {
			labels[label.Key] = label.Value
		}
	}

	annotations := res.GetAnnotations()
	annotations = utils.GenerateAnnotations(annotations, gvr, res.GetName())

	// If destination object type is specified, then replicate the object to that type.
	destGvr, destKind := utils.GetDestinationGvrAndKind(rc.Destination, gvr, res.GetKind())

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": destGvr.GroupVersion().String(),
			"kind":       destKind,
			"metadata": map[string]interface{}{
				"name":        res.GetName(),
				"labels":      labels,
				"namespace":   rc.Destination.Namespace,
				"annotations": annotations,
			},
		},
	}
	spec, ok := res.UnstructuredContent()["spec"]
	if ok {
		obj.Object["spec"] = spec
	}

	if destObj, err := rc.RemoteClient.Resource(destGvr).Namespace(rc.Destination.Namespace).Get(context.TODO(), obj.GetName(),
		metav1.GetOptions{}); destObj != nil && err == nil {
		return nil
	}

	log.Infof("Replication is enabled for the resource: %v, creating...", hierarchy)
	// If the object was successfully replicated, we need to patch the source and remote generation ID.
	destObj, err := rc.RemoteClient.Resource(destGvr).Namespace(rc.Destination.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	if err == nil && destObj != nil {
		if name, err := patchStatusObject(rc, gvr, destGvr, rc.Destination.Namespace, res.GetName(), destObj.GetName(),
			res.GetGeneration(), destObj.GetGeneration()); err != nil {
			log.Errorf("CR %q status patch failed with an error: %v", name, err)
			return err
		}
	}
	return err
}

func updateObject(gvr schema.GroupVersionResource, res *unstructured.Unstructured, children utils.Children,
	hierarchy string, rc utils.ReplicationConfigSpec) error {

	// If destination object type is specified, then replicate the object to that type.
	destGvr, _ := utils.GetDestinationGvrAndKind(rc.Destination, gvr, "")

	oldObject, err := rc.RemoteClient.Resource(destGvr).Namespace(rc.Destination.Namespace).Get(context.TODO(), res.GetName(), metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		log.Infof("Resource %s not found, creating instead", res.GetName())
		return createObject(gvr, res, children, hierarchy, rc)
	}

	spec, ok := res.UnstructuredContent()["spec"]
	if rc.Source.Object.Hierarchical && !rc.Destination.Hierarchical && ok {
		// Ignore relationships.
		utils.DeleteChildGvkFields(spec.(map[string]interface{}), children)
	}
	oldObject.UnstructuredContent()["spec"] = spec

	annotations := oldObject.GetAnnotations()
	annotations = utils.GenerateAnnotations(annotations, gvr, res.GetName())
	oldObject.SetAnnotations(annotations)

	log.Infof("Replication is enabled for the resource: %v, updating...", hierarchy)
	// If the object was successfully replicated, we need to patch the source and remote generation ID.
	destObj, err := rc.RemoteClient.Resource(destGvr).Namespace(rc.Destination.Namespace).Update(context.TODO(), oldObject, metav1.UpdateOptions{})
	if err == nil && destObj != nil {
		if name, err := patchStatusObject(rc, gvr, destGvr, rc.Destination.Namespace, res.GetName(), destObj.GetName(),
			res.GetGeneration(), destObj.GetGeneration()); err != nil {
			log.Errorf("CR %q status patch failed with an error: %v", name, err)
			return err
		}
	}
	return err
}

func patchStatusObject(rc utils.ReplicationConfigSpec, srcGvr, destGvr schema.GroupVersionResource,
	ns, srcObj, destObj string, srcGeneration, destGeneration int64) (string, error) {

	var (
		objName, namespace string
		client             dynamic.Interface
		gvr                schema.GroupVersionResource
	)

	switch rc.StatusEndpoint {
	case utils.Source:
		objName = srcObj
		client = rc.LocalClient
		gvr = srcGvr

	case utils.Destination:
		objName = destObj
		namespace = ns
		client = rc.RemoteClient
		gvr = destGvr

	default:
		return "", nil
	}

	patchBytes, err := utils.CreatePatch(srcGeneration, destGeneration)
	if err != nil {
		log.Errorf("Could not create patch for the CR(%q) %v", objName, err)
		return objName, err
	}

	log.Debugf("Patching status of CR %q: %v", objName, string(patchBytes))

	_, err = client.Resource(gvr).Namespace(namespace).Patch(context.TODO(), objName, types.MergePatchType, patchBytes, metav1.PatchOptions{}, "status")
	if err != nil {
		log.Errorf("Resource %s patching failed with an error: %v", objName, err)
		return objName, err
	}
	return objName, nil
}

func processEvents(h *RemoteHandler, eventType, hierarchy string, res *unstructured.Unstructured, children utils.Children,
	repInfo utils.ReplicationConfigSpec, replicationEnabledNode bool) error {

	switch eventType {
	case utils.Create:
		err := createObject(h.Gvr, res, children, hierarchy, repInfo)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Errorf("Resource %v create failed with an error: %v", hierarchy, err)
			return err
		}
		if repInfo.Source.Object.Hierarchical && replicationEnabledNode {
			// Replicate children only if the obj exactly matches the source object.
			newRepInfo := repInfo
			newRepInfo.Destination.IsChild = true
			if err := ReplicateChildren(h.Gvr, res, newRepInfo); err != nil {
				log.Errorf("Children replication failed for the resource %v: %v", hierarchy, err)
				return err
			}
		}
	case utils.Update:
		err := updateObject(h.Gvr, res, children, hierarchy, repInfo)
		if err != nil {
			log.Errorf("Resource %v update failed with an error: %v", hierarchy, err)
			return err
		}
	}
	return nil
}

// For every object notification received, replicator replicates in the following order.
// 1. If the ResourceType is of interest, simply replicate the object.
// 2. If only the object is of interest and if it is a part of a graph, then replicate the object and its immediate children.
// 3. If the oject is not a part of a graph, then replicate only that object.
func Replicator(obj interface{}, h *RemoteHandler, eventType string) error {
	res := obj.(*unstructured.Unstructured)

	var (
		rc        utils.ReplicationConfigSpec
		repObj    utils.ReplicationObject
		hierarchy string
	)

	// Verify if the type is of interest.
	// If the ResourceType is enabled for replication, simply replicate the object.
	if repConfMap, replicationEnabledResourceType := utils.ReplicationEnabledGVR[h.Gvr]; replicationEnabledResourceType {
		for _, rc = range repConfMap {
			if err := processEvents(h, eventType, res.GetName(), res, nil, rc, false); err != nil {
				return err
			}
		}
		return nil
	}

	labels := res.GetLabels()
	parents := utils.GVRToParentHierarchy[h.Gvr]
	children := utils.GVRToChildren[h.Gvr]

	name := res.GetName()
	if objName, ok := labels[utils.DisplayNameKey]; ok {
		name = objName
	}
	repObj = utils.GetReplicationObject(res.GroupVersionKind().Group, res.GetKind(), name)

	// Verify if the obj is of interest.
	repConfMap, replicationEnabledNode := utils.ReplicationEnabledNode[repObj]
	if replicationEnabledNode {
		for _, rc = range repConfMap {
			hierarchy = name
			if parents != nil {
				hierarchy = utils.GetNodeHierarchy(parents, labels, strings.Join([]string{h.Gvr.Resource, h.Gvr.Group}, "."))
			}
			if err := processEvents(h, eventType, hierarchy, res, children, rc, true); err != nil {
				return err
			}
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
	gvr := utils.GetGVRFromCrdType(immediateParent, utils.CRDTypeToCrdVersion[immediateParent])

	// Get parent information from the object's labels to verify if the object's immediate parent is of interest.
	opts := metav1.ListOptions{LabelSelector: utils.GetParentLabels(parents, labels)}
	c, err := h.LocalClient.Resource(gvr).List(context.TODO(), opts)
	if err != nil {
		log.Errorf("error getting child objects of %v: %v", repObj.Name, err)
		return err
	}

	for _, item := range c.Items {
		name := item.GetName()
		if objName, ok := item.GetLabels()[utils.DisplayNameKey]; ok {
			name = objName
		}
		repObj := utils.GetReplicationObject(item.GroupVersionKind().Group, item.GetKind(), name)
		if repInfoMap, replicationEnabledNode := utils.ReplicationEnabledNode[repObj]; replicationEnabledNode {
			for _, rc = range repInfoMap {
				hierarchy = name
				if rc.Source.Object.Hierarchical {
					hierarchy = utils.GetNodeHierarchy(parents, labels, strings.Join([]string{h.Gvr.Resource, h.Gvr.Group}, "."))
				}
				replicate = true
				break
			}
		}
	}
	if replicate {
		newRepInfo := rc
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
func ReplicateNode(confName string, rc utils.ReplicationConfigSpec) error {

	var (
		children      utils.Children
		labels        []string
		labelSelector string
	)

	// If the source kind is "Type", replicate all the objects of that type.
	if rc.Source.Kind == utils.Type {
		t := rc.Source.Type
		gvr := schema.GroupVersionResource{Group: t.Group, Version: t.Version, Resource: utils.GetGroupResourceName(t.Kind)}

		// Start controller if not started.
		controllers.GvrCh <- gvr

		// Add the entry to the ReplicationEnabledGVR map.
		utils.ConstructMapReplicationEnabledGVR(gvr, confName, rc)

		if err := ReplicateAllObjectsOfType(gvr, rc, children); err != nil {
			return err
		}
		return nil
	}

	// If the source kind is "Object", then:
	// 1. If the source is hierarchical, replicate the object and its immediate children.
	// 2. If not, replicate only the object.
	t := rc.Source.Object
	gvr := schema.GroupVersionResource{Group: t.Group, Version: t.Version, Resource: utils.GetGroupResourceName(t.Kind)}
	children = utils.GVRToChildren[gvr]

	// Start controller if not started.
	controllers.GvrCh <- gvr

	// Add the entry to ReplicationEnabledNode Map.
	repObject := utils.GetReplicationObject(rc.Source.Object.Group, rc.Source.Object.Kind, rc.Source.Object.Name)
	utils.ConstructMapReplicationEnabledNode(repObject, confName, rc)

	if rc.Source.Object.Hierarchical {
		for _, label := range rc.Source.Object.Hierarchy.Labels {
			labels = append(labels, label.Key+"="+label.Value)
		}
		labelSelector = strings.Join(labels, ",")
	}

	list, err := rc.LocalClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		log.Errorf("Error getting objects for gvr %v: %v", gvr, err)
		return err
	}

	for _, item := range list.Items {
		if !rc.Source.Object.Hierarchical {
			if rc.Source.Object.Name != item.GetName() {
				continue
			}
		}
		err = createObject(gvr, &item, children, item.GetName(), rc)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Errorf("Resource %v creation failed with an error: %v", item.GetName(), err)
			return err
		}
		if rc.Source.Object.Hierarchical {
			if err := ReplicateChildren(gvr, &item, rc); err != nil {
				log.Errorf("Children replication failed for replication object %v: %v", &item, err)
				return err
			}
		}
	}
	return nil
}

// ReplicationChildren replicates the immediate children of desired node.
func ReplicateChildren(gvr schema.GroupVersionResource, repEnabledNode *unstructured.Unstructured, rc utils.ReplicationConfigSpec) error {

	children := utils.GVRToChildren[gvr]
	parents := utils.GVRToParentHierarchy[gvr]
	labels := repEnabledNode.GetLabels()

	opts := utils.GetNodeLabels(parents, labels, strings.Join([]string{gvr.Resource, gvr.Group}, "."))
	for child := range children {
		childGvr := utils.GetGVRFromCrdType(child, utils.CRDTypeToCrdVersion[child])
		parents = utils.GVRToParentHierarchy[childGvr]
		c, err := rc.LocalClient.Resource(childGvr).List(context.TODO(), metav1.ListOptions{LabelSelector: opts})
		if err != nil && !errors.IsNotFound(err) {
			log.Errorf("Error getting child objects of %v: %v", repEnabledNode.GetName(), err)
			return err
		}

		for _, item := range c.Items {
			hierarchy := utils.GetNodeHierarchy(parents, item.GetLabels(), child)
			rc.Destination.IsChild = true
			err := createObject(childGvr, &item, children, hierarchy, rc)
			if err != nil && !errors.IsAlreadyExists(err) {
				log.Errorf("Error creating resource, skipping %v: %v", hierarchy, err)
				return err
			}
		}
	}
	return nil
}

func ReplicateAllObjectsOfType(gvr schema.GroupVersionResource, rc utils.ReplicationConfigSpec, children utils.Children) error {
	list, err := rc.LocalClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("Error getting objects for gvr %v: %v", gvr, err)
		return err
	}

	for _, item := range list.Items {
		err = createObject(gvr, &item, children, item.GetName(), rc)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Errorf("Resource %v creation failed with an error: %v", item.GetName(), err)
			return err
		}
	}
	return nil
}
