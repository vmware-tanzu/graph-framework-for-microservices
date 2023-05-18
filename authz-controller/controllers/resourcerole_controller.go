/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	log "github.com/sirupsen/logrus"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"authz-controller/pkg/utils"
	auth_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authorization.nexus.org/v1"
)

// ResourceRoleReconciler reconciles a ResourceRole object
type ResourceRoleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=resourceroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=resourceroles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=resourceroles/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ResourceRoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		resourceRole auth_nexus_org.ResourceRole
		err          error
	)

	eventType := utils.Upsert
	if err = r.Get(ctx, req.NamespacedName, &resourceRole); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch ResourceRole with name %s", req.Name)
			return ctrl.Result{}, err
		}
		eventType = utils.Delete
	}
	log.Debugf("Received event %s for ResourceRole: Name %s", eventType, req.Name)

	switch eventType {
	case utils.Upsert:
		rules := constructResourceRolePolicyRules(resourceRole)

		existingClusterRole := rbacv1.ClusterRole{}
		if err = r.Get(ctx, types.NamespacedName{Name: resourceRole.Name}, &existingClusterRole); err != nil {
			if errors.IsNotFound(err) {
				// If it doesn't exist, just create it
				return createClusterRole(ctx, r.Client, resourceRole.Kind, resourceRole.ObjectMeta, rules)
			}
			log.Errorf("Failed to get ClusterRole for the equivalent ResourceRole %q: %v", resourceRole.Name, err)
			return ctrl.Result{}, err
		}

		return updateClusterRole(ctx,
			r.Client,
			existingClusterRole,
			resourceRole.Labels,
			resourceRole.Annotations,
			rules)

	case utils.Delete:
		deleteRoleFromHierarchicalMap(resourceRole.Name)
	}

	return ctrl.Result{}, nil
}

func constructResourceRolePolicyRules(resourceRole auth_nexus_org.ResourceRole) (rules []rbacv1.PolicyRule) {
	deleteRoleFromHierarchicalMap(resourceRole.Name)
	for _, r := range resourceRole.Spec.Rules {
		rule := constructPolicyRule(resourceRole.Name, r.Hierarchical, r.Resource, nil, r.Verbs)
		rules = append(rules, rule)
	}

	log.Debugf("Policy rules %v for role with name %q", rules, resourceRole.Name)
	return rules
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceRoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&auth_nexus_org.ResourceRole{}).
		Complete(r)
}
