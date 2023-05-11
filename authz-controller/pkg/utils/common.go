package utils

import (
	"strings"
	"sync"
)

type EventType string

const (
	Upsert EventType = "Upsert"
	Delete EventType = "Delete"

	DISPLAY_NAME_LABEL = "nexus/display_name"
)

var (
	CRDTypeToChildren      = make(map[string]Children)
	crdTypeToChildrenMutex = &sync.Mutex{}

	RoleToHierarchicalCRDTypes   = make(map[string]map[string][]string)
	RoleToHierarchicalTypesMutex = &sync.Mutex{}
)

type NexusAnnotation struct {
	Name      string                     `json:"name,omitempty"`
	Hierarchy []string                   `json:"hierarchy,omitempty"`
	Children  map[string]NodeHelperChild `json:"children,omitempty"`
}

type Children map[string]NodeHelperChild

type NodeHelperChild struct {
	FieldName    string `json:"fieldName"`
	FieldNameGvk string `json:"fieldNameGvk"`
	IsNamed      bool   `json:"isNamed"`
}

// ResourceType : This model used to carry APIGroups and
//resource names for all the CRD
type ResourceType struct {
	APIGroups []string `json:"group"`
	Resources []string `json:"kind"`
}

// GetChildrenByCRDType returns all the children of that parent CRD type
func GetChildrenByCRDType(crdType string) Children {
	crdTypeToChildrenMutex.Lock()
	defer crdTypeToChildrenMutex.Unlock()

	return CRDTypeToChildren[crdType]
}

func ConstructMapCRDTypeToChildren(eventType EventType, crdType string, children Children) {
	crdTypeToChildrenMutex.Lock()
	defer crdTypeToChildrenMutex.Unlock()

	if eventType == Delete {
		delete(CRDTypeToChildren, crdType)
		return
	}

	CRDTypeToChildren[crdType] = children
}

// SetParentCRDTypeToChildren set the children crd types to the role name.
// that is useful in case of any new CRD event must be appended to the role rules.
func SetParentCRDTypeToChildren(roleName, crdType string, childrenResourceTypes []string) {
	RoleToHierarchicalTypesMutex.Lock()
	defer RoleToHierarchicalTypesMutex.Unlock()

	if RoleToHierarchicalCRDTypes[roleName] == nil {
		RoleToHierarchicalCRDTypes[roleName] = make(map[string][]string)
	}

	RoleToHierarchicalCRDTypes[roleName][crdType] = childrenResourceTypes
}

func GetRolesToHierarchicalCRDMap() map[string]map[string][]string {
	RoleToHierarchicalTypesMutex.Lock()
	defer RoleToHierarchicalTypesMutex.Unlock()

	copyParentCRDTypeToChildrenTypes := make(map[string]map[string][]string)
	for role, parentInfo := range RoleToHierarchicalCRDTypes {
		copyParentCRDTypeToChildrenTypes[role] = parentInfo
	}

	return copyParentCRDTypeToChildrenTypes
}

func DeleteCRDTypeFromRoleMap(crdType string) {
	RoleToHierarchicalTypesMutex.Lock()
	defer RoleToHierarchicalTypesMutex.Unlock()

	for _, crdTypesMap := range RoleToHierarchicalCRDTypes {
		delete(crdTypesMap, crdType)
	}
}

// SplitCRDType splits the <roots.root.helloworld.com> into
// resource: <roots> and group: <root.helloworld.com>
func SplitCRDType(crdType string) (string, string) {
	parts := strings.Split(crdType, ".")
	apiGroup := strings.Join(parts[1:], ".")

	return parts[0], apiGroup
}

func GetCrdType(kind, groupName string) string {
	return GetGroupResourceName(kind) + "." + groupName // eg roots.root.helloworld.com
}

func GetGroupResourceName(kind string) string {
	return strings.ToLower(ToPlural(kind)) // eg roots
}

func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ParentExists(resources, apiGroups []string, parent string) bool {
	parentResource, parentApiGroup := SplitCRDType(parent)

	if ContainsString(resources, parentResource) &&
		ContainsString(apiGroups, parentApiGroup) {
		return true
	}

	return false
}
