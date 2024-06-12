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

package managementvmwareorg

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	managementvmwareorgv1 "vmware/build/apis/management.vmware.org/v1"
	nexus_client "vmware/build/nexus-client"
)

// LeaderReconciler reconciles a Leader object
type LeaderReconciler struct {
	client.Client
	NexusClient *nexus_client.Clientset
	Scheme      *runtime.Scheme
}

//+kubebuilder:rbac:groups=management.vmware.org.test-app-local.com,resources=leaders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=management.vmware.org.test-app-local.com,resources=leaders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=management.vmware.org.test-app-local.com,resources=leaders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Leader object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *LeaderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	fmt.Println()
	fmt.Println()
	_ = log.FromContext(ctx)

	fmt.Println("New event for Leader reconciler occured")
	leader, err := r.NexusClient.Management().GetLeaderByName(context.TODO(), req.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Object is removed\n")
			return ctrl.Result{}, nil
		}
		fmt.Printf("Error in getting leader obj: %v\n", err)
		return ctrl.Result{}, err
	}
	fmt.Printf("Received event for leader node: %s\n", leader.DisplayName())

	role, err := leader.GetRole(context.TODO())
	if err != nil {
		fmt.Printf("Error in getting role obj: %v\n", err)
		return ctrl.Result{}, err
	}

	if role == nil {
		fmt.Printf("Current role is nil, updating role to default\n")
		defaultRole, err := r.NexusClient.OrgchartRoot("default").GetExecutiveRole(context.TODO(), "default-executive-role")
		if err != nil {
			fmt.Printf("Error in getting default role: %v\n", err)
			return ctrl.Result{}, err
		}
		err = leader.LinkRole(context.TODO(), defaultRole)
		if err != nil {
			return ctrl.Result{}, nil
		}
		role, err = leader.GetRole(context.TODO())
		if err != nil {
			return ctrl.Result{}, nil
		}
	}
	fmt.Printf("Leader's %v role is: %v\n", leader.DisplayName(), role.DisplayName())
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LeaderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managementvmwareorgv1.Leader{}).
		Complete(r)
}
