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

// TODO this logic will have to be reviewed https://jira.eng.vmware.com/browse/NPT-264

//InstanceRoleReconciler reconciles a InstanceRole object
type InstanceRoleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=instanceroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=instanceroles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=instanceroles/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *InstanceRoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		instanceRole auth_nexus_org.InstanceRole
		err          error
	)
	eventType := utils.Upsert
	if err = r.Get(ctx, req.NamespacedName, &instanceRole); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch InstanceRole with name %s", req.Name)
			return ctrl.Result{}, err
		}
		eventType = utils.Delete
	}
	log.Debugf("Received event %s for InstanceRole: Name %s", eventType, req.Name)

	switch eventType {
	case utils.Upsert:
		rules := constructInstanceRolePolicyRules(instanceRole)

		existingClusterRole := rbacv1.ClusterRole{}
		if err = r.Get(ctx, types.NamespacedName{Name: instanceRole.Name}, &existingClusterRole); err != nil {
			if errors.IsNotFound(err) {
				// If it doesn't exist, just create it
				return createClusterRole(ctx, r.Client, instanceRole.Kind, instanceRole.ObjectMeta, rules)
			}
			log.Errorf("Failed to get ClusterRole for the equivalent InstanceRole %q: %v", instanceRole.Name, err)
			return ctrl.Result{}, err
		}

		return updateClusterRole(ctx,
			r.Client,
			existingClusterRole,
			instanceRole.Labels,
			instanceRole.Annotations,
			rules)
	}

	return ctrl.Result{}, nil
}

func constructInstanceRolePolicyRules(instanceRole auth_nexus_org.InstanceRole) []rbacv1.PolicyRule {
	rules := make([]rbacv1.PolicyRule, 0)
	for _, r := range instanceRole.Spec.Rules {
		rule := constructPolicyRule(instanceRole.Name, false, r.Instance.Resource, []string{r.Instance.Name}, r.Verbs)
		rules = append(rules, rule)
	}

	log.Debugf("Policy rules %v for role with name %q", rules, instanceRole.Name)
	return rules
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceRoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&auth_nexus_org.InstanceRole{}).
		Complete(r)
}
