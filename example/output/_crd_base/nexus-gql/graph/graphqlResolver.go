package graph

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
    vConfigFooType *nexus_client.ConfigFooType
    vConfigDomain *nexus_client.ConfigDomain
    vGnsRandomGnsData *nexus_client.GnsRandomGnsData
    vGnsGns *nexus_client.GnsGns
    vGnsBarLink *nexus_client.GnsBarLink
    vGnsDns *nexus_client.GnsDns
    vGnsAdditionalGnsData *nexus_client.GnsAdditionalGnsData
    vServicegroupSvcGroup *nexus_client.ServicegroupSvcGroup
    vPolicypkgAdditionalPolicyData *nexus_client.PolicypkgAdditionalPolicyData
    vPolicypkgAccessControlPolicy *nexus_client.PolicypkgAccessControlPolicy
    vPolicypkgACPConfig *nexus_client.PolicypkgACPConfig
    vPolicypkgVMpolicy *nexus_client.PolicypkgVMpolicy
    vPolicypkgRandomPolicyData *nexus_client.PolicypkgRandomPolicyData
    
}

func getParentName(parentLabels map[string]interface{}, key string) string {
    if v, ok := parentLabels[key]; ok && v != nil {
	    return v.(string)
	}
	return ""
}

//////////////////////////////////////
// Nexus K8sAPIEndpointConfig
//////////////////////////////////////
func getK8sAPIEndpointConfig() *rest.Config {
	filePath := os.Getenv("KUBECONFIG")
	if filePath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", filePath)
		if err != nil {
			return err
		}
		return config
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
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
	    log.Errorf("[getRootResolver]Error getting Root node %s", err)
        return nil, nil
	}
	c.vRootRoot = vRoot
	dn := vRoot.DisplayName()
parentLabels := map[string]interface{}{"roots.root.tsm.tanzu.vmware.com":dn}

	ret := &model.RootRoot {
	Id: &dn,
	ParentLabels: parentLabels,
	}
	log.Debug("[getRootResolver]Output Root object %+v", ret)
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
		log.Errorf("[getRootRootqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
	   log.Errorf("[getRootRootqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
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
		log.Errorf("[getRootRootqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
        log.Errorf("[getRootRootqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
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
		log.Errof("[getRootRootqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getRootRootqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getRootRootqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getRootRootqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getConfigConfigqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
	   log.Errorf("[getConfigConfigqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
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
		log.Errorf("[getConfigConfigqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
        log.Errorf("[getConfigConfigqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
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
		log.Errof("[getConfigConfigqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getConfigConfigqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getConfigConfigqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getConfigConfigqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}



//////////////////////////////////////
// CustomQuery Resolver for Node: FooType in PKG: Config
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getConfigFooTypequeryServiceTableResolver(obj *model.ConfigFooType, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getConfigFooTypequeryServiceVersionTableResolver(obj *model.ConfigFooType, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getConfigFooTypequeryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getConfigFooTypequeryServiceTSResolver(obj *model.ConfigFooType, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getConfigFooTypequeryIncomingAPIsResolver(obj *model.ConfigFooType, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getConfigFooTypequeryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getConfigFooTypequeryOutgoingAPIsResolver(obj *model.ConfigFooType, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getConfigFooTypequeryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getConfigFooTypequeryIncomingTCPResolver(obj *model.ConfigFooType, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getConfigFooTypequeryOutgoingTCPResolver(obj *model.ConfigFooType, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getConfigFooTypequeryServiceTopologyResolver(obj *model.ConfigFooType, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}



//////////////////////////////////////
// CustomQuery Resolver for Node: Domain in PKG: Config
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getConfigDomainqueryServiceTableResolver(obj *model.ConfigDomain, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getConfigDomainqueryServiceVersionTableResolver(obj *model.ConfigDomain, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getConfigDomainqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getConfigDomainqueryServiceTSResolver(obj *model.ConfigDomain, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigDomainqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getConfigDomainqueryIncomingAPIsResolver(obj *model.ConfigDomain, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getConfigDomainqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getConfigDomainqueryOutgoingAPIsResolver(obj *model.ConfigDomain, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getConfigDomainqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getConfigDomainqueryIncomingTCPResolver(obj *model.ConfigDomain, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getConfigDomainqueryOutgoingTCPResolver(obj *model.ConfigDomain, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getConfigDomainqueryServiceTopologyResolver(obj *model.ConfigDomain, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
func (c *resolverConfig) getRootResolver(id *string) ([]*model.GnsRandomGnsData, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	var vRandomGnsDataList []*model.GnsRandomGnsData
	if id != nil && *id != "" {
	    log.Debugf("[getRootResolver]Id: %q", *id)
		vRandomGnsData, err := nc.GetGnsRandomGnsData(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting RandomGnsData node %q: %s", *id, err)
			return nil, nil
		}
		c.vGnsRandomGnsData = vRandomGnsData
		dn := vRandomGnsData.DisplayName()
parentLabels := map[string]interface{}{"randomgnsdatas.gns.tsm.tanzu.vmware.com":dn}

		ret := &model.GnsRandomGnsData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vRandomGnsDataList = append(vRandomGnsDataList, ret)
		log.Debugf("[getRootResolver]Output RandomGnsData objects %+v", vRandomGnsDataList)
		return vRandomGnsDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all RandomGnsDatas")

	vRandomGnsDataListObj, err := nc.RandomGnsData().ListRandomGnsDatas(context.TODO(), metav1.ListOptions{})
	for _,i := range vRandomGnsDataListObj{
		vRandomGnsData, err := nc.GetGnsRandomGnsData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting RandomGnsData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vGnsRandomGnsData = vRandomGnsData
		dn := vRandomGnsData.DisplayName()
parentLabels := map[string]interface{}{"randomgnsdatas.gns.tsm.tanzu.vmware.com":dn}

		ret := &model.GnsRandomGnsData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vRandomGnsDataList = append(vRandomGnsDataList, ret)
	}

	log.Debugf("[getRootResolver]Output RandomGnsData objects %v", vRandomGnsDataList)
	return vRandomGnsDataList, nil
}


//////////////////////////////////////
// CustomQuery Resolver for Node: RandomGnsData in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceTableResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceVersionTableResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getGnsRandomGnsDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceTSResolver(obj *model.GnsRandomGnsData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsRandomGnsDataqueryIncomingAPIsResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getGnsRandomGnsDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsRandomGnsDataqueryOutgoingAPIsResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getGnsRandomGnsDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsRandomGnsDataqueryIncomingTCPResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsRandomGnsDataqueryOutgoingTCPResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceTopologyResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getGnsGnsqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
	   log.Errorf("[getGnsGnsqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
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
		log.Errorf("[getGnsGnsqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
        log.Errorf("[getGnsGnsqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
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
		log.Errof("[getGnsGnsqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getGnsGnsqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getGnsGnsqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getGnsGnsqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}





//////////////////////////////////////
// Singleton Resolver for Parent Node
// PKG: Gns, NODE: Gns
//////////////////////////////////////
func (c *resolverConfig) getRootResolver() (*model.GnsBarLink, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
        return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	vBarLink, err := nc.GetGnsBarLink(context.TODO())
	if err != nil {
	    log.Errorf("[getRootResolver]Error getting BarLink node %s", err)
        return nil, nil
	}
	c.vGnsBarLink = vBarLink
	dn := vBarLink.DisplayName()
parentLabels := map[string]interface{}{"barlinks.gns.tsm.tanzu.vmware.com":dn}

	ret := &model.GnsBarLink {
	Id: &dn,
	ParentLabels: parentLabels,
	}
	log.Debug("[getRootResolver]Output BarLink object %+v", ret)
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
		log.Errorf("[getGnsBarLinkqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
	   log.Errorf("[getGnsBarLinkqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
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
		log.Errorf("[getGnsBarLinkqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
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
        log.Errorf("[getGnsBarLinkqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
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
		log.Errof("[getGnsBarLinkqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getGnsBarLinkqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getGnsBarLinkqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
		log.Errorf("[getGnsBarLinkqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}



//////////////////////////////////////
// CustomQuery Resolver for Node: Dns in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsDnsqueryServiceTableResolver(obj *model.GnsDns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsDnsqueryServiceVersionTableResolver(obj *model.GnsDns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getGnsDnsqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsDnsqueryServiceTSResolver(obj *model.GnsDns, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsDnsqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsDnsqueryIncomingAPIsResolver(obj *model.GnsDns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getGnsDnsqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsDnsqueryOutgoingAPIsResolver(obj *model.GnsDns, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getGnsDnsqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsDnsqueryIncomingTCPResolver(obj *model.GnsDns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsDnsqueryOutgoingTCPResolver(obj *model.GnsDns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsDnsqueryServiceTopologyResolver(obj *model.GnsDns, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
func (c *resolverConfig) getRootResolver(id *string) ([]*model.GnsAdditionalGnsData, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	var vAdditionalGnsDataList []*model.GnsAdditionalGnsData
	if id != nil && *id != "" {
	    log.Debugf("[getRootResolver]Id: %q", *id)
		vAdditionalGnsData, err := nc.GetGnsAdditionalGnsData(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting AdditionalGnsData node %q: %s", *id, err)
			return nil, nil
		}
		c.vGnsAdditionalGnsData = vAdditionalGnsData
		dn := vAdditionalGnsData.DisplayName()
parentLabels := map[string]interface{}{"additionalgnsdatas.gns.tsm.tanzu.vmware.com":dn}

		ret := &model.GnsAdditionalGnsData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vAdditionalGnsDataList = append(vAdditionalGnsDataList, ret)
		log.Debugf("[getRootResolver]Output AdditionalGnsData objects %+v", vAdditionalGnsDataList)
		return vAdditionalGnsDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all AdditionalGnsDatas")

	vAdditionalGnsDataListObj, err := nc.AdditionalGnsData().ListAdditionalGnsDatas(context.TODO(), metav1.ListOptions{})
	for _,i := range vAdditionalGnsDataListObj{
		vAdditionalGnsData, err := nc.GetGnsAdditionalGnsData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting AdditionalGnsData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vGnsAdditionalGnsData = vAdditionalGnsData
		dn := vAdditionalGnsData.DisplayName()
parentLabels := map[string]interface{}{"additionalgnsdatas.gns.tsm.tanzu.vmware.com":dn}

		ret := &model.GnsAdditionalGnsData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vAdditionalGnsDataList = append(vAdditionalGnsDataList, ret)
	}

	log.Debugf("[getRootResolver]Output AdditionalGnsData objects %v", vAdditionalGnsDataList)
	return vAdditionalGnsDataList, nil
}


//////////////////////////////////////
// CustomQuery Resolver for Node: AdditionalGnsData in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceTableResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceVersionTableResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getGnsAdditionalGnsDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceTSResolver(obj *model.GnsAdditionalGnsData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsAdditionalGnsDataqueryIncomingAPIsResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getGnsAdditionalGnsDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsAdditionalGnsDataqueryOutgoingAPIsResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getGnsAdditionalGnsDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsAdditionalGnsDataqueryIncomingTCPResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsAdditionalGnsDataqueryOutgoingTCPResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceTopologyResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}







//////////////////////////////////////
// CustomQuery Resolver for Node: SvcGroup in PKG: Servicegroup
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceTableResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceVersionTableResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getServicegroupSvcGroupqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceTSResolver(obj *model.ServicegroupSvcGroup, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getServicegroupSvcGroupqueryIncomingAPIsResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getServicegroupSvcGroupqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getServicegroupSvcGroupqueryOutgoingAPIsResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getServicegroupSvcGroupqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getServicegroupSvcGroupqueryIncomingTCPResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getServicegroupSvcGroupqueryOutgoingTCPResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceTopologyResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
// PKG: Policypkg, NODE: Policypkg
//////////////////////////////////////
func (c *resolverConfig) getRootResolver(id *string) ([]*model.PolicypkgAdditionalPolicyData, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	var vAdditionalPolicyDataList []*model.PolicypkgAdditionalPolicyData
	if id != nil && *id != "" {
	    log.Debugf("[getRootResolver]Id: %q", *id)
		vAdditionalPolicyData, err := nc.GetPolicypkgAdditionalPolicyData(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting AdditionalPolicyData node %q: %s", *id, err)
			return nil, nil
		}
		c.vPolicypkgAdditionalPolicyData = vAdditionalPolicyData
		dn := vAdditionalPolicyData.DisplayName()
parentLabels := map[string]interface{}{"additionalpolicydatas.policypkg.tsm.tanzu.vmware.com":dn}

		ret := &model.PolicypkgAdditionalPolicyData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vAdditionalPolicyDataList = append(vAdditionalPolicyDataList, ret)
		log.Debugf("[getRootResolver]Output AdditionalPolicyData objects %+v", vAdditionalPolicyDataList)
		return vAdditionalPolicyDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all AdditionalPolicyDatas")

	vAdditionalPolicyDataListObj, err := nc.AdditionalPolicyData().ListAdditionalPolicyDatas(context.TODO(), metav1.ListOptions{})
	for _,i := range vAdditionalPolicyDataListObj{
		vAdditionalPolicyData, err := nc.GetPolicypkgAdditionalPolicyData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting AdditionalPolicyData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vPolicypkgAdditionalPolicyData = vAdditionalPolicyData
		dn := vAdditionalPolicyData.DisplayName()
parentLabels := map[string]interface{}{"additionalpolicydatas.policypkg.tsm.tanzu.vmware.com":dn}

		ret := &model.PolicypkgAdditionalPolicyData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vAdditionalPolicyDataList = append(vAdditionalPolicyDataList, ret)
	}

	log.Debugf("[getRootResolver]Output AdditionalPolicyData objects %v", vAdditionalPolicyDataList)
	return vAdditionalPolicyDataList, nil
}


//////////////////////////////////////
// CustomQuery Resolver for Node: AdditionalPolicyData in PKG: Policypkg
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryServiceTableResolver(obj *model.PolicypkgAdditionalPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAdditionalPolicyDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryServiceVersionTableResolver(obj *model.PolicypkgAdditionalPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getPolicypkgAdditionalPolicyDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryServiceTSResolver(obj *model.PolicypkgAdditionalPolicyData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicypkgAdditionalPolicyDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryIncomingAPIsResolver(obj *model.PolicypkgAdditionalPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getPolicypkgAdditionalPolicyDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryOutgoingAPIsResolver(obj *model.PolicypkgAdditionalPolicyData, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getPolicypkgAdditionalPolicyDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryIncomingTCPResolver(obj *model.PolicypkgAdditionalPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAdditionalPolicyDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryOutgoingTCPResolver(obj *model.PolicypkgAdditionalPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAdditionalPolicyDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicypkgAdditionalPolicyDataqueryServiceTopologyResolver(obj *model.PolicypkgAdditionalPolicyData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAdditionalPolicyDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}







//////////////////////////////////////
// CustomQuery Resolver for Node: AccessControlPolicy in PKG: Policypkg
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryServiceTableResolver(obj *model.PolicypkgAccessControlPolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAccessControlPolicyqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryServiceVersionTableResolver(obj *model.PolicypkgAccessControlPolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getPolicypkgAccessControlPolicyqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryServiceTSResolver(obj *model.PolicypkgAccessControlPolicy, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicypkgAccessControlPolicyqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryIncomingAPIsResolver(obj *model.PolicypkgAccessControlPolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getPolicypkgAccessControlPolicyqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryOutgoingAPIsResolver(obj *model.PolicypkgAccessControlPolicy, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getPolicypkgAccessControlPolicyqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryIncomingTCPResolver(obj *model.PolicypkgAccessControlPolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAccessControlPolicyqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryOutgoingTCPResolver(obj *model.PolicypkgAccessControlPolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAccessControlPolicyqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicypkgAccessControlPolicyqueryServiceTopologyResolver(obj *model.PolicypkgAccessControlPolicy, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgAccessControlPolicyqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}



//////////////////////////////////////
// CustomQuery Resolver for Node: ACPConfig in PKG: Policypkg
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicypkgACPConfigqueryServiceTableResolver(obj *model.PolicypkgACPConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgACPConfigqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicypkgACPConfigqueryServiceVersionTableResolver(obj *model.PolicypkgACPConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getPolicypkgACPConfigqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicypkgACPConfigqueryServiceTSResolver(obj *model.PolicypkgACPConfig, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicypkgACPConfigqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicypkgACPConfigqueryIncomingAPIsResolver(obj *model.PolicypkgACPConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getPolicypkgACPConfigqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicypkgACPConfigqueryOutgoingAPIsResolver(obj *model.PolicypkgACPConfig, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getPolicypkgACPConfigqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicypkgACPConfigqueryIncomingTCPResolver(obj *model.PolicypkgACPConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgACPConfigqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicypkgACPConfigqueryOutgoingTCPResolver(obj *model.PolicypkgACPConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgACPConfigqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicypkgACPConfigqueryServiceTopologyResolver(obj *model.PolicypkgACPConfig, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgACPConfigqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}













//////////////////////////////////////
// CustomQuery Resolver for Node: VMpolicy in PKG: Policypkg
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicypkgVMpolicyqueryServiceTableResolver(obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgVMpolicyqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicypkgVMpolicyqueryServiceVersionTableResolver(obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getPolicypkgVMpolicyqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicypkgVMpolicyqueryServiceTSResolver(obj *model.PolicypkgVMpolicy, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicypkgVMpolicyqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicypkgVMpolicyqueryIncomingAPIsResolver(obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getPolicypkgVMpolicyqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicypkgVMpolicyqueryOutgoingAPIsResolver(obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getPolicypkgVMpolicyqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicypkgVMpolicyqueryIncomingTCPResolver(obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgVMpolicyqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicypkgVMpolicyqueryOutgoingTCPResolver(obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgVMpolicyqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicypkgVMpolicyqueryServiceTopologyResolver(obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgVMpolicyqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
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
// PKG: Policypkg, NODE: Policypkg
//////////////////////////////////////
func (c *resolverConfig) getRootResolver(id *string) ([]*model.PolicypkgRandomPolicyData, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	var vRandomPolicyDataList []*model.PolicypkgRandomPolicyData
	if id != nil && *id != "" {
	    log.Debugf("[getRootResolver]Id: %q", *id)
		vRandomPolicyData, err := nc.GetPolicypkgRandomPolicyData(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting RandomPolicyData node %q: %s", *id, err)
			return nil, nil
		}
		c.vPolicypkgRandomPolicyData = vRandomPolicyData
		dn := vRandomPolicyData.DisplayName()
parentLabels := map[string]interface{}{"randompolicydatas.policypkg.tsm.tanzu.vmware.com":dn}

		ret := &model.PolicypkgRandomPolicyData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vRandomPolicyDataList = append(vRandomPolicyDataList, ret)
		log.Debugf("[getRootResolver]Output RandomPolicyData objects %+v", vRandomPolicyDataList)
		return vRandomPolicyDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all RandomPolicyDatas")

	vRandomPolicyDataListObj, err := nc.RandomPolicyData().ListRandomPolicyDatas(context.TODO(), metav1.ListOptions{})
	for _,i := range vRandomPolicyDataListObj{
		vRandomPolicyData, err := nc.GetPolicypkgRandomPolicyData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting RandomPolicyData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vPolicypkgRandomPolicyData = vRandomPolicyData
		dn := vRandomPolicyData.DisplayName()
parentLabels := map[string]interface{}{"randompolicydatas.policypkg.tsm.tanzu.vmware.com":dn}

		ret := &model.PolicypkgRandomPolicyData {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vRandomPolicyDataList = append(vRandomPolicyDataList, ret)
	}

	log.Debugf("[getRootResolver]Output RandomPolicyData objects %v", vRandomPolicyDataList)
	return vRandomPolicyDataList, nil
}


//////////////////////////////////////
// CustomQuery Resolver for Node: RandomPolicyData in PKG: Policypkg
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryServiceTableResolver(obj *model.PolicypkgRandomPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgRandomPolicyDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryServiceVersionTableResolver(obj *model.PolicypkgRandomPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
	   log.Errorf("[getPolicypkgRandomPolicyDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
	   return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryServiceTSResolver(obj *model.PolicypkgRandomPolicyData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicypkgRandomPolicyDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryIncomingAPIsResolver(obj *model.PolicypkgRandomPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
        log.Errorf("[getPolicypkgRandomPolicyDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
        return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryOutgoingAPIsResolver(obj *model.PolicypkgRandomPolicyData, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: *timeInterval})
	if err != nil {
		log.Errof("[getPolicypkgRandomPolicyDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryIncomingTCPResolver(obj *model.PolicypkgRandomPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgRandomPolicyDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryOutgoingTCPResolver(obj *model.PolicypkgRandomPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgRandomPolicyDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicypkgRandomPolicyDataqueryServiceTopologyResolver(obj *model.PolicypkgRandomPolicyData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData,error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters,TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicypkgRandomPolicyDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret,nil
}






























































































































