package graph

import (
	"context"
	"encoding/json"
	"fmt"
	nexus_client "nexustempmodule/nexus-client"
	"nexustempmodule/nexus-gql/graph/model"

	qm "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/query-manager"
	libgrpc "gitlab.eng.vmware.com/nsx-allspark_users/lib-go/grpc"
	"k8s.io/client-go/rest"
)

var c resolverConfig

type resolverConfig struct {
	vRootRoot *nexus_client.RootRoot
	vRootNonNexusType *nexus_client.RootNonNexusType
	vConfigConfig *nexus_client.ConfigConfig
	vGnsRandomGnsData *nexus_client.GnsRandomGnsData
	vGnsGns *nexus_client.GnsGns
	vGnsDns *nexus_client.GnsDns
	vGnsAdditionalGnsData *nexus_client.GnsAdditionalGnsData
	vServicegroupSvcGroup *nexus_client.ServicegroupSvcGroup
	vPolicypkgAdditionalPolicyData *nexus_client.PolicypkgAdditionalPolicyData
	vPolicypkgAccessControlPolicy *nexus_client.PolicypkgAccessControlPolicy
	vPolicypkgACPConfig *nexus_client.PolicypkgACPConfig
	vPolicypkgVMpolicy *nexus_client.PolicypkgVMpolicy
	vPolicypkgRandomPolicyData *nexus_client.PolicypkgRandomPolicyData
	
}

//////////////////////////////////////
// Nexus K8sAPIEndpointConfig
//////////////////////////////////////
func getK8sAPIEndpointConfig() *rest.Config {

	var config *rest.Config
	config = &rest.Config{
		Host: "http://localhost:45192",
	}

	return config
}

//////////////////////////////////////
// GRPC SERVER CONFIG
//////////////////////////////////////
func grpcServer() qm.ServerClient{
	addr := "localhost:45781"
	conn, err := libgrpc.ClientConn(addr, libgrpc.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to query-manager server, err: %v", err)
	}
	return qm.NewServerClient(conn)
}

//////////////////////////////////////
// Resolver for Parent Node: Root
//////////////////////////////////////
func (c *resolverConfig) getRootResolver() (*model.RootRoot, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nc, err := nexus_client.NewForConfig(k8sApiConfig)
	Id := ""
	if err != nil {
		panic(err)
	}
	vRoot, err := nc.GetRootRoot(context.TODO())
	if err != nil {
		panic(err)
	}
	c.vRootRoot = vRoot
	
	ret := &model.RootRoot {
		Id: &Id,
		Name: &vRoot.Spec.Name,
		NonStructFoo: &vRoot.Spec.NonStructFoo,
	}
	return ret, nil
}
//////////////////////////////////////
// Resolver for Parent Node: Gns
//////////////////////////////////////
func (c *resolverConfig) getRootResolver() (*model.GnsDns, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nc, err := nexus_client.NewForConfig(k8sApiConfig)
	Id := ""
	if err != nil {
		panic(err)
	}
	vDns, err := nc.GetGnsDns(context.TODO())
	if err != nil {
		panic(err)
	}
	c.vGnsDns = vDns
	
	ret := &model.GnsDns {
		Id: &Id,
	}
	return ret, nil
}

//////////////////////////////////////
// Child/Link Node : Config
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootConfigResolver() (*model.ConfigConfig, error) {
	vConfigConfig, err := c.vRootRoot.GetConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vConfigConfig = vConfigConfig
	
	ret := &model.ConfigConfig {
		Id: &Id,
	}
	return ret, nil
}

//////////////////////////////////////
// Child/Link Node : Foolink
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootFoolinkResolver() (*model.ConfigConfig, error) {
	vConfigConfig, err := c.vRootRoot.GetFoolink(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vConfigConfig = vConfigConfig
	
	ret := &model.ConfigConfig {
		Id: &Id,
	}
	return ret, nil
}

//////////////////////////////////////
// Children/Links Node : Foochildren
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootFoochildrenResolver(id *string) ([]*model.ConfigConfig, error) {
	var vConfigConfigList []*model.ConfigConfig
	if id != nil && *id != "" {
		Id := *id
		vConfigConfig, err := c.vRootRoot.GetFoochildren(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ConfigConfig {
		Id: &Id,
	}
		vConfigConfigList = append(vRootList, ret)
		return vConfigConfigList, nil
	}
	for i := range c.vRootRoot.Spec.FoochildrenGvk {
		Id := i
		vConfigConfig, err := c.vRootRoot.GetFoochildren(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ConfigConfig {
		Id: &Id,
	}
		vConfigConfigList = append(vRootList, ret)
	}
	return vConfigConfigList, nil
}

//////////////////////////////////////
// Children/Links Node : Foolinks
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootFoolinksResolver(id *string) ([]*model.ConfigConfig, error) {
	var vConfigConfigList []*model.ConfigConfig
	if id != nil && *id != "" {
		Id := *id
		vConfigConfig, err := c.vRootRoot.GetFoolinks(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ConfigConfig {
		Id: &Id,
	}
		vConfigConfigList = append(vRootList, ret)
		return vConfigConfigList, nil
	}
	for i := range c.vRootRoot.Spec.FoolinksGvk {
		Id := i
		vConfigConfig, err := c.vRootRoot.GetFoolinks(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ConfigConfig {
		Id: &Id,
	}
		vConfigConfigList = append(vRootList, ret)
	}
	return vConfigConfigList, nil
}

//////////////////////////////////////
// CustomField: CustomBar of CustomType: Root
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootCustomBarResolver() (*model.RootBar, error) {
	vRoot := c.vRootRoot
	
	ret := &model.RootBar {
		Foo: &vRoot.Spec.Foo,
	}
	return ret, nil
}
//////////////////////////////////////
// CustomField: Bar of CustomType: NonNexusType
// Resolver for NonNexusType
//////////////////////////////////////
func (c *resolverConfig) getRootNonNexusTypeBarResolver() (*model.RootBar, error) {
	vRoot := c.vRootNonNexusType
	
	ret := &model.RootBar {
		Foo: &vRoot.Spec.Foo,
	}
	return ret, nil
}
//////////////////////////////////////
// Child/Link Node : GNS
// Resolver for Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigGNSResolver() (*model.GnsGns, error) {
	vGnsGns, err := c.vConfigConfig.GetGNS(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vGnsGns = vGnsGns
	
	ret := &model.GnsGns {
		Id: &Id,
		Domain: &vGns.Spec.Domain,
		UseSharedGateway: &vGns.Spec.UseSharedGateway,
		Meta: &vGns.Spec.Meta,
	}
	return ret, nil
}

//////////////////////////////////////
// Child/Link Node : DNS
// Resolver for Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigDNSResolver() (*model.GnsDns, error) {
	vGnsDns, err := c.vConfigConfig.GetDNS(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vGnsDns = vGnsDns
	
	ret := &model.GnsDns {
		Id: &Id,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomField: TestValMarkers of CustomType: Config
// Resolver for Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigTestValMarkersResolver() (*model.ConfigTestValMarkers, error) {
	vConfig := c.vConfigConfig
	
	ret := &model.ConfigTestValMarkers {
		MyStr: &vConfig.Spec.MyStr,
		MyInt: &vConfig.Spec.MyInt,
	}
	return ret, nil
}
//////////////////////////////////////
// CustomField: Description of CustomType: RandomGnsData
// Resolver for RandomGnsData
//////////////////////////////////////
func (c *resolverConfig) getGnsRandomGnsDataDescriptionResolver() (*model.GnsRandomDescription, error) {
	vGns := c.vGnsRandomGnsData
	
	ret := &model.GnsRandomDescription {
		DiscriptionA: &vGns.Spec.DiscriptionA,
		DiscriptionB: &vGns.Spec.DiscriptionB,
		DiscriptionC: &vGns.Spec.DiscriptionC,
		DiscriptionD: &vGns.Spec.DiscriptionD,
	}
	return ret, nil
}
//////////////////////////////////////
// Child/Link Node : Dns
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsDnsResolver() (*model.GnsDns, error) {
	vGnsDns, err := c.vGnsGns.GetDns(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vGnsDns = vGnsDns
	
	ret := &model.GnsDns {
		Id: &Id,
	}
	return ret, nil
}

//////////////////////////////////////
// Children/Links Node : GnsServiceGroups
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsGnsServiceGroupsResolver(id *string) ([]*model.ServicegroupSvcGroup, error) {
	var vServicegroupSvcGroupList []*model.ServicegroupSvcGroup
	if id != nil && *id != "" {
		Id := *id
		vServicegroupSvcGroup, err := c.vGnsGns.GetGnsServiceGroups(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ServicegroupSvcGroup {
		Id: &Id,
		DisplayName: &vServicegroup.Spec.DisplayName,
		Description: &vServicegroup.Spec.Description,
		Color: &vServicegroup.Spec.Color,
	}
		vServicegroupSvcGroupList = append(vGnsList, ret)
		return vServicegroupSvcGroupList, nil
	}
	for i := range c.vGnsGns.Spec.GnsServiceGroupsGvk {
		Id := i
		vServicegroupSvcGroup, err := c.vGnsGns.GetGnsServiceGroups(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ServicegroupSvcGroup {
		Id: &Id,
		DisplayName: &vServicegroup.Spec.DisplayName,
		Description: &vServicegroup.Spec.Description,
		Color: &vServicegroup.Spec.Color,
	}
		vServicegroupSvcGroupList = append(vGnsList, ret)
	}
	return vServicegroupSvcGroupList, nil
}

//////////////////////////////////////
// CustomField: Description of CustomType: Gns
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsDescriptionResolver() (*model.GnsDescription, error) {
	vGns := c.vGnsGns
	
	ret := &model.GnsDescription {
		Color: &vGns.Spec.Color,
		Version: &vGns.Spec.Version,
		ProjectId: &vGns.Spec.ProjectId,
	}
	return ret, nil
}
//////////////////////////////////////
// CustomField: Description of CustomType: AdditionalGnsData
// Resolver for AdditionalGnsData
//////////////////////////////////////
func (c *resolverConfig) getGnsAdditionalGnsDataDescriptionResolver() (*model.GnsAdditionalDescription, error) {
	vGns := c.vGnsAdditionalGnsData
	
	ret := &model.GnsAdditionalDescription {
		DiscriptionA: &vGns.Spec.DiscriptionA,
		DiscriptionB: &vGns.Spec.DiscriptionB,
		DiscriptionC: &vGns.Spec.DiscriptionC,
		DiscriptionD: &vGns.Spec.DiscriptionD,
	}
	return ret, nil
}
//////////////////////////////////////
// CustomField: Description of CustomType: AdditionalPolicyData
// Resolver for AdditionalPolicyData
//////////////////////////////////////
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataDescriptionResolver() (*model.PolicypkgAdditionalDescription, error) {
	vPolicypkg := c.vPolicypkgAdditionalPolicyData
	
	ret := &model.PolicypkgAdditionalDescription {
		DiscriptionA: &vPolicypkg.Spec.DiscriptionA,
		DiscriptionB: &vPolicypkg.Spec.DiscriptionB,
		DiscriptionC: &vPolicypkg.Spec.DiscriptionC,
		DiscriptionD: &vPolicypkg.Spec.DiscriptionD,
	}
	return ret, nil
}
//////////////////////////////////////
// Children/Links Node : SourceSvcGroups
// Resolver for ACPConfig
//////////////////////////////////////
func (c *resolverConfig) getPolicypkgACPConfigSourceSvcGroupsResolver(id *string) ([]*model.ServicegroupSvcGroup, error) {
	var vServicegroupSvcGroupList []*model.ServicegroupSvcGroup
	if id != nil && *id != "" {
		Id := *id
		vServicegroupSvcGroup, err := c.vPolicypkgACPConfig.GetSourceSvcGroups(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ServicegroupSvcGroup {
		Id: &Id,
		DisplayName: &vServicegroup.Spec.DisplayName,
		Description: &vServicegroup.Spec.Description,
		Color: &vServicegroup.Spec.Color,
	}
		vServicegroupSvcGroupList = append(vACPConfigList, ret)
		return vServicegroupSvcGroupList, nil
	}
	for i := range c.vPolicypkgACPConfig.Spec.SourceSvcGroupsGvk {
		Id := i
		vServicegroupSvcGroup, err := c.vPolicypkgACPConfig.GetSourceSvcGroups(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		
		ret := &model.ServicegroupSvcGroup {
		Id: &Id,
		DisplayName: &vServicegroup.Spec.DisplayName,
		Description: &vServicegroup.Spec.Description,
		Color: &vServicegroup.Spec.Color,
	}
		vServicegroupSvcGroupList = append(vACPConfigList, ret)
	}
	return vServicegroupSvcGroupList, nil
}

//////////////////////////////////////
// CustomField: Description of CustomType: RandomPolicyData
// Resolver for RandomPolicyData
//////////////////////////////////////
func (c *resolverConfig) getPolicypkgRandomPolicyDataDescriptionResolver() (*model.PolicypkgRandomDescription, error) {
	vPolicypkg := c.vPolicypkgRandomPolicyData
	
	ret := &model.PolicypkgRandomDescription {
		DiscriptionA: &vPolicypkg.Spec.DiscriptionA,
		DiscriptionB: &vPolicypkg.Spec.DiscriptionB,
		DiscriptionC: &vPolicypkg.Spec.DiscriptionC,
		DiscriptionD: &vPolicypkg.Spec.DiscriptionD,
	}
	return ret, nil
}