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
	"api-gw/pkg/common"
	"api-gw/pkg/model"
	"context"

	log "github.com/sirupsen/logrus"

	tenantv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantruntime.nexus.vmware.com/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TenantReconciler reconciles a Datamodels object
type TenantRuntimeReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	K8sclient     kubernetes.Interface
	GrpcConnector *model.ConnectorObject
}

//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile

func (r *TenantRuntimeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var TenantRuntime tenantv1.Tenant

	eventType := model.Upsert
	if err := r.Get(ctx, req.NamespacedName, &TenantRuntime); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch tenant node with name %s", req.Name)
			return ctrl.Result{}, err
		}
		eventType = model.Delete
	}
	log.Debugf("Received event %s for tenantruntime node: %s", eventType, req.Name)
	if eventType == model.Upsert {
		status, msg := common.GetTenantStatus(TenantRuntime.Status.AppStatus)
		log.Infof("Received event %s for tenantruntime node: %s with TenantName %v", eventType, req.Name, TenantRuntime.Labels[common.DISPLAY_NAME])
		tenantState, ok := common.GetTenantState(TenantRuntime.Labels[common.DISPLAY_NAME])
		if ok {
			common.AddTenantState(TenantRuntime.Labels[common.DISPLAY_NAME], common.TenantState{
				Status:        status,
				Message:       msg,
				CreationStart: tenantState.CreationStart,
				SKU:           tenantState.SKU,
			})
		}
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantRuntimeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenantv1.Tenant{}).
		Complete(r)
}
