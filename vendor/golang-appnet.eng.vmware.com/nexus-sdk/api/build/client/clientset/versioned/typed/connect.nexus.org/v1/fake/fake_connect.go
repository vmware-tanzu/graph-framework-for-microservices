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

	connectnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeConnects implements ConnectInterface
type FakeConnects struct {
	Fake *FakeConnectNexusV1
}

var connectsResource = schema.GroupVersionResource{Group: "connect.nexus.org", Version: "v1", Resource: "connects"}

var connectsKind = schema.GroupVersionKind{Group: "connect.nexus.org", Version: "v1", Kind: "Connect"}

// Get takes name of the connect, and returns the corresponding connect object, and an error if there is any.
func (c *FakeConnects) Get(ctx context.Context, name string, options v1.GetOptions) (result *connectnexusorgv1.Connect, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(connectsResource, name), &connectnexusorgv1.Connect{})
	if obj == nil {
		return nil, err
	}
	return obj.(*connectnexusorgv1.Connect), err
}

// List takes label and field selectors, and returns the list of Connects that match those selectors.
func (c *FakeConnects) List(ctx context.Context, opts v1.ListOptions) (result *connectnexusorgv1.ConnectList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(connectsResource, connectsKind, opts), &connectnexusorgv1.ConnectList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &connectnexusorgv1.ConnectList{ListMeta: obj.(*connectnexusorgv1.ConnectList).ListMeta}
	for _, item := range obj.(*connectnexusorgv1.ConnectList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested connects.
func (c *FakeConnects) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(connectsResource, opts))
}

// Create takes the representation of a connect and creates it.  Returns the server's representation of the connect, and an error, if there is any.
func (c *FakeConnects) Create(ctx context.Context, connect *connectnexusorgv1.Connect, opts v1.CreateOptions) (result *connectnexusorgv1.Connect, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(connectsResource, connect), &connectnexusorgv1.Connect{})
	if obj == nil {
		return nil, err
	}
	return obj.(*connectnexusorgv1.Connect), err
}

// Update takes the representation of a connect and updates it. Returns the server's representation of the connect, and an error, if there is any.
func (c *FakeConnects) Update(ctx context.Context, connect *connectnexusorgv1.Connect, opts v1.UpdateOptions) (result *connectnexusorgv1.Connect, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(connectsResource, connect), &connectnexusorgv1.Connect{})
	if obj == nil {
		return nil, err
	}
	return obj.(*connectnexusorgv1.Connect), err
}

// Delete takes name of the connect and deletes it. Returns an error if one occurs.
func (c *FakeConnects) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(connectsResource, name), &connectnexusorgv1.Connect{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeConnects) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(connectsResource, listOpts)

	_, err := c.Fake.Invokes(action, &connectnexusorgv1.ConnectList{})
	return err
}

// Patch applies the patch and returns the patched connect.
func (c *FakeConnects) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *connectnexusorgv1.Connect, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(connectsResource, name, pt, data, subresources...), &connectnexusorgv1.Connect{})
	if obj == nil {
		return nil, err
	}
	return obj.(*connectnexusorgv1.Connect), err
}
