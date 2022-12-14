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
	v1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/gns.tsm.tanzu.vmware.com/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RandomGnsDataLister helps list RandomGnsDatas.
// All objects returned here must be treated as read-only.
type RandomGnsDataLister interface {
	// List lists all RandomGnsDatas in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.RandomGnsData, err error)
	// Get retrieves the RandomGnsData from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.RandomGnsData, error)
	RandomGnsDataListerExpansion
}

// randomGnsDataLister implements the RandomGnsDataLister interface.
type randomGnsDataLister struct {
	indexer cache.Indexer
}

// NewRandomGnsDataLister returns a new RandomGnsDataLister.
func NewRandomGnsDataLister(indexer cache.Indexer) RandomGnsDataLister {
	return &randomGnsDataLister{indexer: indexer}
}

// List lists all RandomGnsDatas in the indexer.
func (s *randomGnsDataLister) List(selector labels.Selector) (ret []*v1.RandomGnsData, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RandomGnsData))
	})
	return ret, err
}

// Get retrieves the RandomGnsData from the index for a given name.
func (s *randomGnsDataLister) Get(name string) (*v1.RandomGnsData, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("randomgnsdata"), name)
	}
	return obj.(*v1.RandomGnsData), nil
}
