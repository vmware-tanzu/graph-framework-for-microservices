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

// FakeProgressiveUpgradeConfigs implements ProgressiveUpgradeConfigInterface
type FakeProgressiveUpgradeConfigs struct {
	Fake *FakeGlobalTsmV1
}

var progressiveupgradeconfigsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "progressiveupgradeconfigs"}

var progressiveupgradeconfigsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ProgressiveUpgradeConfig"}

// Get takes name of the progressiveUpgradeConfig, and returns the corresponding progressiveUpgradeConfig object, and an error if there is any.
func (c *FakeProgressiveUpgradeConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(progressiveupgradeconfigsResource, name), &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig), err
}

// List takes label and field selectors, and returns the list of ProgressiveUpgradeConfigs that match those selectors.
func (c *FakeProgressiveUpgradeConfigs) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(progressiveupgradeconfigsResource, progressiveupgradeconfigsKind, opts), &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfigList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfigList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested progressiveUpgradeConfigs.
func (c *FakeProgressiveUpgradeConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(progressiveupgradeconfigsResource, opts))
}

// Create takes the representation of a progressiveUpgradeConfig and creates it.  Returns the server's representation of the progressiveUpgradeConfig, and an error, if there is any.
func (c *FakeProgressiveUpgradeConfigs) Create(ctx context.Context, progressiveUpgradeConfig *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(progressiveupgradeconfigsResource, progressiveUpgradeConfig), &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig), err
}

// Update takes the representation of a progressiveUpgradeConfig and updates it. Returns the server's representation of the progressiveUpgradeConfig, and an error, if there is any.
func (c *FakeProgressiveUpgradeConfigs) Update(ctx context.Context, progressiveUpgradeConfig *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(progressiveupgradeconfigsResource, progressiveUpgradeConfig), &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig), err
}

// Delete takes name of the progressiveUpgradeConfig and deletes it. Returns an error if one occurs.
func (c *FakeProgressiveUpgradeConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(progressiveupgradeconfigsResource, name), &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeProgressiveUpgradeConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(progressiveupgradeconfigsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfigList{})
	return err
}

// Patch applies the patch and returns the patched progressiveUpgradeConfig.
func (c *FakeProgressiveUpgradeConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(progressiveupgradeconfigsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeConfig), err
}