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

	// ReplicationObject to ReplicationConfig Spec (ReplicationObject{Group, Kind, Name} => {"conf1": ReplicationConfigSpec{}, "conf2": ReplicationConfigSpec{})
	ReplicationEnabledNode      = make(map[ReplicationObject]map[string]ReplicationConfigSpec)
	replicationEnabledNodeMutex = &sync.Mutex{}

	// CRDType to ReplicationConfig Spec (root.vmware.org => {"conf1": ReplicationConfigSpec{}, "conf2": ReplicationConfigSpec{})
	ReplicationEnabledCRDType      = make(map[string]map[string]ReplicationConfigSpec)
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

func ConstructMapReplicationEnabledNode(repObj ReplicationObject, name string, spec ReplicationConfigSpec) {
	replicationEnabledNodeMutex.Lock()
	defer replicationEnabledNodeMutex.Unlock()

	if ReplicationEnabledNode[repObj] == nil {
		ReplicationEnabledNode[repObj] = make(map[string]ReplicationConfigSpec)
	}

	ReplicationEnabledNode[repObj][name] = spec
}

func ConstructMapReplicationEnabledCRDType(crdType string, name string, spec ReplicationConfigSpec) {
	replicationEnabledCRDTypeMutex.Lock()
	defer replicationEnabledCRDTypeMutex.Unlock()

	if ReplicationEnabledCRDType[crdType] == nil {
		ReplicationEnabledCRDType[crdType] = make(map[string]ReplicationConfigSpec)
	}

	ReplicationEnabledCRDType[crdType][name] = spec
}
