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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/admin.nexus.org/v1"
	apinexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/api.nexus.org/v1"
	apigatewaynexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/apigateway.nexus.org/v1"
	authenticationnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.org/v1"
	confignexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.org/v1"
	connectnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
	routenexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/route.nexus.org/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=admin.nexus.org, Version=v1
	case v1.SchemeGroupVersion.WithResource("proxyrules"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.AdminNexus().V1().ProxyRules().Informer()}, nil

		// Group=api.nexus.org, Version=v1
	case apinexusorgv1.SchemeGroupVersion.WithResource("nexuses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.ApiNexus().V1().Nexuses().Informer()}, nil

		// Group=apigateway.nexus.org, Version=v1
	case apigatewaynexusorgv1.SchemeGroupVersion.WithResource("apigateways"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.ApigatewayNexus().V1().ApiGateways().Informer()}, nil

		// Group=authentication.nexus.org, Version=v1
	case authenticationnexusorgv1.SchemeGroupVersion.WithResource("oidcs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.AuthenticationNexus().V1().OIDCs().Informer()}, nil

		// Group=config.nexus.org, Version=v1
	case confignexusorgv1.SchemeGroupVersion.WithResource("configs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.ConfigNexus().V1().Configs().Informer()}, nil

		// Group=connect.nexus.org, Version=v1
	case connectnexusorgv1.SchemeGroupVersion.WithResource("connects"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.ConnectNexus().V1().Connects().Informer()}, nil
	case connectnexusorgv1.SchemeGroupVersion.WithResource("nexusendpoints"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.ConnectNexus().V1().NexusEndpoints().Informer()}, nil
	case connectnexusorgv1.SchemeGroupVersion.WithResource("replicationconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.ConnectNexus().V1().ReplicationConfigs().Informer()}, nil

		// Group=route.nexus.org, Version=v1
	case routenexusorgv1.SchemeGroupVersion.WithResource("routes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.RouteNexus().V1().Routes().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
