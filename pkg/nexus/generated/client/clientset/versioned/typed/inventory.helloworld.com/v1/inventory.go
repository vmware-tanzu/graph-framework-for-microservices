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

package v1

import (
	"context"
	"time"

	v1 "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/apis/inventory.helloworld.com/v1"
	scheme "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// InventoriesGetter has a method to return a InventoryInterface.
// A group's client should implement this interface.
type InventoriesGetter interface {
	Inventories(namespace string) InventoryInterface
}

// InventoryInterface has methods to work with Inventory resources.
type InventoryInterface interface {
	Create(ctx context.Context, inventory *v1.Inventory, opts metav1.CreateOptions) (*v1.Inventory, error)
	Update(ctx context.Context, inventory *v1.Inventory, opts metav1.UpdateOptions) (*v1.Inventory, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Inventory, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.InventoryList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Inventory, err error)
	InventoryExpansion
}

// inventories implements InventoryInterface
type inventories struct {
	client rest.Interface
	ns     string
}

// newInventories returns a Inventories
func newInventories(c *InventoryHelloworldV1Client, namespace string) *inventories {
	return &inventories{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the inventory, and returns the corresponding inventory object, and an error if there is any.
func (c *inventories) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Inventory, err error) {
	result = &v1.Inventory{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("inventories").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Inventories that match those selectors.
func (c *inventories) List(ctx context.Context, opts metav1.ListOptions) (result *v1.InventoryList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.InventoryList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("inventories").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested inventories.
func (c *inventories) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("inventories").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a inventory and creates it.  Returns the server's representation of the inventory, and an error, if there is any.
func (c *inventories) Create(ctx context.Context, inventory *v1.Inventory, opts metav1.CreateOptions) (result *v1.Inventory, err error) {
	result = &v1.Inventory{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("inventories").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventory).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a inventory and updates it. Returns the server's representation of the inventory, and an error, if there is any.
func (c *inventories) Update(ctx context.Context, inventory *v1.Inventory, opts metav1.UpdateOptions) (result *v1.Inventory, err error) {
	result = &v1.Inventory{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("inventories").
		Name(inventory.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(inventory).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the inventory and deletes it. Returns an error if one occurs.
func (c *inventories) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("inventories").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *inventories) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("inventories").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched inventory.
func (c *inventories) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Inventory, err error) {
	result = &v1.Inventory{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("inventories").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
