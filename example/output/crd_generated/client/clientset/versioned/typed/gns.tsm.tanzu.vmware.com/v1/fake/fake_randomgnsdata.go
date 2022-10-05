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
	gnstsmtanzuvmwarecomv1 "nexustempmodule/apis/gns.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRandomGnsDatas implements RandomGnsDataInterface
type FakeRandomGnsDatas struct {
	Fake *FakeGnsTsmV1
}

var randomgnsdatasResource = schema.GroupVersionResource{Group: "gns.tsm.tanzu.vmware.com", Version: "v1", Resource: "randomgnsdatas"}

var randomgnsdatasKind = schema.GroupVersionKind{Group: "gns.tsm.tanzu.vmware.com", Version: "v1", Kind: "RandomGnsData"}

// Get takes name of the randomGnsData, and returns the corresponding randomGnsData object, and an error if there is any.
func (c *FakeRandomGnsDatas) Get(ctx context.Context, name string, options v1.GetOptions) (result *gnstsmtanzuvmwarecomv1.RandomGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(randomgnsdatasResource, name), &gnstsmtanzuvmwarecomv1.RandomGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.RandomGnsData), err
}

// List takes label and field selectors, and returns the list of RandomGnsDatas that match those selectors.
func (c *FakeRandomGnsDatas) List(ctx context.Context, opts v1.ListOptions) (result *gnstsmtanzuvmwarecomv1.RandomGnsDataList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(randomgnsdatasResource, randomgnsdatasKind, opts), &gnstsmtanzuvmwarecomv1.RandomGnsDataList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &gnstsmtanzuvmwarecomv1.RandomGnsDataList{ListMeta: obj.(*gnstsmtanzuvmwarecomv1.RandomGnsDataList).ListMeta}
	for _, item := range obj.(*gnstsmtanzuvmwarecomv1.RandomGnsDataList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested randomGnsDatas.
func (c *FakeRandomGnsDatas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(randomgnsdatasResource, opts))
}

// Create takes the representation of a randomGnsData and creates it.  Returns the server's representation of the randomGnsData, and an error, if there is any.
func (c *FakeRandomGnsDatas) Create(ctx context.Context, randomGnsData *gnstsmtanzuvmwarecomv1.RandomGnsData, opts v1.CreateOptions) (result *gnstsmtanzuvmwarecomv1.RandomGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(randomgnsdatasResource, randomGnsData), &gnstsmtanzuvmwarecomv1.RandomGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.RandomGnsData), err
}

// Update takes the representation of a randomGnsData and updates it. Returns the server's representation of the randomGnsData, and an error, if there is any.
func (c *FakeRandomGnsDatas) Update(ctx context.Context, randomGnsData *gnstsmtanzuvmwarecomv1.RandomGnsData, opts v1.UpdateOptions) (result *gnstsmtanzuvmwarecomv1.RandomGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(randomgnsdatasResource, randomGnsData), &gnstsmtanzuvmwarecomv1.RandomGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.RandomGnsData), err
}

// Delete takes name of the randomGnsData and deletes it. Returns an error if one occurs.
func (c *FakeRandomGnsDatas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(randomgnsdatasResource, name), &gnstsmtanzuvmwarecomv1.RandomGnsData{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRandomGnsDatas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(randomgnsdatasResource, listOpts)

	_, err := c.Fake.Invokes(action, &gnstsmtanzuvmwarecomv1.RandomGnsDataList{})
	return err
}

// Patch applies the patch and returns the patched randomGnsData.
func (c *FakeRandomGnsDatas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *gnstsmtanzuvmwarecomv1.RandomGnsData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(randomgnsdatasResource, name, pt, data, subresources...), &gnstsmtanzuvmwarecomv1.RandomGnsData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.RandomGnsData), err
}