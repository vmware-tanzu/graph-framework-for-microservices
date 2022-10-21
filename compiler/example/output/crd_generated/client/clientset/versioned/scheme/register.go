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

package scheme

import (
	configtsmv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/config.tsm.tanzu.vmware.com/v1"
	gnstsmv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/gns.tsm.tanzu.vmware.com/v1"
	policypkgtsmv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/policypkg.tsm.tanzu.vmware.com/v1"
	roottsmv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/root.tsm.tanzu.vmware.com/v1"
	servicegrouptsmv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/servicegroup.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var Scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(Scheme)
var ParameterCodec = runtime.NewParameterCodec(Scheme)
var localSchemeBuilder = runtime.SchemeBuilder{
	configtsmv1.AddToScheme,
	gnstsmv1.AddToScheme,
	policypkgtsmv1.AddToScheme,
	roottsmv1.AddToScheme,
	servicegrouptsmv1.AddToScheme,
}

// AddToScheme adds all types of this clientset into the given scheme. This allows composition
// of clientsets, like in:
//
//	import (
//	  "k8s.io/client-go/kubernetes"
//	  clientsetscheme "k8s.io/client-go/kubernetes/scheme"
//	  aggregatorclientsetscheme "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme"
//	)
//
//	kclientset, _ := kubernetes.NewForConfig(c)
//	_ = aggregatorclientsetscheme.AddToScheme(clientsetscheme.Scheme)
//
// After this, RawExtensions in Kubernetes types will serialize kube-aggregator types
// correctly.
var AddToScheme = localSchemeBuilder.AddToScheme

func init() {
	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	utilruntime.Must(AddToScheme(Scheme))
}
