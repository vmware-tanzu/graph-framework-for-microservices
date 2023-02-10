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

// FakeNetworkAttachmentDefinitions implements NetworkAttachmentDefinitionInterface
type FakeNetworkAttachmentDefinitions struct {
	Fake *FakeGlobalTsmV1
}

var networkattachmentdefinitionsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "networkattachmentdefinitions"}

var networkattachmentdefinitionsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "NetworkAttachmentDefinition"}

// Get takes name of the networkAttachmentDefinition, and returns the corresponding networkAttachmentDefinition object, and an error if there is any.
func (c *FakeNetworkAttachmentDefinitions) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(networkattachmentdefinitionsResource, name), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition), err
}

// List takes label and field selectors, and returns the list of NetworkAttachmentDefinitions that match those selectors.
func (c *FakeNetworkAttachmentDefinitions) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(networkattachmentdefinitionsResource, networkattachmentdefinitionsKind, opts), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested networkAttachmentDefinitions.
func (c *FakeNetworkAttachmentDefinitions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(networkattachmentdefinitionsResource, opts))
}

// Create takes the representation of a networkAttachmentDefinition and creates it.  Returns the server's representation of the networkAttachmentDefinition, and an error, if there is any.
func (c *FakeNetworkAttachmentDefinitions) Create(ctx context.Context, networkAttachmentDefinition *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(networkattachmentdefinitionsResource, networkAttachmentDefinition), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition), err
}

// Update takes the representation of a networkAttachmentDefinition and updates it. Returns the server's representation of the networkAttachmentDefinition, and an error, if there is any.
func (c *FakeNetworkAttachmentDefinitions) Update(ctx context.Context, networkAttachmentDefinition *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(networkattachmentdefinitionsResource, networkAttachmentDefinition), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition), err
}

// Delete takes name of the networkAttachmentDefinition and deletes it. Returns an error if one occurs.
func (c *FakeNetworkAttachmentDefinitions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(networkattachmentdefinitionsResource, name), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNetworkAttachmentDefinitions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(networkattachmentdefinitionsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinitionList{})
	return err
}

// Patch applies the patch and returns the patched networkAttachmentDefinition.
func (c *FakeNetworkAttachmentDefinitions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(networkattachmentdefinitionsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.NetworkAttachmentDefinition), err
}