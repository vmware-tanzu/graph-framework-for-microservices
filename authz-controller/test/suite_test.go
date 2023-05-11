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

package test_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	coreV1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	api_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/api.nexus.org/v1"
	auth_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authorization.nexus.org/v1"
	config_nexus_org "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.org/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	ctx         context.Context
	nexusClient *nexus_client.Clientset
	authNode    *nexus_client.AuthorizationAuthorization
	scheme      *runtime.Scheme
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Authz Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	nexusClient = nexus_client.NewFakeClient()
	ctx = context.Background()
	scheme = addScheme()
})

var _ = AfterSuite(func() {
	clearNexusAPI()
})

func createParentNodes() *nexus_client.AuthorizationAuthorization {
	nexusapi, err := nexusClient.AddApiNexus(ctx, &api_nexus_org.Nexus{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nexus",
		},
		Spec: api_nexus_org.NexusSpec{},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(nexusapi.DisplayName()).To(Equal("nexus"))

	config, err := nexusapi.AddConfig(ctx, &config_nexus_org.Config{
		ObjectMeta: metav1.ObjectMeta{
			Name: "config",
		},
		Spec: config_nexus_org.ConfigSpec{},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(config.DisplayName()).To(Equal("config"))

	authz, err := config.AddAuthorization(ctx, &auth_nexus_org.Authorization{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: auth_nexus_org.AuthorizationSpec{},
	})

	Expect(err).NotTo(HaveOccurred())
	Expect(authz.DisplayName()).To(Equal("default"))

	return authz
}

func addScheme() *runtime.Scheme {
	scheme = runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = coreV1.AddToScheme(scheme)
	_ = auth_nexus_org.AddToScheme(scheme)
	_ = apiextensionsv1.AddToScheme(scheme)
	return scheme
}

func clearNexusAPI() {
	_ = nexusClient.DeleteApiNexus(ctx, "nexus")
}
