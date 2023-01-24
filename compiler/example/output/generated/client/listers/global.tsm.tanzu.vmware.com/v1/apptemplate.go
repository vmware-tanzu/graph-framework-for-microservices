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

// AppTemplateLister helps list AppTemplates.
// All objects returned here must be treated as read-only.
type AppTemplateLister interface {
	// List lists all AppTemplates in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.AppTemplate, err error)
	// Get retrieves the AppTemplate from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.AppTemplate, error)
	AppTemplateListerExpansion
}

// appTemplateLister implements the AppTemplateLister interface.
type appTemplateLister struct {
	indexer cache.Indexer
}

// NewAppTemplateLister returns a new AppTemplateLister.
func NewAppTemplateLister(indexer cache.Indexer) AppTemplateLister {
	return &appTemplateLister{indexer: indexer}
}

// List lists all AppTemplates in the indexer.
func (s *appTemplateLister) List(selector labels.Selector) (ret []*v1.AppTemplate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.AppTemplate))
	})
	return ret, err
}

// Get retrieves the AppTemplate from the index for a given name.
func (s *appTemplateLister) Get(name string) (*v1.AppTemplate, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("apptemplate"), name)
	}
	return obj.(*v1.AppTemplate), nil
}