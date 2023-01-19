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

	servicegrouptsmtanzuvmwarecomv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/servicegroup.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSvcGroups implements SvcGroupInterface
type FakeSvcGroups struct {
	Fake *FakeServicegroupTsmV1
}

var svcgroupsResource = schema.GroupVersionResource{Group: "servicegroup.tsm.tanzu.vmware.com", Version: "v1", Resource: "svcgroups"}

var svcgroupsKind = schema.GroupVersionKind{Group: "servicegroup.tsm.tanzu.vmware.com", Version: "v1", Kind: "SvcGroup"}

// Get takes name of the svcGroup, and returns the corresponding svcGroup object, and an error if there is any.
func (c *FakeSvcGroups) Get(ctx context.Context, name string, options v1.GetOptions) (result *servicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(svcgroupsResource, name), &servicegrouptsmtanzuvmwarecomv1.SvcGroup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicegrouptsmtanzuvmwarecomv1.SvcGroup), err
}

// List takes label and field selectors, and returns the list of SvcGroups that match those selectors.
func (c *FakeSvcGroups) List(ctx context.Context, opts v1.ListOptions) (result *servicegrouptsmtanzuvmwarecomv1.SvcGroupList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(svcgroupsResource, svcgroupsKind, opts), &servicegrouptsmtanzuvmwarecomv1.SvcGroupList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &servicegrouptsmtanzuvmwarecomv1.SvcGroupList{ListMeta: obj.(*servicegrouptsmtanzuvmwarecomv1.SvcGroupList).ListMeta}
	for _, item := range obj.(*servicegrouptsmtanzuvmwarecomv1.SvcGroupList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested svcGroups.
func (c *FakeSvcGroups) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(svcgroupsResource, opts))
}

// Create takes the representation of a svcGroup and creates it.  Returns the server's representation of the svcGroup, and an error, if there is any.
func (c *FakeSvcGroups) Create(ctx context.Context, svcGroup *servicegrouptsmtanzuvmwarecomv1.SvcGroup, opts v1.CreateOptions) (result *servicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(svcgroupsResource, svcGroup), &servicegrouptsmtanzuvmwarecomv1.SvcGroup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicegrouptsmtanzuvmwarecomv1.SvcGroup), err
}

// Update takes the representation of a svcGroup and updates it. Returns the server's representation of the svcGroup, and an error, if there is any.
func (c *FakeSvcGroups) Update(ctx context.Context, svcGroup *servicegrouptsmtanzuvmwarecomv1.SvcGroup, opts v1.UpdateOptions) (result *servicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(svcgroupsResource, svcGroup), &servicegrouptsmtanzuvmwarecomv1.SvcGroup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicegrouptsmtanzuvmwarecomv1.SvcGroup), err
}

// Delete takes name of the svcGroup and deletes it. Returns an error if one occurs.
func (c *FakeSvcGroups) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(svcgroupsResource, name), &servicegrouptsmtanzuvmwarecomv1.SvcGroup{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSvcGroups) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(svcgroupsResource, listOpts)

	_, err := c.Fake.Invokes(action, &servicegrouptsmtanzuvmwarecomv1.SvcGroupList{})
	return err
}

// Patch applies the patch and returns the patched svcGroup.
func (c *FakeSvcGroups) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *servicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(svcgroupsResource, name, pt, data, subresources...), &servicegrouptsmtanzuvmwarecomv1.SvcGroup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*servicegrouptsmtanzuvmwarecomv1.SvcGroup), err
}
