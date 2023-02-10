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

// FakeProgressiveUpgrades implements ProgressiveUpgradeInterface
type FakeProgressiveUpgrades struct {
	Fake *FakeGlobalTsmV1
}

var progressiveupgradesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "progressiveupgrades"}

var progressiveupgradesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "ProgressiveUpgrade"}

// Get takes name of the progressiveUpgrade, and returns the corresponding progressiveUpgrade object, and an error if there is any.
func (c *FakeProgressiveUpgrades) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgrade, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(progressiveupgradesResource, name), &globaltsmtanzuvmwarecomv1.ProgressiveUpgrade{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgrade), err
}

// List takes label and field selectors, and returns the list of ProgressiveUpgrades that match those selectors.
func (c *FakeProgressiveUpgrades) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgradeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(progressiveupgradesResource, progressiveupgradesKind, opts), &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgradeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested progressiveUpgrades.
func (c *FakeProgressiveUpgrades) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(progressiveupgradesResource, opts))
}

// Create takes the representation of a progressiveUpgrade and creates it.  Returns the server's representation of the progressiveUpgrade, and an error, if there is any.
func (c *FakeProgressiveUpgrades) Create(ctx context.Context, progressiveUpgrade *globaltsmtanzuvmwarecomv1.ProgressiveUpgrade, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgrade, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(progressiveupgradesResource, progressiveUpgrade), &globaltsmtanzuvmwarecomv1.ProgressiveUpgrade{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgrade), err
}

// Update takes the representation of a progressiveUpgrade and updates it. Returns the server's representation of the progressiveUpgrade, and an error, if there is any.
func (c *FakeProgressiveUpgrades) Update(ctx context.Context, progressiveUpgrade *globaltsmtanzuvmwarecomv1.ProgressiveUpgrade, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgrade, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(progressiveupgradesResource, progressiveUpgrade), &globaltsmtanzuvmwarecomv1.ProgressiveUpgrade{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgrade), err
}

// Delete takes name of the progressiveUpgrade and deletes it. Returns an error if one occurs.
func (c *FakeProgressiveUpgrades) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(progressiveupgradesResource, name), &globaltsmtanzuvmwarecomv1.ProgressiveUpgrade{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeProgressiveUpgrades) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(progressiveupgradesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.ProgressiveUpgradeList{})
	return err
}

// Patch applies the patch and returns the patched progressiveUpgrade.
func (c *FakeProgressiveUpgrades) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.ProgressiveUpgrade, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(progressiveupgradesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.ProgressiveUpgrade{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.ProgressiveUpgrade), err
}