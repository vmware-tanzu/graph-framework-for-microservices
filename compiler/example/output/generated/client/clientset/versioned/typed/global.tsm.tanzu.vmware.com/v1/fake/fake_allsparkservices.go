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

// FakeAllSparkServiceses implements AllSparkServicesInterface
type FakeAllSparkServiceses struct {
	Fake *FakeGlobalTsmV1
}

var allsparkservicesesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "allsparkserviceses"}

var allsparkservicesesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "AllSparkServices"}

// Get takes name of the allSparkServices, and returns the corresponding allSparkServices object, and an error if there is any.
func (c *FakeAllSparkServiceses) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.AllSparkServices, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(allsparkservicesesResource, name), &globaltsmtanzuvmwarecomv1.AllSparkServices{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AllSparkServices), err
}

// List takes label and field selectors, and returns the list of AllSparkServiceses that match those selectors.
func (c *FakeAllSparkServiceses) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.AllSparkServicesList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(allsparkservicesesResource, allsparkservicesesKind, opts), &globaltsmtanzuvmwarecomv1.AllSparkServicesList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.AllSparkServicesList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.AllSparkServicesList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.AllSparkServicesList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested allSparkServiceses.
func (c *FakeAllSparkServiceses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(allsparkservicesesResource, opts))
}

// Create takes the representation of a allSparkServices and creates it.  Returns the server's representation of the allSparkServices, and an error, if there is any.
func (c *FakeAllSparkServiceses) Create(ctx context.Context, allSparkServices *globaltsmtanzuvmwarecomv1.AllSparkServices, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.AllSparkServices, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(allsparkservicesesResource, allSparkServices), &globaltsmtanzuvmwarecomv1.AllSparkServices{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AllSparkServices), err
}

// Update takes the representation of a allSparkServices and updates it. Returns the server's representation of the allSparkServices, and an error, if there is any.
func (c *FakeAllSparkServiceses) Update(ctx context.Context, allSparkServices *globaltsmtanzuvmwarecomv1.AllSparkServices, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.AllSparkServices, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(allsparkservicesesResource, allSparkServices), &globaltsmtanzuvmwarecomv1.AllSparkServices{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AllSparkServices), err
}

// Delete takes name of the allSparkServices and deletes it. Returns an error if one occurs.
func (c *FakeAllSparkServiceses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(allsparkservicesesResource, name), &globaltsmtanzuvmwarecomv1.AllSparkServices{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAllSparkServiceses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(allsparkservicesesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.AllSparkServicesList{})
	return err
}

// Patch applies the patch and returns the patched allSparkServices.
func (c *FakeAllSparkServiceses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.AllSparkServices, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(allsparkservicesesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.AllSparkServices{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AllSparkServices), err
}