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
	v1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/apis/config.tsm.tanzu.vmware.com/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DomainLister helps list Domains.
// All objects returned here must be treated as read-only.
type DomainLister interface {
	// List lists all Domains in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Domain, err error)
	// Get retrieves the Domain from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Domain, error)
	DomainListerExpansion
}

// domainLister implements the DomainLister interface.
type domainLister struct {
	indexer cache.Indexer
}

// NewDomainLister returns a new DomainLister.
func NewDomainLister(indexer cache.Indexer) DomainLister {
	return &domainLister{indexer: indexer}
}

// List lists all Domains in the indexer.
func (s *domainLister) List(selector labels.Selector) (ret []*v1.Domain, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Domain))
	})
	return ret, err
}

// Get retrieves the Domain from the index for a given name.
func (s *domainLister) Get(name string) (*v1.Domain, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("domain"), name)
	}
	return obj.(*v1.Domain), nil
}
