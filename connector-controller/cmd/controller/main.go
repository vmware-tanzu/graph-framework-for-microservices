/*
Copyright 2020, 2021 The Flux authors

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
	"context"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/fluxcd/pkg/runtime/logger"
	flag "github.com/spf13/pflag"
	nxcontrollers "gitlab.eng.vmware.com/nexus/controller/controllers"
	"gitlab.eng.vmware.com/nexus/controller/helpers"
	nxv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.vmware.com/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	// +kubebuilder:scaffold:imports
)

var (
	scheme     = runtime.NewScheme()
	setupLog   = ctrl.Log.WithName("setup")
	logOptions logger.Options
)

func init() {
	utilruntime.Must(nxv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	if err := run(); err != nil {
		switch err {
		case context.Canceled:
			// not considered error
		default:
			log.Fatalf("could not run Nexus Controller app: %v", err)
		}
	}
}

func run() error {
	var (
		metricsAddr          string
		probeAddr            string
		enableLeaderElection bool
	)
	logOptions.BindFlags(flag.CommandLine)
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	flag.Parse()
	ctrl.SetLogger(logger.NewLogger(logOptions))

	outerK8sClient, err := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		setupLog.Error(err, "unable to get outer k8s client")
		return err
	}

	nexusK8s, err := helpers.GetNexusConfig()
	if err != nil {
		setupLog.Error(err, "unable to get nexus k8s config")
		return err
	}

	mgr, err := ctrl.NewManager(nexusK8s, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		Logger:                 ctrl.Log,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "7b10c258.controller.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		return err
	}
	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}
	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}

	if err = (&nxcontrollers.NexusConnectorReconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		K8sClient: outerK8sClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ApiCollaborationSpace")
		return err
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	errGrp, errCtx := errgroup.WithContext(ctrl.SetupSignalHandler())

	errGrp.Go(func() error {
		setupLog.Info("starting Nexus API server manager")
		return mgr.Start(errCtx)
	})

	return errGrp.Wait()
}
