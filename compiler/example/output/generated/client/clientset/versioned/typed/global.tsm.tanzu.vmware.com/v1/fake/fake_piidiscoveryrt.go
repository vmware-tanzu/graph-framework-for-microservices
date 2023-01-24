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

// FakePiiDiscoveryRTs implements PiiDiscoveryRTInterface
type FakePiiDiscoveryRTs struct {
	Fake *FakeGlobalTsmV1
}

var piidiscoveryrtsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "piidiscoveryrts"}

var piidiscoveryrtsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "PiiDiscoveryRT"}

// Get takes name of the piiDiscoveryRT, and returns the corresponding piiDiscoveryRT object, and an error if there is any.
func (c *FakePiiDiscoveryRTs) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.PiiDiscoveryRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(piidiscoveryrtsResource, name), &globaltsmtanzuvmwarecomv1.PiiDiscoveryRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.PiiDiscoveryRT), err
}

// List takes label and field selectors, and returns the list of PiiDiscoveryRTs that match those selectors.
func (c *FakePiiDiscoveryRTs) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.PiiDiscoveryRTList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(piidiscoveryrtsResource, piidiscoveryrtsKind, opts), &globaltsmtanzuvmwarecomv1.PiiDiscoveryRTList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.PiiDiscoveryRTList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.PiiDiscoveryRTList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.PiiDiscoveryRTList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested piiDiscoveryRTs.
func (c *FakePiiDiscoveryRTs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(piidiscoveryrtsResource, opts))
}

// Create takes the representation of a piiDiscoveryRT and creates it.  Returns the server's representation of the piiDiscoveryRT, and an error, if there is any.
func (c *FakePiiDiscoveryRTs) Create(ctx context.Context, piiDiscoveryRT *globaltsmtanzuvmwarecomv1.PiiDiscoveryRT, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.PiiDiscoveryRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(piidiscoveryrtsResource, piiDiscoveryRT), &globaltsmtanzuvmwarecomv1.PiiDiscoveryRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.PiiDiscoveryRT), err
}

// Update takes the representation of a piiDiscoveryRT and updates it. Returns the server's representation of the piiDiscoveryRT, and an error, if there is any.
func (c *FakePiiDiscoveryRTs) Update(ctx context.Context, piiDiscoveryRT *globaltsmtanzuvmwarecomv1.PiiDiscoveryRT, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.PiiDiscoveryRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(piidiscoveryrtsResource, piiDiscoveryRT), &globaltsmtanzuvmwarecomv1.PiiDiscoveryRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.PiiDiscoveryRT), err
}

// Delete takes name of the piiDiscoveryRT and deletes it. Returns an error if one occurs.
func (c *FakePiiDiscoveryRTs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(piidiscoveryrtsResource, name), &globaltsmtanzuvmwarecomv1.PiiDiscoveryRT{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePiiDiscoveryRTs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(piidiscoveryrtsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.PiiDiscoveryRTList{})
	return err
}

// Patch applies the patch and returns the patched piiDiscoveryRT.
func (c *FakePiiDiscoveryRTs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.PiiDiscoveryRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(piidiscoveryrtsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.PiiDiscoveryRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.PiiDiscoveryRT), err
}