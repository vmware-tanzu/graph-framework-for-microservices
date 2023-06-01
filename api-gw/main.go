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

package main

import (
	"api-gw/internal/tenant/registration"
	"api-gw/pkg/client"
	"api-gw/pkg/common"
	"api-gw/pkg/config"
	"api-gw/pkg/envoy"
	"api-gw/pkg/model"
	"api-gw/pkg/openapi/api"
	"api-gw/pkg/openapi/declarative"
	"api-gw/pkg/utils"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	reg_svc "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/registration-service/global"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	"os"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"strconv"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	authnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.vmware.com/v1"
	corev1 "k8s.io/api/core/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"api-gw/controllers"

	adminnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/admin.nexus.vmware.com/v1"
	apigatewaynexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/apigateway.nexus.vmware.com/v1"

	middleware_nexus_org_v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/domain.nexus.vmware.com/v1"
	routenexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/route.nexus.vmware.com/v1"
	tenantv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	tenantruntimev1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantruntime.nexus.vmware.com/v1"

	//+kubebuilder:scaffold:imports

	"api-gw/pkg/server/echo_server"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	apiextensionsv1.AddToScheme(scheme)
	authnexusv1.AddToScheme(scheme)
	apiregistrationv1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
	routenexusorgv1.AddToScheme(scheme)
	apigatewaynexusorgv1.AddToScheme(scheme)
	adminnexusorgv1.AddToScheme(scheme)
	middleware_nexus_org_v1.AddToScheme(scheme)
	tenantv1.AddToScheme(scheme)
	tenantruntimev1.AddToScheme(scheme)
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	k8sConfig := ctrl.GetConfigOrDie()
	k8sClientSet, err := kubernetes.NewForConfig(k8sConfig)

	if err != nil {
		log.Fatalf("Failed to create K8sclient: %v", err)
	}

	gsRoutes, err := config.LoadStaticUrlsConfig("/config/staticRoutes.yaml")
	if err != nil {
		log.Warnf("Error loading config: %v\n", err)
	}
	config.GlobalStaticRouteConfig = gsRoutes

	nexusClientSet, err := nexus_client.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("Failed to create nexusclient: %v", err)
	}
	// Setup log level
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "ERROR"
	}
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Failed to configure logging: %v\n", err)
	}
	log.SetLevel(lvl)

	common.Mode = os.Getenv("GATEWAY_MODE")
	log.Infof("Gateway Mode: %s", common.Mode)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02T15:04:05Z07:00"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	conf, err := config.LoadConfig("/config/api-gw-config")
	if err != nil {
		log.Warnf("Error loading config: %v\n", err)
	}
	config.Cfg = conf

	if common.IsModeAdmin() {
		skuConfig, err := config.LoadSKUConfig("/config/skuconfigmap")
		if err != nil {
			log.Warnf("Error loading  skuconfig: %v\n", err)
		}
		config.SKUConfig = skuConfig

	}

	stopCh := make(chan struct{})

	if common.IsModeAdmin() {
		utils.VersionCalls = []*model.ConnectorObject{
			{
				Service:    "global-registration-service:60000",
				Protocol:   "grpc",
				Connection: nil,
			},
			{
				Service:    fmt.Sprintf("%s:60000/ui/version", common.GlobalUISvcName),
				Protocol:   "http",
				Connection: nil,
			},
		}
		for _, v := range utils.VersionCalls {
			err := v.InitConnection()
			if err != nil {
				log.Errorf("could not create connection for : %s", v.Service)
			}
		}
	}

	log.Infoln("Init Echo Server")
	// Start server
	echo_server.InitEcho(stopCh, conf, k8sClientSet, nexusClientSet)

	if conf.BackendService != "" {
		echo_server.WatchForOpenApiSpecChanges(stopCh, declarative.OpenApiSpecDir, declarative.OpenApiSpecFile)
	}

	if conf.EnableNexusRuntime {
		InitManager(metricsAddr, probeAddr, enableLeaderElection, stopCh, lvl)
	}

	select {}
}

func InitManager(metricsAddr string, probeAddr string, enableLeaderElection bool, stopCh chan struct{}, lvl log.Level) {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "7b10c258.api-gw.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	if err = (&controllers.CustomResourceDefinitionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		StopCh: stopCh,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CustomResourceDefinition")
		os.Exit(1)
	}

	if err = (&controllers.OidcConfigReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "OidcConfig")
		os.Exit(1)
	}

	if err = (&controllers.CORSReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CORSreconciler")
		os.Exit(1)
	}
	baseConfig, err := rest.InClusterConfig()
	if err != nil {
		setupLog.Error(err, "unable to create manager for base k8s")
		os.Exit(1)
	}

	baseClient, err := runtimeclient.New(baseConfig, runtimeclient.Options{})

	baseNamespace := os.Getenv("NAMESPACE")
	ingressControllerName := "ingress-nginx-controller"
	if os.Getenv("INGRESS_CONTROLLER_NAME") != "" {
		ingressControllerName = os.Getenv("INGRESS_CONTROLLER_NAME")
	}

	// Adding nginx base server , needed because of the / requirement.
	defaultBackendservice := os.Getenv("DEFAULT_BACKEND_SERVICE_NAME")
	defaultBackendPort, _ := strconv.Atoi(os.Getenv("DEFAULT_BACKEND_SERVICE_PORT"))
	if err = (&controllers.RouteReconciler{
		Client:                mgr.GetClient(),
		BaseClient:            baseClient,
		Scheme:                mgr.GetScheme(),
		BaseNamespace:         baseNamespace,
		IngressControllerName: ingressControllerName,
		DefaultBackend:        defaultBackendservice,
		DefaultBackendPort:    int32(defaultBackendPort),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Route")
		os.Exit(1)
	}

	if err = (&controllers.ProxyRuleReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ProxyRule")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	// Create new dynamic client for kubernetes
	if err = client.New(ctrl.GetConfigOrDie()); err != nil {
		setupLog.Error(err, "unable to set up dynamic client")
		os.Exit(1)
	}

	if err = client.NewNexusClient(ctrl.GetConfigOrDie()); err != nil {
		setupLog.Error(err, "unable to set up nexus client")
		os.Exit(1)
	}

	if err = (&controllers.DatamodelReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Dynamic: client.Client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Datamodel")
		os.Exit(1)
	}

	// Create new openapi3 schema
	api.Recreate()
	go api.DatamodelUpdateNotification()

	common.SSLEnabled = os.Getenv("SSL_ENABLED")
	log.Infof("SSL CertsEnabled: %s", common.SSLEnabled)

	log.Infoln("Init xDS server")
	if jwt, upstreams, headerUpstreams, err := utils.GetEnvoyInitParams(); err != nil {
		log.Errorf("error getting envoy init params: %s\n", err)
		// start with a blank envoy config and let the controllers reconcile the envoy state
		if err = envoy.Init(nil, nil, nil, lvl); err != nil {
			log.Fatalf("error initializing envoy in main(): %s", err)
		}
	} else {
		if err = envoy.Init(jwt, upstreams, headerUpstreams, lvl); err != nil {
			panic(err)
		}
	}
	log.Infoln("successfully initialized xDS server")

	if common.IsModeAdmin() {
		//Fetch CSPPermissionName and CSP ServiceID
		err := common.SetCSPVariables()
		if err != nil {
			setupLog.Info("CSPVariables passed as empty")
		} else {
			setupLog.Info(fmt.Sprintf("CSPVariables set as %s=%s , %s=%s", common.CSPPermissionName, common.CSP_PERMISSION_NAME, common.CSPServiceID, common.CSP_SERVICE_ID))
		}
		common.SetCSPPermissionOrg()
		grpcConnector := model.ConnectorObject{
			Service:  "global-registration-service:30031",
			Protocol: "grpc",
		}
		err = grpcConnector.InitConnection()
		if err != nil {
			setupLog.Error(err, "unable to reconcile TenantConfig")
			os.Exit(1)
		}

		reg_client := reg_svc.NewGlobalRegistrationClient(grpcConnector.Connection)

		if err := registration.InitTenantConfig(reg_client); err != nil {
			setupLog.Error(err, "unable to reconcile TenantConfig")
			os.Exit(1)
		}
		if err := registration.InitTenantRuntimeCache(reg_client); err != nil {
			setupLog.Error(err, "unable to start Tenant Cache")
			os.Exit(1)
		}
		if err := common.InitAdminDatamodelCache(); err != nil {
			setupLog.Error(err, "unable to start User Cache")
			os.Exit(1)
		}
		if err = (&controllers.TenantReconciler{
			Client:        mgr.GetClient(),
			Scheme:        mgr.GetScheme(),
			K8sclient:     client.CoreClient,
			GrpcConnector: &grpcConnector,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Tenant Config")
			os.Exit(1)
		}
		if err = (&controllers.TenantRuntimeReconciler{
			Client:        mgr.GetClient(),
			Scheme:        mgr.GetScheme(),
			K8sclient:     client.CoreClient,
			GrpcConnector: &grpcConnector,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Tenant runtime")
			os.Exit(1)
		}
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
