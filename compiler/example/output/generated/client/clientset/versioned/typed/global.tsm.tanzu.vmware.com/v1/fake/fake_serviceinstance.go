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

// FakeServiceInstances implements ServiceInstanceInterface
type FakeServiceInstances struct {
	Fake *FakeGlobalTsmV1
}

var serviceinstancesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "serviceinstances"}

var serviceinstancesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ServiceInstance"}

// Get takes name of the serviceInstance, and returns the corresponding serviceInstance object, and an error if there is any.
func (c *FakeServiceInstances) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(serviceinstancesResource, name), &globaltsmtanzuvmwarecomv1.ServiceInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstance), err
}

// List takes label and field selectors, and returns the list of ServiceInstances that match those selectors.
func (c *FakeServiceInstances) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstanceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(serviceinstancesResource, serviceinstancesKind, opts), &globaltsmtanzuvmwarecomv1.ServiceInstanceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ServiceInstanceList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ServiceInstanceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceInstances.
func (c *FakeServiceInstances) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(serviceinstancesResource, opts))
}

// Create takes the representation of a serviceInstance and creates it.  Returns the server's representation of the serviceInstance, and an error, if there is any.
func (c *FakeServiceInstances) Create(ctx context.Context, serviceInstance *globaltsmtanzuvmwarecomv1.ServiceInstance, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(serviceinstancesResource, serviceInstance), &globaltsmtanzuvmwarecomv1.ServiceInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstance), err
}

// Update takes the representation of a serviceInstance and updates it. Returns the server's representation of the serviceInstance, and an error, if there is any.
func (c *FakeServiceInstances) Update(ctx context.Context, serviceInstance *globaltsmtanzuvmwarecomv1.ServiceInstance, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(serviceinstancesResource, serviceInstance), &globaltsmtanzuvmwarecomv1.ServiceInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstance), err
}

// Delete takes name of the serviceInstance and deletes it. Returns an error if one occurs.
func (c *FakeServiceInstances) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(serviceinstancesResource, name), &globaltsmtanzuvmwarecomv1.ServiceInstance{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceInstances) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(serviceinstancesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ServiceInstanceList{})
	return err
}

// Patch applies the patch and returns the patched serviceInstance.
func (c *FakeServiceInstances) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ServiceInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(serviceinstancesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ServiceInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceInstance), err
}