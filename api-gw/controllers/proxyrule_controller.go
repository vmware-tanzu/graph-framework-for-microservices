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
	"api-gw/pkg/envoy"
	"api-gw/pkg/model"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	adminnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/admin.nexus.vmware.com/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProxyRuleReconciler reconciles a ProxyRule object
type ProxyRuleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=admin.nexus.vmware.com.api-gw.com,resources=proxyrules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=admin.nexus.vmware.com.api-gw.com,resources=proxyrules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=admin.nexus.vmware.com.api-gw.com,resources=proxyrules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ProxyRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var proxyRule adminnexusorgv1.ProxyRule
	eventType := model.Upsert
	// TODO use the shim layer to fetch the object
	if err := r.Get(ctx, req.NamespacedName, &proxyRule); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Errorf("Error while trying to fetch ProxyRule node with name %s: %s", req.Name, err)
			return ctrl.Result{}, err
		}
		eventType = model.Delete
	}
	log.Debugf("Received event %s for ProxyRule node: Name=%s", eventType, req.NamespacedName.Name)

	switch eventType {
	case model.Delete:
		// calling delete on both since we don't know the match type of the deleted object
		// is there a way to get the spec of the deleted object?
		err := envoy.DeleteUpstream(req.NamespacedName.Name)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error deleting envoy upstream: %s", err)
		}
		err = envoy.DeleteHeaderUpstream(req.NamespacedName.Name)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error deleting envoy header upstream: %s", err)
		}
		log.Debugf("deleted proxy rule %s", req.NamespacedName.Name)
	case model.Upsert:
		if proxyRule.Spec.MatchCondition.Type == "jwt" {
			err := envoy.AddUpstream(req.NamespacedName.Name, &envoy.UpstreamConfig{
				Name:          proxyRule.Name,
				Host:          proxyRule.Spec.Upstream.Host,
				Port:          proxyRule.Spec.Upstream.Port,
				JwtClaimKey:   proxyRule.Spec.MatchCondition.Key,
				JwtClaimValue: proxyRule.Spec.MatchCondition.Value,
			})
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("error adding envoy upstream: %s", err)
			}
		} else if proxyRule.Spec.MatchCondition.Type == "header" {
			err := envoy.AddHeaderUpstream(req.NamespacedName.Name, &envoy.HeaderMatchedUpstream{
				Name:        proxyRule.Name,
				Host:        proxyRule.Spec.Upstream.Host,
				Port:        proxyRule.Spec.Upstream.Port,
				HeaderName:  proxyRule.Spec.MatchCondition.Key,
				HeaderValue: proxyRule.Spec.MatchCondition.Value,
			})
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("error adding envoy header upstream: %s", err)
			}
		} else {
			return ctrl.Result{}, fmt.Errorf("match type %s not supported", proxyRule.Spec.MatchCondition.Type)
		}
		log.Debugf("updated proxy rule %s", proxyRule.Name)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProxyRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&adminnexusorgv1.ProxyRule{}).
		Complete(r)
}
