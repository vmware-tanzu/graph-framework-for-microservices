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
	vGnsEmptyData *nexus_client.GnsEmptyData
	
}

//////////////////////////////////////
// Nexus K8sAPIEndpointConfig
//////////////////////////////////////
func getK8sAPIEndpointConfig() *rest.Config {

	filePath := os.Getenv("KUBECONFIG")
	if filePath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", filePath)
		if err != nil {
			return nil
		}
		return config
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil
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
	if err != nil {
		panic(err)
	}
	vRoot, err := nc.GetRootRoot(context.TODO())
	if err != nil {
		panic(err)
	}
	c.vRootRoot = vRoot
	
	vDisplayName := string(vRoot.Spec.DisplayName)

	ret := &model.RootRoot {
	DisplayName: &vDisplayName,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: RootRoot
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getRootRootqueryServiceTableResolver(startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getRootRootqueryServiceVersionTableResolver(startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getRootRootqueryServiceTSResolver(svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getRootRootqueryIncomingAPIsResolver(startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getRootRootqueryOutgoingAPIsResolver(startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getRootRootqueryIncomingTCPResolver(startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getRootRootqueryOutgoingTCPResolver(startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getRootRootqueryServiceTopologyResolver(startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		fmt.Printf("Failed to getMetrics, err: %v", err)
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}


//////////////////////////////////////
// Child/Link Node : Config Config
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootConfigResolver() (*model.ConfigConfig, error) {
	vConfig, err := c.vRootRoot.GetConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	c.vConfigConfig = vConfig
	vConfigName := string(vConfig.Spec.ConfigName)
FooA, _ := json.Marshal(vConfig.Spec.FooA)
FooAData := string(FooA)
FooMap, _ := json.Marshal(vConfig.Spec.FooMap)
FooMapData := string(FooMap)
FooD, _ := json.Marshal(vConfig.Spec.FooD)
FooDData := string(FooD)
vXYZPort := int(vConfig.Spec.XYZPort)
vABCHost := string(vConfig.Spec.ABCHost)

	ret := &model.ConfigConfig {
	ConfigName: &vConfigName,
	FooA: &FooAData,
	FooMap: &FooMapData,
	FooD: &FooDData,
	XYZPort: &vXYZPort,
	ABCHost: &vABCHost,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomField: CustomBar of CustomType: Root
// Resolver for Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootCustomBarResolver() (*model.RootBar, error) {
	vRoot := c.vRootRoot
	vName := string(vRoot.Spec.CustomBar.Name)

	ret := &model.RootBar {
	Name: &vName,
	}
	return ret, nil
}
//////////////////////////////////////
// Child/Link Node : GNS Gns
// Resolver for Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigGNSResolver() (*model.GnsGns, error) {
	vGns, err := c.vConfigConfig.GetGNS(context.TODO())
	if err != nil {
		panic(err)
	}
	c.vGnsGns = vGns
	vDomain := string(vGns.Spec.Domain)
vUseSharedGateway := bool(vGns.Spec.UseSharedGateway)
vInstance := string(vGns.Spec.Instance)

	ret := &model.GnsGns {
	Domain: &vDomain,
	UseSharedGateway: &vUseSharedGateway,
	Instance: &vInstance,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomField: Cluster of CustomType: Config
// Resolver for Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigClusterResolver() (*model.ConfigCluster, error) {
	vConfig := c.vConfigConfig
	vName := string(vConfig.Spec.Cluster.Name)
vMyID := int(vConfig.Spec.Cluster.MyID)

	ret := &model.ConfigCluster {
	Name: &vName,
	MyID: &vMyID,
	}
	return ret, nil
}
//////////////////////////////////////
// Child/Link Node : FooLink Bar
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooLinkResolver() (*model.GnsBar, error) {
	vBar, err := c.vGnsGns.GetFooLink(context.TODO())
	if err != nil {
		panic(err)
	}
	c.vGnsBar = vBar
	vName := string(vBar.Spec.Name)

	ret := &model.GnsBar {
	Name: &vName,
	}
	return ret, nil
}

//////////////////////////////////////
// Child/Link Node : FooChild Bar
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooChildResolver() (*model.GnsBar, error) {
	vBar, err := c.vGnsGns.GetFooChild(context.TODO())
	if err != nil {
		panic(err)
	}
	c.vGnsBar = vBar
	vName := string(vBar.Spec.Name)

	ret := &model.GnsBar {
	Name: &vName,
	}
	return ret, nil
}

//////////////////////////////////////
// Children/Links Node : FooLinks
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooLinksResolver(id *string) ([]*model.GnsBar, error) {
	var vGnsBarList []*model.GnsBar
	if id != nil && *id != "" {
		vBar, err := c.vGnsGns.GetFooLinks(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		vName := string(vBar.Spec.Name)

		ret := &model.GnsBar {
	Name: &vName,
	}
		vGnsBarList = append(vGnsBarList, ret)
		return vGnsBarList, nil
	}
	for i := range c.vGnsGns.Spec.FooLinksGvk {
		vBar, err := c.vGnsGns.GetFooLinks(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		vName := string(vBar.Spec.Name)

		ret := &model.GnsBar {
	Name: &vName,
	}
		vGnsBarList = append(vGnsBarList, ret)
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
		vBar, err := c.vGnsGns.GetFooChildren(context.TODO(), *id)
		if err != nil {
			panic(err)
		}
		vName := string(vBar.Spec.Name)

		ret := &model.GnsBar {
	Name: &vName,
	}
		vGnsBarList = append(vGnsBarList, ret)
		return vGnsBarList, nil
	}
	for i := range c.vGnsGns.Spec.FooChildrenGvk {
		vBar, err := c.vGnsGns.GetFooChildren(context.TODO(), i)
		if err != nil {
			panic(err)
		}
		vName := string(vBar.Spec.Name)

		ret := &model.GnsBar {
	Name: &vName,
	}
		vGnsBarList = append(vGnsBarList, ret)
	}
	return vGnsBarList, nil
	
}

//////////////////////////////////////
// CustomField: Mydesc of CustomType: Gns
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsMydescResolver() (*model.GnsDescription, error) {
	vGns := c.vGnsGns
	vColor := string(vGns.Spec.Mydesc.Color)
vVersion := string(vGns.Spec.Mydesc.Version)
vProjectID := string(vGns.Spec.Mydesc.ProjectID)
vInstance := string(vGns.Spec.Mydesc.Instance)

	ret := &model.GnsDescription {
	Color: &vColor,
	Version: &vVersion,
	ProjectID: &vProjectID,
	Instance: &vInstance,
	}
	return ret, nil
}
//////////////////////////////////////
// CustomField: HostPort of CustomType: Gns
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsHostPortResolver() (*model.GnsHostPort, error) {
	vGns := c.vGnsGns
	vHost := string(vGns.Spec.HostPort.Host)
vPort := int(vGns.Spec.HostPort.Port)

	ret := &model.GnsHostPort {
	Host: &vHost,
	Port: &vPort,
	}
	return ret, nil
}
//////////////////////////////////////
// CustomField: TestArray of CustomType: Gns
// Resolver for Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsTestArrayResolver() (*model.GnsEmptyData, error) {
	ret := &model.GnsEmptyData {}
	return ret, nil
}