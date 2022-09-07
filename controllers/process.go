package controllers

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"authz-controller/pkg/utils"
)

func (r *CustomResourceDefinitionReconciler) ProcessAnnotation(ctx context.Context, crdType string,
	annotations map[string]string, eventType utils.EventType) error {
	n := utils.NexusAnnotation{}

	if eventType != utils.Delete {
		apiInfo, ok := annotations["nexus"]
		if !ok {
			return nil
		}

		// unmarshall to nexus annotation struct
		err := json.Unmarshal([]byte(apiInfo), &n)
		if err != nil {
			log.Errorf("Error unmarshaling Nexus annotation %v\n", err)
			return err
		}

		log.Debugf("NexusAnnotation %v for crdType %q", n, crdType)
	}

	children := make(map[string]utils.NodeHelperChild)
	if n.Children != nil {
		children = n.Children
	}

	// Store Children information for a given CRD Type.
	utils.ConstructMapCRDTypeToChildren(eventType, crdType, children)

	if eventType == utils.Delete {
		utils.DeleteCRDTypeFromRoleMap(crdType)
		return nil
	}

	// check crd type to be appended to the k8s role rules.
	checkCRDTypeToAddToClusterRole(r.Client, ctx, crdType, n.Hierarchy)

	return nil
}

/*
 checkCRDTypeToAddToClusterRole iterates through all the roles in the map, and it's parent CRDs
 against the crd type from the event notification received. if the crd type parent matched
 with the map value and the child crd type not already in k8s ClusterRole rules, if so append the
 crd type's <resources> and <apiGroups> to the rules and update the k8s ClusterRole.
 Do this recursively ...

 roleName - k8s ClusterRole name
        Ex: root-admin-role

 parentToChildrenCRDMap - indicates the map of parent CRD to list of children
         Ex: map[roots.root.helloworld.com] = []Children{
                                                     projects.project.helloworld.com,
                                                     configs.config.helloworld.com
                                                    }

 childrenCRDTypes - indicates all the children CRDs which is captured during the ResourceRole event handler.
*/
func checkCRDTypeToAddToClusterRole(client client.Client, ctx context.Context, crdType string, hierarchy []string) {
	for roleName, parentToChildrenCRDMap := range utils.GetRolesToHierarchicalCRDMap() {
		var (
			role          rbacv1.ClusterRole
			foundExisting bool
		)

		log.Debugf("Check CRD type %q to be added to the Role %q", crdType, roleName)

		if err := client.Get(ctx, types.NamespacedName{Name: roleName}, &role); err != nil {
			log.Errorf("Error getting role object: %q with err %v", roleName, err)
			continue
		}

		for _, parent := range hierarchy {
			childrenCRDTypes, exists := parentToChildrenCRDMap[parent] // eg parent = roots.root.helloworld.com
			if !exists {
				log.Debugf("Parent %q not matched, skip: %q", parent, roleName)
				// RoleMap doesn't contain the parent from hierarchy list.
				continue
			}

			// if crd type is newly created, not added already to the resource during the ResourceRole event handler,
			// to be added during the CRD event notification.
			foundExisting = utils.ContainsString(childrenCRDTypes, crdType)
			if !foundExisting {
				log.Debugf("CRD type %q should be added to role rules %v", crdType, childrenCRDTypes)
				// bind it to the rules
				for i := range role.Rules {
					existingRule := &role.Rules[i]

					log.Debugf("Parent %q and Existing Rules %v", parent, existingRule)
					// if parent exists in the list, then append the child to it.
					if utils.ParentExists(existingRule.Resources, existingRule.APIGroups, parent) {
						childResource, childApiGroup := utils.SplitCRDType(crdType)
						if !utils.ContainsString(existingRule.Resources, childResource) {
							existingRule.Resources = append(existingRule.Resources, childResource)
							log.Debugf("Appended resources %v", existingRule.Resources)
						}
						if !utils.ContainsString(existingRule.APIGroups, childApiGroup) {
							existingRule.APIGroups = append(existingRule.APIGroups, childApiGroup)
							log.Debugf("Appended APIGroups %v", existingRule.APIGroups)
						}
					}
				}
			}
		}

		if !foundExisting {
			if err := client.Update(ctx, &role); err != nil {
				log.Errorf("Couldn't update role %q with CRD %q", role.Name, crdType)
				continue
			}
		}
	}
}
