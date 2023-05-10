package utils

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	/* GVR to ParentInfo.
	Eg: {Group: vmware.org, Version: v1, Resource: configs} => []string{root.vmware.org, project.vmware.org}) */
	GVRToParentHierarchy      = make(map[schema.GroupVersionResource][]string)
	gvrToParentHierarchyMutex = &sync.Mutex{}

	/* GVR to ChildrenInfo.
	Eg: {Group: vmware.org, Version: v1, Resource: roots} => Children{}) */
	GVRToChildren      = make(map[schema.GroupVersionResource]Children)
	gvrToChildrenMutex = &sync.Mutex{}

	/* ReplicationObject to ReplicationConfig Spec
	Eg: {Group: vmware.org, Kind: Root, Name: foo} => {"conf1": ReplicationConfigSpec{}, "conf2": ReplicationConfigSpec{}) */
	ReplicationEnabledNode      = make(map[ReplicationObject]map[string]ReplicationConfigSpec)
	replicationEnabledNodeMutex = &sync.Mutex{}

	/* GVR to ReplicationConfig Spec
	Eg: {Group: vmware.org, Version: v1, Resource: roots} => {"conf1": ReplicationConfigSpec{}, "conf2": ReplicationConfigSpec{}) */
	ReplicationEnabledGVR      = make(map[schema.GroupVersionResource]map[string]ReplicationConfigSpec)
	replicationEnabledGVRMutex = &sync.Mutex{}

	/* CRD Type to CRD version, assuming that there exists CRD of only one version.
	Eg: {roots.vmware.org} => v1 */
	CRDTypeToCrdVersion      = make(map[string]string)
	crdTypeTocrdVersionMutex = &sync.Mutex{}
)

func ConstructCRDTypeToCrdVersion(eventType EventType, crdType, crdVersion string) {
	crdTypeTocrdVersionMutex.Lock()
	defer crdTypeTocrdVersionMutex.Unlock()

	if eventType == Delete {
		delete(CRDTypeToCrdVersion, crdType)
		return
	}

	CRDTypeToCrdVersion[crdType] = crdVersion
}

func ConstructMapGVRToParentHierarchy(eventType EventType, gvr schema.GroupVersionResource, hierarchy []string) {
	gvrToParentHierarchyMutex.Lock()
	defer gvrToParentHierarchyMutex.Unlock()

	if eventType == Delete {
		delete(GVRToParentHierarchy, gvr)
		return
	}

	GVRToParentHierarchy[gvr] = hierarchy
}

func GetParents(gvr schema.GroupVersionResource) []string {
	gvrToParentHierarchyMutex.Lock()
	defer gvrToParentHierarchyMutex.Unlock()
	return GVRToParentHierarchy[gvr]
}

func ConstructMapGVRToChildren(eventType EventType, gvr schema.GroupVersionResource, children Children) {
	gvrToChildrenMutex.Lock()
	defer gvrToChildrenMutex.Unlock()

	if eventType == Delete {
		delete(GVRToChildren, gvr)
		return
	}

	GVRToChildren[gvr] = children
}

func GetChildren(gvr schema.GroupVersionResource) Children {
	gvrToChildrenMutex.Lock()
	defer gvrToChildrenMutex.Unlock()
	return GVRToChildren[gvr]
}

func ConstructMapReplicationEnabledNode(repObj ReplicationObject, name string, spec ReplicationConfigSpec) {
	replicationEnabledNodeMutex.Lock()
	defer replicationEnabledNodeMutex.Unlock()

	if ReplicationEnabledNode[repObj] == nil {
		ReplicationEnabledNode[repObj] = make(map[string]ReplicationConfigSpec)
	}

	ReplicationEnabledNode[repObj][name] = spec
}

func ConstructMapReplicationEnabledGVR(gvr schema.GroupVersionResource, name string, spec ReplicationConfigSpec) {
	replicationEnabledGVRMutex.Lock()
	defer replicationEnabledGVRMutex.Unlock()

	if ReplicationEnabledGVR[gvr] == nil {
		ReplicationEnabledGVR[gvr] = make(map[string]ReplicationConfigSpec)
	}

	ReplicationEnabledGVR[gvr][name] = spec
}
