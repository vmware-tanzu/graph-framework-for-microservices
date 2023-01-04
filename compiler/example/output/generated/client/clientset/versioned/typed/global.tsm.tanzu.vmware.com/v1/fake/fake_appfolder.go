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

// FakeAppFolders implements AppFolderInterface
type FakeAppFolders struct {
	Fake *FakeGlobalTsmV1
}

var appfoldersResource = schema.GroupVersionResource{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Resource: "appfolders"}

var appfoldersKind = schema.GroupVersionKind{Group: "global.tsm.tanzu.vmware.com", Version: "v1", Kind: "AppFolder"}

// Get takes name of the appFolder, and returns the corresponding appFolder object, and an error if there is any.
func (c *FakeAppFolders) Get(ctx context.Context, name string, options v1.GetOptions) (result *globaltsmtanzuvmwarecomv1.AppFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(appfoldersResource, name), &globaltsmtanzuvmwarecomv1.AppFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppFolder), err
}

// List takes label and field selectors, and returns the list of AppFolders that match those selectors.
func (c *FakeAppFolders) List(ctx context.Context, opts v1.ListOptions) (result *globaltsmtanzuvmwarecomv1.AppFolderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(appfoldersResource, appfoldersKind, opts), &globaltsmtanzuvmwarecomv1.AppFolderList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &globaltsmtanzuvmwarecomv1.AppFolderList{ListMeta: obj.(*globaltsmtanzuvmwarecomv1.AppFolderList).ListMeta}
	for _, item := range obj.(*globaltsmtanzuvmwarecomv1.AppFolderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested appFolders.
func (c *FakeAppFolders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(appfoldersResource, opts))
}

// Create takes the representation of a appFolder and creates it.  Returns the server's representation of the appFolder, and an error, if there is any.
func (c *FakeAppFolders) Create(ctx context.Context, appFolder *globaltsmtanzuvmwarecomv1.AppFolder, opts v1.CreateOptions) (result *globaltsmtanzuvmwarecomv1.AppFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(appfoldersResource, appFolder), &globaltsmtanzuvmwarecomv1.AppFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppFolder), err
}

// Update takes the representation of a appFolder and updates it. Returns the server's representation of the appFolder, and an error, if there is any.
func (c *FakeAppFolders) Update(ctx context.Context, appFolder *globaltsmtanzuvmwarecomv1.AppFolder, opts v1.UpdateOptions) (result *globaltsmtanzuvmwarecomv1.AppFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(appfoldersResource, appFolder), &globaltsmtanzuvmwarecomv1.AppFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppFolder), err
}

// Delete takes name of the appFolder and deletes it. Returns an error if one occurs.
func (c *FakeAppFolders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(appfoldersResource, name), &globaltsmtanzuvmwarecomv1.AppFolder{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAppFolders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(appfoldersResource, listOpts)

	_, err := c.Fake.Invokes(action, &globaltsmtanzuvmwarecomv1.AppFolderList{})
	return err
}

// Patch applies the patch and returns the patched appFolder.
func (c *FakeAppFolders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *globaltsmtanzuvmwarecomv1.AppFolder, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(appfoldersResource, name, pt, data, subresources...), &globaltsmtanzuvmwarecomv1.AppFolder{})
	if obj == nil {
		return nil, err
	}
	return obj.(*globaltsmtanzuvmwarecomv1.AppFolder), err
}