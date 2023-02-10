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
	globaltsmtanzuvmwarecomv1 "nexustempmodule/apis/global.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeServices implements ServiceInterface
type FakeServices struct {
	Fake *FakeGlobalTsmV1
}

var servicesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "services"}

var servicesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "Service"}

// Get takes name of the service, and returns the corresponding service object, and an error if there is any.
func (c *FakeServices) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(servicesResource, name), &globaltsmtanzuvmwarecomv1.Service{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Service), err
}

// List takes label and field selectors, and returns the list of Services that match those selectors.
func (c *FakeServices) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(servicesResource, servicesKind, opts), &globaltsmtanzuvmwarecomv1.ServiceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ServiceList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ServiceList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested services.
func (c *FakeServices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(servicesResource, opts))
}

// Create takes the representation of a service and creates it.  Returns the server's representation of the service, and an error, if there is any.
func (c *FakeServices) Create(ctx context.Context, service *globaltsmtanzuvmwarecomv1.Service, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(servicesResource, service), &globaltsmtanzuvmwarecomv1.Service{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Service), err
}

// Update takes the representation of a service and updates it. Returns the server's representation of the service, and an error, if there is any.
func (c *FakeServices) Update(ctx context.Context, service *globaltsmtanzuvmwarecomv1.Service, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(servicesResource, service), &globaltsmtanzuvmwarecomv1.Service{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Service), err
}

// Delete takes name of the service and deletes it. Returns an error if one occurs.
func (c *FakeServices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(servicesResource, name), &globaltsmtanzuvmwarecomv1.Service{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(servicesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ServiceList{})
	return err
}

// Patch applies the patch and returns the patched service.
func (c *FakeServices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(servicesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.Service{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Service), err
}