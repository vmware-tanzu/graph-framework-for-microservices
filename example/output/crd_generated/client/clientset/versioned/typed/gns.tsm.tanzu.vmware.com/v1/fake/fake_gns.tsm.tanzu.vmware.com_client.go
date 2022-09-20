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
	v1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/client/clientset/versioned/typed/gns.tsm.tanzu.vmware.com/v1"

	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeGnsTsmV1 struct {
	*testing.Fake
}

func (c *FakeGnsTsmV1) Bars() v1.BarInterface {
	return &FakeBars{c}
}

func (c *FakeGnsTsmV1) EmptyDatas() v1.EmptyDataInterface {
	return &FakeEmptyDatas{c}
}

func (c *FakeGnsTsmV1) Gnses() v1.GnsInterface {
	return &FakeGnses{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeGnsTsmV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
