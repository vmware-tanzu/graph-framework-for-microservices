package utils

import (
	"sync"
)

var (
	// CRD Type to ParentInfo (config.vmware.org => []string{root.vmware.org, project.vmware.org})
	CRDTypeToParentHierarchy      = make(map[string][]string)
	crdTypeToParentHierarchyMutex = &sync.Mutex{}

	// CRD Type to ChilrenInfo (root.vmware.org => Children{})
	CRDTypeToChildren      = make(map[string]Children)
	crdTypeToChildrenMutex = &sync.Mutex{}

	// ReplicationObject to ReplicationConfig Spec (ReplicationObject{Group, Kind, Name} => ReplicationConfigSpec{})
	ReplicationEnabledNode      = make(map[ReplicationObject]ReplicationConfigSpec)
	replicationEnabledNodeMutex = &sync.Mutex{}

	// CRDType to ReplicationConfig Spec (root.vmware.org => ReplicationConfigSpec{})
	ReplicationEnabledCRDType      = make(map[string]ReplicationConfigSpec)
	replicationEnabledCRDTypeMutex = &sync.Mutex{}
)

func ConstructMapCRDTypeToParentHierarchy(eventType EventType, crdType string, hierarchy []string) {
	crdTypeToParentHierarchyMutex.Lock()
	defer crdTypeToParentHierarchyMutex.Unlock()

	if eventType == Delete {
		delete(CRDTypeToParentHierarchy, crdType)
	}

	CRDTypeToParentHierarchy[crdType] = hierarchy
}

func ConstructMapCRDTypeToChildren(eventType EventType, crdType string, children Children) {
	crdTypeToChildrenMutex.Lock()
	defer crdTypeToChildrenMutex.Unlock()

	if eventType == Delete {
		delete(CRDTypeToChildren, crdType)
	}

	CRDTypeToChildren[crdType] = children
}

func ConstructMapReplicationEnabledNode(repObj ReplicationObject, spec ReplicationConfigSpec) {
	replicationEnabledNodeMutex.Lock()
	defer replicationEnabledNodeMutex.Unlock()

	ReplicationEnabledNode[repObj] = spec
}

func ConstructMapReplicationEnabledCRDType(crdType string, spec ReplicationConfigSpec) {
	replicationEnabledCRDTypeMutex.Lock()
	defer replicationEnabledCRDTypeMutex.Unlock()

	ReplicationEnabledCRDType[crdType] = spec
}
