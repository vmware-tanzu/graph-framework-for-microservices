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

// FakeACPConfigs implements ACPConfigInterface
type FakeACPConfigs struct {
	Fake *FakePolicypkgTsmV1
}

var acpconfigsResource = schema.GroupVersionResource{Group: "policypkg.tsm.tanzu.vmware.com", Version: "v1", Resource: "acpconfigs"}

var acpconfigsKind = schema.GroupVersionKind{Group: "policypkg.tsm.tanzu.vmware.com", Version: "v1", Kind: "ACPConfig"}

// Get takes name of the aCPConfig, and returns the corresponding aCPConfig object, and an error if there is any.
func (c *FakeACPConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *policypkgtsmtanzuvmwarecomv1.ACPConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(acpconfigsResource, name), &policypkgtsmtanzuvmwarecomv1.ACPConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.ACPConfig), err
}

// List takes label and field selectors, and returns the list of ACPConfigs that match those selectors.
func (c *FakeACPConfigs) List(ctx context.Context, opts v1.ListOptions) (result *policypkgtsmtanzuvmwarecomv1.ACPConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(acpconfigsResource, acpconfigsKind, opts), &policypkgtsmtanzuvmwarecomv1.ACPConfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &policypkgtsmtanzuvmwarecomv1.ACPConfigList{ListMeta: obj.(*policypkgtsmtanzuvmwarecomv1.ACPConfigList).ListMeta}
	for _, item := range obj.(*policypkgtsmtanzuvmwarecomv1.ACPConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested aCPConfigs.
func (c *FakeACPConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(acpconfigsResource, opts))
}

// Create takes the representation of a aCPConfig and creates it.  Returns the server's representation of the aCPConfig, and an error, if there is any.
func (c *FakeACPConfigs) Create(ctx context.Context, aCPConfig *policypkgtsmtanzuvmwarecomv1.ACPConfig, opts v1.CreateOptions) (result *policypkgtsmtanzuvmwarecomv1.ACPConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(acpconfigsResource, aCPConfig), &policypkgtsmtanzuvmwarecomv1.ACPConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.ACPConfig), err
}

// Update takes the representation of a aCPConfig and updates it. Returns the server's representation of the aCPConfig, and an error, if there is any.
func (c *FakeACPConfigs) Update(ctx context.Context, aCPConfig *policypkgtsmtanzuvmwarecomv1.ACPConfig, opts v1.UpdateOptions) (result *policypkgtsmtanzuvmwarecomv1.ACPConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(acpconfigsResource, aCPConfig), &policypkgtsmtanzuvmwarecomv1.ACPConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.ACPConfig), err
}

// Delete takes name of the aCPConfig and deletes it. Returns an error if one occurs.
func (c *FakeACPConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(acpconfigsResource, name), &policypkgtsmtanzuvmwarecomv1.ACPConfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeACPConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(acpconfigsResource, listOpts)

	_, err := c.Fake.Invokes(action, &policypkgtsmtanzuvmwarecomv1.ACPConfigList{})
	return err
}

// Patch applies the patch and returns the patched aCPConfig.
func (c *FakeACPConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *policypkgtsmtanzuvmwarecomv1.ACPConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(acpconfigsResource, name, pt, data, subresources...), &policypkgtsmtanzuvmwarecomv1.ACPConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.ACPConfig), err
}
