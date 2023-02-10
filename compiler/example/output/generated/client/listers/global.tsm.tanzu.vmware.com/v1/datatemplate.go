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
	v1 "nexustempmodule/apis/global.tsm.tanzu.vmware.com/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DataTemplateLister helps list DataTemplates.
// All objects returned here must be treated as read-only.
type DataTemplateLister interface {
	// List lists all DataTemplates in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.DataTemplate, err error)
	// Get retrieves the DataTemplate from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.DataTemplate, error)
	DataTemplateListerExpansion
}

// dataTemplateLister implements the DataTemplateLister interface.
type dataTemplateLister struct {
	indexer cache.Indexer
}

// NewDataTemplateLister returns a new DataTemplateLister.
func NewDataTemplateLister(indexer cache.Indexer) DataTemplateLister {
	return &dataTemplateLister{indexer: indexer}
}

// List lists all DataTemplates in the indexer.
func (s *dataTemplateLister) List(selector labels.Selector) (ret []*v1.DataTemplate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.DataTemplate))
	})
	return ret, err
}

// Get retrieves the DataTemplate from the index for a given name.
func (s *dataTemplateLister) Get(name string) (*v1.DataTemplate, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("datatemplate"), name)
	}
	return obj.(*v1.DataTemplate), nil
}