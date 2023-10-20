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
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	networkingAPI "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	routenexusorgv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/route.nexus.vmware.com/v1"
)

// RouteReconciler reconciles a Route object
type RouteReconciler struct {
	client.Client
	BaseClient            client.Client
	Scheme                *runtime.Scheme
	BaseNamespace         string
	IngressControllerName string
	DefaultBackend        string
	DefaultBackendPort    int32
}

const LocalAPIServerNamespace = "default"
const apiServiceNameIngress = "nginx-ingress"
const customAPIVersion = "v1"

func (r *RouteReconciler) DeleteBaseObjects(baseObjectName string) error {
	ingressList := &networkingAPI.IngressList{}
	ingressListopts := []client.ListOption{
		client.InNamespace(r.BaseNamespace),
		client.MatchingLabels{"baseExtensionObject": baseObjectName},
	}
	r.BaseClient.List(context.Background(), ingressList, ingressListopts...)
	for _, ingressToDelete := range ingressList.Items {
		fmt.Printf("ingress object for deletion: %s\n", ingressToDelete.Name)
		if err := r.BaseClient.Delete(context.Background(), &ingressToDelete); err != nil {
			if !errors.IsNotFound(err) {
				ctrl.Log.Error(err, fmt.Sprintf("could not delete ingress object  %s due to", ingressToDelete.Name))
				return err
			}
		}

	}
	// local client objects
	apiServiceList := apiregistrationv1.APIServiceList{}
	apiServiceListOpts := []client.ListOption{
		client.MatchingLabels{"baseExtensionObject": baseObjectName},
	}
	r.Client.List(context.Background(), &apiServiceList, apiServiceListOpts...)
	for _, apiserviceToDelete := range apiServiceList.Items {
		fmt.Printf("apiservice object for deletion: %s\n", apiserviceToDelete.Name)
		if err := r.Client.Delete(context.Background(), &apiserviceToDelete); err != nil {
			if !errors.IsNotFound(err) {
				ctrl.Log.Error(err, fmt.Sprintf("could not delete apiservice object %s due to", apiserviceToDelete.Name))
				return err
			}
		}
	}
	return nil
}

//+kubebuilder:rbac:groups=route.nexus.vmware.com.api-operator.com,resources=routes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=route.nexus.vmware.com.api-operator.com,resources=routes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=route.nexus.vmware.com.api-operator.com,resources=routes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile

func (r *RouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var apiExtensionsobject routenexusorgv1.Route
	baseObjectName := req.NamespacedName.Name
	if err := r.Get(ctx, req.NamespacedName, &apiExtensionsobject); err != nil {
		if errors.IsNotFound(err) {
			err = r.DeleteBaseObjects(baseObjectName)
			if err != nil {
				ctrl.Log.Error(err, "Could not delete ingress/ apiservice objects")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		ctrl.Log.Error(err, fmt.Sprintf("Could not get %s apiextension object due to", baseObjectName))
		return ctrl.Result{}, err
	}
	ctrl.Log.Info(fmt.Sprintf("Recieved create event for %s", baseObjectName))
	var portNumber int32
	ingressClass := fmt.Sprintf("nginx-%s", r.BaseNamespace)
	ingressName := apiExtensionsobject.Spec.Resource.Name
	uriPath := apiExtensionsobject.Spec.Uri
	serviceName := apiExtensionsobject.Spec.Service.Name
	// Create Port object to use default value if not provided.
	portObject := apiExtensionsobject.Spec.Service.Port
	if &portObject == nil {
		portNumber = 80
	} else {
		portNumber = apiExtensionsobject.Spec.Service.Port
	}

	port := networkingAPI.ServiceBackendPort{
		Number: portNumber,
	}

	// Add Https annotation if backend is https service
	ingressAnnotations := map[string]string{
		"nginx.ingress.kubernetes.io/rewrite-target": "/$1",
	}
	if strings.ToLower(apiExtensionsobject.Spec.Service.Scheme) == "https" {
		ingressAnnotations["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
	}
	Prefix := networkingAPI.PathTypePrefix
	Exact := networkingAPI.PathTypeExact
	//Base object label
	baseObjectLabel := map[string]string{
		"baseExtensionObject": baseObjectName,
	}
	PathsConfig := []networkingAPI.HTTPIngressPath{
		{
			Path:     fmt.Sprintf("/apis/%s/%s/?(.+)", ingressName, customAPIVersion),
			PathType: &Prefix,
			Backend: networkingAPI.IngressBackend{
				Service: &networkingAPI.IngressServiceBackend{
					Name: serviceName,
					Port: port,
				},
			},
		},
	}
	if uriPath == "/" {
		PathsConfig = append(PathsConfig, networkingAPI.HTTPIngressPath{
			Path:     fmt.Sprintf("/apis/%s/%s", ingressName, customAPIVersion),
			PathType: &Exact,
			Backend: networkingAPI.IngressBackend{
				Service: &networkingAPI.IngressServiceBackend{
					Name: serviceName,
					Port: port,
				},
			},
		})
	} else {
		defaultPort := networkingAPI.ServiceBackendPort{
			Number: r.DefaultBackendPort,
		}
		PathsConfig = append(PathsConfig, networkingAPI.HTTPIngressPath{
			Path:     fmt.Sprintf("/apis/%s/%s", ingressName, customAPIVersion),
			PathType: &Exact,
			Backend: networkingAPI.IngressBackend{
				Service: &networkingAPI.IngressServiceBackend{
					Name: r.DefaultBackend,
					Port: defaultPort,
				},
			},
		})
	}

	ingressObject := networkingAPI.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: ingressAnnotations,
			Namespace:   r.BaseNamespace,
			Name:        ingressName,
			Labels:      baseObjectLabel,
		},
		Spec: networkingAPI.IngressSpec{
			IngressClassName: &ingressClass,
			Rules: []networkingAPI.IngressRule{
				{
					IngressRuleValue: networkingAPI.IngressRuleValue{
						HTTP: &networkingAPI.HTTPIngressRuleValue{
							Paths: PathsConfig,
						},
					},
				},
			},
		},
	}
	fmt.Printf("Creating ingress object: %v\n", ingressObject)
	if err := r.BaseClient.Get(context.Background(), types.NamespacedName{
		Namespace: ingressObject.Namespace,
		Name:      ingressObject.Name,
	}, &ingressObject); err != nil {
		if errors.IsNotFound(err) {
			if err := r.BaseClient.Create(context.Background(), &ingressObject); err != nil {
				ctrl.Log.Error(err, fmt.Sprintf("could not create the ingress object %s on base cluster", ingressObject.Name))
				return ctrl.Result{}, err
			}
		} else {
			ctrl.Log.Error(err, fmt.Sprintf("could not get the ingress object %s on base cluster", ingressObject.Name))
			return ctrl.Result{}, err
		}
	}
	if err := r.BaseClient.Update(context.Background(), &ingressObject); err != nil {
		ctrl.Log.Error(err, fmt.Sprintf("could not update the ingress object %s on base cluster", ingressObject.Name))
		return ctrl.Result{}, err
	}

	//service object creation if not present...
	serviceObject := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiServiceNameIngress,
			Namespace: LocalAPIServerNamespace,
		},
		Spec: corev1.ServiceSpec{
			ExternalName: fmt.Sprintf("%s.%s.svc", r.IngressControllerName, r.BaseNamespace),
			Type:         corev1.ServiceTypeExternalName,
		},
	}
	if err := r.Client.Get(context.Background(), types.NamespacedName{
		Namespace: serviceObject.Namespace,
		Name:      serviceObject.Name,
	}, &serviceObject); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Client.Create(context.Background(), &serviceObject); err != nil {
				ctrl.Log.Error(err, fmt.Sprintf("could not internal service object %s", apiServiceNameIngress))
				return ctrl.Result{}, err
			}
		} else {
			ctrl.Log.Error(err, fmt.Sprintf("could not get the internal service object %s on base cluster", apiServiceNameIngress))
			return ctrl.Result{}, err
		}
	}

	apiServiceObject := apiregistrationv1.APIService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s.%s", customAPIVersion, ingressName),
			Namespace: LocalAPIServerNamespace,
			Labels:    baseObjectLabel,
		},
		Spec: apiregistrationv1.APIServiceSpec{
			GroupPriorityMinimum: 100,
			Service: &apiregistrationv1.ServiceReference{
				Namespace: LocalAPIServerNamespace,
				Name:      apiServiceNameIngress,
			},
			Group:                 ingressName,
			Version:               customAPIVersion,
			VersionPriority:       100,
			InsecureSkipTLSVerify: true,
		},
	}
	err := r.Client.Get(context.Background(), types.NamespacedName{
		Name:      apiServiceObject.Name,
		Namespace: apiServiceObject.Namespace}, &apiServiceObject)
	if err != nil {
		if errors.IsNotFound(err) {
			if err := r.Client.Create(context.Background(), &apiServiceObject); err != nil {
				ctrl.Log.Error(err, "could not create apiservice")
				return ctrl.Result{}, err
			}
		} else {
			ctrl.Log.Error(err, fmt.Sprintf("could not get the internal apiservice object %s on base cluster", apiServiceObject.Name))
			return ctrl.Result{}, err
		}
	}
	if err := r.Client.Update(context.Background(), &apiServiceObject); err != nil {
		ctrl.Log.Error(err, "could not update apiservice")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&routenexusorgv1.Route{}).
		Complete(r)
}
