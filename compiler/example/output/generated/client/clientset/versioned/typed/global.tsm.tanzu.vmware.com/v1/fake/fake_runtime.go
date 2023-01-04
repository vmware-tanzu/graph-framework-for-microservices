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

// FakeRuntimes implements RuntimeInterface
type FakeRuntimes struct {
	Fake *FakeGlobalTsmV1
}

var runtimesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "runtimes"}

var runtimesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "Runtime"}

// Get takes name of the runtime, and returns the corresponding runtime object, and an error if there is any.
func (c *FakeRuntimes) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.Runtime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(runtimesResource, name), &globaltsmtanzuvmwarecomv1.Runtime{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Runtime), err
}

// List takes label and field selectors, and returns the list of Runtimes that match those selectors.
func (c *FakeRuntimes) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.RuntimeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(runtimesResource, runtimesKind, opts), &globaltsmtanzuvmwarecomv1.RuntimeList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.RuntimeList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.RuntimeList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.RuntimeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested runtimes.
func (c *FakeRuntimes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(runtimesResource, opts))
}

// Create takes the representation of a runtime and creates it.  Returns the server's representation of the runtime, and an error, if there is any.
func (c *FakeRuntimes) Create(ctx context.Context, runtime *globaltsmtanzuvmwarecomv1.Runtime, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.Runtime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(runtimesResource, runtime), &globaltsmtanzuvmwarecomv1.Runtime{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Runtime), err
}

// Update takes the representation of a runtime and updates it. Returns the server's representation of the runtime, and an error, if there is any.
func (c *FakeRuntimes) Update(ctx context.Context, runtime *globaltsmtanzuvmwarecomv1.Runtime, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.Runtime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(runtimesResource, runtime), &globaltsmtanzuvmwarecomv1.Runtime{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Runtime), err
}

// Delete takes name of the runtime and deletes it. Returns an error if one occurs.
func (c *FakeRuntimes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(runtimesResource, name), &globaltsmtanzuvmwarecomv1.Runtime{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRuntimes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(runtimesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.RuntimeList{})
	return err
}

// Patch applies the patch and returns the patched runtime.
func (c *FakeRuntimes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.Runtime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(runtimesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.Runtime{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.Runtime), err
}