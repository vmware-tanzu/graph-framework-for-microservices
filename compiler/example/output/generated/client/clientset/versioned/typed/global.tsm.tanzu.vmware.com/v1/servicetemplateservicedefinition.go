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
	v1 "nexustempmodule/apis/global.tsm.tanzu.vmware.com/v1"
	scheme "nexustempmodule/client/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ServiceTemplateServiceDefinitionsGetter has a method to return a ServiceTemplateServiceDefinitionInterface.
// A group's client should implement this interface.
type ServiceTemplateServiceDefinitionsGetter interface {
	ServiceTemplateServiceDefinitions() ServiceTemplateServiceDefinitionInterface
}

// ServiceTemplateServiceDefinitionInterface has methods to work with ServiceTemplateServiceDefinition resources.
type ServiceTemplateServiceDefinitionInterface interface {
	Create(ctx context.Context, serviceTemplateServiceDefinition *v1.ServiceTemplateServiceDefinition, opts metav1.CreateOptions) (*v1.ServiceTemplateServiceDefinition, error)
	Update(ctx context.Context, serviceTemplateServiceDefinition *v1.ServiceTemplateServiceDefinition, opts metav1.UpdateOptions) (*v1.ServiceTemplateServiceDefinition, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ServiceTemplateServiceDefinition, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ServiceTemplateServiceDefinitionList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServiceTemplateServiceDefinition, err error)
	ServiceTemplateServiceDefinitionExpansion
}

// serviceTemplateServiceDefinitions implements ServiceTemplateServiceDefinitionInterface
type serviceTemplateServiceDefinitions struct {
	client rest.Interface
}

// newServiceTemplateServiceDefinitions returns a ServiceTemplateServiceDefinitions
func newServiceTemplateServiceDefinitions(c *GlobalTsmV1Client) *serviceTemplateServiceDefinitions {
	return &serviceTemplateServiceDefinitions{
		client: c.RESTClient(),
	}
}

// Get takes name of the serviceTemplateServiceDefinition, and returns the corresponding serviceTemplateServiceDefinition object, and an error if there is any.
func (c *serviceTemplateServiceDefinitions) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ServiceTemplateServiceDefinition, err error) {
	result = &v1.ServiceTemplateServiceDefinition{}
	err = c.client.Get().
		Resource("servicetemplateservicedefinitions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServiceTemplateServiceDefinitions that match those selectors.
func (c *serviceTemplateServiceDefinitions) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ServiceTemplateServiceDefinitionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ServiceTemplateServiceDefinitionList{}
	err = c.client.Get().
		Resource("servicetemplateservicedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serviceTemplateServiceDefinitions.
func (c *serviceTemplateServiceDefinitions) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("servicetemplateservicedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a serviceTemplateServiceDefinition and creates it.  Returns the server's representation of the serviceTemplateServiceDefinition, and an error, if there is any.
func (c *serviceTemplateServiceDefinitions) Create(ctx context.Context, serviceTemplateServiceDefinition *v1.ServiceTemplateServiceDefinition, opts metav1.CreateOptions) (result *v1.ServiceTemplateServiceDefinition, err error) {
	result = &v1.ServiceTemplateServiceDefinition{}
	err = c.client.Post().
		Resource("servicetemplateservicedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serviceTemplateServiceDefinition).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a serviceTemplateServiceDefinition and updates it. Returns the server's representation of the serviceTemplateServiceDefinition, and an error, if there is any.
func (c *serviceTemplateServiceDefinitions) Update(ctx context.Context, serviceTemplateServiceDefinition *v1.ServiceTemplateServiceDefinition, opts metav1.UpdateOptions) (result *v1.ServiceTemplateServiceDefinition, err error) {
	result = &v1.ServiceTemplateServiceDefinition{}
	err = c.client.Put().
		Resource("servicetemplateservicedefinitions").
		Name(serviceTemplateServiceDefinition.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(serviceTemplateServiceDefinition).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the serviceTemplateServiceDefinition and deletes it. Returns an error if one occurs.
func (c *serviceTemplateServiceDefinitions) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("servicetemplateservicedefinitions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serviceTemplateServiceDefinitions) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("servicetemplateservicedefinitions").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched serviceTemplateServiceDefinition.
func (c *serviceTemplateServiceDefinitions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ServiceTemplateServiceDefinition, err error) {
	result = &v1.ServiceTemplateServiceDefinition{}
	err = c.client.Patch(pt).
		Resource("servicetemplateservicedefinitions").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}