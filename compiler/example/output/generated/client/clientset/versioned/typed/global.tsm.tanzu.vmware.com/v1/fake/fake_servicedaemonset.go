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

// FakeServiceDaemonSets implements ServiceDaemonSetInterface
type FakeServiceDaemonSets struct {
	Fake *FakeGlobalTsmV1
}

var servicedaemonsetsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "servicedaemonsets"}

var servicedaemonsetsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ServiceDaemonSet"}

// Get takes name of the serviceDaemonSet, and returns the corresponding serviceDaemonSet object, and an error if there is any.
func (c *FakeServiceDaemonSets) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDaemonSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(servicedaemonsetsResource, name), &globaltsmtanzuvmwarecomv1.ServiceDaemonSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDaemonSet), err
}

// List takes label and field selectors, and returns the list of ServiceDaemonSets that match those selectors.
func (c *FakeServiceDaemonSets) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDaemonSetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(servicedaemonsetsResource, servicedaemonsetsKind, opts), &globaltsmtanzuvmwarecomv1.ServiceDaemonSetList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ServiceDaemonSetList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ServiceDaemonSetList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ServiceDaemonSetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceDaemonSets.
func (c *FakeServiceDaemonSets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(servicedaemonsetsResource, opts))
}

// Create takes the representation of a serviceDaemonSet and creates it.  Returns the server's representation of the serviceDaemonSet, and an error, if there is any.
func (c *FakeServiceDaemonSets) Create(ctx context.Context, serviceDaemonSet *globaltsmtanzuvmwarecomv1.ServiceDaemonSet, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDaemonSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(servicedaemonsetsResource, serviceDaemonSet), &globaltsmtanzuvmwarecomv1.ServiceDaemonSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDaemonSet), err
}

// Update takes the representation of a serviceDaemonSet and updates it. Returns the server's representation of the serviceDaemonSet, and an error, if there is any.
func (c *FakeServiceDaemonSets) Update(ctx context.Context, serviceDaemonSet *globaltsmtanzuvmwarecomv1.ServiceDaemonSet, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ServiceDaemonSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(servicedaemonsetsResource, serviceDaemonSet), &globaltsmtanzuvmwarecomv1.ServiceDaemonSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDaemonSet), err
}

// Delete takes name of the serviceDaemonSet and deletes it. Returns an error if one occurs.
func (c *FakeServiceDaemonSets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(servicedaemonsetsResource, name), &globaltsmtanzuvmwarecomv1.ServiceDaemonSet{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceDaemonSets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(servicedaemonsetsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ServiceDaemonSetList{})
	return err
}

// Patch applies the patch and returns the patched serviceDaemonSet.
func (c *FakeServiceDaemonSets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ServiceDaemonSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(servicedaemonsetsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ServiceDaemonSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ServiceDaemonSet), err
}