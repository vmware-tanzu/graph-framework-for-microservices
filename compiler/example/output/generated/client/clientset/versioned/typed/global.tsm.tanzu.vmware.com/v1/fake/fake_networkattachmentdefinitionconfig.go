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

// FakeNetworkAttachmentDefinitionConfigs implements NetworkAttachmentDefinitionConfigInterface
type FakeNetworkAttachmentDefinitionConfigs struct {
	Fake *FakeGlobalTsmV1
}

var networkattachmentdefinitionconfigsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "networkattachmentdefinitionconfigs"}

var networkattachmentdefinitionconfigsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "NetworkAttachmentDefinitionConfig"}

// Get takes name of the networkAttachmentDefinitionConfig, and returns the corresponding networkAttachmentDefinitionConfig object, and an error if there is any.
func (c *FakeNetworkAttachmentDefinitionConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(networkattachmentdefinitionconfigsResource, name), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig), err
}

// List takes label and field selectors, and returns the list of NetworkAttachmentDefinitionConfigs that match those selectors.
func (c *FakeNetworkAttachmentDefinitionConfigs) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(networkattachmentdefinitionconfigsResource, networkattachmentdefinitionconfigsKind, opts), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfigList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfigList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested networkAttachmentDefinitionConfigs.
func (c *FakeNetworkAttachmentDefinitionConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(networkattachmentdefinitionconfigsResource, opts))
}

// Create takes the representation of a networkAttachmentDefinitionConfig and creates it.  Returns the server's representation of the networkAttachmentDefinitionConfig, and an error, if there is any.
func (c *FakeNetworkAttachmentDefinitionConfigs) Create(ctx context.Context, networkAttachmentDefinitionConfig *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(networkattachmentdefinitionconfigsResource, networkAttachmentDefinitionConfig), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig), err
}

// Update takes the representation of a networkAttachmentDefinitionConfig and updates it. Returns the server's representation of the networkAttachmentDefinitionConfig, and an error, if there is any.
func (c *FakeNetworkAttachmentDefinitionConfigs) Update(ctx context.Context, networkAttachmentDefinitionConfig *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(networkattachmentdefinitionconfigsResource, networkAttachmentDefinitionConfig), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig), err
}

// Delete takes name of the networkAttachmentDefinitionConfig and deletes it. Returns an error if one occurs.
func (c *FakeNetworkAttachmentDefinitionConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(networkattachmentdefinitionconfigsResource, name), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNetworkAttachmentDefinitionConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(networkattachmentdefinitionconfigsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfigList{})
	return err
}

// Patch applies the patch and returns the patched networkAttachmentDefinitionConfig.
func (c *FakeNetworkAttachmentDefinitionConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(networkattachmentdefinitionconfigsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionConfig), err
}