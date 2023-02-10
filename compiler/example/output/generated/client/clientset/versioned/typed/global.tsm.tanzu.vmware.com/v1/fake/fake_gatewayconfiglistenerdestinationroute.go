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

// FakeGatewayConfigListenerDestinationRoutes implements GatewayConfigListenerDestinationRouteInterface
type FakeGatewayConfigListenerDestinationRoutes struct {
	Fake *FakeGlobalTsmV1
}

var gatewayconfiglistenerdestinationroutesResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "gatewayconfiglistenerdestinationroutes"}

var gatewayconfiglistenerdestinationroutesKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "GatewayConfigListenerDestinationRoute"}

// Get takes name of the gatewayConfigListenerDestinationRoute, and returns the corresponding gatewayConfigListenerDestinationRoute object, and an error if there is any.
func (c *FakeGatewayConfigListenerDestinationRoutes) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(gatewayconfiglistenerdestinationroutesResource, name), &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute), err
}

// List takes label and field selectors, and returns the list of GatewayConfigListenerDestinationRoutes that match those selectors.
func (c *FakeGatewayConfigListenerDestinationRoutes) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRouteList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(gatewayconfiglistenerdestinationroutesResource, gatewayconfiglistenerdestinationroutesKind, opts), &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRouteList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRouteList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRouteList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRouteList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gatewayConfigListenerDestinationRoutes.
func (c *FakeGatewayConfigListenerDestinationRoutes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(gatewayconfiglistenerdestinationroutesResource, opts))
}

// Create takes the representation of a gatewayConfigListenerDestinationRoute and creates it.  Returns the server's representation of the gatewayConfigListenerDestinationRoute, and an error, if there is any.
func (c *FakeGatewayConfigListenerDestinationRoutes) Create(ctx context.Context, gatewayConfigListenerDestinationRoute *globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(gatewayconfiglistenerdestinationroutesResource, gatewayConfigListenerDestinationRoute), &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute), err
}

// Update takes the representation of a gatewayConfigListenerDestinationRoute and updates it. Returns the server's representation of the gatewayConfigListenerDestinationRoute, and an error, if there is any.
func (c *FakeGatewayConfigListenerDestinationRoutes) Update(ctx context.Context, gatewayConfigListenerDestinationRoute *globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(gatewayconfiglistenerdestinationroutesResource, gatewayConfigListenerDestinationRoute), &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute), err
}

// Delete takes name of the gatewayConfigListenerDestinationRoute and deletes it. Returns an error if one occurs.
func (c *FakeGatewayConfigListenerDestinationRoutes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(gatewayconfiglistenerdestinationroutesResource, name), &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGatewayConfigListenerDestinationRoutes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(gatewayconfiglistenerdestinationroutesResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRouteList{})
	return err
}

// Patch applies the patch and returns the patched gatewayConfigListenerDestinationRoute.
func (c *FakeGatewayConfigListenerDestinationRoutes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(gatewayconfiglistenerdestinationroutesResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.GatewayConfigListenerDestinationRoute), err
}