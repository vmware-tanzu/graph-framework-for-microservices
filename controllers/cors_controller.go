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
	"api-gw/pkg/model"
	"context"

	log "github.com/sirupsen/logrus"
	domain_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/domain.nexus.vmware.com/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// DatamodelReconciler reconciles a Datamodels object
type CORSReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile

func (r *CORSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var CORSConfig domain_nexus_org.CORSConfig
	eventType := model.Upsert
	if err := r.Get(ctx, req.NamespacedName, &CORSConfig); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch cors node with name %s", req.Name)
			return ctrl.Result{}, err
		}
	}
	log.Debugf("Received event %s for cors node: %s", eventType, req.Name)

	model.CorsChan <- model.CorsNodeEvent{Cors: CORSConfig, Type: eventType}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CORSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&domain_nexus_org.CORSConfig{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				delete(model.CorsConfigOrigins, e.Object.GetName())
				_, ok := model.CorsConfigHeaders[e.Object.GetName()]
				if ok {
					delete(model.CorsConfigHeaders, e.Object.GetName())
				}
				model.CorsChan <- model.CorsNodeEvent{Cors: domain_nexus_org.CORSConfig{}, Type: model.Delete}
				return false
			},
		}).
		Complete(r)
}
