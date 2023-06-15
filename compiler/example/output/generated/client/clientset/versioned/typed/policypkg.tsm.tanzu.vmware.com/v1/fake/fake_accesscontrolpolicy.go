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
	policypkgtsmtanzuvmwarecomv1 "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/apis/policypkg.tsm.tanzu.vmware.com/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAccessControlPolicies implements AccessControlPolicyInterface
type FakeAccessControlPolicies struct {
	Fake *FakePolicypkgTsmV1
}

var accesscontrolpoliciesResource = schema.GroupVersionResource{Group: "policypkg.tsm.tanzu.vmware.com", Version: "v1", Resource: "accesscontrolpolicies"}

var accesscontrolpoliciesKind = schema.GroupVersionKind{Group: "policypkg.tsm.tanzu.vmware.com", Version: "v1", Kind: "AccessControlPolicy"}

// Get takes name of the accessControlPolicy, and returns the corresponding accessControlPolicy object, and an error if there is any.
func (c *FakeAccessControlPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *policypkgtsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(accesscontrolpoliciesResource, name), &policypkgtsmtanzuvmwarecomv1.AccessControlPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.AccessControlPolicy), err
}

// List takes label and field selectors, and returns the list of AccessControlPolicies that match those selectors.
func (c *FakeAccessControlPolicies) List(ctx context.Context, opts v1.ListOptions) (result *policypkgtsmtanzuvmwarecomv1.AccessControlPolicyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(accesscontrolpoliciesResource, accesscontrolpoliciesKind, opts), &policypkgtsmtanzuvmwarecomv1.AccessControlPolicyList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &policypkgtsmtanzuvmwarecomv1.AccessControlPolicyList{ListMeta: obj.(*policypkgtsmtanzuvmwarecomv1.AccessControlPolicyList).ListMeta}
	for _, item := range obj.(*policypkgtsmtanzuvmwarecomv1.AccessControlPolicyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested accessControlPolicies.
func (c *FakeAccessControlPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(accesscontrolpoliciesResource, opts))
}

// Create takes the representation of a accessControlPolicy and creates it.  Returns the server's representation of the accessControlPolicy, and an error, if there is any.
func (c *FakeAccessControlPolicies) Create(ctx context.Context, accessControlPolicy *policypkgtsmtanzuvmwarecomv1.AccessControlPolicy, opts v1.CreateOptions) (result *policypkgtsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(accesscontrolpoliciesResource, accessControlPolicy), &policypkgtsmtanzuvmwarecomv1.AccessControlPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.AccessControlPolicy), err
}

// Update takes the representation of a accessControlPolicy and updates it. Returns the server's representation of the accessControlPolicy, and an error, if there is any.
func (c *FakeAccessControlPolicies) Update(ctx context.Context, accessControlPolicy *policypkgtsmtanzuvmwarecomv1.AccessControlPolicy, opts v1.UpdateOptions) (result *policypkgtsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(accesscontrolpoliciesResource, accessControlPolicy), &policypkgtsmtanzuvmwarecomv1.AccessControlPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.AccessControlPolicy), err
}

// Delete takes name of the accessControlPolicy and deletes it. Returns an error if one occurs.
func (c *FakeAccessControlPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(accesscontrolpoliciesResource, name, opts), &policypkgtsmtanzuvmwarecomv1.AccessControlPolicy{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAccessControlPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(accesscontrolpoliciesResource, listOpts)

	_, err := c.Fake.Invokes(action, &policypkgtsmtanzuvmwarecomv1.AccessControlPolicyList{})
	return err
}

// Patch applies the patch and returns the patched accessControlPolicy.
func (c *FakeAccessControlPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *policypkgtsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(accesscontrolpoliciesResource, name, pt, data, subresources...), &policypkgtsmtanzuvmwarecomv1.AccessControlPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*policypkgtsmtanzuvmwarecomv1.AccessControlPolicy), err
}
