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

// ProgressiveUpgradeConfigsGetter has a method to return a ProgressiveUpgradeConfigInterface.
// A group's client should implement this interface.
type ProgressiveUpgradeConfigsGetter interface {
	ProgressiveUpgradeConfigs() ProgressiveUpgradeConfigInterface
}

// ProgressiveUpgradeConfigInterface has methods to work with ProgressiveUpgradeConfig resources.
type ProgressiveUpgradeConfigInterface interface {
	Create(ctx context.Context, progressiveUpgradeConfig *v1.ProgressiveUpgradeConfig, opts metav1.CreateOptions) (*v1.ProgressiveUpgradeConfig, error)
	Update(ctx context.Context, progressiveUpgradeConfig *v1.ProgressiveUpgradeConfig, opts metav1.UpdateOptions) (*v1.ProgressiveUpgradeConfig, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ProgressiveUpgradeConfig, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ProgressiveUpgradeConfigList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProgressiveUpgradeConfig, err error)
	ProgressiveUpgradeConfigExpansion
}

// progressiveUpgradeConfigs implements ProgressiveUpgradeConfigInterface
type progressiveUpgradeConfigs struct {
	client rest.Interface
}

// newProgressiveUpgradeConfigs returns a ProgressiveUpgradeConfigs
func newProgressiveUpgradeConfigs(c *GlobalTsmV1Client) *progressiveUpgradeConfigs {
	return &progressiveUpgradeConfigs{
		client: c.RESTClient(),
	}
}

// Get takes name of the progressiveUpgradeConfig, and returns the corresponding progressiveUpgradeConfig object, and an error if there is any.
func (c *progressiveUpgradeConfigs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ProgressiveUpgradeConfig, err error) {
	result = &v1.ProgressiveUpgradeConfig{}
	err = c.client.Get().
		Resource("progressiveupgradeconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ProgressiveUpgradeConfigs that match those selectors.
func (c *progressiveUpgradeConfigs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ProgressiveUpgradeConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ProgressiveUpgradeConfigList{}
	err = c.client.Get().
		Resource("progressiveupgradeconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested progressiveUpgradeConfigs.
func (c *progressiveUpgradeConfigs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("progressiveupgradeconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a progressiveUpgradeConfig and creates it.  Returns the server's representation of the progressiveUpgradeConfig, and an error, if there is any.
func (c *progressiveUpgradeConfigs) Create(ctx context.Context, progressiveUpgradeConfig *v1.ProgressiveUpgradeConfig, opts metav1.CreateOptions) (result *v1.ProgressiveUpgradeConfig, err error) {
	result = &v1.ProgressiveUpgradeConfig{}
	err = c.client.Post().
		Resource("progressiveupgradeconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(progressiveUpgradeConfig).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a progressiveUpgradeConfig and updates it. Returns the server's representation of the progressiveUpgradeConfig, and an error, if there is any.
func (c *progressiveUpgradeConfigs) Update(ctx context.Context, progressiveUpgradeConfig *v1.ProgressiveUpgradeConfig, opts metav1.UpdateOptions) (result *v1.ProgressiveUpgradeConfig, err error) {
	result = &v1.ProgressiveUpgradeConfig{}
	err = c.client.Put().
		Resource("progressiveupgradeconfigs").
		Name(progressiveUpgradeConfig.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(progressiveUpgradeConfig).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the progressiveUpgradeConfig and deletes it. Returns an error if one occurs.
func (c *progressiveUpgradeConfigs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("progressiveupgradeconfigs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *progressiveUpgradeConfigs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("progressiveupgradeconfigs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched progressiveUpgradeConfig.
func (c *progressiveUpgradeConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ProgressiveUpgradeConfig, err error) {
	result = &v1.ProgressiveUpgradeConfig{}
	err = c.client.Patch(pt).
		Resource("progressiveupgradeconfigs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}