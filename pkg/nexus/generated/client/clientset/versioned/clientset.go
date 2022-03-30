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

package versioned

import (
	"fmt"

	confighelloworldv1 "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/client/clientset/versioned/typed/config.helloworld.com/v1"
	inventoryhelloworldv1 "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/client/clientset/versioned/typed/inventory.helloworld.com/v1"
	nexushelloworldv1 "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/client/clientset/versioned/typed/nexus.helloworld.com/v1"
	roothelloworldv1 "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/client/clientset/versioned/typed/root.helloworld.com/v1"
	runtimehelloworldv1 "gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/client/clientset/versioned/typed/runtime.helloworld.com/v1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	ConfigHelloworldV1() confighelloworldv1.ConfigHelloworldV1Interface
	InventoryHelloworldV1() inventoryhelloworldv1.InventoryHelloworldV1Interface
	NexusHelloworldV1() nexushelloworldv1.NexusHelloworldV1Interface
	RootHelloworldV1() roothelloworldv1.RootHelloworldV1Interface
	RuntimeHelloworldV1() runtimehelloworldv1.RuntimeHelloworldV1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	configHelloworldV1    *confighelloworldv1.ConfigHelloworldV1Client
	inventoryHelloworldV1 *inventoryhelloworldv1.InventoryHelloworldV1Client
	nexusHelloworldV1     *nexushelloworldv1.NexusHelloworldV1Client
	rootHelloworldV1      *roothelloworldv1.RootHelloworldV1Client
	runtimeHelloworldV1   *runtimehelloworldv1.RuntimeHelloworldV1Client
}

// ConfigHelloworldV1 retrieves the ConfigHelloworldV1Client
func (c *Clientset) ConfigHelloworldV1() confighelloworldv1.ConfigHelloworldV1Interface {
	return c.configHelloworldV1
}

// InventoryHelloworldV1 retrieves the InventoryHelloworldV1Client
func (c *Clientset) InventoryHelloworldV1() inventoryhelloworldv1.InventoryHelloworldV1Interface {
	return c.inventoryHelloworldV1
}

// NexusHelloworldV1 retrieves the NexusHelloworldV1Client
func (c *Clientset) NexusHelloworldV1() nexushelloworldv1.NexusHelloworldV1Interface {
	return c.nexusHelloworldV1
}

// RootHelloworldV1 retrieves the RootHelloworldV1Client
func (c *Clientset) RootHelloworldV1() roothelloworldv1.RootHelloworldV1Interface {
	return c.rootHelloworldV1
}

// RuntimeHelloworldV1 retrieves the RuntimeHelloworldV1Client
func (c *Clientset) RuntimeHelloworldV1() runtimehelloworldv1.RuntimeHelloworldV1Interface {
	return c.runtimeHelloworldV1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.configHelloworldV1, err = confighelloworldv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.inventoryHelloworldV1, err = inventoryhelloworldv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.nexusHelloworldV1, err = nexushelloworldv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.rootHelloworldV1, err = roothelloworldv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.runtimeHelloworldV1, err = runtimehelloworldv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.configHelloworldV1 = confighelloworldv1.NewForConfigOrDie(c)
	cs.inventoryHelloworldV1 = inventoryhelloworldv1.NewForConfigOrDie(c)
	cs.nexusHelloworldV1 = nexushelloworldv1.NewForConfigOrDie(c)
	cs.rootHelloworldV1 = roothelloworldv1.NewForConfigOrDie(c)
	cs.runtimeHelloworldV1 = runtimehelloworldv1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.configHelloworldV1 = confighelloworldv1.New(c)
	cs.inventoryHelloworldV1 = inventoryhelloworldv1.New(c)
	cs.nexusHelloworldV1 = nexushelloworldv1.New(c)
	cs.rootHelloworldV1 = roothelloworldv1.New(c)
	cs.runtimeHelloworldV1 = runtimehelloworldv1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
