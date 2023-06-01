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
	gnstsmtanzuvmwarecomv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/gns.tsm.tanzu.vmware.com/v1"
	versioned "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/client/clientset/versioned"
	internalinterfaces "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/client/informers/externalversions/internalinterfaces"
	v1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/client/listers/gns.tsm.tanzu.vmware.com/v1"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GnsInformer provides access to a shared informer and lister for
// Gnses.
type GnsInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.GnsLister
}

type gnsInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewGnsInformer constructs a new informer for Gns type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGnsInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGnsInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredGnsInformer constructs a new informer for Gns type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGnsInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GnsTsmV1().Gnses().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GnsTsmV1().Gnses().Watch(context.TODO(), options)
			},
		},
		&gnstsmtanzuvmwarecomv1.Gns{},
		resyncPeriod,
		indexers,
	)
}

func (f *gnsInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGnsInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gnsInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&gnstsmtanzuvmwarecomv1.Gns{}, f.defaultInformer)
}

func (f *gnsInformer) Lister() v1.GnsLister {
	return v1.NewGnsLister(f.Informer().GetIndexer())
}
