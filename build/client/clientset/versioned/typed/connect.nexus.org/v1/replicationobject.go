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

	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/connect.nexus.org/v1"
	scheme "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ReplicationObjectsGetter has a method to return a ReplicationObjectInterface.
// A group's client should implement this interface.
type ReplicationObjectsGetter interface {
	ReplicationObjects() ReplicationObjectInterface
}

// ReplicationObjectInterface has methods to work with ReplicationObject resources.
type ReplicationObjectInterface interface {
	Create(ctx context.Context, replicationObject *v1.ReplicationObject, opts metav1.CreateOptions) (*v1.ReplicationObject, error)
	Update(ctx context.Context, replicationObject *v1.ReplicationObject, opts metav1.UpdateOptions) (*v1.ReplicationObject, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ReplicationObject, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ReplicationObjectList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ReplicationObject, err error)
	ReplicationObjectExpansion
}

// replicationObjects implements ReplicationObjectInterface
type replicationObjects struct {
	client rest.Interface
}

// newReplicationObjects returns a ReplicationObjects
func newReplicationObjects(c *ConnectNexusV1Client) *replicationObjects {
	return &replicationObjects{
		client: c.RESTClient(),
	}
}

// Get takes name of the replicationObject, and returns the corresponding replicationObject object, and an error if there is any.
func (c *replicationObjects) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ReplicationObject, err error) {
	result = &v1.ReplicationObject{}
	err = c.client.Get().
		Resource("replicationobjects").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ReplicationObjects that match those selectors.
func (c *replicationObjects) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ReplicationObjectList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ReplicationObjectList{}
	err = c.client.Get().
		Resource("replicationobjects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested replicationObjects.
func (c *replicationObjects) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("replicationobjects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a replicationObject and creates it.  Returns the server's representation of the replicationObject, and an error, if there is any.
func (c *replicationObjects) Create(ctx context.Context, replicationObject *v1.ReplicationObject, opts metav1.CreateOptions) (result *v1.ReplicationObject, err error) {
	result = &v1.ReplicationObject{}
	err = c.client.Post().
		Resource("replicationobjects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(replicationObject).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a replicationObject and updates it. Returns the server's representation of the replicationObject, and an error, if there is any.
func (c *replicationObjects) Update(ctx context.Context, replicationObject *v1.ReplicationObject, opts metav1.UpdateOptions) (result *v1.ReplicationObject, err error) {
	result = &v1.ReplicationObject{}
	err = c.client.Put().
		Resource("replicationobjects").
		Name(replicationObject.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(replicationObject).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the replicationObject and deletes it. Returns an error if one occurs.
func (c *replicationObjects) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("replicationobjects").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *replicationObjects) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("replicationobjects").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched replicationObject.
func (c *replicationObjects) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ReplicationObject, err error) {
	result = &v1.ReplicationObject{}
	err = c.client.Patch(pt).
		Resource("replicationobjects").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
