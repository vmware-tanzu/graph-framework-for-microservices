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
	v1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/apis/gns.tsm.tanzu.vmware.com/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// AdditionalGnsDataLister helps list AdditionalGnsDatas.
// All objects returned here must be treated as read-only.
type AdditionalGnsDataLister interface {
	// List lists all AdditionalGnsDatas in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.AdditionalGnsData, err error)
	// Get retrieves the AdditionalGnsData from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.AdditionalGnsData, error)
	AdditionalGnsDataListerExpansion
}

// additionalGnsDataLister implements the AdditionalGnsDataLister interface.
type additionalGnsDataLister struct {
	indexer cache.Indexer
}

// NewAdditionalGnsDataLister returns a new AdditionalGnsDataLister.
func NewAdditionalGnsDataLister(indexer cache.Indexer) AdditionalGnsDataLister {
	return &additionalGnsDataLister{indexer: indexer}
}

// List lists all AdditionalGnsDatas in the indexer.
func (s *additionalGnsDataLister) List(selector labels.Selector) (ret []*v1.AdditionalGnsData, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.AdditionalGnsData))
	})
	return ret, err
}

// Get retrieves the AdditionalGnsData from the index for a given name.
func (s *additionalGnsDataLister) Get(name string) (*v1.AdditionalGnsData, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("additionalgnsdata"), name)
	}
	return obj.(*v1.AdditionalGnsData), nil
}
