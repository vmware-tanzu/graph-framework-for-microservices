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

	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/admin.nexus.org/v1"
	scheme "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/clientset/versioned/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ProxyRulesGetter has a method to return a ProxyRuleInterface.
// A group's client should implement this interface.
type ProxyRulesGetter interface {
	ProxyRules() ProxyRuleInterface
}

// ProxyRuleInterface has methods to work with ProxyRule resources.
type ProxyRuleInterface interface {
	Create(ctx context.Context, proxyRule *v1.ProxyRule, opts metav1.CreateOptions) (*v1.ProxyRule, error)
	Update(ctx context.Context, proxyRule *v1.ProxyRule, opts metav1.UpdateOptions) (*v1.ProxyRule, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ProxyRule, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ProxyRuleList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProxyRule, err error)
	ProxyRuleExpansion
}

// proxyRules implements ProxyRuleInterface
type proxyRules struct {
	client rest.Interface
}

// newProxyRules returns a ProxyRules
func newProxyRules(c *AdminNexusV1Client) *proxyRules {
	return &proxyRules{
		client: c.RESTClient(),
	}
}

// Get takes name of the proxyRule, and returns the corresponding proxyRule object, and an error if there is any.
func (c *proxyRules) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ProxyRule, err error) {
	result = &v1.ProxyRule{}
	err = c.client.Get().
		Resource("proxyrules").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ProxyRules that match those selectors.
func (c *proxyRules) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ProxyRuleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ProxyRuleList{}
	err = c.client.Get().
		Resource("proxyrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested proxyRules.
func (c *proxyRules) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("proxyrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a proxyRule and creates it.  Returns the server's representation of the proxyRule, and an error, if there is any.
func (c *proxyRules) Create(ctx context.Context, proxyRule *v1.ProxyRule, opts metav1.CreateOptions) (result *v1.ProxyRule, err error) {
	result = &v1.ProxyRule{}
	err = c.client.Post().
		Resource("proxyrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(proxyRule).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a proxyRule and updates it. Returns the server's representation of the proxyRule, and an error, if there is any.
func (c *proxyRules) Update(ctx context.Context, proxyRule *v1.ProxyRule, opts metav1.UpdateOptions) (result *v1.ProxyRule, err error) {
	result = &v1.ProxyRule{}
	err = c.client.Put().
		Resource("proxyrules").
		Name(proxyRule.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(proxyRule).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the proxyRule and deletes it. Returns an error if one occurs.
func (c *proxyRules) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("proxyrules").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *proxyRules) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("proxyrules").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched proxyRule.
func (c *proxyRules) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProxyRule, err error) {
	result = &v1.ProxyRule{}
	err = c.client.Patch(pt).
		Resource("proxyrules").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
