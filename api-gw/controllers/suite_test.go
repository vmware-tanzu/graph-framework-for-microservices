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
	"os"
	"path/filepath"
	"testing"
	"time"

	adminnexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/admin.nexus.vmware.com/v1"
	apinexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/api.nexus.vmware.com/v1"
	authenticationnexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/authentication.nexus.vmware.com/v1"
	confignexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/config.nexus.vmware.com/v1"
	domain_nexus_org_v1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/domain.nexus.vmware.com/v1"
	routenexusv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/route.nexus.vmware.com/v1"
	tenantv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	tenantruntimev1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/tenantruntime.nexus.vmware.com/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg           *rest.Config
	k8sClient     client.Client
	dynamicClient dynamic.Interface
	testEnv       *envtest.Environment
	ctx           context.Context
	cancel        context.CancelFunc
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	Expect(os.Setenv("TEST_ASSET_KUBE_APISERVER", "../test/bin/kube-apiserver")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_ETCD", "../test/bin/etcd")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_KUBECTL", "../test/bin/kubectl")).To(Succeed())

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())

	var err error
	testEnv = &envtest.Environment{}
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())

	// Install CRDs
	opts := envtest.CRDInstallOptions{
		Paths: []string{filepath.Join("..", "test", "crds", "bases")},
	}
	crds, err := envtest.InstallCRDs(cfg, opts)
	Expect(err).NotTo(HaveOccurred())

	err = envtest.WaitForCRDs(cfg, crds, envtest.CRDInstallOptions{})
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	err = apiextensions.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = authenticationnexusv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = adminnexusv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = routenexusv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = apinexusv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = confignexusv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = domain_nexus_org_v1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = tenantv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = tenantruntimev1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	dynamicClient, err = dynamic.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	Expect(dynamicClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	// DatamodelReconciler
	err = (&DatamodelReconciler{
		Client:  k8sManager.GetClient(),
		Scheme:  k8sManager.GetScheme(),
		Dynamic: dynamicClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	// CustomResourceDefinitionReconciler
	err = (&CustomResourceDefinitionReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	// OidcConfigReconciler
	err = (&OidcConfigReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	// ProxyRuleReconciler
	err = (&ProxyRuleReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	// RouteReconciler
	err = (&RouteReconciler{
		Client:     k8sManager.GetClient(),
		Scheme:     k8sManager.GetScheme(),
		BaseClient: k8sManager.GetClient(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&CORSReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&TenantReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		GrpcConnector: &model.ConnectorObject{},
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&TenantRuntimeReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		GrpcConnector: &model.ConnectorObject{},
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()
}, 60)

var _ = AfterSuite(func() {
	// https://github.com/kubernetes-sigs/controller-runtime/issues/1571
	cancel()
	err := testEnv.Stop()
	if err != nil {
		time.Sleep(4 * time.Second)
	}
	err = testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())

	if envoy.XDSListener != nil {
		envoy.XDSListener.Close()
	}
	Expect(os.Unsetenv("TEST_ASSET_KUBE_APISERVER")).To(Succeed())
	Expect(os.Unsetenv("TEST_ASSET_ETCD")).To(Succeed())
	Expect(os.Unsetenv("TEST_ASSET_KUBECTL")).To(Succeed())
})
