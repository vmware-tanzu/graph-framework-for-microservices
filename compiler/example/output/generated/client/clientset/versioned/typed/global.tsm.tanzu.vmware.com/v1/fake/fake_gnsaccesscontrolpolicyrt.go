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

// FakeGnsAccessControlPolicyRTs implements GnsAccessControlPolicyRTInterface
type FakeGnsAccessControlPolicyRTs struct {
	Fake *FakeGlobalTsmV1
}

var gnsaccesscontrolpolicyrtsResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "gnsaccesscontrolpolicyrts"}

var gnsaccesscontrolpolicyrtsKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "GnsAccessControlPolicyRT"}

// Get takes name of the gnsAccessControlPolicyRT, and returns the corresponding gnsAccessControlPolicyRT object, and an error if there is any.
func (c *FakeGnsAccessControlPolicyRTs) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(gnsaccesscontrolpolicyrtsResource, name), &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT), err
}

// List takes label and field selectors, and returns the list of GnsAccessControlPolicyRTs that match those selectors.
func (c *FakeGnsAccessControlPolicyRTs) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRTList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(gnsaccesscontrolpolicyrtsResource, gnsaccesscontrolpolicyrtsKind, opts), &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRTList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRTList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRTList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRTList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gnsAccessControlPolicyRTs.
func (c *FakeGnsAccessControlPolicyRTs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(gnsaccesscontrolpolicyrtsResource, opts))
}

// Create takes the representation of a gnsAccessControlPolicyRT and creates it.  Returns the server's representation of the gnsAccessControlPolicyRT, and an error, if there is any.
func (c *FakeGnsAccessControlPolicyRTs) Create(ctx context.Context, gnsAccessControlPolicyRT *globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(gnsaccesscontrolpolicyrtsResource, gnsAccessControlPolicyRT), &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT), err
}

// Update takes the representation of a gnsAccessControlPolicyRT and updates it. Returns the server's representation of the gnsAccessControlPolicyRT, and an error, if there is any.
func (c *FakeGnsAccessControlPolicyRTs) Update(ctx context.Context, gnsAccessControlPolicyRT *globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(gnsaccesscontrolpolicyrtsResource, gnsAccessControlPolicyRT), &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT), err
}

// Delete takes name of the gnsAccessControlPolicyRT and deletes it. Returns an error if one occurs.
func (c *FakeGnsAccessControlPolicyRTs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(gnsaccesscontrolpolicyrtsResource, name), &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGnsAccessControlPolicyRTs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(gnsaccesscontrolpolicyrtsResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRTList{})
	return err
}

// Patch applies the patch and returns the patched gnsAccessControlPolicyRT.
func (c *FakeGnsAccessControlPolicyRTs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(gnsaccesscontrolpolicyrtsResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GnsAccessControlPolicyRT), err
}