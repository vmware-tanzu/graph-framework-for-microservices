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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	configtsmtanzuvmwarecomv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/config.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDomains implements DomainInterface
type FakeDomains struct {
	Fake *FakeConfigTsmV1
}

var domainsResource = schema.GroupVersionResource{Group: "config.tsm.tanzu.vmware.com", Version: "v1", Resource: "domains"}

var domainsKind = schema.GroupVersionKind{Group: "config.tsm.tanzu.vmware.com", Version: "v1", Kind: "Domain"}

// Get takes name of the domain, and returns the corresponding domain object, and an error if there is any.
func (c *FakeDomains) Get(ctx context.Context, name string, options v1.GetOptions) (result *configtsmtanzuvmwarecomv1.Domain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(domainsResource, name), &configtsmtanzuvmwarecomv1.Domain{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.Domain), err
}

// List takes label and field selectors, and returns the list of Domains that match those selectors.
func (c *FakeDomains) List(ctx context.Context, opts v1.ListOptions) (result *configtsmtanzuvmwarecomv1.DomainList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(domainsResource, domainsKind, opts), &configtsmtanzuvmwarecomv1.DomainList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &configtsmtanzuvmwarecomv1.DomainList{ListMeta: obj.(*configtsmtanzuvmwarecomv1.DomainList).ListMeta}
	for _, item := range obj.(*configtsmtanzuvmwarecomv1.DomainList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested domains.
func (c *FakeDomains) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(domainsResource, opts))
}

// Create takes the representation of a domain and creates it.  Returns the server's representation of the domain, and an error, if there is any.
func (c *FakeDomains) Create(ctx context.Context, domain *configtsmtanzuvmwarecomv1.Domain, opts v1.CreateOptions) (result *configtsmtanzuvmwarecomv1.Domain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(domainsResource, domain), &configtsmtanzuvmwarecomv1.Domain{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.Domain), err
}

// Update takes the representation of a domain and updates it. Returns the server's representation of the domain, and an error, if there is any.
func (c *FakeDomains) Update(ctx context.Context, domain *configtsmtanzuvmwarecomv1.Domain, opts v1.UpdateOptions) (result *configtsmtanzuvmwarecomv1.Domain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(domainsResource, domain), &configtsmtanzuvmwarecomv1.Domain{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.Domain), err
}

// Delete takes name of the domain and deletes it. Returns an error if one occurs.
func (c *FakeDomains) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(domainsResource, name, opts), &configtsmtanzuvmwarecomv1.Domain{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDomains) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(domainsResource, listOpts)

	_, err := c.Fake.Invokes(action, &configtsmtanzuvmwarecomv1.DomainList{})
	return err
}

// Patch applies the patch and returns the patched domain.
func (c *FakeDomains) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *configtsmtanzuvmwarecomv1.Domain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(domainsResource, name, pt, data, subresources...), &configtsmtanzuvmwarecomv1.Domain{})
	if obj == nil {
		return nil, err
	}
	return obj.(*configtsmtanzuvmwarecomv1.Domain), err
}
