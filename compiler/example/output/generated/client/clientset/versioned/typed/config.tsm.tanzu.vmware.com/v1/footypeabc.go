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
	v1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/config.tsm.tanzu.vmware.com/v1"
	scheme "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/client/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// FooTypeABCsGetter has a method to return a FooTypeABCInterface.
// A group's client should implement this interface.
type FooTypeABCsGetter interface {
	FooTypeABCs() FooTypeABCInterface
}

// FooTypeABCInterface has methods to work with FooTypeABC resources.
type FooTypeABCInterface interface {
	Create(ctx context.Context, fooTypeABC *v1.FooTypeABC, opts metav1.CreateOptions) (*v1.FooTypeABC, error)
	Update(ctx context.Context, fooTypeABC *v1.FooTypeABC, opts metav1.UpdateOptions) (*v1.FooTypeABC, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.FooTypeABC, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.FooTypeABCList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.FooTypeABC, err error)
	FooTypeABCExpansion
}

// fooTypeABCs implements FooTypeABCInterface
type fooTypeABCs struct {
	client rest.Interface
}

// newFooTypeABCs returns a FooTypeABCs
func newFooTypeABCs(c *ConfigTsmV1Client) *fooTypeABCs {
	return &fooTypeABCs{
		client: c.RESTClient(),
	}
}

// Get takes name of the fooTypeABC, and returns the corresponding fooTypeABC object, and an error if there is any.
func (c *fooTypeABCs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.FooTypeABC, err error) {
	result = &v1.FooTypeABC{}
	err = c.client.Get().
		Resource("footypeabcs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FooTypeABCs that match those selectors.
func (c *fooTypeABCs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.FooTypeABCList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.FooTypeABCList{}
	err = c.client.Get().
		Resource("footypeabcs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fooTypeABCs.
func (c *fooTypeABCs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("footypeabcs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fooTypeABC and creates it.  Returns the server's representation of the fooTypeABC, and an error, if there is any.
func (c *fooTypeABCs) Create(ctx context.Context, fooTypeABC *v1.FooTypeABC, opts metav1.CreateOptions) (result *v1.FooTypeABC, err error) {
	result = &v1.FooTypeABC{}
	err = c.client.Post().
		Resource("footypeabcs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fooTypeABC).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fooTypeABC and updates it. Returns the server's representation of the fooTypeABC, and an error, if there is any.
func (c *fooTypeABCs) Update(ctx context.Context, fooTypeABC *v1.FooTypeABC, opts metav1.UpdateOptions) (result *v1.FooTypeABC, err error) {
	result = &v1.FooTypeABC{}
	err = c.client.Put().
		Resource("footypeabcs").
		Name(fooTypeABC.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fooTypeABC).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fooTypeABC and deletes it. Returns an error if one occurs.
func (c *fooTypeABCs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("footypeabcs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fooTypeABCs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("footypeabcs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fooTypeABC.
func (c *fooTypeABCs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.FooTypeABC, err error) {
	result = &v1.FooTypeABC{}
	err = c.client.Patch(pt).
		Resource("footypeabcs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
