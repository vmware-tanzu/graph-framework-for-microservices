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

// ResourceRoleBindingReconciler reconciles ResourceRoleBinding object
type ResourceRoleBindingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=resourcerolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=resourcerolebindings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=authorization.nexus.org.authz-controller.com,resources=resourcerolebindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ResourceRoleBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		resourceRoleBinding auth_nexus_org.ResourceRoleBinding
		err                 error
	)

	eventType := utils.Upsert
	if err = r.Get(ctx, req.NamespacedName, &resourceRoleBinding); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch ResourceRoleBinding with name %s", req.Name)
			return ctrl.Result{}, err
		}
		eventType = utils.Delete
	}
	log.Debugf("Received event %s for ResourceRoleBinding: Name %s", eventType, req.Name)

	switch eventType {
	case utils.Upsert:
		user, subjects := constructSubjectsAndUser(resourceRoleBinding.Spec.RoleGvk,
			resourceRoleBinding.Spec.UsersGvk, resourceRoleBinding.Spec.GroupsGvk)

		log.Debugf("ClusterRoleBinding (%q): Subjects: %v and User: %v from ResourceRoleBinding", req.Name, subjects, user)

		existingClusterRoleBinding := rbacv1.ClusterRoleBinding{}
		if err = r.Get(ctx, types.NamespacedName{Name: resourceRoleBinding.Name}, &existingClusterRoleBinding); err != nil {
			if errors.IsNotFound(err) {
				// If it doesn't exist, just create it
				return createClusterRoleBinding(ctx, r.Client, resourceRoleBinding.Kind, resourceRoleBinding.ObjectMeta,
					user, subjects)
			}
			log.Errorf("Failed to get ClusterRoleBinding for the equivalent ResourceRoleBinding %q: %v", resourceRoleBinding.Name, err)
			return ctrl.Result{}, err
		}

		return updateClusterRoleBinding(ctx,
			r.Client,
			existingClusterRoleBinding,
			resourceRoleBinding.Labels,
			resourceRoleBinding.Annotations,
			user,
			subjects,
		)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceRoleBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&auth_nexus_org.ResourceRoleBinding{}).
		Complete(r)
}
