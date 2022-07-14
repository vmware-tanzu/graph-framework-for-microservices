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

package nxcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	logger "github.com/sirupsen/logrus"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"

	nxv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
)

const (
	Upsert = "Upsert"
	Delete = "Delete"
	// connector
	runtimeNamespace           = "NAMESPACE"
	nexusConnectorImageVersion = "NEXUS_CONNECTOR_VERSION"
	healthEndpoint             = "/healthz"
	readyEndpoint              = "/readyz"
	volumeMountName            = "config"
	volumeMountPath            = "/config"
	kubeconfig                 = "KUBECONFIG"
	kubeconfigPath             = "/config/kubeconfig"
	// remote endpoint
	remoteEndpointHost = "REMOTE_ENDPOINT_HOST"
	remoteEndpointPort = "REMOTE_ENDPOINT_PORT"
	remoteEndpointCert = "REMOTE_ENDPOINT_CERT"
	kubeconfigLocal    = "connector-kubeconfig-local"
	// connector - init container
	initContainerName  = "check-nexus-proxy-container"
	initContainerImage = "gcr.io/mesh7-public-images/tools:latest"
)

// ApiCollaborationSpaceReconciler reconciles a ApiCollaborationSpace object
type NexusConnectorReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	K8sClient *kubernetes.Clientset
}

// +kubebuilder:rbac:groups=config.mazinger.com.design-controllers.com,resources=apicollaborationspaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=config.mazinger.com.design-controllers.com,resources=apicollaborationspaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=config.mazinger.com.design-controllers.com,resources=apicollaborationspaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApiCollaborationSpace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *NexusConnectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	logger.Infof("request received from endpoint: %q\n", req.Name)

	var endpoint nxv1.NexusEndpoint
	eventType := Upsert
	if err := r.Get(ctx, req.NamespacedName, &endpoint); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		eventType = Delete
	}
	logger.Infof("Request received from endpoint: %q\n with eventType: %q", req.Name, eventType)

	name := "nexus-connector-" + req.Name
	namespace := os.Getenv(runtimeNamespace)
	if eventType == Delete {
		deletePolicy := metav1.DeletePropagationForeground
		deleteOptions := metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}
		err := r.K8sClient.AppsV1().Deployments(namespace).Delete(ctx, name, deleteOptions)
		if err != nil {
			logger.Warnf(err.Error())
		}
		return ctrl.Result{}, nil
	}

	if endpoint.Spec.Host == "" || endpoint.Spec.Port == "" {
		return ctrl.Result{}, fmt.Errorf("endpoint Host/Port is empty")
	}

	connectorImage := os.Getenv(nexusConnectorImageVersion)
	if connectorImage == "" {
		return ctrl.Result{}, fmt.Errorf("env var NEXUS_CONNECTOR_VERSION is missing")
	}

	// config map
	err := r.createConfigMap(ctx, "connector-kubeconfig-local", namespace)
	if err != nil {
		logger.Errorf(err.Error())
		return ctrl.Result{}, err
	}

	err = r.createDeployment(ctx, name, namespace, connectorImage, &endpoint)
	if err != nil {
		logger.Errorf(err.Error())
		return ctrl.Result{}, err
	}
	err = r.createService(ctx, name, namespace)
	if err != nil {
		logger.Errorf(err.Error())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
func int32Ptr(i int32) *int32 {
	return &i
}

// SetupWithManager sets up the controller with the Manager.
func (r *NexusConnectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nxv1.NexusEndpoint{}).
		Complete(r)
}

func (r *NexusConnectorReconciler) createService(ctx context.Context, name, namespace string) error {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     80,
					Name:     "http",
					TargetPort: intstr.IntOrString{
						IntVal: 80,
					},
				},
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     443,
					Name:     "https",
					TargetPort: intstr.IntOrString{
						IntVal: 443,
					},
				},
			},
			Selector: map[string]string{
				"control-plane": name,
			},
		},
	}

	s, err := r.K8sClient.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		_, err = r.K8sClient.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
	} else {
		s.Spec = service.Spec
		_, err = r.K8sClient.CoreV1().Services(namespace).Update(ctx, s, metav1.UpdateOptions{})
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
	}
	return nil
}
func (r *NexusConnectorReconciler) createConfigMap(ctx context.Context, name, namespace string) error {
	var data = `
  connector-config: |
    dispatcher:
        workerTTL: "15s"
        maxWorkerCount: 100
        closeRequestsQueueSize: 15
        eventProcessedQueueSize: 100
    ignoredNamespaces:
        matchNames:
            - "kube-public"
            - "kube-system"
            - "kube-node-lease"
            - "istio-system"
            - "ibm-system"
            - "ibm-operators"
            - "ibm-cert-store"
  kubeconfig: |
    current-context: localapiserver
    apiVersion: v1
    kind: Config
    clusters:
    - cluster:
        api-version: v1
        server: http://nexus-proxy-container:80
        insecure-skip-tls-verify: true
      name: localapiserver
    contexts:
    - context:
        cluster: localapiserver
      name: localapiserver

`
	jsonData, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		return err
	}
	dataMap := map[string]string{}
	err = json.Unmarshal(jsonData, &dataMap)
	if err != nil {
		return err
	}

	cm, err := r.K8sClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		configMap := &apiv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: dataMap,
		}
		_, err = r.K8sClient.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		cm.Data = dataMap
		_, err = r.K8sClient.CoreV1().ConfigMaps(namespace).Update(ctx, cm, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *NexusConnectorReconciler) createDeployment(ctx context.Context, name, namespace, connectorImage string, endpoint *nxv1.NexusEndpoint) error {
	initContainerCommand := `#!/bin/bash
	set -x
	URL="http://nexus-proxy-container/api/v1/namespaces"
	max_retries=20
	counter=0
	while [[ $counter -lt $max_retries ]]; do
		  status=$(curl -s -o /dev/null -I -w "%{http_code}" -XGET $URL)
		  if [ $status == "200" ]; then
			  echo "$URL is reachable"
			  exit 0
		  else
			  counter=$((counter +1))
			  sleep 5
		  fi
	done`
	allowPrivilegeEscalation := false
	terminationGracePeriodSeconds := int64(10)
	runAsUser := int64(0)
	runAsGroup := int64(0)

	labels := make(map[string]string)
	labels["control-plane"] = "connector"
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						"kubectl.kubernetes.io/default-container": "manager",
					},
				},
				Spec: apiv1.PodSpec{
					InitContainers: []apiv1.Container{
						{
							Name:    initContainerName,
							Image:   initContainerImage,
							Command: []string{"/bin/bash", "-c", initContainerCommand},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: connectorImage,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							ImagePullPolicy: apiv1.PullIfNotPresent,
							SecurityContext: &apiv1.SecurityContext{
								AllowPrivilegeEscalation: &allowPrivilegeEscalation,
							},
							LivenessProbe: &apiv1.Probe{
								InitialDelaySeconds: int32(15),
								PeriodSeconds:       int32(20),
								ProbeHandler: apiv1.ProbeHandler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: healthEndpoint,
										Port: intstr.IntOrString{
											IntVal: 8081,
										},
									},
								},
							},
							ReadinessProbe: &apiv1.Probe{
								InitialDelaySeconds: int32(5),
								PeriodSeconds:       int32(10),
								ProbeHandler: apiv1.ProbeHandler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path: readyEndpoint,
										Port: intstr.IntOrString{
											IntVal: 8081,
										},
									},
								},
							},
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU: resource.Quantity{
										Format: "500m", // TODO to be taken from cm
									},
									apiv1.ResourceMemory: resource.Quantity{
										Format: "128Mi", // TODO to be taken from cm
									},
								},
								Requests: apiv1.ResourceList{
									apiv1.ResourceCPU: resource.Quantity{
										Format: "10m", // TODO to be taken from cm
									},
									apiv1.ResourceMemory: resource.Quantity{
										Format: "64Mi", // TODO to be taken from cm
									},
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      volumeMountName,
									MountPath: volumeMountPath,
								},
							},

							Env: []apiv1.EnvVar{
								{
									Name:  kubeconfig,
									Value: kubeconfigPath,
								},
								{
									Name:  remoteEndpointHost,
									Value: endpoint.Spec.Host,
								},
								{
									Name:  remoteEndpointPort,
									Value: endpoint.Spec.Port,
								},
								{
									Name:  remoteEndpointCert,
									Value: endpoint.Spec.Cert,
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: volumeMountName,
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: kubeconfigLocal,
									},
								},
							},
						},
					},
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					SecurityContext: &apiv1.PodSecurityContext{
						RunAsUser:  &runAsUser,
						RunAsGroup: &runAsGroup,
					},
				},
			},
		},
	}

	dep, err := r.K8sClient.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		_, err = r.K8sClient.AppsV1().Deployments(namespace).Create(ctx, deploy, metav1.CreateOptions{})
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
		logger.Infof("Nexus deploymemt: %q created\n", name)
	} else {
		dep.Spec = deploy.Spec
		_, err = r.K8sClient.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
		logger.Infof("Nexus deploymemt: %q updated\n", name)
	}

	return nil
}
