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

// GatewayConfigLister helps list GatewayConfigs.
// All objects returned here must be treated as read-only.
type GatewayConfigLister interface {
	// List lists all GatewayConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.GatewayConfig, err error)
	// Get retrieves the GatewayConfig from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.GatewayConfig, error)
	GatewayConfigListerExpansion
}

// gatewayConfigLister implements the GatewayConfigLister interface.
type gatewayConfigLister struct {
	indexer cache.Indexer
}

// NewGatewayConfigLister returns a new GatewayConfigLister.
func NewGatewayConfigLister(indexer cache.Indexer) GatewayConfigLister {
	return &gatewayConfigLister{indexer: indexer}
}

// List lists all GatewayConfigs in the indexer.
func (s *gatewayConfigLister) List(selector labels.Selector) (ret []*v1.GatewayConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.GatewayConfig))
	})
	return ret, err
}

// Get retrieves the GatewayConfig from the index for a given name.
func (s *gatewayConfigLister) Get(name string) (*v1.GatewayConfig, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("gatewayconfig"), name)
	}
	return obj.(*v1.GatewayConfig), nil
}