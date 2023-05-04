/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	adminnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/admin.nexus.vmware.com/v1"
	apinexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/api.nexus.vmware.com/v1"
	apigatewaynexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/apigateway.nexus.vmware.com/v1"
	authenticationnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.vmware.com/v1"
	commonnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/common.nexus.vmware.com/v1"
	confignexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.vmware.com/v1"
	connectnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.vmware.com/v1"
	domainnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/domain.nexus.vmware.com/v1"
	routenexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/route.nexus.vmware.com/v1"
	runtimenexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/runtime.nexus.vmware.com/v1"
	tenantconfignexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	tenantruntimenexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantruntime.nexus.vmware.com/v1"
	usernexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/user.nexus.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)

var localSchemeBuilder = runtime.SchemeBuilder{
	adminnexusv1.AddToScheme,
	apinexusv1.AddToScheme,
	apigatewaynexusv1.AddToScheme,
	authenticationnexusv1.AddToScheme,
	commonnexusv1.AddToScheme,
	confignexusv1.AddToScheme,
	connectnexusv1.AddToScheme,
	domainnexusv1.AddToScheme,
	routenexusv1.AddToScheme,
	runtimenexusv1.AddToScheme,
	tenantconfignexusv1.AddToScheme,
	tenantruntimenexusv1.AddToScheme,
	usernexusv1.AddToScheme,
}

// AddToScheme adds all types of this clientset into the given scheme. This allows composition
// of clientsets, like in:
//
//   import (
//     "k8s.io/client-go/kubernetes"
//     clientsetscheme "k8s.io/client-go/kubernetes/scheme"
//     aggregatorclientsetscheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"
//   )
//
//   kclientset, _ := kubernetes.NewForConfig(c)
//   _ = aggregatorclientsetscheme.AddToScheme(clientsetscheme.Scheme)
//
// After this, RawExtensions in Kubernetes types will serialize kube-aggregator types
// correctly.
var AddToScheme = localSchemeBuilder.AddToScheme

func init() {
	v1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})
	utilruntime.Must(AddToScheme(scheme))
}
