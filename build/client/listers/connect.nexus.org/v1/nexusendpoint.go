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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// NexusEndpointLister helps list NexusEndpoints.
// All objects returned here must be treated as read-only.
type NexusEndpointLister interface {
	// List lists all NexusEndpoints in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.NexusEndpoint, err error)
	// Get retrieves the NexusEndpoint from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.NexusEndpoint, error)
	NexusEndpointListerExpansion
}

// nexusEndpointLister implements the NexusEndpointLister interface.
type nexusEndpointLister struct {
	indexer cache.Indexer
}

// NewNexusEndpointLister returns a new NexusEndpointLister.
func NewNexusEndpointLister(indexer cache.Indexer) NexusEndpointLister {
	return &nexusEndpointLister{indexer: indexer}
}

// List lists all NexusEndpoints in the indexer.
func (s *nexusEndpointLister) List(selector labels.Selector) (ret []*v1.NexusEndpoint, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NexusEndpoint))
	})
	return ret, err
}

// Get retrieves the NexusEndpoint from the index for a given name.
func (s *nexusEndpointLister) Get(name string) (*v1.NexusEndpoint, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("nexusendpoint"), name)
	}
	return obj.(*v1.NexusEndpoint), nil
}
