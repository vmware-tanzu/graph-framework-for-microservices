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
	"io"
	"net/http"
	net_http "net/http"
	"os"

	"plugin"

	logger "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const defaultPort = "8080"

var PLUGIN_PATH = os.Getenv("PLUGIN_PATH")

// DatamodelReconciler reconciles a Datamodels object
type DatamodelReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	StopCh  chan struct{}
	Dynamic dynamic.Interface
}

//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apiextensions.k8s.io.api-gw.com,resources=datamodels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile

func DownloadFile(url string, filename string) error {
	url = fmt.Sprintf("%s", url)

	var (
		retries int = 5
		resp    *net_http.Response
		err     error
	)
	for retries > 0 {
		resp, err = net_http.Get(url)
		if err != nil {
			logger.Errorf("Retrying download of file : %s due to :%s...", url, err)
			retries -= 1
		} else {
			break
		}
	}
	if resp != nil {
		defer resp.Body.Close()
		out, _ := os.Create(filename)
		defer out.Close()
		_, err := io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("URL is not valid")

}

func Startserver(stopCh chan struct{}, graphqlBuildplugin string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	if _, err := os.Stat(graphqlBuildplugin); err != nil {
		logger.Errorf("error in checking graphql plugin file %s", graphqlBuildplugin)
		panic(err)
	}
	// Opening graphql plugin file archieved from datamodel image
	pl, err := plugin.Open(graphqlBuildplugin)
	if err != nil {
		logger.Errorf("could not open pluginfile: %s", err)
		panic(err)
	}

	// Lookup init method present
	plsm, err := pl.Lookup("StartHttpServer")
	if err != nil {
		logger.Errorf("could not lookup the InitMethod : %s", err)
		panic(err)
	}
	// Execute the init method for initialising resolvers and typecast to expected format
	plsm.(func())()

	fmt.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	go func() {
		srv := &http.Server{Addr: fmt.Sprintf(":%s", port)}
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("could not start graphqlServer")
			panic(err)
		}
		select {
		case <-stopCh:
			if err := srv.Shutdown(context.TODO()); err != nil {
				logger.Error("could not stop running graphqlServer")
				panic(err)
			}
		}
	}()

}

func (r *DatamodelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	eventType := "Update"

	wholeObject, err := r.Dynamic.Resource(schema.GroupVersionResource{
		Group:    "nexus.org",
		Version:  "v1",
		Resource: "datamodels",
	}).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		eventType = "Delete"
	}
	obj := wholeObject.Object
	spec := obj["spec"].(map[string]interface{})
	if eventType != "Delete" {
		if enableGraphql, ok := spec["enableGraphql"]; ok && enableGraphql.(bool) {
			if graphqlUrl, ok := spec["graphqlPath"]; ok {
				logger.Infof("Downloading graphqlPlugin from URL...: %s", graphqlUrl.(string))
				if err := DownloadFile(graphqlUrl.(string), PLUGIN_PATH); err != nil {
					return ctrl.Result{}, fmt.Errorf("could not download the graphql plugin fro url %s", graphqlUrl.(string))
				}
				logger.Infof("Downloaded graphqlPlugin file from URL.....: %s", graphqlUrl.(string))
				logger.Info("stopping existing graphqlServer.....")
				go func() { r.StopCh <- struct{}{} }()
				logger.Info("restarting graphqlServer with new plugin.....")
				Startserver(r.StopCh, PLUGIN_PATH)
			}
		}
	}
	logger.Infof("Received Datamodel notification for Name %s Type %s", req.Name, eventType)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatamodelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "Datamodel",
		Group:   "nexus.org",
		Version: "v1",
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(u).
		Complete(r)
}
