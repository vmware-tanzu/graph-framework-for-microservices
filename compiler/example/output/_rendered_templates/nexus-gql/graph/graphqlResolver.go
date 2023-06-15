package graph

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	qm "github.com/vmware-tanzu/graph-framework-for-microservices/nexus/generated/query-manager"
	nexus_client "nexustempmodule/nexus-client"
	"nexustempmodule/nexus-gql/graph/model"
)

var c = GrpcClients{
		mtx: sync.Mutex{},
		Clients: map[string]GrpcClient{},
}
var nc *nexus_client.Clientset

func getParentName(parentLabels map[string]interface{}, key string) string {
    if v, ok := parentLabels[key]; ok && v != nil {
	    return v.(string)
	}
	return ""
}

type NodeMetricTypeEnum string
type ServiceMetricTypeEnum string
type ServiceGroupByEnum string
type HTTPMethodEnum string
type EventSeverityEnum string
type AnalyticsMetricEnum string
type AnalyticsSubMetricEnum string
type TrafficDirectionEnum string
type SloDetailsEnum string

//////////////////////////////////////
// Nexus K8sAPIEndpointConfig
//////////////////////////////////////
func getK8sAPIEndpointConfig() *rest.Config {
    var (
		config *rest.Config
		err    error
	)
	filePath := os.Getenv("KUBECONFIG")
	if filePath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", filePath)
		if err != nil {
			return nil
		}
	} else {
	    config, err = rest.InClusterConfig()
	    if err != nil {
		    return nil
	    }
	}
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(200, 300)
	return config
}
//////////////////////////////////////
// Singleton Resolver for Parent Node
// PKG: Root, NODE: Root
//////////////////////////////////////
func getRootResolver() (*model.RootRoot, error) {
	if nc == nil {
		k8sApiConfig := getK8sAPIEndpointConfig()
		nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get k8s client config: %s", err)
		}
		nc = nexusClient
		nc.SubscribeAll()
		log.Debugf("Subscribed to all nodes in datamodel")
	}

	vRoot, err := nc.GetRootRoot(context.TODO())
	if err != nil {
		log.Errorf("[getRootResolver]Error getting Root node %s", err)
		return nil, nil
	}
	dn := vRoot.DisplayName()
parentLabels := map[string]interface{}{"roots.root.tsm.tanzu.vmware.com":dn}

	ret := &model.RootRoot {
	Id: &dn,
	ParentLabels: parentLabels,
	}
	log.Debugf("[getRootResolver]Output Root object %+v", ret)
	return ret, nil
}
// Custom query
func getConfigConfigQueryExampleResolver(obj *model.ConfigConfig,  StartTime *string, EndTime *string, Interval *string, IsServiceDeployment *bool, StartVal *int,) (*model.NexusGraphqlResponse, error) {
	parentLabels := make(map[string]string)
	if obj != nil {
		for k, v := range obj.ParentLabels {
			val, ok := v.(string)
			if ok {
				parentLabels[k] = val
			}
		}
	}
	query := &graphql.GraphQLQuery{
		Query: "QueryExample",
		UserProvidedArgs: map[string]string{
			"StartTime": pointerToString(StartTime),
			"EndTime": pointerToString(EndTime),
			"Interval": pointerToString(Interval),
			"IsServiceDeployment": pointerToString(IsServiceDeployment),
			"StartVal": pointerToString(StartVal),
		},
		Hierarchy: parentLabels,
	}

	resp, err := c.Request("query-manager:6000", nexus.GraphQLQueryApi, query)
	if err != nil {
		return nil, err
	}
	return resp.(*model.NexusGraphqlResponse), nil
}
// Custom query
func getGnsGnsqueryGns1Resolver(obj *model.GnsGns,  StartTime *string, EndTime *string, Interval *string, IsServiceDeployment *bool, StartVal *int,) (*model.NexusGraphqlResponse, error) {
	parentLabels := make(map[string]string)
	if obj != nil {
		for k, v := range obj.ParentLabels {
			val, ok := v.(string)
			if ok {
				parentLabels[k] = val
			}
		}
	}
	query := &graphql.GraphQLQuery{
		Query: "queryGns1",
		UserProvidedArgs: map[string]string{
			"StartTime": pointerToString(StartTime),
			"EndTime": pointerToString(EndTime),
			"Interval": pointerToString(Interval),
			"IsServiceDeployment": pointerToString(IsServiceDeployment),
			"StartVal": pointerToString(StartVal),
		},
		Hierarchy: parentLabels,
	}

	resp, err := c.Request("nexus-query-responder:15000", nexus.GraphQLQueryApi, query)
	if err != nil {
		return nil, err
	}
	return resp.(*model.NexusGraphqlResponse), nil
}
// Custom query
func getGnsGnsqueryGnsQM1Resolver(obj *model.GnsGns, ) (*model.TimeSeriesData, error) {
	parentLabels := make(map[string]string)
	if obj != nil {
		for k, v := range obj.ParentLabels {
			val, ok := v.(string)
			if ok {
				parentLabels[k] = val
			}
		}
	}
	metricArgs := &qm.MetricArg{
		QueryType: "/queryGnsQM1",
		Hierarchy: parentLabels,
		UserProvidedArgs: map[string]string{
		},
	}
	resp, err := c.Request("query-manager:15002", nexus.GetMetricsApi, metricArgs)
	if err != nil {
		return nil, err
	}
	return resp.(*model.TimeSeriesData), nil
}
// Custom query
func getGnsGnsqueryGnsQMResolver(obj *model.GnsGns,  StartTime *string, EndTime *string, TimeInterval *string, SomeUserArg1 *string, SomeUserArg2 *int, SomeUserArg3 *bool,) (*model.TimeSeriesData, error) {
	parentLabels := make(map[string]string)
	if obj != nil {
		for k, v := range obj.ParentLabels {
			val, ok := v.(string)
			if ok {
				parentLabels[k] = val
			}
		}
	}
	metricArgs := &qm.MetricArg{
		QueryType: "/queryGnsQM",
		StartTime: *StartTime,
		EndTime: *EndTime,
		TimeInterval: *TimeInterval,
		Hierarchy: parentLabels,
		UserProvidedArgs: map[string]string{
			"SomeUserArg1": pointerToString(SomeUserArg1),
			"SomeUserArg2": pointerToString(SomeUserArg2),
			"SomeUserArg3": pointerToString(SomeUserArg3),
		},
	}
	resp, err := c.Request("query-manager:15003", nexus.GetMetricsApi, metricArgs)
	if err != nil {
		return nil, err
	}
	return resp.(*model.TimeSeriesData), nil
}
// Custom query
func getPolicypkgVMpolicyqueryGns1Resolver(obj *model.PolicypkgVMpolicy,  StartTime *string, EndTime *string, Interval *string, IsServiceDeployment *bool, StartVal *int,) (*model.NexusGraphqlResponse, error) {
	parentLabels := make(map[string]string)
	if obj != nil {
		for k, v := range obj.ParentLabels {
			val, ok := v.(string)
			if ok {
				parentLabels[k] = val
			}
		}
	}
	query := &graphql.GraphQLQuery{
		Query: "queryGns1",
		UserProvidedArgs: map[string]string{
			"StartTime": pointerToString(StartTime),
			"EndTime": pointerToString(EndTime),
			"Interval": pointerToString(Interval),
			"IsServiceDeployment": pointerToString(IsServiceDeployment),
			"StartVal": pointerToString(StartVal),
		},
		Hierarchy: parentLabels,
	}

	resp, err := c.Request("nexus-query-responder:15000", nexus.GraphQLQueryApi, query)
	if err != nil {
		return nil, err
	}
	return resp.(*model.NexusGraphqlResponse), nil
}
// Custom query
func getPolicypkgVMpolicyqueryGnsQM1Resolver(obj *model.PolicypkgVMpolicy, ) (*model.TimeSeriesData, error) {
	parentLabels := make(map[string]string)
	if obj != nil {
		for k, v := range obj.ParentLabels {
			val, ok := v.(string)
			if ok {
				parentLabels[k] = val
			}
		}
	}
	metricArgs := &qm.MetricArg{
		QueryType: "/queryGnsQM1",
		Hierarchy: parentLabels,
		UserProvidedArgs: map[string]string{
		},
	}
	resp, err := c.Request("query-manager:15002", nexus.GetMetricsApi, metricArgs)
	if err != nil {
		return nil, err
	}
	return resp.(*model.TimeSeriesData), nil
}
// Custom query
func getPolicypkgVMpolicyqueryGnsQMResolver(obj *model.PolicypkgVMpolicy,  StartTime *string, EndTime *string, TimeInterval *string, SomeUserArg1 *string, SomeUserArg2 *int, SomeUserArg3 *bool,) (*model.TimeSeriesData, error) {
	parentLabels := make(map[string]string)
	if obj != nil {
		for k, v := range obj.ParentLabels {
			val, ok := v.(string)
			if ok {
				parentLabels[k] = val
			}
		}
	}
	metricArgs := &qm.MetricArg{
		QueryType: "/queryGnsQM",
		StartTime: *StartTime,
		EndTime: *EndTime,
		TimeInterval: *TimeInterval,
		Hierarchy: parentLabels,
		UserProvidedArgs: map[string]string{
			"SomeUserArg1": pointerToString(SomeUserArg1),
			"SomeUserArg2": pointerToString(SomeUserArg2),
			"SomeUserArg3": pointerToString(SomeUserArg3),
		},
	}
	resp, err := c.Request("query-manager:15003", nexus.GetMetricsApi, metricArgs)
	if err != nil {
		return nil, err
	}
	return resp.(*model.TimeSeriesData), nil
}
//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: Config Node: Root PKG: Root
//////////////////////////////////////
func getRootRootConfigResolver(obj *model.RootRoot, id *string) (*model.ConfigConfig, error) {
	log.Debugf("[getRootRootConfigResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
	     log.Debugf("[getRootRootConfigResolver]Id %q", *id)
		vConfig, err := nc.RootRoot().GetConfig(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootRootConfigResolver]Error getting Config node %q : %s", *id, err)
			return &model.ConfigConfig{}, nil
		}
		dn := vConfig.DisplayName()
parentLabels := map[string]interface{}{"configs.config.tsm.tanzu.vmware.com":dn}
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
vCuOption := string(vConfig.Spec.CuOption)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConfigConfig {
	Id: &dn,
	ParentLabels: parentLabels,
	MyStr0: &MyStr0Data,
	MyStr1: &MyStr1Data,
	MyStr2: &MyStr2Data,
	XYZPort: &XYZPortData,
	ABCHost: &ABCHostData,
	ClusterNamespaces: &ClusterNamespacesData,
	TestValMarkers: &TestValMarkersData,
	Instance: &vInstance,
	CuOption: &vCuOption,
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
	dn := vConfig.DisplayName()
parentLabels := map[string]interface{}{"configs.config.tsm.tanzu.vmware.com":dn}
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
vCuOption := string(vConfig.Spec.CuOption)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.ConfigConfig {
	Id: &dn,
	ParentLabels: parentLabels,
	MyStr0: &MyStr0Data,
	MyStr1: &MyStr1Data,
	MyStr2: &MyStr2Data,
	XYZPort: &XYZPortData,
	ABCHost: &ABCHostData,
	ClusterNamespaces: &ClusterNamespacesData,
	TestValMarkers: &TestValMarkersData,
	Instance: &vInstance,
	CuOption: &vCuOption,
	}

	log.Debugf("[getRootRootConfigResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: GNS Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigGNSResolver(obj *model.ConfigConfig, id *string) (*model.GnsGns, error) {
	log.Debugf("[getConfigConfigGNSResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
	     log.Debugf("[getConfigConfigGNSResolver]Id %q", *id)
		vGns, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetGNS(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigGNSResolver]Error getting GNS node %q : %s", *id, err)
			return &model.GnsGns{}, nil
		}
		dn := vGns.DisplayName()
parentLabels := map[string]interface{}{"gnses.gns.tsm.tanzu.vmware.com":dn}
vDomain := string(vGns.Spec.Domain)
vUseSharedGateway := bool(vGns.Spec.UseSharedGateway)
Annotations, _ := json.Marshal(vGns.Spec.Annotations)
AnnotationsData := string(Annotations)
TargetPort, _ := json.Marshal(vGns.Spec.TargetPort)
TargetPortData := string(TargetPort)
Description, _ := json.Marshal(vGns.Spec.Description)
DescriptionData := string(Description)
vMeta := string(vGns.Spec.Meta)
IntOrString, _ := json.Marshal(vGns.Spec.IntOrString)
IntOrStringData := string(IntOrString)
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
ServiceSegmentRef, _ := json.Marshal(vGns.Spec.ServiceSegmentRef)
ServiceSegmentRefData := string(ServiceSegmentRef)
ServiceSegmentRefPointer, _ := json.Marshal(vGns.Spec.ServiceSegmentRefPointer)
ServiceSegmentRefPointerData := string(ServiceSegmentRefPointer)
ServiceSegmentRefs, _ := json.Marshal(vGns.Spec.ServiceSegmentRefs)
ServiceSegmentRefsData := string(ServiceSegmentRefs)
ServiceSegmentRefMap, _ := json.Marshal(vGns.Spec.ServiceSegmentRefMap)
ServiceSegmentRefMapData := string(ServiceSegmentRefMap)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.GnsGns {
	Id: &dn,
	ParentLabels: parentLabels,
	Domain: &vDomain,
	UseSharedGateway: &vUseSharedGateway,
	Annotations: &AnnotationsData,
	TargetPort: &TargetPortData,
	Description: &DescriptionData,
	Meta: &vMeta,
	IntOrString: &IntOrStringData,
	OtherDescription: &OtherDescriptionData,
	MapPointer: &MapPointerData,
	SlicePointer: &SlicePointerData,
	WorkloadSpec: &WorkloadSpecData,
	DifferentSpec: &DifferentSpecData,
	ServiceSegmentRef: &ServiceSegmentRefData,
	ServiceSegmentRefPointer: &ServiceSegmentRefPointerData,
	ServiceSegmentRefs: &ServiceSegmentRefsData,
	ServiceSegmentRefMap: &ServiceSegmentRefMapData,
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
	dn := vGns.DisplayName()
parentLabels := map[string]interface{}{"gnses.gns.tsm.tanzu.vmware.com":dn}
vDomain := string(vGns.Spec.Domain)
vUseSharedGateway := bool(vGns.Spec.UseSharedGateway)
Annotations, _ := json.Marshal(vGns.Spec.Annotations)
AnnotationsData := string(Annotations)
TargetPort, _ := json.Marshal(vGns.Spec.TargetPort)
TargetPortData := string(TargetPort)
Description, _ := json.Marshal(vGns.Spec.Description)
DescriptionData := string(Description)
vMeta := string(vGns.Spec.Meta)
IntOrString, _ := json.Marshal(vGns.Spec.IntOrString)
IntOrStringData := string(IntOrString)
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
ServiceSegmentRef, _ := json.Marshal(vGns.Spec.ServiceSegmentRef)
ServiceSegmentRefData := string(ServiceSegmentRef)
ServiceSegmentRefPointer, _ := json.Marshal(vGns.Spec.ServiceSegmentRefPointer)
ServiceSegmentRefPointerData := string(ServiceSegmentRefPointer)
ServiceSegmentRefs, _ := json.Marshal(vGns.Spec.ServiceSegmentRefs)
ServiceSegmentRefsData := string(ServiceSegmentRefs)
ServiceSegmentRefMap, _ := json.Marshal(vGns.Spec.ServiceSegmentRefMap)
ServiceSegmentRefMapData := string(ServiceSegmentRefMap)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.GnsGns {
	Id: &dn,
	ParentLabels: parentLabels,
	Domain: &vDomain,
	UseSharedGateway: &vUseSharedGateway,
	Annotations: &AnnotationsData,
	TargetPort: &TargetPortData,
	Description: &DescriptionData,
	Meta: &vMeta,
	IntOrString: &IntOrStringData,
	OtherDescription: &OtherDescriptionData,
	MapPointer: &MapPointerData,
	SlicePointer: &SlicePointerData,
	WorkloadSpec: &WorkloadSpecData,
	DifferentSpec: &DifferentSpecData,
	ServiceSegmentRef: &ServiceSegmentRefData,
	ServiceSegmentRefPointer: &ServiceSegmentRefPointerData,
	ServiceSegmentRefs: &ServiceSegmentRefsData,
	ServiceSegmentRefMap: &ServiceSegmentRefMapData,
	}

	log.Debugf("[getConfigConfigGNSResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Singleton)
// FieldName: DNS Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigDNSResolver(obj *model.ConfigConfig) (*model.GnsDns, error) {
	log.Debugf("[getConfigConfigDNSResolver]Parent Object %+v", obj)
	vDns, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetDNS(context.TODO())
	if err != nil {
	    log.Errorf("[getConfigConfigDNSResolver]Error getting Config node %s", err)
        return &model.GnsDns{}, nil
    }
	dn := vDns.DisplayName()
parentLabels := map[string]interface{}{"dnses.gns.tsm.tanzu.vmware.com":dn}

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.GnsDns {
	Id: &dn,
	ParentLabels: parentLabels,
	}

    log.Debugf("[getConfigConfigDNSResolver]Output object %+v", ret)
	return ret, nil
}
//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: VMPPolicies Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigVMPPoliciesResolver(obj *model.ConfigConfig, id *string) (*model.PolicypkgVMpolicy, error) {
	log.Debugf("[getConfigConfigVMPPoliciesResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
	     log.Debugf("[getConfigConfigVMPPoliciesResolver]Id %q", *id)
		vVMpolicy, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetVMPPolicies(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigVMPPoliciesResolver]Error getting VMPPolicies node %q : %s", *id, err)
			return &model.PolicypkgVMpolicy{}, nil
		}
		dn := vVMpolicy.DisplayName()
parentLabels := map[string]interface{}{"vmpolicies.policypkg.tsm.tanzu.vmware.com":dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.PolicypkgVMpolicy {
	Id: &dn,
	ParentLabels: parentLabels,
	}

		log.Debugf("[getConfigConfigVMPPoliciesResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getConfigConfigVMPPoliciesResolver]Id is empty, process all VMPPoliciess")
	vVMpolicyParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
	    log.Errorf("[getConfigConfigVMPPoliciesResolver]Failed to get parent node %s", err)
        return &model.PolicypkgVMpolicy{}, nil
    }
	vVMpolicy, err := vVMpolicyParent.GetVMPPolicies(context.TODO())
	if err != nil {
	    log.Errorf("[getConfigConfigVMPPoliciesResolver]Error getting VMPPolicies node %s", err)
        return &model.PolicypkgVMpolicy{}, nil
    }
	dn := vVMpolicy.DisplayName()
parentLabels := map[string]interface{}{"vmpolicies.policypkg.tsm.tanzu.vmware.com":dn}

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.PolicypkgVMpolicy {
	Id: &dn,
	ParentLabels: parentLabels,
	}

	log.Debugf("[getConfigConfigVMPPoliciesResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: Domain Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigDomainResolver(obj *model.ConfigConfig, id *string) (*model.ConfigDomain, error) {
	log.Debugf("[getConfigConfigDomainResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
	     log.Debugf("[getConfigConfigDomainResolver]Id %q", *id)
		vDomain, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetDomain(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigDomainResolver]Error getting Domain node %q : %s", *id, err)
			return &model.ConfigDomain{}, nil
		}
		dn := vDomain.DisplayName()
parentLabels := map[string]interface{}{"domains.config.tsm.tanzu.vmware.com":dn}
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
		ret := &model.ConfigDomain {
	Id: &dn,
	ParentLabels: parentLabels,
	PointPort: &PointPortData,
	PointMap: &PointMapData,
	PointSlice: &PointSliceData,
	SliceOfPoints: &SliceOfPointsData,
	SliceOfArrPoints: &SliceOfArrPointsData,
	MapOfArrsPoints: &MapOfArrsPointsData,
	PointStruct: &PointStructData,
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
	dn := vDomain.DisplayName()
parentLabels := map[string]interface{}{"domains.config.tsm.tanzu.vmware.com":dn}
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
	ret := &model.ConfigDomain {
	Id: &dn,
	ParentLabels: parentLabels,
	PointPort: &PointPortData,
	PointMap: &PointMapData,
	PointSlice: &PointSliceData,
	SliceOfPoints: &SliceOfPointsData,
	SliceOfArrPoints: &SliceOfArrPointsData,
	MapOfArrsPoints: &MapOfArrsPointsData,
	PointStruct: &PointStructData,
	}

	log.Debugf("[getConfigConfigDomainResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: SvcGrpInfo Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigSvcGrpInfoResolver(obj *model.ConfigConfig, id *string) (*model.ServicegroupSvcGroupLinkInfo, error) {
	log.Debugf("[getConfigConfigSvcGrpInfoResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
	     log.Debugf("[getConfigConfigSvcGrpInfoResolver]Id %q", *id)
		vSvcGroupLinkInfo, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetSvcGrpInfo(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigSvcGrpInfoResolver]Error getting SvcGrpInfo node %q : %s", *id, err)
			return &model.ServicegroupSvcGroupLinkInfo{}, nil
		}
		dn := vSvcGroupLinkInfo.DisplayName()
parentLabels := map[string]interface{}{"svcgrouplinkinfos.servicegroup.tsm.tanzu.vmware.com":dn}
vClusterName := string(vSvcGroupLinkInfo.Spec.ClusterName)
vDomainName := string(vSvcGroupLinkInfo.Spec.DomainName)
vServiceName := string(vSvcGroupLinkInfo.Spec.ServiceName)
vServiceType := string(vSvcGroupLinkInfo.Spec.ServiceType)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ServicegroupSvcGroupLinkInfo {
	Id: &dn,
	ParentLabels: parentLabels,
	ClusterName: &vClusterName,
	DomainName: &vDomainName,
	ServiceName: &vServiceName,
	ServiceType: &vServiceType,
	}

		log.Debugf("[getConfigConfigSvcGrpInfoResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getConfigConfigSvcGrpInfoResolver]Id is empty, process all SvcGrpInfos")
	vSvcGroupLinkInfoParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
	    log.Errorf("[getConfigConfigSvcGrpInfoResolver]Failed to get parent node %s", err)
        return &model.ServicegroupSvcGroupLinkInfo{}, nil
    }
	vSvcGroupLinkInfo, err := vSvcGroupLinkInfoParent.GetSvcGrpInfo(context.TODO())
	if err != nil {
	    log.Errorf("[getConfigConfigSvcGrpInfoResolver]Error getting SvcGrpInfo node %s", err)
        return &model.ServicegroupSvcGroupLinkInfo{}, nil
    }
	dn := vSvcGroupLinkInfo.DisplayName()
parentLabels := map[string]interface{}{"svcgrouplinkinfos.servicegroup.tsm.tanzu.vmware.com":dn}
vClusterName := string(vSvcGroupLinkInfo.Spec.ClusterName)
vDomainName := string(vSvcGroupLinkInfo.Spec.DomainName)
vServiceName := string(vSvcGroupLinkInfo.Spec.ServiceName)
vServiceType := string(vSvcGroupLinkInfo.Spec.ServiceType)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.ServicegroupSvcGroupLinkInfo {
	Id: &dn,
	ParentLabels: parentLabels,
	ClusterName: &vClusterName,
	DomainName: &vDomainName,
	ServiceName: &vServiceName,
	ServiceType: &vServiceType,
	}

	log.Debugf("[getConfigConfigSvcGrpInfoResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: FooExample Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigFooExampleResolver(obj *model.ConfigConfig, id *string) ([]*model.ConfigFooTypeABC, error) {
	log.Debugf("[getConfigConfigFooExampleResolver]Parent Object %+v", obj)
	var vConfigFooTypeABCList []*model.ConfigFooTypeABC
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigFooExampleResolver]Id %q", *id)
		vFooTypeABC, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetFooExample(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigFooExampleResolver]Error getting FooExample node %q : %s", *id, err)
            return vConfigFooTypeABCList, nil
        }
		dn := vFooTypeABC.DisplayName()
parentLabels := map[string]interface{}{"footypeabcs.config.tsm.tanzu.vmware.com":dn}
FooA, _ := json.Marshal(vFooTypeABC.Spec.FooA)
FooAData := string(FooA)
FooB, _ := json.Marshal(vFooTypeABC.Spec.FooB)
FooBData := string(FooB)
FooD, _ := json.Marshal(vFooTypeABC.Spec.FooD)
FooDData := string(FooD)
FooF, _ := json.Marshal(vFooTypeABC.Spec.FooF)
FooFData := string(FooF)

        for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.ConfigFooTypeABC {
	Id: &dn,
	ParentLabels: parentLabels,
	FooA: &FooAData,
	FooB: &FooBData,
	FooD: &FooDData,
	FooF: &FooFData,
	}
		vConfigFooTypeABCList = append(vConfigFooTypeABCList, ret)

		log.Debugf("[getConfigConfigFooExampleResolver]Output FooExample objects %v", vConfigFooTypeABCList)

		return vConfigFooTypeABCList, nil
	}

	log.Debug("[getConfigConfigFooExampleResolver]Id is empty, process all FooExamples")

	vFooTypeABCParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
	    log.Errorf("[getConfigConfigFooExampleResolver]Error getting parent node %s", err)
        return vConfigFooTypeABCList, nil
    }
	vFooTypeABCAllObj, err := vFooTypeABCParent.GetAllFooExample(context.TODO())
	if err != nil {
	    log.Errorf("[getConfigConfigFooExampleResolver]Error getting FooExample objects %s", err)
        return vConfigFooTypeABCList, nil
    }
	for _, i := range vFooTypeABCAllObj {
		vFooTypeABC, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetFooExample(context.TODO(), i.DisplayName())
		if err != nil {
	        log.Errorf("[getConfigConfigFooExampleResolver]Error getting FooExample node %q : %s", i.DisplayName(), err)
            continue
		}
		dn := vFooTypeABC.DisplayName()
parentLabels := map[string]interface{}{"footypeabcs.config.tsm.tanzu.vmware.com":dn}
FooA, _ := json.Marshal(vFooTypeABC.Spec.FooA)
FooAData := string(FooA)
FooB, _ := json.Marshal(vFooTypeABC.Spec.FooB)
FooBData := string(FooB)
FooD, _ := json.Marshal(vFooTypeABC.Spec.FooD)
FooDData := string(FooD)
FooF, _ := json.Marshal(vFooTypeABC.Spec.FooF)
FooFData := string(FooF)

		for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.ConfigFooTypeABC {
	Id: &dn,
	ParentLabels: parentLabels,
	FooA: &FooAData,
	FooB: &FooBData,
	FooD: &FooDData,
	FooF: &FooFData,
	}
		vConfigFooTypeABCList = append(vConfigFooTypeABCList, ret)
	}

	log.Debugf("[getConfigConfigFooExampleResolver]Output FooExample objects %v", vConfigFooTypeABCList)

	return vConfigFooTypeABCList, nil
}

//////////////////////////////////////
// LINKS RESOLVER
// FieldName: ACPPolicies Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigACPPoliciesResolver(obj *model.ConfigConfig, id *string) ([]*model.PolicypkgAccessControlPolicy, error) {
	log.Debugf("[getConfigConfigACPPoliciesResolver]Parent Object %+v", obj)
	var vPolicypkgAccessControlPolicyList []*model.PolicypkgAccessControlPolicy
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigACPPoliciesResolver]Id %q", *id)
		vAccessControlPolicyParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
		if err != nil {
			log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting ACPPolicies %q : %s", *id, err)
			return vPolicypkgAccessControlPolicyList, nil
		}
		vAccessControlPolicy, err := vAccessControlPolicyParent.GetACPPolicies(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting ACPPolicies %q : %s", *id, err)
			return vPolicypkgAccessControlPolicyList, nil
		}
		dn := vAccessControlPolicy.DisplayName()
parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com":dn}

        for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.PolicypkgAccessControlPolicy {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vPolicypkgAccessControlPolicyList = append(vPolicypkgAccessControlPolicyList, ret)

		log.Debugf("[getConfigConfigACPPoliciesResolver]Output ACPPolicies objects %v", vPolicypkgAccessControlPolicyList)

		return vPolicypkgAccessControlPolicyList, nil
	}

	log.Debug("[getConfigConfigACPPoliciesResolver]Id is empty, process all ACPPoliciess")

	vAccessControlPolicyParent, err := nc.RootRoot().GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com"))
	if err != nil {
	    log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting parent node %s", err)
        return vPolicypkgAccessControlPolicyList, nil
    }
	vAccessControlPolicyAllObj, err := vAccessControlPolicyParent.GetAllACPPolicies(context.TODO())
	if err != nil {
	    log.Errorf("[getConfigConfigACPPoliciesResolver]Error getting ACPPolicies %s", err)
        return vPolicypkgAccessControlPolicyList, nil
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
parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com":dn}

		for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.PolicypkgAccessControlPolicy {
	Id: &dn,
	ParentLabels: parentLabels,
	}
		vPolicypkgAccessControlPolicyList = append(vPolicypkgAccessControlPolicyList, ret)
	}
	log.Debugf("[getConfigConfigACPPoliciesResolver]List of ACPPolicies object %v", vPolicypkgAccessControlPolicyList)
	return vPolicypkgAccessControlPolicyList, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: GnsAccessControlPolicy Node: Gns PKG: Gns
//////////////////////////////////////
func getGnsGnsGnsAccessControlPolicyResolver(obj *model.GnsGns, id *string) (*model.PolicypkgAccessControlPolicy, error) {
	log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
	     log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Id %q", *id)
		vAccessControlPolicy, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GetGnsAccessControlPolicy(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getGnsGnsGnsAccessControlPolicyResolver]Error getting GnsAccessControlPolicy node %q : %s", *id, err)
			return &model.PolicypkgAccessControlPolicy{}, nil
		}
		dn := vAccessControlPolicy.DisplayName()
parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com":dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.PolicypkgAccessControlPolicy {
	Id: &dn,
	ParentLabels: parentLabels,
	}

		log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getGnsGnsGnsAccessControlPolicyResolver]Id is empty, process all GnsAccessControlPolicys")
	vAccessControlPolicyParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GetGNS(context.TODO(), getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com"))
	if err != nil {
	    log.Errorf("[getGnsGnsGnsAccessControlPolicyResolver]Failed to get parent node %s", err)
        return &model.PolicypkgAccessControlPolicy{}, nil
    }
	vAccessControlPolicy, err := vAccessControlPolicyParent.GetGnsAccessControlPolicy(context.TODO())
	if err != nil {
	    log.Errorf("[getGnsGnsGnsAccessControlPolicyResolver]Error getting GnsAccessControlPolicy node %s", err)
        return &model.PolicypkgAccessControlPolicy{}, nil
    }
	dn := vAccessControlPolicy.DisplayName()
parentLabels := map[string]interface{}{"accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com":dn}

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.PolicypkgAccessControlPolicy {
	Id: &dn,
	ParentLabels: parentLabels,
	}

	log.Debugf("[getGnsGnsGnsAccessControlPolicyResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Singleton)
// FieldName: FooChild Node: Gns PKG: Gns
//////////////////////////////////////
func getGnsGnsFooChildResolver(obj *model.GnsGns) (*model.GnsBarChild, error) {
	log.Debugf("[getGnsGnsFooChildResolver]Parent Object %+v", obj)
	vBarChild, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GetFooChild(context.TODO())
	if err != nil {
	    log.Errorf("[getGnsGnsFooChildResolver]Error getting Gns node %s", err)
        return &model.GnsBarChild{}, nil
    }
	dn := vBarChild.DisplayName()
parentLabels := map[string]interface{}{"barchilds.gns.tsm.tanzu.vmware.com":dn}
vName := string(vBarChild.Spec.Name)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.GnsBarChild {
	Id: &dn,
	ParentLabels: parentLabels,
	Name: &vName,
	}

    log.Debugf("[getGnsGnsFooChildResolver]Output object %+v", ret)
	return ret, nil
}
//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: PolicyConfigs Node: AccessControlPolicy PKG: Policypkg
//////////////////////////////////////
func getPolicypkgAccessControlPolicyPolicyConfigsResolver(obj *model.PolicypkgAccessControlPolicy, id *string) ([]*model.PolicypkgACPConfig, error) {
	log.Debugf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Parent Object %+v", obj)
	var vPolicypkgACPConfigList []*model.PolicypkgACPConfig
	if id != nil && *id != "" {
		log.Debugf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Id %q", *id)
		vACPConfig, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Error getting PolicyConfigs node %q : %s", *id, err)
            return vPolicypkgACPConfigList, nil
        }
		dn := vACPConfig.DisplayName()
parentLabels := map[string]interface{}{"acpconfigs.policypkg.tsm.tanzu.vmware.com":dn}
vDisplayName := string(vACPConfig.Spec.DisplayName)
vGns := string(vACPConfig.Spec.Gns)
vDescription := string(vACPConfig.Spec.Description)
Tags, _ := json.Marshal(vACPConfig.Spec.Tags)
TagsData := string(Tags)
vProjectId := string(vACPConfig.Spec.ProjectId)
Conditions, _ := json.Marshal(vACPConfig.Spec.Conditions)
ConditionsData := string(Conditions)

        for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.PolicypkgACPConfig {
	Id: &dn,
	ParentLabels: parentLabels,
	DisplayName: &vDisplayName,
	Gns: &vGns,
	Description: &vDescription,
	Tags: &TagsData,
	ProjectId: &vProjectId,
	Conditions: &ConditionsData,
	}
		vPolicypkgACPConfigList = append(vPolicypkgACPConfigList, ret)

		log.Debugf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Output PolicyConfigs objects %v", vPolicypkgACPConfigList)

		return vPolicypkgACPConfigList, nil
	}

	log.Debug("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Id is empty, process all PolicyConfigss")

	vACPConfigParent, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GetGnsAccessControlPolicy(context.TODO(), getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com"))
	if err != nil {
	    log.Errorf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Error getting parent node %s", err)
        return vPolicypkgACPConfigList, nil
    }
	vACPConfigAllObj, err := vACPConfigParent.GetAllPolicyConfigs(context.TODO())
	if err != nil {
	    log.Errorf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Error getting PolicyConfigs objects %s", err)
        return vPolicypkgACPConfigList, nil
    }
	for _, i := range vACPConfigAllObj {
		vACPConfig, err := nc.RootRoot().Config(getParentName(obj.ParentLabels, "configs.config.tsm.tanzu.vmware.com")).GNS(getParentName(obj.ParentLabels, "gnses.gns.tsm.tanzu.vmware.com")).GnsAccessControlPolicy(getParentName(obj.ParentLabels, "accesscontrolpolicies.policypkg.tsm.tanzu.vmware.com")).GetPolicyConfigs(context.TODO(), i.DisplayName())
		if err != nil {
	        log.Errorf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Error getting PolicyConfigs node %q : %s", i.DisplayName(), err)
            continue
		}
		dn := vACPConfig.DisplayName()
parentLabels := map[string]interface{}{"acpconfigs.policypkg.tsm.tanzu.vmware.com":dn}
vDisplayName := string(vACPConfig.Spec.DisplayName)
vGns := string(vACPConfig.Spec.Gns)
vDescription := string(vACPConfig.Spec.Description)
Tags, _ := json.Marshal(vACPConfig.Spec.Tags)
TagsData := string(Tags)
vProjectId := string(vACPConfig.Spec.ProjectId)
Conditions, _ := json.Marshal(vACPConfig.Spec.Conditions)
ConditionsData := string(Conditions)

		for k, v := range obj.ParentLabels {
            parentLabels[k] = v
        }
		ret := &model.PolicypkgACPConfig {
	Id: &dn,
	ParentLabels: parentLabels,
	DisplayName: &vDisplayName,
	Gns: &vGns,
	Description: &vDescription,
	Tags: &TagsData,
	ProjectId: &vProjectId,
	Conditions: &ConditionsData,
	}
		vPolicypkgACPConfigList = append(vPolicypkgACPConfigList, ret)
	}

	log.Debugf("[getPolicypkgAccessControlPolicyPolicyConfigsResolver]Output PolicyConfigs objects %v", vPolicypkgACPConfigList)

	return vPolicypkgACPConfigList, nil
}

