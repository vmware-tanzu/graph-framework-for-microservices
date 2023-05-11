package controllers

import (
	"context"
	log "github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"authz-controller/pkg/utils"
	auth_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authorization.nexus.org/v1"
)

const (
	// Kind
	ClusterRoleBinding = "ClusterRoleBinding"
	ClusterRole        = "ClusterRole"

	// RBACAPIVersion
	RBACAPIVersion = "rbac.authorization.k8s.io/v1"

	// RBACAPIGroup
	RBACAPIGroup = "rbac.authorization.k8s.io"

	// Kind
	User  = "User"
	Group = "Group"
)

func metaData(meta *metav1.ObjectMeta, ownerRef metav1.OwnerReference) *metav1.ObjectMeta {
	return &metav1.ObjectMeta{
		Name:            meta.Name,
		Labels:          meta.Labels,
		Annotations:     meta.Annotations,
		OwnerReferences: []metav1.OwnerReference{ownerRef},
	}
}

func createOwnerReference(apiVersion, kind, name string, uid types.UID) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		UID:        uid,
	}
}

func constructClusterRoleBinding(objectMeta metav1.ObjectMeta, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject, ownerRef metav1.OwnerReference) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       ClusterRoleBinding,
			APIVersion: RBACAPIVersion,
		},
		ObjectMeta: *metaData(&objectMeta, ownerRef),
		Subjects:   subjects,
		RoleRef:    roleRef,
	}
}

func notExists(allTypes map[string]bool, value string) bool {
	if _, v := allTypes[value]; !v {
		allTypes[value] = true
		return true
	}

	return false
}

// getResourceTypesForChildren called recursively to hold all the resources and apiGroups.
func getResourceTypesForChildren(crdType string, resourceType *utils.ResourceType,
	allApiGroups, allResources map[string]bool, childrenResourceTypes []string) {
	children := utils.GetChildrenByCRDType(crdType)

	log.Debugf("Children %v of parent %q", children, crdType)
	for k := range children {
		resource, apiGroup := utils.SplitCRDType(k)

		if notExists(allApiGroups, apiGroup) {
			resourceType.APIGroups = append(resourceType.APIGroups, apiGroup)
		}

		if notExists(allResources, resource) {
			resourceType.Resources = append(resourceType.Resources, resource)
		}

		childrenResourceTypes = append(childrenResourceTypes, k)
		getResourceTypesForChildren(k, resourceType, allApiGroups, allResources, childrenResourceTypes)
	}
}

// deleteRoleFromHierarchicalMap contains being called on ResourceRole delete event which deletes the role from the map.
func deleteRoleFromHierarchicalMap(roleName string) {
	utils.RoleToHierarchicalTypesMutex.Lock()
	defer utils.RoleToHierarchicalTypesMutex.Unlock()

	delete(utils.RoleToHierarchicalCRDTypes, roleName)
}

func collectResourceTypes(roleName, crdType, resource, apiGroup string, resourceType *utils.ResourceType) *utils.ResourceType {
	childrenResourceTypes := []string{}
	allApiGroups := map[string]bool{
		apiGroup: true,
	}
	allResources := map[string]bool{
		resource: true,
	}

	getResourceTypesForChildren(crdType, resourceType, allApiGroups, allResources, childrenResourceTypes)
	utils.SetParentCRDTypeToChildren(roleName, crdType, childrenResourceTypes)

	return resourceType
}

// getResourcesTypes holds all the resources and apiGroups based the CRD configured in ResourceRole.
func getResourcesTypes(hierarchical bool, roleName, apiGroup, kind string) *utils.ResourceType {
	resource := utils.GetGroupResourceName(kind)
	resourceType := &utils.ResourceType{
		APIGroups: []string{apiGroup},
		Resources: []string{resource},
	}

	// If hierarchical, collect all the children CRD type of `apiGroup/kind'
	if hierarchical {
		crdType := utils.GetCrdType(kind, apiGroup)
		log.Debugf("Get hierarchical %t resource types of crdType %q for Role: %q", hierarchical, crdType, roleName)
		resourceType = collectResourceTypes(roleName, crdType, resource, apiGroup, resourceType)
	}

	log.Debugf("ResourceType: %v for Role: %q ", resourceType, roleName)

	return resourceType
}

func updateClusterRole(ctx context.Context, client client.Client,
	existingClusterRole rbacv1.ClusterRole, labels, annotations map[string]string,
	rules []rbacv1.PolicyRule) (ctrl.Result, error) {

	existingClusterRole.Labels = labels
	existingClusterRole.Annotations = annotations
	existingClusterRole.Rules = rules

	if err := client.Update(ctx, &existingClusterRole); err != nil {
		log.Errorf("Failed to update ClusterRole (%+v) %v", existingClusterRole, err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func createClusterRole(ctx context.Context, client client.Client, kind string, meta metav1.ObjectMeta,
	rules []rbacv1.PolicyRule) (ctrl.Result, error) {
	ownerRef := createOwnerReference(auth_nexus_org.SchemeGroupVersion.String(),
		kind, meta.Name, meta.UID)
	clusterRole := constructClusterRole(meta, rules, ownerRef)

	if err := client.Create(ctx, &clusterRole); err != nil {
		log.Errorf("Error creating ClusterRole %+v with error %v", clusterRole, err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// convertVerbs converts the list of nexus type verbs to k8s verbs types that apply to ALL the ResourceKinds.
func convertVerbs(verbs []auth_nexus_org.Verb) (toVerbs []string) {
	for _, v := range verbs {
		toVerbs = append(toVerbs, string(v))
	}
	return
}

func constructClusterRole(objectMeta metav1.ObjectMeta, rules []rbacv1.PolicyRule, ownerRef metav1.OwnerReference) rbacv1.ClusterRole {
	return rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       ClusterRole,
			APIVersion: RBACAPIVersion,
		},
		ObjectMeta: *metaData(&objectMeta, ownerRef),
		Rules:      rules,
	}
}

// constructSubjectsAndUser
//1. construct the k8s role binding from nexus role binding that points to the nexus role being created already.
//2. construct the k8s object user identities.
func constructSubjectsAndUser(roleGvk *auth_nexus_org.Link, usersGvk, groupGvk map[string]auth_nexus_org.Link) (rbacv1.RoleRef, []rbacv1.Subject) {
	roleRef := rbacv1.RoleRef{}
	if roleGvk != nil {
		roleRef = rbacv1.RoleRef{
			APIGroup: RBACAPIGroup,
			Kind:     ClusterRole,
			Name:     roleGvk.Name,
		}
	}

	subjects := []rbacv1.Subject{}
	for _, u := range usersGvk {
		subjects = append(subjects, rbacv1.Subject{
			APIGroup: RBACAPIGroup,
			Kind:     User,
			Name:     u.Name,
		})
	}

	for _, u := range groupGvk {
		subjects = append(subjects, rbacv1.Subject{
			APIGroup: RBACAPIGroup,
			Kind:     Group,
			Name:     u.Name,
		})
	}

	return roleRef, subjects
}

func updateClusterRoleBinding(ctx context.Context, client client.Client,
	existingClusterRoleBinding rbacv1.ClusterRoleBinding,
	labels, annotations map[string]string,
	user rbacv1.RoleRef, subjects []rbacv1.Subject) (ctrl.Result, error) {

	existingClusterRoleBinding.Labels = labels
	existingClusterRoleBinding.Annotations = annotations
	existingClusterRoleBinding.Subjects = subjects
	existingClusterRoleBinding.RoleRef = user

	if err := client.Update(ctx, &existingClusterRoleBinding); err != nil {
		log.Errorf("Failed to update clusterRoleBinding (%q) %v", existingClusterRoleBinding.Name, err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func createClusterRoleBinding(ctx context.Context, client client.Client, kind string, meta metav1.ObjectMeta,
	user rbacv1.RoleRef, subjects []rbacv1.Subject) (ctrl.Result, error) {

	ownerRef := createOwnerReference(auth_nexus_org.SchemeGroupVersion.String(),
		kind, meta.Name, meta.UID)
	clusterRoleBinding := constructClusterRoleBinding(meta, user, subjects, ownerRef)

	if err := client.Create(ctx, clusterRoleBinding); err != nil {
		log.Errorf("Error creating ClusterRoleBinding %+v with error %v", clusterRoleBinding, err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// constructPolicyRule construct the information that describes a policy rule.
func constructPolicyRule(resourceRoleName string, hierarchical bool, nexusResourceType auth_nexus_org.ResourceType,
	resourceName []string, verbs []auth_nexus_org.Verb) (rule rbacv1.PolicyRule) {
	resourceType := getResourcesTypes(hierarchical, resourceRoleName, nexusResourceType.Group, nexusResourceType.Kind)
	return rbacv1.PolicyRule{
		Verbs:         convertVerbs(verbs),
		APIGroups:     resourceType.APIGroups,
		Resources:     resourceType.Resources,
		ResourceNames: resourceName,
	}
}
