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

	gnstsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generatedapis/gns.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDnses implements DnsInterface
type FakeDnses struct {
	Fake *FakeGnsTsmV1
}

var dnsesResource = schema.GroupVersionResource{Group: "gns.tsm.tanzu.vmware.com", Version: "v1", Resource: "dnses"}

var dnsesKind = schema.GroupVersionKind{Group: "gns.tsm.tanzu.vmware.com", Version: "v1", Kind: "Dns"}

// Get takes name of the dns, and returns the corresponding dns object, and an error if there is any.
func (c *FakeDnses) Get(ctx context.Context, name string, options v1.GetOptions) (result *gnstsmtanzuvmwarecomv1.Dns, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(dnsesResource, name), &gnstsmtanzuvmwarecomv1.Dns{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.Dns), err
}

// List takes label and field selectors, and returns the list of Dnses that match those selectors.
func (c *FakeDnses) List(ctx context.Context, opts v1.ListOptions) (result *gnstsmtanzuvmwarecomv1.DnsList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(dnsesResource, dnsesKind, opts), &gnstsmtanzuvmwarecomv1.DnsList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &gnstsmtanzuvmwarecomv1.DnsList{ListMeta: obj.(*gnstsmtanzuvmwarecomv1.DnsList).ListMeta}
	for _, item := range obj.(*gnstsmtanzuvmwarecomv1.DnsList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dnses.
func (c *FakeDnses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(dnsesResource, opts))
}

// Create takes the representation of a dns and creates it.  Returns the server's representation of the dns, and an error, if there is any.
func (c *FakeDnses) Create(ctx context.Context, dns *gnstsmtanzuvmwarecomv1.Dns, opts v1.CreateOptions) (result *gnstsmtanzuvmwarecomv1.Dns, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(dnsesResource, dns), &gnstsmtanzuvmwarecomv1.Dns{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.Dns), err
}

// Update takes the representation of a dns and updates it. Returns the server's representation of the dns, and an error, if there is any.
func (c *FakeDnses) Update(ctx context.Context, dns *gnstsmtanzuvmwarecomv1.Dns, opts v1.UpdateOptions) (result *gnstsmtanzuvmwarecomv1.Dns, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(dnsesResource, dns), &gnstsmtanzuvmwarecomv1.Dns{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.Dns), err
}

// Delete takes name of the dns and deletes it. Returns an error if one occurs.
func (c *FakeDnses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(dnsesResource, name), &gnstsmtanzuvmwarecomv1.Dns{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDnses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(dnsesResource, listOpts)

	_, err := c.Fake.Invokes(action, &gnstsmtanzuvmwarecomv1.DnsList{})
	return err
}

// Patch applies the patch and returns the patched dns.
func (c *FakeDnses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *gnstsmtanzuvmwarecomv1.Dns, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(dnsesResource, name, pt, data, subresources...), &gnstsmtanzuvmwarecomv1.Dns{})
	if obj == nil {
		return nil, err
	}
	return obj.(*gnstsmtanzuvmwarecomv1.Dns), err
}