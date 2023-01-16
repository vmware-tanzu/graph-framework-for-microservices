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

	policypkgtsmtanzuvmwarecomv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/policypkg.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRandomPolicyDatas implements RandomPolicyDataInterface
type FakeRandomPolicyDatas struct {
	Fake *FakePolicypkgTsmV1
}

var randompolicydatasResource = schema.GroupVersionResource{Group: "policypkg.tsm.tanzu.vmware.com", Version: "v1", Resource: "randompolicydatas"}

var randompolicydatasKind = schema.GroupVersionKind{Group: "policypkg.tsm.tanzu.vmware.com", Version: "v1", Kind: "RandomPolicyData"}

// Get takes name of the randomPolicyData, and returns the corresponding randomPolicyData object, and an error if there is any.
func (c *FakeRandomPolicyDatas) Get(ctx context.Context, name string, options v1.GetOptions) (result *policypkgtsmtanzuvmwarecomv1.RandomPolicyData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(randompolicydatasResource, name), &policypkgtsmtanzuvmwarecomv1.RandomPolicyData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.RandomPolicyData), err
}

// List takes label and field selectors, and returns the list of RandomPolicyDatas that match those selectors.
func (c *FakeRandomPolicyDatas) List(ctx context.Context, opts v1.ListOptions) (result *policypkgtsmtanzuvmwarecomv1.RandomPolicyDataList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(randompolicydatasResource, randompolicydatasKind, opts), &policypkgtsmtanzuvmwarecomv1.RandomPolicyDataList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &policypkgtsmtanzuvmwarecomv1.RandomPolicyDataList{ListMeta: obj.(*policypkgtsmtanzuvmwarecomv1.RandomPolicyDataList).ListMeta}
	for _, item := range obj.(*policypkgtsmtanzuvmwarecomv1.RandomPolicyDataList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested randomPolicyDatas.
func (c *FakeRandomPolicyDatas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(randompolicydatasResource, opts))
}

// Create takes the representation of a randomPolicyData and creates it.  Returns the server's representation of the randomPolicyData, and an error, if there is any.
func (c *FakeRandomPolicyDatas) Create(ctx context.Context, randomPolicyData *policypkgtsmtanzuvmwarecomv1.RandomPolicyData, opts v1.CreateOptions) (result *policypkgtsmtanzuvmwarecomv1.RandomPolicyData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(randompolicydatasResource, randomPolicyData), &policypkgtsmtanzuvmwarecomv1.RandomPolicyData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.RandomPolicyData), err
}

// Update takes the representation of a randomPolicyData and updates it. Returns the server's representation of the randomPolicyData, and an error, if there is any.
func (c *FakeRandomPolicyDatas) Update(ctx context.Context, randomPolicyData *policypkgtsmtanzuvmwarecomv1.RandomPolicyData, opts v1.UpdateOptions) (result *policypkgtsmtanzuvmwarecomv1.RandomPolicyData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(randompolicydatasResource, randomPolicyData), &policypkgtsmtanzuvmwarecomv1.RandomPolicyData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.RandomPolicyData), err
}

// Delete takes name of the randomPolicyData and deletes it. Returns an error if one occurs.
func (c *FakeRandomPolicyDatas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(randompolicydatasResource, name), &policypkgtsmtanzuvmwarecomv1.RandomPolicyData{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRandomPolicyDatas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(randompolicydatasResource, listOpts)

	_, err := c.Fake.Invokes(action, &policypkgtsmtanzuvmwarecomv1.RandomPolicyDataList{})
	return err
}

// Patch applies the patch and returns the patched randomPolicyData.
func (c *FakeRandomPolicyDatas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *policypkgtsmtanzuvmwarecomv1.RandomPolicyData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(randompolicydatasResource, name, pt, data, subresources...), &policypkgtsmtanzuvmwarecomv1.RandomPolicyData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.RandomPolicyData), err
}
