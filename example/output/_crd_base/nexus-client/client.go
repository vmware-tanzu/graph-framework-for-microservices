package nexus_client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	baseClientset "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/_crd_generated/client/clientset/versioned"

	baseconfigtsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/_crd_generated/apis/config.tsm.tanzu.vmware.com/v1"
	basegnstsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/_crd_generated/apis/gns.tsm.tanzu.vmware.com/v1"
	basepolicytsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/_crd_generated/apis/policy.tsm.tanzu.vmware.com/v1"
	baseroottsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/_crd_generated/apis/root.tsm.tanzu.vmware.com/v1"
	baseservicegrouptsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/_crd_generated/apis/servicegroup.tsm.tanzu.vmware.com/v1"
)

type Clientset struct {
	baseClient        *baseClientset.Clientset
	rootTsmV1         *RootTsmV1
	configTsmV1       *ConfigTsmV1
	gnsTsmV1          *GnsTsmV1
	servicegroupTsmV1 *ServicegroupTsmV1
	policyTsmV1       *PolicyTsmV1
}

func NewForConfig(config *rest.Config) (*Clientset, error) {
	baseClient, err := baseClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client := &Clientset{}
	client.baseClient = baseClient
	client.rootTsmV1 = newRootTsmV1(client)
	client.configTsmV1 = newConfigTsmV1(client)
	client.gnsTsmV1 = newGnsTsmV1(client)
	client.servicegroupTsmV1 = newServicegroupTsmV1(client)
	client.policyTsmV1 = newPolicyTsmV1(client)

	return client, nil
}

func (c *Clientset) RootTsmV1() *RootTsmV1 {
	return c.rootTsmV1
}
func (c *Clientset) ConfigTsmV1() *ConfigTsmV1 {
	return c.configTsmV1
}
func (c *Clientset) GnsTsmV1() *GnsTsmV1 {
	return c.gnsTsmV1
}
func (c *Clientset) ServicegroupTsmV1() *ServicegroupTsmV1 {
	return c.servicegroupTsmV1
}
func (c *Clientset) PolicyTsmV1() *PolicyTsmV1 {
	return c.policyTsmV1
}

type RootTsmV1 struct {
	roots *rootRootTsmV1
}

func newRootTsmV1(client *Clientset) *RootTsmV1 {
	return &RootTsmV1{
		roots: &rootRootTsmV1{
			client: client,
		},
	}
}

type rootRootTsmV1 struct {
	client *Clientset
}

func (obj *RootTsmV1) Roots() *rootRootTsmV1 {
	return obj.roots
}

type ConfigTsmV1 struct {
	configs *configConfigTsmV1
}

func newConfigTsmV1(client *Clientset) *ConfigTsmV1 {
	return &ConfigTsmV1{
		configs: &configConfigTsmV1{
			client: client,
		},
	}
}

type configConfigTsmV1 struct {
	client *Clientset
}

func (obj *ConfigTsmV1) Configs() *configConfigTsmV1 {
	return obj.configs
}

type GnsTsmV1 struct {
	gnses *gnsGnsTsmV1
	dnses *dnsGnsTsmV1
}

func newGnsTsmV1(client *Clientset) *GnsTsmV1 {
	return &GnsTsmV1{
		gnses: &gnsGnsTsmV1{
			client: client,
		},
		dnses: &dnsGnsTsmV1{
			client: client,
		},
	}
}

type gnsGnsTsmV1 struct {
	client *Clientset
}

func (obj *GnsTsmV1) Gnses() *gnsGnsTsmV1 {
	return obj.gnses
}

type dnsGnsTsmV1 struct {
	client *Clientset
}

func (obj *GnsTsmV1) Dnses() *dnsGnsTsmV1 {
	return obj.dnses
}

type ServicegroupTsmV1 struct {
	svcgroups *svcgroupServicegroupTsmV1
}

func newServicegroupTsmV1(client *Clientset) *ServicegroupTsmV1 {
	return &ServicegroupTsmV1{
		svcgroups: &svcgroupServicegroupTsmV1{
			client: client,
		},
	}
}

type svcgroupServicegroupTsmV1 struct {
	client *Clientset
}

func (obj *ServicegroupTsmV1) SvcGroups() *svcgroupServicegroupTsmV1 {
	return obj.svcgroups
}

type PolicyTsmV1 struct {
	accesscontrolpolicies *accesscontrolpolicyPolicyTsmV1
	acpconfigs            *acpconfigPolicyTsmV1
}

func newPolicyTsmV1(client *Clientset) *PolicyTsmV1 {
	return &PolicyTsmV1{
		accesscontrolpolicies: &accesscontrolpolicyPolicyTsmV1{
			client: client,
		},
		acpconfigs: &acpconfigPolicyTsmV1{
			client: client,
		},
	}
}

type accesscontrolpolicyPolicyTsmV1 struct {
	client *Clientset
}

func (obj *PolicyTsmV1) AccessControlPolicies() *accesscontrolpolicyPolicyTsmV1 {
	return obj.accesscontrolpolicies
}

type acpconfigPolicyTsmV1 struct {
	client *Clientset
}

func (obj *PolicyTsmV1) ACPConfigs() *acpconfigPolicyTsmV1 {
	return obj.acpconfigs
}

func (obj *rootRootTsmV1) Get(ctx context.Context, name string, options metav1.GetOptions) (result *baseroottsmtanzuvmwarecomv1.Root, err error) {
	result, err = obj.client.baseClient.RootTsmV1().Roots().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	if result.Spec.ConfigGvk.Name != "" {
		field, err := obj.client.ConfigTsmV1().Configs().Get(ctx, result.Spec.ConfigGvk.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.Config = *field
	}

	return
}

func (obj *configConfigTsmV1) Get(ctx context.Context, name string, options metav1.GetOptions) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	result, err = obj.client.baseClient.ConfigTsmV1().Configs().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	if result.Spec.GNSGvk.Name != "" {
		field, err := obj.client.GnsTsmV1().Gnses().Get(ctx, result.Spec.GNSGvk.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.GNS = *field
	}

	for k, v := range result.Spec.PolicyGvk {
		obj, err := obj.client.PolicyTsmV1().AccessControlPolicies().Get(ctx, v.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.Policy[k] = *obj
	}

	return
}

func (obj *gnsGnsTsmV1) Get(ctx context.Context, name string, options metav1.GetOptions) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	result, err = obj.client.baseClient.GnsTsmV1().Gnses().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	for k, v := range result.Spec.GnsServiceGroupsGvk {
		obj, err := obj.client.Service_groupTsmV1().SvcGroups().Get(ctx, v.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.GnsServiceGroups[k] = *obj
	}

	if result.Spec.GnsAccessControlPolicyGvk.Name != "" {
		field, err := obj.client.PolicyTsmV1().AccessControlPolicies().Get(ctx, result.Spec.GnsAccessControlPolicyGvk.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.GnsAccessControlPolicy = *field
	}

	if result.Spec.DnsGvk.Name != "" {
		field, err := obj.client.GnsTsmV1().Dnses().Get(ctx, result.Spec.DnsGvk.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.Dns = *field
	}

	return
}

func (obj *dnsGnsTsmV1) Get(ctx context.Context, name string, options metav1.GetOptions) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	result, err = obj.client.baseClient.GnsTsmV1().Dnses().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	return
}

func (obj *svcgroupServicegroupTsmV1) Get(ctx context.Context, name string, options metav1.GetOptions) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	result, err = obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	for k, v := range result.Spec.ServicesGvk {
		obj, err := obj.client.Core_v1TsmV1().Services().Get(ctx, v.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.Services[k] = *obj
	}

	return
}

func (obj *accesscontrolpolicyPolicyTsmV1) Get(ctx context.Context, name string, options metav1.GetOptions) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	result, err = obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	for k, v := range result.Spec.PolicyConfigsGvk {
		obj, err := obj.client.PolicyTsmV1().ACPConfigs().Get(ctx, v.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.PolicyConfigs[k] = *obj
	}

	return
}

func (obj *acpconfigPolicyTsmV1) Get(ctx context.Context, name string, options metav1.GetOptions) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	for k, v := range result.Spec.DestSvcGroupsGvk {
		obj, err := obj.client.Service_groupTsmV1().SvcGroups().Get(ctx, v.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.DestSvcGroups[k] = *obj
	}

	for k, v := range result.Spec.SourceSvcGroupsGvk {
		obj, err := obj.client.Service_groupTsmV1().SvcGroups().Get(ctx, v.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.SourceSvcGroups[k] = *obj
	}

	return
}
