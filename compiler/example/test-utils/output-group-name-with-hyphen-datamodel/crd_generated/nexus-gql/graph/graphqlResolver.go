package graph

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	qm "github.com/vmware-tanzu/graph-framework-for-microservices/nexus/generated/query-manager"
	nexus_client "../../example/test-utils/output-group-name-with-hyphen-datamodel/crd_generated/nexus-client"
	"../../example/test-utils/output-group-name-with-hyphen-datamodel/crd_generated/nexus-gql/graph/model"
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
// Non Singleton Resolver for Parent Node
// PKG: Root, NODE: Root
//////////////////////////////////////
func getRootResolver(id *string) ([]*model.RootRoot, error) {
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

	var vRootList []*model.RootRoot
	if id != nil && *id != "" {
		log.Debugf("[getRootResolver]Id: %q", *id)
		vRoot, err := nc.GetRootRoot(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting Root node %q: %s", *id, err)
			return nil, nil
		}
		dn := vRoot.DisplayName()
parentLabels := map[string]interface{}{"roots.root.tsm-tanzu.vmware.com":dn}
vSomeRootData := string(vRoot.Spec.SomeRootData)

		ret := &model.RootRoot {
	Id: &dn,
	ParentLabels: parentLabels,
	SomeRootData: &vSomeRootData,
	}
		vRootList = append(vRootList, ret)
		log.Debugf("[getRootResolver]Output Root objects %+v", vRootList)
		return vRootList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all Roots")

	vRootListObj, err := nc.Root().ListRoots(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("[getRootResolver]Error getting Root node %s", err)
		return nil, nil
	}
	for _,i := range vRootListObj{
		vRoot, err := nc.GetRootRoot(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting Root node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vRoot.DisplayName()
parentLabels := map[string]interface{}{"roots.root.tsm-tanzu.vmware.com":dn}
vSomeRootData := string(vRoot.Spec.SomeRootData)

		ret := &model.RootRoot {
	Id: &dn,
	ParentLabels: parentLabels,
	SomeRootData: &vSomeRootData,
	}
		vRootList = append(vRootList, ret)
	}

	log.Debugf("[getRootResolver]Output Root objects %v", vRootList)
	return vRootList, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Singleton)
// FieldName: Project Node: Root PKG: Root
//////////////////////////////////////
func getRootRootProjectResolver(obj *model.RootRoot) (*model.ProjectProject, error) {
	log.Debugf("[getRootRootProjectResolver]Parent Object %+v", obj)
	vProject, err := nc.RootRoot(getParentName(obj.ParentLabels, "roots.root.tsm-tanzu.vmware.com")).GetProject(context.TODO())
	if err != nil {
	    log.Errorf("[getRootRootProjectResolver]Error getting Root node %s", err)
        return &model.ProjectProject{}, nil
    }
	dn := vProject.DisplayName()
parentLabels := map[string]interface{}{"projects.project.tsm-tanzu.vmware.com":dn}
vKey := string(vProject.Spec.Key)
vField1 := string(vProject.Spec.Field1)
vField2 := int(vProject.Spec.Field2)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.ProjectProject {
	Id: &dn,
	ParentLabels: parentLabels,
	Key: &vKey,
	Field1: &vField1,
	Field2: &vField2,
	}

    log.Debugf("[getRootRootProjectResolver]Output object %+v", ret)
	return ret, nil
}
//////////////////////////////////////
// CHILD RESOLVER (Singleton)
// FieldName: Config Node: Project PKG: Project
//////////////////////////////////////
func getProjectProjectConfigResolver(obj *model.ProjectProject) (*model.ConfigConfig, error) {
	log.Debugf("[getProjectProjectConfigResolver]Parent Object %+v", obj)
	vConfig, err := nc.RootRoot(getParentName(obj.ParentLabels, "roots.root.tsm-tanzu.vmware.com")).Project().GetConfig(context.TODO())
	if err != nil {
	    log.Errorf("[getProjectProjectConfigResolver]Error getting Project node %s", err)
        return &model.ConfigConfig{}, nil
    }
	dn := vConfig.DisplayName()
parentLabels := map[string]interface{}{"configs.config.tsm-tanzu.vmware.com":dn}
vFieldX := string(vConfig.Spec.FieldX)
vFieldY := int(vConfig.Spec.FieldY)
MyStructField, _ := json.Marshal(vConfig.Spec.MyStructField)
MyStructFieldData := string(MyStructField)

    for k, v := range obj.ParentLabels {
        parentLabels[k] = v
    }
	ret := &model.ConfigConfig {
	Id: &dn,
	ParentLabels: parentLabels,
	FieldX: &vFieldX,
	FieldY: &vFieldY,
	MyStructField: &MyStructFieldData,
	}

    log.Debugf("[getProjectProjectConfigResolver]Output object %+v", ret)
	return ret, nil
}
