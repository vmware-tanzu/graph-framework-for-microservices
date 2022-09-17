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
	vConfigConfig *nexus_client.ConfigConfig
	vGnsGns *nexus_client.GnsGns
	vGnsBar *nexus_client.GnsBar
	
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
	
	
	return ret, nil
}

//////////////////////////////////////
// Child/Link Node : Config
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootConfigResolver() (*model.ConfigConfig, error) {
	vConfig, err := c.vRootRoot.GetConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vConfigConfig = vConfig
	
	
	return ret, nil
}

//////////////////////////////////////
// CustomField: CustomBar of CustomType: Root
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootCustomBarResolver() (*model.RootBar, error) {
	vRoot := c.vRootRoot
	
	
	return ret, nil
}
//////////////////////////////////////
// Child/Link Node : GNS
// Resolver for Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigGNSResolver() (*model.GnsGns, error) {
	vGNS, err := c.vConfigConfig.GetGNS(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vGnsGns = vGNS
	
	
	return ret, nil
}

//////////////////////////////////////
// CustomField: Cluster of CustomType: Config
// Resolver for Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigClusterResolver() (*model.ConfigCluster, error) {
	vConfig := c.vConfigConfig
	
	
	return ret, nil
}
//////////////////////////////////////
// Child/Link Node : FooLink
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooLinkResolver() (*model.GnsBar, error) {
	vFooLink, err := c.vGnsGns.GetFooLink(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vGnsBar = vFooLink
	
	
	return ret, nil
}

//////////////////////////////////////
// Child/Link Node : FooChild
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooChildResolver() (*model.GnsBar, error) {
	vFooChild, err := c.vGnsGns.GetFooChild(context.TODO())
	if err != nil {
		panic(err)
	}
	Id := ""
	c.vGnsBar = vFooChild
	
	
	return ret, nil
}

//////////////////////////////////////
// Children/Links Node : FooLinks
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooLinksResolver(id *string) ([]*model.GnsBar, error) {
	var vGnsBarList []*model.GnsBar
	if id != nil && *id != "" {
		Id := *id
		vGnsBar, err := c.vGnsGns.GetFooLinks(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		
		
		vGnsBarList = append(vGnsList, ret)
		return vGnsBarList, nil
	}
	for i := range c.vGnsGns.Spec.FooLinksGvk {
		Id := i
		vGnsBar, err := c.vGnsGns.GetFooLinks(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		
		
		vGnsBarList = append(vGnsList, ret)
	}
	return vGnsBarList, nil
}

//////////////////////////////////////
// Children/Links Node : FooChildren
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooChildrenResolver(id *string) ([]*model.GnsBar, error) {
	var vGnsBarList []*model.GnsBar
	if id != nil && *id != "" {
		Id := *id
		vGnsBar, err := c.vGnsGns.GetFooChildren(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		
		
		vGnsBarList = append(vGnsList, ret)
		return vGnsBarList, nil
	}
	for i := range c.vGnsGns.Spec.FooChildrenGvk {
		Id := i
		vGnsBar, err := c.vGnsGns.GetFooChildren(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		
		
		vGnsBarList = append(vGnsList, ret)
	}
	return vGnsBarList, nil
}

//////////////////////////////////////
// CustomField: Mydesc of CustomType: Gns
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsMydescResolver() (*model.GnsDescription, error) {
	vGns := c.vGnsGns
	
	
	return ret, nil
}
//////////////////////////////////////
// CustomField: HostPort of CustomType: Gns
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsHostPortResolver() (*model.GnsHostPort, error) {
	vGns := c.vGnsGns
	
	
	return ret, nil
}