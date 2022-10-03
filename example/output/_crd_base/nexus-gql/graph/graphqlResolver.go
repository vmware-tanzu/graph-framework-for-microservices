package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
	vRootRoot                   *nexus_client.RootRoot
	vConfigConfig               *nexus_client.ConfigConfig
	vConfigFooType              *nexus_client.ConfigFooType
	vConfigDomain               *nexus_client.ConfigDomain
	vGnsRandomGnsData           *nexus_client.GnsRandomGnsData
	vGnsGns                     *nexus_client.GnsGns
	vGnsBarLink                 *nexus_client.GnsBarLink
	vGnsDns                     *nexus_client.GnsDns
	vGnsAdditionalGnsData       *nexus_client.GnsAdditionalGnsData
	vServicegroupSvcGroup       *nexus_client.ServicegroupSvcGroup
	vPolicyAdditionalPolicyData *nexus_client.PolicyAdditionalPolicyData
	vPolicyAccessControlPolicy  *nexus_client.PolicyAccessControlPolicy
	vPolicyACPConfig            *nexus_client.PolicyACPConfig
	vPolicyVMpolicy             *nexus_client.PolicyVMpolicy
	vPolicyRandomPolicyData     *nexus_client.PolicyRandomPolicyData
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
func grpcServer() qm.ServerClient {
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
	parentLabels := map[string]interface{}{"roots.root.tsm.tanzu.vmware.com": dn}

	ret := &model.RootRoot{
		Id:           &dn,
		ParentLabels: parentLabels,
	}
	log.Debug("[getRootResolver]Output Root object %+v", ret)
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: Root in PKG: Root
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getRootRootqueryServiceTableResolver(obj *model.RootRoot, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getRootRootqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getRootRootqueryServiceVersionTableResolver(obj *model.RootRoot, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getRootRootqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getRootRootqueryServiceTSResolver(obj *model.RootRoot, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getRootRootqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getRootRootqueryIncomingAPIsResolver(obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getRootRootqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getRootRootqueryOutgoingAPIsResolver(obj *model.RootRoot, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getRootRootqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getRootRootqueryIncomingTCPResolver(obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getRootRootqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getRootRootqueryOutgoingTCPResolver(obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getRootRootqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getRootRootqueryServiceTopologyResolver(obj *model.RootRoot, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getRootRootqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: Config in PKG: Config
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getConfigConfigqueryServiceTableResolver(obj *model.ConfigConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigConfigqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getConfigConfigqueryServiceVersionTableResolver(obj *model.ConfigConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigConfigqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getConfigConfigqueryServiceTSResolver(obj *model.ConfigConfig, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigConfigqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getConfigConfigqueryIncomingAPIsResolver(obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigConfigqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getConfigConfigqueryOutgoingAPIsResolver(obj *model.ConfigConfig, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigConfigqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getConfigConfigqueryIncomingTCPResolver(obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigConfigqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getConfigConfigqueryOutgoingTCPResolver(obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigConfigqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getConfigConfigqueryServiceTopologyResolver(obj *model.ConfigConfig, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigConfigqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: FooType in PKG: Config
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getConfigFooTypequeryServiceTableResolver(obj *model.ConfigFooType, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getConfigFooTypequeryServiceVersionTableResolver(obj *model.ConfigFooType, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getConfigFooTypequeryServiceTSResolver(obj *model.ConfigFooType, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getConfigFooTypequeryIncomingAPIsResolver(obj *model.ConfigFooType, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getConfigFooTypequeryOutgoingAPIsResolver(obj *model.ConfigFooType, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getConfigFooTypequeryIncomingTCPResolver(obj *model.ConfigFooType, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getConfigFooTypequeryOutgoingTCPResolver(obj *model.ConfigFooType, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getConfigFooTypequeryServiceTopologyResolver(obj *model.ConfigFooType, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigFooTypequeryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: Domain in PKG: Config
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getConfigDomainqueryServiceTableResolver(obj *model.ConfigDomain, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getConfigDomainqueryServiceVersionTableResolver(obj *model.ConfigDomain, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getConfigDomainqueryServiceTSResolver(obj *model.ConfigDomain, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigDomainqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getConfigDomainqueryIncomingAPIsResolver(obj *model.ConfigDomain, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigDomainqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getConfigDomainqueryOutgoingAPIsResolver(obj *model.ConfigDomain, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getConfigDomainqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getConfigDomainqueryIncomingTCPResolver(obj *model.ConfigDomain, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getConfigDomainqueryOutgoingTCPResolver(obj *model.ConfigDomain, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getConfigDomainqueryServiceTopologyResolver(obj *model.ConfigDomain, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getConfigDomainqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
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
		parentLabels := map[string]interface{}{"randomgnsdatas.gns.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vRandomGnsData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.GnsRandomGnsData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
		}
		vRandomGnsDataList = append(vRandomGnsDataList, ret)
		log.Debugf("[getRootResolver]Output RandomGnsData objects %+v", vRandomGnsDataList)
		return vRandomGnsDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all RandomGnsDatas")

	vRandomGnsDataListObj, err := nc.RandomGnsData().ListRandomGnsDatas(context.TODO(), metav1.ListOptions{})
	for _, i := range vRandomGnsDataListObj {
		vRandomGnsData, err := nc.GetGnsRandomGnsData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting RandomGnsData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vGnsRandomGnsData = vRandomGnsData
		dn := vRandomGnsData.DisplayName()
		parentLabels := map[string]interface{}{"randomgnsdatas.gns.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vRandomGnsData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.GnsRandomGnsData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
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
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceTableResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceVersionTableResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceTSResolver(obj *model.GnsRandomGnsData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsRandomGnsDataqueryIncomingAPIsResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsRandomGnsDataqueryOutgoingAPIsResolver(obj *model.GnsRandomGnsData, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsRandomGnsDataqueryIncomingTCPResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsRandomGnsDataqueryOutgoingTCPResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsRandomGnsDataqueryServiceTopologyResolver(obj *model.GnsRandomGnsData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsRandomGnsDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: Gns in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsGnsqueryServiceTableResolver(obj *model.GnsGns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsGnsqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsGnsqueryServiceVersionTableResolver(obj *model.GnsGns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsGnsqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsGnsqueryServiceTSResolver(obj *model.GnsGns, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsGnsqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsGnsqueryIncomingAPIsResolver(obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsGnsqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsGnsqueryOutgoingAPIsResolver(obj *model.GnsGns, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsGnsqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsGnsqueryIncomingTCPResolver(obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsGnsqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsGnsqueryOutgoingTCPResolver(obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsGnsqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsGnsqueryServiceTopologyResolver(obj *model.GnsGns, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsGnsqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
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
	parentLabels := map[string]interface{}{"barlinks.gns.tsm.tanzu.vmware.com": dn}
	vName := string(vBarLink.Spec.Name)

	ret := &model.GnsBarLink{
		Id:           &dn,
		ParentLabels: parentLabels,
		Name:         &vName,
	}
	log.Debug("[getRootResolver]Output BarLink object %+v", ret)
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: BarLink in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsBarLinkqueryServiceTableResolver(obj *model.GnsBarLink, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsBarLinkqueryServiceVersionTableResolver(obj *model.GnsBarLink, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsBarLinkqueryServiceTSResolver(obj *model.GnsBarLink, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsBarLinkqueryIncomingAPIsResolver(obj *model.GnsBarLink, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsBarLinkqueryOutgoingAPIsResolver(obj *model.GnsBarLink, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsBarLinkqueryIncomingTCPResolver(obj *model.GnsBarLink, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsBarLinkqueryOutgoingTCPResolver(obj *model.GnsBarLink, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsBarLinkqueryServiceTopologyResolver(obj *model.GnsBarLink, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsBarLinkqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: Dns in PKG: Gns
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getGnsDnsqueryServiceTableResolver(obj *model.GnsDns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsDnsqueryServiceVersionTableResolver(obj *model.GnsDns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsDnsqueryServiceTSResolver(obj *model.GnsDns, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsDnsqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsDnsqueryIncomingAPIsResolver(obj *model.GnsDns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsDnsqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsDnsqueryOutgoingAPIsResolver(obj *model.GnsDns, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsDnsqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsDnsqueryIncomingTCPResolver(obj *model.GnsDns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsDnsqueryOutgoingTCPResolver(obj *model.GnsDns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsDnsqueryServiceTopologyResolver(obj *model.GnsDns, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsDnsqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
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
		parentLabels := map[string]interface{}{"additionalgnsdatas.gns.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vAdditionalGnsData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.GnsAdditionalGnsData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
		}
		vAdditionalGnsDataList = append(vAdditionalGnsDataList, ret)
		log.Debugf("[getRootResolver]Output AdditionalGnsData objects %+v", vAdditionalGnsDataList)
		return vAdditionalGnsDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all AdditionalGnsDatas")

	vAdditionalGnsDataListObj, err := nc.AdditionalGnsData().ListAdditionalGnsDatas(context.TODO(), metav1.ListOptions{})
	for _, i := range vAdditionalGnsDataListObj {
		vAdditionalGnsData, err := nc.GetGnsAdditionalGnsData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting AdditionalGnsData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vGnsAdditionalGnsData = vAdditionalGnsData
		dn := vAdditionalGnsData.DisplayName()
		parentLabels := map[string]interface{}{"additionalgnsdatas.gns.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vAdditionalGnsData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.GnsAdditionalGnsData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
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
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceTableResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceVersionTableResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceTSResolver(obj *model.GnsAdditionalGnsData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getGnsAdditionalGnsDataqueryIncomingAPIsResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getGnsAdditionalGnsDataqueryOutgoingAPIsResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getGnsAdditionalGnsDataqueryIncomingTCPResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getGnsAdditionalGnsDataqueryOutgoingTCPResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getGnsAdditionalGnsDataqueryServiceTopologyResolver(obj *model.GnsAdditionalGnsData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getGnsAdditionalGnsDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: SvcGroup in PKG: Servicegroup
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceTableResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceVersionTableResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceTSResolver(obj *model.ServicegroupSvcGroup, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getServicegroupSvcGroupqueryIncomingAPIsResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getServicegroupSvcGroupqueryOutgoingAPIsResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getServicegroupSvcGroupqueryIncomingTCPResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getServicegroupSvcGroupqueryOutgoingTCPResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getServicegroupSvcGroupqueryServiceTopologyResolver(obj *model.ServicegroupSvcGroup, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getServicegroupSvcGroupqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// Non Singleton Resolver for Parent Node
// PKG: Policy, NODE: Policy
//////////////////////////////////////
func (c *resolverConfig) getRootResolver(id *string) ([]*model.PolicyAdditionalPolicyData, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	var vAdditionalPolicyDataList []*model.PolicyAdditionalPolicyData
	if id != nil && *id != "" {
		log.Debugf("[getRootResolver]Id: %q", *id)
		vAdditionalPolicyData, err := nc.GetPolicyAdditionalPolicyData(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting AdditionalPolicyData node %q: %s", *id, err)
			return nil, nil
		}
		c.vPolicyAdditionalPolicyData = vAdditionalPolicyData
		dn := vAdditionalPolicyData.DisplayName()
		parentLabels := map[string]interface{}{"additionalpolicydatas.policypkg.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vAdditionalPolicyData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.PolicyAdditionalPolicyData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
		}
		vAdditionalPolicyDataList = append(vAdditionalPolicyDataList, ret)
		log.Debugf("[getRootResolver]Output AdditionalPolicyData objects %+v", vAdditionalPolicyDataList)
		return vAdditionalPolicyDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all AdditionalPolicyDatas")

	vAdditionalPolicyDataListObj, err := nc.AdditionalPolicyData().ListAdditionalPolicyDatas(context.TODO(), metav1.ListOptions{})
	for _, i := range vAdditionalPolicyDataListObj {
		vAdditionalPolicyData, err := nc.GetPolicyAdditionalPolicyData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting AdditionalPolicyData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vPolicyAdditionalPolicyData = vAdditionalPolicyData
		dn := vAdditionalPolicyData.DisplayName()
		parentLabels := map[string]interface{}{"additionalpolicydatas.policypkg.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vAdditionalPolicyData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.PolicyAdditionalPolicyData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
		}
		vAdditionalPolicyDataList = append(vAdditionalPolicyDataList, ret)
	}

	log.Debugf("[getRootResolver]Output AdditionalPolicyData objects %v", vAdditionalPolicyDataList)
	return vAdditionalPolicyDataList, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: AdditionalPolicyData in PKG: Policy
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryServiceTableResolver(obj *model.PolicyAdditionalPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryServiceVersionTableResolver(obj *model.PolicyAdditionalPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryServiceTSResolver(obj *model.PolicyAdditionalPolicyData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryIncomingAPIsResolver(obj *model.PolicyAdditionalPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryOutgoingAPIsResolver(obj *model.PolicyAdditionalPolicyData, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryIncomingTCPResolver(obj *model.PolicyAdditionalPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryOutgoingTCPResolver(obj *model.PolicyAdditionalPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicyAdditionalPolicyDataqueryServiceTopologyResolver(obj *model.PolicyAdditionalPolicyData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAdditionalPolicyDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: AccessControlPolicy in PKG: Policy
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicyAccessControlPolicyqueryServiceTableResolver(obj *model.PolicyAccessControlPolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicyAccessControlPolicyqueryServiceVersionTableResolver(obj *model.PolicyAccessControlPolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicyAccessControlPolicyqueryServiceTSResolver(obj *model.PolicyAccessControlPolicy, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicyAccessControlPolicyqueryIncomingAPIsResolver(obj *model.PolicyAccessControlPolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicyAccessControlPolicyqueryOutgoingAPIsResolver(obj *model.PolicyAccessControlPolicy, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicyAccessControlPolicyqueryIncomingTCPResolver(obj *model.PolicyAccessControlPolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicyAccessControlPolicyqueryOutgoingTCPResolver(obj *model.PolicyAccessControlPolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicyAccessControlPolicyqueryServiceTopologyResolver(obj *model.PolicyAccessControlPolicy, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyAccessControlPolicyqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: ACPConfig in PKG: Policy
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicyACPConfigqueryServiceTableResolver(obj *model.PolicyACPConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicyACPConfigqueryServiceVersionTableResolver(obj *model.PolicyACPConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicyACPConfigqueryServiceTSResolver(obj *model.PolicyACPConfig, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicyACPConfigqueryIncomingAPIsResolver(obj *model.PolicyACPConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicyACPConfigqueryOutgoingAPIsResolver(obj *model.PolicyACPConfig, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicyACPConfigqueryIncomingTCPResolver(obj *model.PolicyACPConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicyACPConfigqueryOutgoingTCPResolver(obj *model.PolicyACPConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicyACPConfigqueryServiceTopologyResolver(obj *model.PolicyACPConfig, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyACPConfigqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: VMpolicy in PKG: Policy
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicyVMpolicyqueryServiceTableResolver(obj *model.PolicyVMpolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicyVMpolicyqueryServiceVersionTableResolver(obj *model.PolicyVMpolicy, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicyVMpolicyqueryServiceTSResolver(obj *model.PolicyVMpolicy, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicyVMpolicyqueryIncomingAPIsResolver(obj *model.PolicyVMpolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicyVMpolicyqueryOutgoingAPIsResolver(obj *model.PolicyVMpolicy, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicyVMpolicyqueryIncomingTCPResolver(obj *model.PolicyVMpolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicyVMpolicyqueryOutgoingTCPResolver(obj *model.PolicyVMpolicy, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicyVMpolicyqueryServiceTopologyResolver(obj *model.PolicyVMpolicy, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyVMpolicyqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// Non Singleton Resolver for Parent Node
// PKG: Policy, NODE: Policy
//////////////////////////////////////
func (c *resolverConfig) getRootResolver(id *string) ([]*model.PolicyRandomPolicyData, error) {
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	var vRandomPolicyDataList []*model.PolicyRandomPolicyData
	if id != nil && *id != "" {
		log.Debugf("[getRootResolver]Id: %q", *id)
		vRandomPolicyData, err := nc.GetPolicyRandomPolicyData(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting RandomPolicyData node %q: %s", *id, err)
			return nil, nil
		}
		c.vPolicyRandomPolicyData = vRandomPolicyData
		dn := vRandomPolicyData.DisplayName()
		parentLabels := map[string]interface{}{"randompolicydatas.policypkg.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vRandomPolicyData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.PolicyRandomPolicyData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
		}
		vRandomPolicyDataList = append(vRandomPolicyDataList, ret)
		log.Debugf("[getRootResolver]Output RandomPolicyData objects %+v", vRandomPolicyDataList)
		return vRandomPolicyDataList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all RandomPolicyDatas")

	vRandomPolicyDataListObj, err := nc.RandomPolicyData().ListRandomPolicyDatas(context.TODO(), metav1.ListOptions{})
	for _, i := range vRandomPolicyDataListObj {
		vRandomPolicyData, err := nc.GetPolicyRandomPolicyData(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting RandomPolicyData node %q : %s", i.DisplayName(), err)
			continue
		}
		c.vPolicyRandomPolicyData = vRandomPolicyData
		dn := vRandomPolicyData.DisplayName()
		parentLabels := map[string]interface{}{"randompolicydatas.policypkg.tsm.tanzu.vmware.com": dn}
		Description, _ := json.Marshal(vRandomPolicyData.Spec.Description)
		DescriptionData := string(Description)

		ret := &model.PolicyRandomPolicyData{
			Id:           &dn,
			ParentLabels: parentLabels,
			Description:  &DescriptionData,
		}
		vRandomPolicyDataList = append(vRandomPolicyDataList, ret)
	}

	log.Debugf("[getRootResolver]Output RandomPolicyData objects %v", vRandomPolicyDataList)
	return vRandomPolicyDataList, nil
}

//////////////////////////////////////
// CustomQuery Resolver for Node: RandomPolicyData in PKG: Policy
//////////////////////////////////////

// Resolver for queryServiceTable
func (c *resolverConfig) getPolicyRandomPolicyDataqueryServiceTableResolver(obj *model.PolicyRandomPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryServiceTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceVersionTable
func (c *resolverConfig) getPolicyRandomPolicyDataqueryServiceVersionTableResolver(obj *model.PolicyRandomPolicyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTable", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryServiceVersionTableResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTS
func (c *resolverConfig) getPolicyRandomPolicyDataqueryServiceTSResolver(obj *model.PolicyRandomPolicyData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceMetricSeries", StartTime: *startTime, EndTime: *endTime, Metric: *svcMetric, Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryServiceTSResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingAPIs
func (c *resolverConfig) getPolicyRandomPolicyDataqueryIncomingAPIsResolver(obj *model.PolicyRandomPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryIncomingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingAPIs
func (c *resolverConfig) getPolicyRandomPolicyDataqueryOutgoingAPIsResolver(obj *model.PolicyRandomPolicyData, startTime *string, endTime, destinationService, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingAPIs", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: *timeInterval})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryOutgoingAPIsResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryIncomingTCP
func (c *resolverConfig) getPolicyRandomPolicyDataqueryIncomingTCPResolver(obj *model.PolicyRandomPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/IncomingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryIncomingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryOutgoingTCP
func (c *resolverConfig) getPolicyRandomPolicyDataqueryOutgoingTCPResolver(obj *model.PolicyRandomPolicyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/OutgoingTCP", StartTime: *startTime, EndTime: *endTime, Metric: "", Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryOutgoingTCPResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

// Resolver for queryServiceTopology
func (c *resolverConfig) getPolicyRandomPolicyDataqueryServiceTopologyResolver(obj *model.PolicyRandomPolicyData, startTime *string, endTime *string, metricStringArray *string) (*model.TimeSeriesData, error) {
	ctx := context.Background()
	var filters = make(map[string]string)
	filters[""] = ""
	resp, err := grpcServer().GetMetrics(ctx, &qm.MetricArg{QueryType: "/ServiceTopology", StartTime: *startTime, EndTime: *endTime, Metric: *metricStringArray, Filters: filters, TimeInterval: ""})
	if err != nil {
		log.Errorf("[getPolicyRandomPolicyDataqueryServiceTopologyResolver]Failed to get metrics, err: %v", err)
		return &model.TimeSeriesData{}, nil
	}
	b, _ := json.Marshal(resp.Data)
	data := string(b)
	ret := &model.TimeSeriesData{
		Data: &data,
	}
	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: Config Node: Root PKG: Root
//////////////////////////////////////
func (c *resolverConfig) getRootRootConfigResolver(obj *model.RootRoot, id *string) (*model.ConfigConfig, error) {
	log.Debugf("[getRootRootConfigResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getRootRootConfigResolver]Id %q", *id)
		vConfig, err := nc.RootRoot().GetConfig(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootRootConfigResolver]Error getting Config node %q : %s", *id, err)
			return &model.ConfigConfig{}, nil
		}
		c.vConfigConfig = vConfig
		dn := vConfig.DisplayName()
		parentLabels := map[string]interface{}{"configs.config.tsm.tanzu.vmware.com": dn}
		MyStr0, _ := json.Marshal(vConfig.Spec.MyStr0)
		MyStr0Data := string(MyStr0)
		MyStr1, _ := json.Marshal(vConfig.Spec.MyStr1)
		MyStr1Data := string(MyStr1)
		MyStr2, _ := json.Marshal(vConfig.Spec.MyStr2)
		MyStr2Data := string(MyStr2)
		XYZPort, _ := json.Marshal(vConfig.Spec.XYZPort)
		XYZPortData := string(XYZPort)
		ABCHost, _ := json.Marshal(vConfig.Spec.ABCHost)
		ABCHostData := string(ABCHost)
		ClusterNamespaces, _ := json.Marshal(vConfig.Spec.ClusterNamespaces)
		ClusterNamespacesData := string(ClusterNamespaces)
		TestValMarkers, _ := json.Marshal(vConfig.Spec.TestValMarkers)
		TestValMarkersData := string(TestValMarkers)
		vInstance := float64(vConfig.Spec.Instance)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConfigConfig{
			Id:                &dn,
			ParentLabels:      parentLabels,
			MyStr0:            &MyStr0Data,
			MyStr1:            &MyStr1Data,
			MyStr2:            &MyStr2Data,
			XYZPort:           &XYZPortData,
			ABCHost:           &ABCHostData,
			ClusterNamespaces: &ClusterNamespacesData,
			TestValMarkers:    &TestValMarkersData,
			Instance:          &vInstance,
		}

		log.Debugf("[getRootRootConfigResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getRootRootConfigResolver]Id is empty, process all Configs")
	vConfigParent, err := nc.GetRootRoot(context.TODO())
	if err != nil {
		log.Errorf("[getRootRootConfigResolver]Failed to get parent node %s", err)
		return &model.ConfigConfig{}, nil
	}
	vConfig, err := vConfigParent.GetConfig(context.TODO())
	if err != nil {
		log.Errorf("[getRootRootConfigResolver]Error getting Config node %s", err)
		return &model.ConfigConfig{}, nil
	}
	c.vConfigConfig = vConfig
	dn := vConfig.DisplayName()
	parentLabels := map[string]interface{}{"configs.config.tsm.tanzu.vmware.com": dn}
	MyStr0, _ := json.Marshal(vConfig.Spec.MyStr0)
	MyStr0Data := string(MyStr0)
	MyStr1, _ := json.Marshal(vConfig.Spec.MyStr1)
	MyStr1Data := string(MyStr1)
	MyStr2, _ := json.Marshal(vConfig.Spec.MyStr2)
	MyStr2Data := string(MyStr2)
	XYZPort, _ := json.Marshal(vConfig.Spec.XYZPort)
	XYZPortData := string(XYZPort)
	ABCHost, _ := json.Marshal(vConfig.Spec.ABCHost)
	ABCHostData := string(ABCHost)
	ClusterNamespaces, _ := json.Marshal(vConfig.Spec.ClusterNamespaces)
	ClusterNamespacesData := string(ClusterNamespaces)
	TestValMarkers, _ := json.Marshal(vConfig.Spec.TestValMarkers)
	TestValMarkersData := string(TestValMarkers)
	vInstance := float64(vConfig.Spec.Instance)

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.ConfigConfig{
		Id:                &dn,
		ParentLabels:      parentLabels,
		MyStr0:            &MyStr0Data,
		MyStr1:            &MyStr1Data,
		MyStr2:            &MyStr2Data,
		XYZPort:           &XYZPortData,
		ABCHost:           &ABCHostData,
		ClusterNamespaces: &ClusterNamespacesData,
		TestValMarkers:    &TestValMarkersData,
		Instance:          &vInstance,
	}

	log.Debugf("[getRootRootConfigResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: GNS Node: Config PKG: Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigGNSResolver(obj *model.ConfigConfig, id *string) (*model.GnsGns, error) {
	log.Debugf("[getConfigConfigGNSResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigGNSResolver]Id %q", *id)
		vGns, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetGNS(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigGNSResolver]Error getting GNS node %q : %s", *id, err)
			return &model.GnsGns{}, nil
		}
		c.vGnsGns = vGns
		dn := vGns.DisplayName()
		parentLabels := map[string]interface{}{"gnses.gns.tsm.tanzu.vmware.com": dn}
		vDomain := string(vGns.Spec.Domain)
		vUseSharedGateway := bool(vGns.Spec.UseSharedGateway)
		Description, _ := json.Marshal(vGns.Spec.Description)
		DescriptionData := string(Description)
		vMeta := string(vGns.Spec.Meta)
		OtherDescription, _ := json.Marshal(vGns.Spec.OtherDescription)
		OtherDescriptionData := string(OtherDescription)
		MapPointer, _ := json.Marshal(vGns.Spec.MapPointer)
		MapPointerData := string(MapPointer)
		SlicePointer, _ := json.Marshal(vGns.Spec.SlicePointer)
		SlicePointerData := string(SlicePointer)
		WorkloadSpec, _ := json.Marshal(vGns.Spec.WorkloadSpec)
		WorkloadSpecData := string(WorkloadSpec)
		DifferentSpec, _ := json.Marshal(vGns.Spec.DifferentSpec)
		DifferentSpecData := string(DifferentSpec)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.GnsGns{
			Id:               &dn,
			ParentLabels:     parentLabels,
			Domain:           &vDomain,
			UseSharedGateway: &vUseSharedGateway,
			Description:      &DescriptionData,
			Meta:             &vMeta,
			OtherDescription: &OtherDescriptionData,
			MapPointer:       &MapPointerData,
			SlicePointer:     &SlicePointerData,
			WorkloadSpec:     &WorkloadSpecData,
			DifferentSpec:    &DifferentSpecData,
		}

		log.Debugf("[getConfigConfigGNSResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getConfigConfigGNSResolver]Id is empty, process all GNSs")
	vGnsParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigGNSResolver]Failed to get parent node %s", err)
		return &model.GnsGns{}, nil
	}
	vGns, err := vGnsParent.GetGNS(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigGNSResolver]Error getting GNS node %s", err)
		return &model.GnsGns{}, nil
	}
	c.vGnsGns = vGns
	dn := vGns.DisplayName()
	parentLabels := map[string]interface{}{"gnses.gns.tsm.tanzu.vmware.com": dn}
	vDomain := string(vGns.Spec.Domain)
	vUseSharedGateway := bool(vGns.Spec.UseSharedGateway)
	Description, _ := json.Marshal(vGns.Spec.Description)
	DescriptionData := string(Description)
	vMeta := string(vGns.Spec.Meta)
	OtherDescription, _ := json.Marshal(vGns.Spec.OtherDescription)
	OtherDescriptionData := string(OtherDescription)
	MapPointer, _ := json.Marshal(vGns.Spec.MapPointer)
	MapPointerData := string(MapPointer)
	SlicePointer, _ := json.Marshal(vGns.Spec.SlicePointer)
	SlicePointerData := string(SlicePointer)
	WorkloadSpec, _ := json.Marshal(vGns.Spec.WorkloadSpec)
	WorkloadSpecData := string(WorkloadSpec)
	DifferentSpec, _ := json.Marshal(vGns.Spec.DifferentSpec)
	DifferentSpecData := string(DifferentSpec)

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.GnsGns{
		Id:               &dn,
		ParentLabels:     parentLabels,
		Domain:           &vDomain,
		UseSharedGateway: &vUseSharedGateway,
		Description:      &DescriptionData,
		Meta:             &vMeta,
		OtherDescription: &OtherDescriptionData,
		MapPointer:       &MapPointerData,
		SlicePointer:     &SlicePointerData,
		WorkloadSpec:     &WorkloadSpecData,
		DifferentSpec:    &DifferentSpecData,
	}

	log.Debugf("[getConfigConfigGNSResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Singleton)
// FieldName: DNS Node: Config PKG: Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigDNSResolver(obj *model.ConfigConfig) (*model.GnsDns, error) {
	log.Debugf("[getConfigConfigDNSResolver]Parent Object %+v", obj)
	vDns, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetDNS(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigDNSResolver]Error getting Config node %s", err)
		return &model.GnsDns{}, nil
	}
	c.vGnsDns = vDns
	dn := vDns.DisplayName()
	parentLabels := map[string]interface{}{"dnses.gns.tsm.tanzu.vmware.com": dn}

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.GnsDns{
		Id:           &dn,
		ParentLabels: parentLabels,
	}

	log.Debugf("[getConfigConfigDNSResolver]Output object %+v", ret)
	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: VMPPolicies Node: Config PKG: Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigVMPPoliciesResolver(obj *model.ConfigConfig, id *string) (*model.PolicyVMpolicy, error) {
	log.Debugf("[getConfigConfigVMPPoliciesResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigVMPPoliciesResolver]Id %q", *id)
		vVMpolicy, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetVMPPolicies(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigVMPPoliciesResolver]Error getting VMPPolicies node %q : %s", *id, err)
			return &model.PolicyVMpolicy{}, nil
		}
		c.vPolicyVMpolicy = vVMpolicy
		dn := vVMpolicy.DisplayName()
		parentLabels := map[string]interface{}{"vmpolicies.policypkg.tsm.tanzu.vmware.com": dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.PolicyVMpolicy{
			Id:           &dn,
			ParentLabels: parentLabels,
		}

		log.Debugf("[getConfigConfigVMPPoliciesResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getConfigConfigVMPPoliciesResolver]Id is empty, process all VMPPoliciess")
	vVMpolicyParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigVMPPoliciesResolver]Failed to get parent node %s", err)
		return &model.PolicyVMpolicy{}, nil
	}
	vVMpolicy, err := vVMpolicyParent.GetVMPPolicies(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigVMPPoliciesResolver]Error getting VMPPolicies node %s", err)
		return &model.PolicyVMpolicy{}, nil
	}
	c.vPolicyVMpolicy = vVMpolicy
	dn := vVMpolicy.DisplayName()
	parentLabels := map[string]interface{}{"vmpolicies.policypkg.tsm.tanzu.vmware.com": dn}

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.PolicyVMpolicy{
		Id:           &dn,
		ParentLabels: parentLabels,
	}

	log.Debugf("[getConfigConfigVMPPoliciesResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: Domain Node: Config PKG: Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigDomainResolver(obj *model.ConfigConfig, id *string) (*model.ConfigDomain, error) {
	log.Debugf("[getConfigConfigDomainResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigDomainResolver]Id %q", *id)
		vDomain, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetDomain(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigDomainResolver]Error getting Domain node %q : %s", *id, err)
			return &model.ConfigDomain{}, nil
		}
		c.vConfigDomain = vDomain
		dn := vDomain.DisplayName()
		parentLabels := map[string]interface{}{"domains.config.tsm.tanzu.vmware.com": dn}
		PointPort, _ := json.Marshal(vDomain.Spec.PointPort)
		PointPortData := string(PointPort)
		PointMap, _ := json.Marshal(vDomain.Spec.PointMap)
		PointMapData := string(PointMap)
		PointSlice, _ := json.Marshal(vDomain.Spec.PointSlice)
		PointSliceData := string(PointSlice)
		SliceOfPoints, _ := json.Marshal(vDomain.Spec.SliceOfPoints)
		SliceOfPointsData := string(SliceOfPoints)
		SliceOfArrPoints, _ := json.Marshal(vDomain.Spec.SliceOfArrPoints)
		SliceOfArrPointsData := string(SliceOfArrPoints)
		MapOfArrsPoints, _ := json.Marshal(vDomain.Spec.MapOfArrsPoints)
		MapOfArrsPointsData := string(MapOfArrsPoints)
		PointStruct, _ := json.Marshal(vDomain.Spec.PointStruct)
		PointStructData := string(PointStruct)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConfigDomain{
			Id:               &dn,
			ParentLabels:     parentLabels,
			PointPort:        &PointPortData,
			PointMap:         &PointMapData,
			PointSlice:       &PointSliceData,
			SliceOfPoints:    &SliceOfPointsData,
			SliceOfArrPoints: &SliceOfArrPointsData,
			MapOfArrsPoints:  &MapOfArrsPointsData,
			PointStruct:      &PointStructData,
		}

		log.Debugf("[getConfigConfigDomainResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getConfigConfigDomainResolver]Id is empty, process all Domains")
	vDomainParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigDomainResolver]Failed to get parent node %s", err)
		return &model.ConfigDomain{}, nil
	}
	vDomain, err := vDomainParent.GetDomain(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigDomainResolver]Error getting Domain node %s", err)
		return &model.ConfigDomain{}, nil
	}
	c.vConfigDomain = vDomain
	dn := vDomain.DisplayName()
	parentLabels := map[string]interface{}{"domains.config.tsm.tanzu.vmware.com": dn}
	PointPort, _ := json.Marshal(vDomain.Spec.PointPort)
	PointPortData := string(PointPort)
	PointMap, _ := json.Marshal(vDomain.Spec.PointMap)
	PointMapData := string(PointMap)
	PointSlice, _ := json.Marshal(vDomain.Spec.PointSlice)
	PointSliceData := string(PointSlice)
	SliceOfPoints, _ := json.Marshal(vDomain.Spec.SliceOfPoints)
	SliceOfPointsData := string(SliceOfPoints)
	SliceOfArrPoints, _ := json.Marshal(vDomain.Spec.SliceOfArrPoints)
	SliceOfArrPointsData := string(SliceOfArrPoints)
	MapOfArrsPoints, _ := json.Marshal(vDomain.Spec.MapOfArrsPoints)
	MapOfArrsPointsData := string(MapOfArrsPoints)
	PointStruct, _ := json.Marshal(vDomain.Spec.PointStruct)
	PointStructData := string(PointStruct)

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.ConfigDomain{
		Id:               &dn,
		ParentLabels:     parentLabels,
		PointPort:        &PointPortData,
		PointMap:         &PointMapData,
		PointSlice:       &PointSliceData,
		SliceOfPoints:    &SliceOfPointsData,
		SliceOfArrPoints: &SliceOfArrPointsData,
		MapOfArrsPoints:  &MapOfArrsPointsData,
		PointStruct:      &PointStructData,
	}

	log.Debugf("[getConfigConfigDomainResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: FooExample Node: Config PKG: Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigFooExampleResolver(obj *model.ConfigConfig, id *string) ([]*model.ConfigFooType, error) {
	log.Debugf("[getConfigConfigFooExampleResolver]Parent Object %+v", obj)
	var vConfigFooTypeList []*model.ConfigFooType
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigFooExampleResolver]Id %q", *id)
		vFooType, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetFooExample(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigFooExampleResolver]Error getting FooExample node %q : %s", *id, err)
			return vConfigFooTypeList, nil
		}
		dn := vFooType.DisplayName()
		parentLabels := map[string]interface{}{"footypes.config.tsm.tanzu.vmware.com": dn}
		FooA, _ := json.Marshal(vFooType.Spec.FooA)
		FooAData := string(FooA)
		FooB, _ := json.Marshal(vFooType.Spec.FooB)
		FooBData := string(FooB)
		FooD, _ := json.Marshal(vFooType.Spec.FooD)
		FooDData := string(FooD)
		FooF, _ := json.Marshal(vFooType.Spec.FooF)
		FooFData := string(FooF)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConfigFooType{
			Id:           &dn,
			ParentLabels: parentLabels,
			FooA:         &FooAData,
			FooB:         &FooBData,
			FooD:         &FooDData,
			FooF:         &FooFData,
		}
		vConfigFooTypeList = append(vConfigFooTypeList, ret)

		log.Debugf("[getConfigConfigFooExampleResolver]Output FooExample objects %v", vConfigFooTypeList)

		return vConfigFooTypeList, nil
	}

	log.Debug("[getConfigConfigFooExampleResolver]Id is empty, process all FooExamples")

	vFooTypeParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigFooExampleResolver]Error getting parent node %s", err)
		return vConfigFooTypeList, nil
	}
	vFooTypeAllObj, err := vFooTypeParent.GetAllFooExample(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigFooExampleResolver]Error getting FooExample objects %s", err)
		return vConfigFooTypeList, nil
	}
	for _, i := range vFooTypeAllObj {
		vFooType, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetFooExample(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConfigConfigFooExampleResolver]Error getting FooExample node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vFooType.DisplayName()
		parentLabels := map[string]interface{}{"footypes.config.tsm.tanzu.vmware.com": dn}
		FooA, _ := json.Marshal(vFooType.Spec.FooA)
		FooAData := string(FooA)
		FooB, _ := json.Marshal(vFooType.Spec.FooB)
		FooBData := string(FooB)
		FooD, _ := json.Marshal(vFooType.Spec.FooD)
		FooDData := string(FooD)
		FooF, _ := json.Marshal(vFooType.Spec.FooF)
		FooFData := string(FooF)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConfigFooType{
			Id:           &dn,
			ParentLabels: parentLabels,
			FooA:         &FooAData,
			FooB:         &FooBData,
			FooD:         &FooDData,
			FooF:         &FooFData,
		}
		vConfigFooTypeList = append(vConfigFooTypeList, ret)
	}

	log.Debugf("[getConfigConfigFooExampleResolver]Output FooExample objects %v", vConfigFooTypeList)

	return vConfigFooTypeList, nil
}

//////////////////////////////////////
// LINKS RESOLVER
// FieldName: ACPPolicies Node: Config PKG: Config
//////////////////////////////////////
func (c *resolverConfig) getConfigConfigACPPoliciesResolver(obj *model.ConfigConfig, id *string) ([]*model.PolicyAccessControlPolicy, error) {
	log.Debugf("[getConfigConfigACPPoliciesResolver]Parent Object %+v", obj)
	var vPolicyAccessControlPolicyList []*model.PolicyAccessControlPolicy
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigACPPoliciesResolver]Id %q", *id)
		vAccessControlPolicyParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
		if err != nil {
			log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting ACPPolicies %q : %s", *id, err)
			return vPolicyAccessControlPolicyList, nil
		}
		vAccessControlPolicy, err := vAccessControlPolicyParent.GetACPPolicies(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting ACPPolicies %q : %s", *id, err)
			return vPolicyAccessControlPolicyList, nil
		}
		dn := vAccessControlPolicy.DisplayName()
		parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com": dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.PolicyAccessControlPolicy{
			Id:           &dn,
			ParentLabels: parentLabels,
		}
		vPolicyAccessControlPolicyList = append(vPolicyAccessControlPolicyList, ret)

		log.Debugf("[getConfigConfigACPPoliciesResolver]Output ACPPolicies objects %v", vPolicyAccessControlPolicyList)

		return vPolicyAccessControlPolicyList, nil
	}

	log.Debug("[getConfigConfigACPPoliciesResolver]Id is empty, process all ACPPoliciess")

	vAccessControlPolicyParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting parent node %s", err)
		return vPolicyAccessControlPolicyList, nil
	}
	vAccessControlPolicyAllObj, err := vAccessControlPolicyParent.GetAllACPPolicies(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting ACPPolicies %s", err)
		return vPolicyAccessControlPolicyList, nil
	}
	for _, i := range vAccessControlPolicyAllObj {
		vAccessControlPolicyParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
		if err != nil {
			log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting parent node %s, skipping...", err)
			continue
		}
		vAccessControlPolicy, err := vAccessControlPolicyParent.GetACPPolicies(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting ACPPolicies node %q : %s, skipping...", i.DisplayName(), err)
			continue
		}
		dn := vAccessControlPolicy.DisplayName()
		parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com": dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.PolicyAccessControlPolicy{
			Id:           &dn,
			ParentLabels: parentLabels,
		}
		vPolicyAccessControlPolicyList = append(vPolicyAccessControlPolicyList, ret)
	}
	log.Debugf("[getConfigConfigACPPoliciesResolver]List of ACPPolicies object %v", vPolicyAccessControlPolicyList)
	return vPolicyAccessControlPolicyList, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: GnsAccessControlPolicy Node: Gns PKG: Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsGnsAccessControlPolicyResolver(obj *model.GnsGns, id *string) (*model.PolicyAccessControlPolicy, error) {
	log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Id %q", *id)
		vAccessControlPolicy, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GetGnsAccessControlPolicy(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getGnsGnsGnsAccessControlPolicyResolver]Error getting GnsAccessControlPolicy node %q : %s", *id, err)
			return &model.PolicyAccessControlPolicy{}, nil
		}
		c.vPolicyAccessControlPolicy = vAccessControlPolicy
		dn := vAccessControlPolicy.DisplayName()
		parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com": dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.PolicyAccessControlPolicy{
			Id:           &dn,
			ParentLabels: parentLabels,
		}

		log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getGnsGnsGnsAccessControlPolicyResolver]Id is empty, process all GnsAccessControlPolicys")
	vAccessControlPolicyParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetGNS(context.TODO(), getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getGnsGnsGnsAccessControlPolicyResolver]Failed to get parent node %s", err)
		return &model.PolicyAccessControlPolicy{}, nil
	}
	vAccessControlPolicy, err := vAccessControlPolicyParent.GetGnsAccessControlPolicy(context.TODO())
	if err != nil {
		log.Errorf("[getGnsGnsGnsAccessControlPolicyResolver]Error getting GnsAccessControlPolicy node %s", err)
		return &model.PolicyAccessControlPolicy{}, nil
	}
	c.vPolicyAccessControlPolicy = vAccessControlPolicy
	dn := vAccessControlPolicy.DisplayName()
	parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com": dn}

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.PolicyAccessControlPolicy{
		Id:           &dn,
		ParentLabels: parentLabels,
	}

	log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: GnsServiceGroups Node: Gns PKG: Gns
//////////////////////////////////////
func (c *resolverConfig) getGnsGnsGnsServiceGroupsResolver(obj *model.GnsGns, id *string) ([]*model.ServicegroupSvcGroup, error) {
	log.Debugf("[getGnsGnsGnsServiceGroupsResolver]Parent Object %+v", obj)
	var vServicegroupSvcGroupList []*model.ServicegroupSvcGroup
	if id != nil && *id != "" {
		log.Debugf("[getGnsGnsGnsServiceGroupsResolver]Id %q", *id)
		vSvcGroup, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GetGnsServiceGroups(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getGnsGnsGnsServiceGroupsResolver]Error getting GnsServiceGroups node %q : %s", *id, err)
			return vServicegroupSvcGroupList, nil
		}
		dn := vSvcGroup.DisplayName()
		parentLabels := map[string]interface{}{"svcgroups.servicegroup.tsm.tanzu.vmware.com": dn}
		vDisplayName := string(vSvcGroup.Spec.DisplayName)
		vDescription := string(vSvcGroup.Spec.Description)
		vColor := string(vSvcGroup.Spec.Color)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ServicegroupSvcGroup{
			Id:           &dn,
			ParentLabels: parentLabels,
			DisplayName:  &vDisplayName,
			Description:  &vDescription,
			Color:        &vColor,
		}
		vServicegroupSvcGroupList = append(vServicegroupSvcGroupList, ret)

		log.Debugf("[getGnsGnsGnsServiceGroupsResolver]Output GnsServiceGroups objects %v", vServicegroupSvcGroupList)

		return vServicegroupSvcGroupList, nil
	}

	log.Debug("[getGnsGnsGnsServiceGroupsResolver]Id is empty, process all GnsServiceGroupss")

	vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetGNS(context.TODO(), getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getGnsGnsGnsServiceGroupsResolver]Error getting parent node %s", err)
		return vServicegroupSvcGroupList, nil
	}
	vSvcGroupAllObj, err := vSvcGroupParent.GetAllGnsServiceGroups(context.TODO())
	if err != nil {
		log.Errorf("[getGnsGnsGnsServiceGroupsResolver]Error getting GnsServiceGroups objects %s", err)
		return vServicegroupSvcGroupList, nil
	}
	for _, i := range vSvcGroupAllObj {
		vSvcGroup, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GetGnsServiceGroups(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getGnsGnsGnsServiceGroupsResolver]Error getting GnsServiceGroups node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vSvcGroup.DisplayName()
		parentLabels := map[string]interface{}{"svcgroups.servicegroup.tsm.tanzu.vmware.com": dn}
		vDisplayName := string(vSvcGroup.Spec.DisplayName)
		vDescription := string(vSvcGroup.Spec.Description)
		vColor := string(vSvcGroup.Spec.Color)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ServicegroupSvcGroup{
			Id:           &dn,
			ParentLabels: parentLabels,
			DisplayName:  &vDisplayName,
			Description:  &vDescription,
			Color:        &vColor,
		}
		vServicegroupSvcGroupList = append(vServicegroupSvcGroupList, ret)
	}

	log.Debugf("[getGnsGnsGnsServiceGroupsResolver]Output GnsServiceGroups objects %v", vServicegroupSvcGroupList)

	return vServicegroupSvcGroupList, nil
}

//////////////////////////////////////
// LINK RESOLVER
// FieldName: DestSvcGroups Node: ACPConfig PKG: Policy
//////////////////////////////////////
func (c *resolverConfig) getPolicyACPConfigDestSvcGroupsResolver(obj *model.PolicyACPConfig) (*model.PolicypkgSvcGroup, error) {
	log.Debugf("[getPolicyACPConfigDestSvcGroupsResolver]Parent Object %+v", obj)
	vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), getParentName(obj.ParentLabels, "acpconfigs.policypkg.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting parent node %s", err)
		return &model.PolicypkgSvcGroup{}, nil
	}
	vSvcGroup, err := vSvcGroupParent.GetDestSvcGroups(context.TODO())
	if err != nil {
		log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting DestSvcGroups object %s", err)
		return &model.PolicypkgSvcGroup{}, nil
	}
	c.vPolicypkgSvcGroup = vSvcGroup

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}

	log.Debugf("[getPolicyACPConfigDestSvcGroupsResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// LINKS RESOLVER
// FieldName: DestSvcGroups Node: ACPConfig PKG: Policy
//////////////////////////////////////
func (c *resolverConfig) getPolicyACPConfigDestSvcGroupsResolver(obj *model.PolicyACPConfig, id *string) ([]*model.PolicypkgSvcGroup, error) {
	log.Debugf("[getPolicyACPConfigDestSvcGroupsResolver]Parent Object %+v", obj)
	var vPolicypkgSvcGroupList []*model.PolicypkgSvcGroup
	if id != nil && *id != "" {
		log.Debugf("[getPolicyACPConfigDestSvcGroupsResolver]Id %q", *id)
		vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), getParentName(obj.ParentLabels, "acpconfigs.policypkg.tsm.tanzu.vmware.com"))
		if err != nil {
			log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting DestSvcGroups %q : %s", *id, err)
			return vPolicypkgSvcGroupList, nil
		}
		vSvcGroup, err := vSvcGroupParent.GetDestSvcGroups(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting DestSvcGroups %q : %s", *id, err)
			return vPolicypkgSvcGroupList, nil
		}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}

		vPolicypkgSvcGroupList = append(vPolicypkgSvcGroupList, ret)

		log.Debugf("[getPolicyACPConfigDestSvcGroupsResolver]Output DestSvcGroups objects %v", vPolicypkgSvcGroupList)

		return vPolicypkgSvcGroupList, nil
	}

	log.Debug("[getPolicyACPConfigDestSvcGroupsResolver]Id is empty, process all DestSvcGroupss")

	vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), getParentName(obj.ParentLabels, "acpconfigs.policypkg.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting parent node %s", err)
		return vPolicypkgSvcGroupList, nil
	}
	vSvcGroupAllObj, err := vSvcGroupParent.GetAllDestSvcGroups(context.TODO())
	if err != nil {
		log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting DestSvcGroups %s", err)
		return vPolicypkgSvcGroupList, nil
	}
	for _, i := range vSvcGroupAllObj {
		vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), getParentName(obj.ParentLabels, "acpconfigs.policypkg.tsm.tanzu.vmware.com"))
		if err != nil {
			log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting parent node %s, skipping...", err)
			continue
		}
		vSvcGroup, err := vSvcGroupParent.GetDestSvcGroups(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getPolicyACPConfigDestSvcGroupsResolver]Error getting DestSvcGroups node %q : %s, skipping...", i.DisplayName(), err)
			continue
		}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}

		vPolicypkgSvcGroupList = append(vPolicypkgSvcGroupList, ret)
	}
	log.Debugf("[getPolicyACPConfigDestSvcGroupsResolver]List of DestSvcGroups object %v", vPolicypkgSvcGroupList)
	return vPolicypkgSvcGroupList, nil
}

//////////////////////////////////////
// LINKS RESOLVER
// FieldName: SourceSvcGroups Node: ACPConfig PKG: Policy
//////////////////////////////////////
func (c *resolverConfig) getPolicyACPConfigSourceSvcGroupsResolver(obj *model.PolicyACPConfig, id *string) ([]*model.ServicegroupSvcGroup, error) {
	log.Debugf("[getPolicyACPConfigSourceSvcGroupsResolver]Parent Object %+v", obj)
	var vServicegroupSvcGroupList []*model.ServicegroupSvcGroup
	if id != nil && *id != "" {
		log.Debugf("[getPolicyACPConfigSourceSvcGroupsResolver]Id %q", *id)
		vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), getParentName(obj.ParentLabels, "acpconfigs.policypkg.tsm.tanzu.vmware.com"))
		if err != nil {
			log.Errorf("[getPolicyACPConfigSourceSvcGroupsResolver]Error getting SourceSvcGroups %q : %s", *id, err)
			return vServicegroupSvcGroupList, nil
		}
		vSvcGroup, err := vSvcGroupParent.GetSourceSvcGroups(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getPolicyACPConfigSourceSvcGroupsResolver]Error getting SourceSvcGroups %q : %s", *id, err)
			return vServicegroupSvcGroupList, nil
		}
		dn := vSvcGroup.DisplayName()
		parentLabels := map[string]interface{}{"svcgroups.servicegroup.tsm.tanzu.vmware.com": dn}
		vDisplayName := string(vSvcGroup.Spec.DisplayName)
		vDescription := string(vSvcGroup.Spec.Description)
		vColor := string(vSvcGroup.Spec.Color)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ServicegroupSvcGroup{
			Id:           &dn,
			ParentLabels: parentLabels,
			DisplayName:  &vDisplayName,
			Description:  &vDescription,
			Color:        &vColor,
		}
		vServicegroupSvcGroupList = append(vServicegroupSvcGroupList, ret)

		log.Debugf("[getPolicyACPConfigSourceSvcGroupsResolver]Output SourceSvcGroups objects %v", vServicegroupSvcGroupList)

		return vServicegroupSvcGroupList, nil
	}

	log.Debug("[getPolicyACPConfigSourceSvcGroupsResolver]Id is empty, process all SourceSvcGroupss")

	vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), getParentName(obj.ParentLabels, "acpconfigs.policypkg.tsm.tanzu.vmware.com"))
	if err != nil {
		log.Errorf("[getPolicyACPConfigSourceSvcGroupsResolver]Error getting parent node %s", err)
		return vServicegroupSvcGroupList, nil
	}
	vSvcGroupAllObj, err := vSvcGroupParent.GetAllSourceSvcGroups(context.TODO())
	if err != nil {
		log.Errorf("[getPolicyACPConfigSourceSvcGroupsResolver]Error getting SourceSvcGroups %s", err)
		return vServicegroupSvcGroupList, nil
	}
	for _, i := range vSvcGroupAllObj {
		vSvcGroupParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), getParentName(obj.ParentLabels, "acpconfigs.policypkg.tsm.tanzu.vmware.com"))
		if err != nil {
			log.Errorf("[getPolicyACPConfigSourceSvcGroupsResolver]Error getting parent node %s, skipping...", err)
			continue
		}
		vSvcGroup, err := vSvcGroupParent.GetSourceSvcGroups(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getPolicyACPConfigSourceSvcGroupsResolver]Error getting SourceSvcGroups node %q : %s, skipping...", i.DisplayName(), err)
			continue
		}
		dn := vSvcGroup.DisplayName()
		parentLabels := map[string]interface{}{"svcgroups.servicegroup.tsm.tanzu.vmware.com": dn}
		vDisplayName := string(vSvcGroup.Spec.DisplayName)
		vDescription := string(vSvcGroup.Spec.Description)
		vColor := string(vSvcGroup.Spec.Color)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ServicegroupSvcGroup{
			Id:           &dn,
			ParentLabels: parentLabels,
			DisplayName:  &vDisplayName,
			Description:  &vDescription,
			Color:        &vColor,
		}
		vServicegroupSvcGroupList = append(vServicegroupSvcGroupList, ret)
	}
	log.Debugf("[getPolicyACPConfigSourceSvcGroupsResolver]List of SourceSvcGroups object %v", vServicegroupSvcGroupList)
	return vServicegroupSvcGroupList, nil
}
