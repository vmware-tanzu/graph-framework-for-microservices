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

	confighelloworldcomv1 "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/apis/config.helloworld.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeConfigs implements ConfigInterface
type FakeConfigs struct {
	Fake *FakeConfigHelloworldV1
	ns   string
}

var configsResource = schema.GroupVersionResource{Group: "config.helloworld.com", Version: "v1", Resource: "configs"}

var configsKind = schema.GroupVersionKind{Group: "config.helloworld.com", Version: "v1", Kind: "Config"}

// Get takes name of the config, and returns the corresponding config object, and an error if there is any.
func (c *FakeConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *confighelloworldcomv1.Config, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(configsResource, c.ns, name), &confighelloworldcomv1.Config{})

	if obj == nil {
		return nil, err
	}
	return obj.(*confighelloworldcomv1.Config), err
}

// List takes label and field selectors, and returns the list of Configs that match those selectors.
func (c *FakeConfigs) List(ctx context.Context, opts v1.ListOptions) (result *confighelloworldcomv1.ConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(configsResource, configsKind, c.ns, opts), &confighelloworldcomv1.ConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &confighelloworldcomv1.ConfigList{ListMeta: obj.(*confighelloworldcomv1.ConfigList).ListMeta}
	for _, item := range obj.(*confighelloworldcomv1.ConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested configs.
func (c *FakeConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(configsResource, c.ns, opts))

}

// Create takes the representation of a config and creates it.  Returns the server's representation of the config, and an error, if there is any.
func (c *FakeConfigs) Create(ctx context.Context, config *confighelloworldcomv1.Config, opts v1.CreateOptions) (result *confighelloworldcomv1.Config, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(configsResource, c.ns, config), &confighelloworldcomv1.Config{})

	if obj == nil {
		return nil, err
	}
	return obj.(*confighelloworldcomv1.Config), err
}

// Update takes the representation of a config and updates it. Returns the server's representation of the config, and an error, if there is any.
func (c *FakeConfigs) Update(ctx context.Context, config *confighelloworldcomv1.Config, opts v1.UpdateOptions) (result *confighelloworldcomv1.Config, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(configsResource, c.ns, config), &confighelloworldcomv1.Config{})

	if obj == nil {
		return nil, err
	}
	return obj.(*confighelloworldcomv1.Config), err
}

// Delete takes name of the config and deletes it. Returns an error if one occurs.
func (c *FakeConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(configsResource, c.ns, name), &confighelloworldcomv1.Config{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(configsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &confighelloworldcomv1.ConfigList{})
	return err
}

// Patch applies the patch and returns the patched config.
func (c *FakeConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *confighelloworldcomv1.Config, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(configsResource, c.ns, name, pt, data, subresources...), &confighelloworldcomv1.Config{})

	if obj == nil {
		return nil, err
	}
	return obj.(*confighelloworldcomv1.Config), err
}
