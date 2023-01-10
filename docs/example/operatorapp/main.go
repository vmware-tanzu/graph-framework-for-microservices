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
	"flag"
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	nexus_client "vmware/build/nexus-client"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	managementvmwareorgv1 "vmware/build/apis/management.vmware.org/v1"

	managementvmwareorgcontrollers "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/docs/example/operatorapp/controllers/management.vmware.org"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	managementvmwareorgv1.AddToScheme(scheme)
	//+kubebuilder:scaffold:scheme
}

func main() {
	config := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	ctrl.SetLogger(zap.New())

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&managementvmwareorgcontrollers.LeaderReconciler{
		Client:      mgr.GetClient(),
		NexusClient: nexusClient,
		Scheme:      mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Leader")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// getK8sAPIEndpointConfig determines K8s API Server endpoint
//
// If "host" is specified in command line argument, connect to API server pointed to by it.
// If not, if "kubeconfig" file is provided as input, connect to API server pointed to by it.
// If not, attempt to read kubeconfig file from home directory and connect to API server pointed to by it.
// If none of these are available, then exit with error.
func getK8sAPIEndpointConfig() *rest.Config {
	var (
		host           *string
		kubeconfigHome string
		config         *rest.Config
		err            error
	)

	host = flag.String("host", "", "portfowarded host to reach the app")
	if home := homedir.HomeDir(); home != "" {
		kubeconfigHome = filepath.Join(home, ".kube", "config")
	}

	flag.Parse()

	if len(*host) > 0 {
		fmt.Println("Connecting to k8s API at host: ", *host)
		config = &rest.Config{
			Host: *host,
		}
	} else if len(kubeconfigHome) > 0 {
		fmt.Println("Connecting to k8s API in kubeconfig in home dir: ", kubeconfigHome)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigHome)
		if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Println("Unable to determing k8s API server endpoint. Exiting application.")
		os.Exit(1)
	}

	return config
}
