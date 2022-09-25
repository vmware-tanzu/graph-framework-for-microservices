package graph

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"

	qm "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/query-manager"
	libgrpc "gitlab.eng.vmware.com/nsx-allspark_users/lib-go/grpc"
	nexus_client "nexustempmodule/nexus-client"
	"nexustempmodule/nexus-gql/graph/model"

)

var c resolverConfig
var nc *nexus_client.Clientset

type resolverConfig struct {
	vRootRoot *nexus_client.RootRoot
    vConfigConfig *nexus_client.ConfigConfig
    vGnsGns *nexus_client.GnsGns
    vGnsBarLink *nexus_client.GnsBarLink
    vGnsBarChild *nexus_client.GnsBarChild
    vGnsBarChildren *nexus_client.GnsBarChildren
    vGnsBarLinks *nexus_client.GnsBarLinks
    
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
// Singleton Resolver for Parent Node
// PKG: Root, NODE: Root
//////////////////////////////////////
func (c *resolverConfig) getRootResolver() (*model.RootRoot, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
        return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	vRoot, err := nc.GetRootRoot(context.TODO())
	if err != nil {
	    log.Errorf("Error getting root node %s", err)
        return nil, fmt.Errorf("failed to get root node: %s", err)
	}
	c.vRootRoot = vRoot
	id := vRoot.DisplayName()
parentLabels := map[string]interface{}{"roots.root.tsm.tanzu.vmware.com":id}
vDisplayName := string(vRoot.Spec.DisplayName)
CustomBar, _ := json.Marshal(vRoot.Spec.CustomBar)
CustomBarData := string(CustomBar)

	ret := &model.RootRoot {
	Id: &id,
	ParentLabels: parentLabels,
	DisplayName: &vDisplayName,
	CustomBar: &CustomBarData,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: Root in PKG: Root
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getRootRootqueryServiceTableResolver(obj *model.RootRoot, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getRootRootqueryServiceVersionTableResolver(obj *model.RootRoot, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getRootRootqueryServiceTSResolver(obj *model.RootRoot, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getRootRootqueryIncomingAPIsResolver(obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getRootRootqueryOutgoingAPIsResolver(obj *model.RootRoot, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getRootRootqueryIncomingTCPResolver(obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getRootRootqueryOutgoingTCPResolver(obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getRootRootqueryServiceTopologyResolver(obj *model.RootRoot, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
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
// CustomQuery Resolver for Node: Config in PKG: Config
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getConfigConfigqueryServiceTableResolver(obj *model.ConfigConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getConfigConfigqueryServiceVersionTableResolver(obj *model.ConfigConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getConfigConfigqueryServiceTSResolver(obj *model.ConfigConfig, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getConfigConfigqueryIncomingAPIsResolver(obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getConfigConfigqueryOutgoingAPIsResolver(obj *model.ConfigConfig, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getConfigConfigqueryIncomingTCPResolver(obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getConfigConfigqueryOutgoingTCPResolver(obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getConfigConfigqueryServiceTopologyResolver(obj *model.ConfigConfig, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
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
// CustomQuery Resolver for Node: Gns in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsGnsqueryServiceTableResolver(obj *model.GnsGns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsGnsqueryServiceVersionTableResolver(obj *model.GnsGns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsGnsqueryServiceTSResolver(obj *model.GnsGns, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsGnsqueryIncomingAPIsResolver(obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsGnsqueryOutgoingAPIsResolver(obj *model.GnsGns, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsGnsqueryIncomingTCPResolver(obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsGnsqueryOutgoingTCPResolver(obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsGnsqueryServiceTopologyResolver(obj *model.GnsGns, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
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
// Non Singleton Resolver for Parent Node
// PKG: Gns, NODE: Gns
//////////////////////////////////////
func (c *resolverConfig) getRootResolver(id *string) (*model.GnsBarLink, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	if id != nil && *id != "" {
		vBarLink, err := nc.GetGnsBarLink(context.TODO(), *id)
		if err != nil {
			log.Errorf("Error getting root node %s", err)
			return nil, fmt.Errorf("failed to get root node: %s", err)
		}
		c.vGnsBarLink = vBarLink
		id := vBarLink.DisplayName()
parentLabels := map[string]interface{}{"barlinks.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarLink.Spec.Name)

		ret := &model.GnsBarLink {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
		return ret, nil
	}
	ret := &model.GnsBarLink{}
	return ret, nil
}


//////////////////////////////////////
// CustomQuery Resolver for Node: BarLink in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsBarLinkqueryServiceTableResolver(obj *model.GnsBarLink, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinkqueryServiceVersionTableResolver(obj *model.GnsBarLink, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinkqueryServiceTSResolver(obj *model.GnsBarLink, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinkqueryIncomingAPIsResolver(obj *model.GnsBarLink, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinkqueryOutgoingAPIsResolver(obj *model.GnsBarLink, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinkqueryIncomingTCPResolver(obj *model.GnsBarLink, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinkqueryOutgoingTCPResolver(obj *model.GnsBarLink, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinkqueryServiceTopologyResolver(obj *model.GnsBarLink, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
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
// CustomQuery Resolver for Node: BarChild in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsBarChildqueryServiceTableResolver(obj *model.GnsBarChild, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildqueryServiceVersionTableResolver(obj *model.GnsBarChild, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildqueryServiceTSResolver(obj *model.GnsBarChild, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildqueryIncomingAPIsResolver(obj *model.GnsBarChild, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildqueryOutgoingAPIsResolver(obj *model.GnsBarChild, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildqueryIncomingTCPResolver(obj *model.GnsBarChild, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildqueryOutgoingTCPResolver(obj *model.GnsBarChild, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildqueryServiceTopologyResolver(obj *model.GnsBarChild, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
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
// CustomQuery Resolver for Node: BarChildren in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsBarChildrenqueryServiceTableResolver(obj *model.GnsBarChildren, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildrenqueryServiceVersionTableResolver(obj *model.GnsBarChildren, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildrenqueryServiceTSResolver(obj *model.GnsBarChildren, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildrenqueryIncomingAPIsResolver(obj *model.GnsBarChildren, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildrenqueryOutgoingAPIsResolver(obj *model.GnsBarChildren, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildrenqueryIncomingTCPResolver(obj *model.GnsBarChildren, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildrenqueryOutgoingTCPResolver(obj *model.GnsBarChildren, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarChildrenqueryServiceTopologyResolver(obj *model.GnsBarChildren, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
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
// Non Singleton Resolver for Parent Node
// PKG: Gns, NODE: Gns
//////////////////////////////////////
func (c *resolverConfig) getRootResolver(id *string) (*model.GnsBarLinks, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	if id != nil && *id != "" {
		vBarLinks, err := nc.GetGnsBarLinks(context.TODO(), *id)
		if err != nil {
			log.Errorf("Error getting root node %s", err)
			return nil, fmt.Errorf("failed to get root node: %s", err)
		}
		c.vGnsBarLinks = vBarLinks
		id := vBarLinks.DisplayName()
parentLabels := map[string]interface{}{"barlinkses.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarLinks.Spec.Name)

		ret := &model.GnsBarLinks {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
		return ret, nil
	}
	ret := &model.GnsBarLinks{}
	return ret, nil
}


//////////////////////////////////////
// CustomQuery Resolver for Node: BarLinks in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsBarLinksqueryServiceTableResolver(obj *model.GnsBarLinks, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinksqueryServiceVersionTableResolver(obj *model.GnsBarLinks, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinksqueryServiceTSResolver(obj *model.GnsBarLinks, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinksqueryIncomingAPIsResolver(obj *model.GnsBarLinks, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinksqueryOutgoingAPIsResolver(obj *model.GnsBarLinks, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinksqueryIncomingTCPResolver(obj *model.GnsBarLinks, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinksqueryOutgoingTCPResolver(obj *model.GnsBarLinks, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
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
func (c *resolverConfig) getGnsBarLinksqueryServiceTopologyResolver(obj *model.GnsBarLinks, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
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
// CHILD RESOLVER (Non Singleton)
// FieldName: Config Node: Root PKG: Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootConfigResolver(obj *model.RootRoot, id *string) (*model.ConfigConfig, error) {
	vConfig, err := nc.RootRoot().GetConfig(context.TODO(), *id)
	if err != nil {
	    log.Errorf("Error getting node %s", err)
        return nil, fmt.Errorf("failed to get node: %s", err)
    }
	c.vConfigConfig = vConfig
	id := vConfig.DisplayName()
parentLabels := map[string]interface{}{"configs.config.tsm.tanzu.vmware.com":id}
vConfigName := string(vConfig.Spec.ConfigName)
Cluster, _ := json.Marshal(vConfig.Spec.Cluster)
ClusterData := string(Cluster)
FooA, _ := json.Marshal(vConfig.Spec.FooA)
FooAData := string(FooA)
FooMap, _ := json.Marshal(vConfig.Spec.FooMap)
FooMapData := string(FooMap)
FooB, _ := json.Marshal(vConfig.Spec.FooB)
FooBData := string(FooB)
FooD, _ := json.Marshal(vConfig.Spec.FooD)
FooDData := string(FooD)
FooF, _ := json.Marshal(vConfig.Spec.FooF)
FooFData := string(FooF)
XYZPort, _ := json.Marshal(vConfig.Spec.XYZPort)
XYZPortData := string(XYZPort)
ABCHost, _ := json.Marshal(vConfig.Spec.ABCHost)
ABCHostData := string(ABCHost)
ClusterNamespaces, _ := json.Marshal(vConfig.Spec.ClusterNamespaces)
ClusterNamespacesData := string(ClusterNamespaces)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.ConfigConfig {
	Id: &id,
	ParentLabels: parentLabels,
	ConfigName: &vConfigName,
	Cluster: &ClusterData,
	FooA: &FooAData,
	FooMap: &FooMapData,
	FooB: &FooBData,
	FooD: &FooDData,
	FooF: &FooFData,
	XYZPort: &XYZPortData,
	ABCHost: &ABCHostData,
	ClusterNamespaces: &ClusterNamespacesData,
	}
	return ret, nil
}






//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: GNS Node: Config PKG: Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigGNSResolver(obj *model.ConfigConfig, id *string) (*model.GnsGns, error) {
	vGns, err := nc.RootRoot().Config(obj.ParentLabels["configs.config.tsm.tanzu.vmware.com"].(string)).GetGNS(context.TODO(), *id)
	if err != nil {
	    log.Errorf("Error getting node %s", err)
        return nil, fmt.Errorf("failed to get node: %s", err)
    }
	c.vGnsGns = vGns
	id := vGns.DisplayName()
parentLabels := map[string]interface{}{"gnses.gns.tsm.tanzu.vmware.com":id}
vDomain := string(vGns.Spec.Domain)
vUseSharedGateway := bool(vGns.Spec.UseSharedGateway)
Mydesc, _ := json.Marshal(vGns.Spec.Mydesc)
MydescData := string(Mydesc)
HostPort, _ := json.Marshal(vGns.Spec.HostPort)
HostPortData := string(HostPort)
Instance, _ := json.Marshal(vGns.Spec.Instance)
InstanceData := string(Instance)
vArray1 := float64(vGns.Spec.Array1)
Array2, _ := json.Marshal(vGns.Spec.Array2)
Array2Data := string(Array2)
Array3, _ := json.Marshal(vGns.Spec.Array3)
Array3Data := string(Array3)
Array4, _ := json.Marshal(vGns.Spec.Array4)
Array4Data := string(Array4)
Array5, _ := json.Marshal(vGns.Spec.Array5)
Array5Data := string(Array5)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.GnsGns {
	Id: &id,
	ParentLabels: parentLabels,
	Domain: &vDomain,
	UseSharedGateway: &vUseSharedGateway,
	Mydesc: &MydescData,
	HostPort: &HostPortData,
	Instance: &InstanceData,
	Array1: &vArray1,
	Array2: &Array2Data,
	Array3: &Array3Data,
	Array4: &Array4Data,
	Array5: &Array5Data,
	}
	return ret, nil
}















//////////////////////////////////////
// CHILD RESOLVER (Singleton)
// FieldName: FooChild Node: Gns PKG: Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooChildResolver(obj *model.GnsGns) (*model.GnsBarChild, error) {
	vBarChild, err := nc.RootRoot().Config(obj.ParentLabels["configs.config.tsm.tanzu.vmware.com"].(string)).GNS(obj.ParentLabels["gnses.gns.tsm.tanzu.vmware.com"].(string)).GetFooChild(context.TODO())
	if err != nil {
	    log.Errorf("Error getting node %s", err)
        return nil, fmt.Errorf("failed to get node: %s", err)
    }
	c.vGnsBarChild = vBarChild
	id := vBarChild.DisplayName()
parentLabels := map[string]interface{}{"barchilds.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarChild.Spec.Name)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.GnsBarChild {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
	return ret, nil
}

//////////////////////////////////////
// LINK RESOLVER (Non Singleton)
// FieldName: FooLink Node: Gns PKG: Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooLinkResolver(obj *model.GnsGns, id *string) (*model.GnsBarLink, error) {
	vBarLinkParent, err := nc.RootRoot().Config(obj.ParentLabels["configs.config.tsm.tanzu.vmware.com"].(string)).GetGNS(context.TODO(), obj.ParentLabels["gnses.gns.tsm.tanzu.vmware.com"].(string))
	if err != nil {
	    log.Errorf("Error getting Parent node details %s", err)
        return nil, fmt.Errorf("failed to get Parent node details: %s", err)
    }
	vBarLink, err := vBarLinkParent.GetFooLink(context.TODO(), *id)
	if err != nil {
	    log.Errorf("Error getting node %s", err)
        return nil, fmt.Errorf("failed to get node: %s", err)
    }
	c.vGnsBarLink = vBarLink
	id := vBarLink.DisplayName()
parentLabels := map[string]interface{}{"barlinks.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarLink.Spec.Name)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.GnsBarLink {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER (Non Singleton
// FieldName: FooChildren Node: Gns PKG: Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooChildrenResolver(obj *model.GnsGns, id *string) ([]*model.GnsBarChildren, error) {
	var vGnsBarChildrenList []*model.GnsBarChildren
	if id != nil && *id != "" {
		vBarChildren, err := nc.RootRoot().Config(obj.ParentLabels["configs.config.tsm.tanzu.vmware.com"].(string)).GNS(obj.ParentLabels["gnses.gns.tsm.tanzu.vmware.com"].(string)).GetFooChildren(context.TODO(), *id)
		if err != nil {
	        log.Errorf("Error getting node %s", err)
            return nil, fmt.Errorf("failed to get node: %s", err)
        }
		id := vBarChildren.DisplayName()
parentLabels := map[string]interface{}{"barchildrens.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarChildren.Spec.Name)

        for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.GnsBarChildren {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
		vGnsBarChildrenList = append(vGnsBarChildrenList, ret)
		return vGnsBarChildrenList, nil
	}
	for i := range c.vGnsGns.Spec.FooChildrenGvk {
		vBarChildren, err := nc.RootRoot().Config(obj.ParentLabels["configs.config.tsm.tanzu.vmware.com"].(string)).GNS(obj.ParentLabels["gnses.gns.tsm.tanzu.vmware.com"].(string)).GetFooChildren(context.TODO(), i)
		if err != nil {
	        log.Errorf("Error getting node %s", err)
            return nil, fmt.Errorf("failed to get node: %s", err)
		}
		id := vBarChildren.DisplayName()
parentLabels := map[string]interface{}{"barchildrens.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarChildren.Spec.Name)

		ret := &model.GnsBarChildren {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
		vGnsBarChildrenList = append(vGnsBarChildrenList, ret)
	}
	return vGnsBarChildrenList, nil
}

//////////////////////////////////////
// LINKS RESOLVER (Singleton/Non Singleton)
// FieldName: FooLinks Node: Gns PKG: Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsFooLinksResolver(obj *model.GnsGns, id *string) ([]*model.GnsBarLinks, error) {
	var vGnsBarLinksList []*model.GnsBarLinks
	if id != nil && *id != "" {
		vBarLinksParent, err := nc.RootRoot().Config(obj.ParentLabels["configs.config.tsm.tanzu.vmware.com"].(string)).GetGNS(context.TODO(), obj.ParentLabels["gnses.gns.tsm.tanzu.vmware.com"].(string))
		if err != nil {
			log.Errorf("Error getting Parent node details %s", err)
			return nil, fmt.Errorf("failed to get Parent node details: %s", err)
		}
		vBarLinks, err := vBarLinksParent.GetFooLinks(context.TODO())
		if err != nil {
			log.Errorf("Error getting node %s", err)
			return nil, fmt.Errorf("failed to get node: %s", err)
		}
		id := vBarLinks.DisplayName()
parentLabels := map[string]interface{}{"barlinkses.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarLinks.Spec.Name)

        for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.GnsBarLinks {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
		vGnsBarLinksList = append(vGnsBarLinksList, ret)
		return vGnsBarLinksList, nil
	}
	for i := range c.vGnsGns.Spec.FooLinksGvk {
		vBarLinksParent, err := nc.RootRoot().Config(obj.ParentLabels["configs.config.tsm.tanzu.vmware.com"].(string)).GetGNS(context.TODO(), obj.ParentLabels["gnses.gns.tsm.tanzu.vmware.com"].(string))
		if err != nil {
			log.Errorf("Error getting Parent node details %s", err)
			return nil, fmt.Errorf("failed to get Parent node details: %s", err)
		}
		vBarLinks, err := vBarLinksParent.GetFooLinks(context.TODO())
		if err != nil {
			log.Errorf("Error getting node %s", err)
			return nil, fmt.Errorf("failed to get node: %s", err)
		}
		id := vBarLinks.DisplayName()
parentLabels := map[string]interface{}{"barlinkses.gns.tsm.tanzu.vmware.com":id}
vName := string(vBarLinks.Spec.Name)

		ret := &model.GnsBarLinks {
	Id: &id,
	ParentLabels: parentLabels,
	Name: &vName,
	}
		vGnsBarLinksList = append(vGnsBarLinksList, ret)
	}
	return vGnsBarLinksList, nil
}


















