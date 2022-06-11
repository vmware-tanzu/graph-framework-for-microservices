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

package v1

import (
	"context"
	time "time"

	connectnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
	versioned "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/clientset/versioned"
	internalinterfaces "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/informers/externalversions/internalinterfaces"
	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/listers/connect.nexus.org/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ConnectInformer provides access to a shared informer and lister for
// Connects.
type ConnectInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ConnectLister
}

type connectInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewConnectInformer constructs a new informer for Connect type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewConnectInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredConnectInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredConnectInformer constructs a new informer for Connect type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredConnectInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConnectNexusV1().Connects().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConnectNexusV1().Connects().Watch(context.TODO(), options)
			},
		},
		&connectnexusorgv1.Connect{},
		resyncPeriod,
		indexers,
	)
}

func (f *connectInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredConnectInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *connectInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&connectnexusorgv1.Connect{}, f.defaultInformer)
}

func (f *connectInformer) Lister() v1.ConnectLister {
	return v1.NewConnectLister(f.Informer().GetIndexer())
}
