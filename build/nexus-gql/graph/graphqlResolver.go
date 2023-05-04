package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"

	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	"golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-gql/graph/model"
)

var c = GrpcClients{
	mtx:     sync.Mutex{},
	Clients: map[string]GrpcClient{},
}
var nc *nexus_client.Clientset

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
// PKG: Api, NODE: Api
//////////////////////////////////////
func getRootResolver(id *string) ([]*model.ApiNexus, error) {
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

	var vNexusList []*model.ApiNexus
	if id != nil && *id != "" {
		log.Debugf("[getRootResolver]Id: %q", *id)
		vNexus, err := nc.GetApiNexus(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting Nexus node %q: %s", *id, err)
			return nil, nil
		}
		dn := vNexus.DisplayName()
		parentLabels := map[string]interface{}{"nexuses.api.nexus.vmware.com": dn}

		ret := &model.ApiNexus{
			Id:           &dn,
			ParentLabels: parentLabels,
		}
		vNexusList = append(vNexusList, ret)
		log.Debugf("[getRootResolver]Output Nexus objects %+v", vNexusList)
		return vNexusList, nil
	}

	log.Debugf("[getRootResolver]Id is empty, process all Nexuss")

	vNexusListObj, err := nc.Api().ListNexuses(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("[getRootResolver]Error getting Nexus node %s", err)
		return nil, nil
	}
	for _, i := range vNexusListObj {
		vNexus, err := nc.GetApiNexus(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRootResolver]Error getting Nexus node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vNexus.DisplayName()
		parentLabels := map[string]interface{}{"nexuses.api.nexus.vmware.com": dn}

		ret := &model.ApiNexus{
			Id:           &dn,
			ParentLabels: parentLabels,
		}
		vNexusList = append(vNexusList, ret)
	}

	log.Debugf("[getRootResolver]Output Nexus objects %v", vNexusList)
	return vNexusList, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: Config Node: Nexus PKG: Api
//////////////////////////////////////
func getApiNexusConfigResolver(obj *model.ApiNexus, id *string) (*model.ConfigConfig, error) {
	log.Debugf("[getApiNexusConfigResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getApiNexusConfigResolver]Id %q", *id)
		vConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetConfig(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApiNexusConfigResolver]Error getting Config node %q : %s", *id, err)
			return &model.ConfigConfig{}, nil
		}
		dn := vConfig.DisplayName()
		parentLabels := map[string]interface{}{"configs.config.nexus.vmware.com": dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConfigConfig{
			Id:           &dn,
			ParentLabels: parentLabels,
		}

		log.Debugf("[getApiNexusConfigResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getApiNexusConfigResolver]Id is empty, process all Configs")
	vConfigParent, err := nc.GetApiNexus(context.TODO(), getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getApiNexusConfigResolver]Failed to get parent node %s", err)
		return &model.ConfigConfig{}, nil
	}
	vConfig, err := vConfigParent.GetConfig(context.TODO())
	if err != nil {
		log.Errorf("[getApiNexusConfigResolver]Error getting Config node %s", err)
		return &model.ConfigConfig{}, nil
	}
	dn := vConfig.DisplayName()
	parentLabels := map[string]interface{}{"configs.config.nexus.vmware.com": dn}

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.ConfigConfig{
		Id:           &dn,
		ParentLabels: parentLabels,
	}

	log.Debugf("[getApiNexusConfigResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Singleton)
// FieldName: Runtime Node: Nexus PKG: Api
//////////////////////////////////////
func getApiNexusRuntimeResolver(obj *model.ApiNexus) (*model.RuntimeRuntime, error) {
	log.Debugf("[getApiNexusRuntimeResolver]Parent Object %+v", obj)
	vRuntime, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetRuntime(context.TODO())
	if err != nil {
		log.Errorf("[getApiNexusRuntimeResolver]Error getting Nexus node %s", err)
		return &model.RuntimeRuntime{}, nil
	}
	dn := vRuntime.DisplayName()
	parentLabels := map[string]interface{}{"runtimes.runtime.nexus.vmware.com": dn}

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.RuntimeRuntime{
		Id:           &dn,
		ParentLabels: parentLabels,
	}

	log.Debugf("[getApiNexusRuntimeResolver]Output object %+v", ret)
	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: Authn Node: ApiGateway PKG: Apigateway
//////////////////////////////////////
func getApigatewayApiGatewayAuthnResolver(obj *model.ApigatewayApiGateway, id *string) (*model.AuthenticationOIDC, error) {
	log.Debugf("[getApigatewayApiGatewayAuthnResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getApigatewayApiGatewayAuthnResolver]Id %q", *id)
		vOIDC, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com")).GetAuthn(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayAuthnResolver]Error getting Authn node %q : %s", *id, err)
			return &model.AuthenticationOIDC{}, nil
		}
		dn := vOIDC.DisplayName()
		parentLabels := map[string]interface{}{"oidcs.authentication.nexus.vmware.com": dn}
		Config, _ := json.Marshal(vOIDC.Spec.Config)
		ConfigData := string(Config)
		ValidationProps, _ := json.Marshal(vOIDC.Spec.ValidationProps)
		ValidationPropsData := string(ValidationProps)
		vJwtClaimUsername := string(vOIDC.Spec.JwtClaimUsername)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.AuthenticationOIDC{
			Id:               &dn,
			ParentLabels:     parentLabels,
			Config:           &ConfigData,
			ValidationProps:  &ValidationPropsData,
			JwtClaimUsername: &vJwtClaimUsername,
		}

		log.Debugf("[getApigatewayApiGatewayAuthnResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getApigatewayApiGatewayAuthnResolver]Id is empty, process all Authns")
	vOIDCParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetApiGateway(context.TODO(), getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getApigatewayApiGatewayAuthnResolver]Failed to get parent node %s", err)
		return &model.AuthenticationOIDC{}, nil
	}
	vOIDC, err := vOIDCParent.GetAuthn(context.TODO())
	if err != nil {
		log.Errorf("[getApigatewayApiGatewayAuthnResolver]Error getting Authn node %s", err)
		return &model.AuthenticationOIDC{}, nil
	}
	dn := vOIDC.DisplayName()
	parentLabels := map[string]interface{}{"oidcs.authentication.nexus.vmware.com": dn}
	Config, _ := json.Marshal(vOIDC.Spec.Config)
	ConfigData := string(Config)
	ValidationProps, _ := json.Marshal(vOIDC.Spec.ValidationProps)
	ValidationPropsData := string(ValidationProps)
	vJwtClaimUsername := string(vOIDC.Spec.JwtClaimUsername)

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.AuthenticationOIDC{
		Id:               &dn,
		ParentLabels:     parentLabels,
		Config:           &ConfigData,
		ValidationProps:  &ValidationPropsData,
		JwtClaimUsername: &vJwtClaimUsername,
	}

	log.Debugf("[getApigatewayApiGatewayAuthnResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: ProxyRules Node: ApiGateway PKG: Apigateway
//////////////////////////////////////
func getApigatewayApiGatewayProxyRulesResolver(obj *model.ApigatewayApiGateway, id *string) ([]*model.AdminProxyRule, error) {
	log.Debugf("[getApigatewayApiGatewayProxyRulesResolver]Parent Object %+v", obj)
	var vAdminProxyRuleList []*model.AdminProxyRule
	if id != nil && *id != "" {
		log.Debugf("[getApigatewayApiGatewayProxyRulesResolver]Id %q", *id)
		vProxyRule, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com")).GetProxyRules(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayProxyRulesResolver]Error getting ProxyRules node %q : %s", *id, err)
			return vAdminProxyRuleList, nil
		}
		dn := vProxyRule.DisplayName()
		parentLabels := map[string]interface{}{"proxyrules.admin.nexus.vmware.com": dn}
		MatchCondition, _ := json.Marshal(vProxyRule.Spec.MatchCondition)
		MatchConditionData := string(MatchCondition)
		Upstream, _ := json.Marshal(vProxyRule.Spec.Upstream)
		UpstreamData := string(Upstream)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.AdminProxyRule{
			Id:             &dn,
			ParentLabels:   parentLabels,
			MatchCondition: &MatchConditionData,
			Upstream:       &UpstreamData,
		}
		vAdminProxyRuleList = append(vAdminProxyRuleList, ret)

		log.Debugf("[getApigatewayApiGatewayProxyRulesResolver]Output ProxyRules objects %v", vAdminProxyRuleList)

		return vAdminProxyRuleList, nil
	}

	log.Debug("[getApigatewayApiGatewayProxyRulesResolver]Id is empty, process all ProxyRuless")

	vProxyRuleParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetApiGateway(context.TODO(), getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getApigatewayApiGatewayProxyRulesResolver]Error getting parent node %s", err)
		return vAdminProxyRuleList, nil
	}
	vProxyRuleAllObj, err := vProxyRuleParent.GetAllProxyRules(context.TODO())
	if err != nil {
		log.Errorf("[getApigatewayApiGatewayProxyRulesResolver]Error getting ProxyRules objects %s", err)
		return vAdminProxyRuleList, nil
	}
	for _, i := range vProxyRuleAllObj {
		vProxyRule, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com")).GetProxyRules(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayProxyRulesResolver]Error getting ProxyRules node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vProxyRule.DisplayName()
		parentLabels := map[string]interface{}{"proxyrules.admin.nexus.vmware.com": dn}
		MatchCondition, _ := json.Marshal(vProxyRule.Spec.MatchCondition)
		MatchConditionData := string(MatchCondition)
		Upstream, _ := json.Marshal(vProxyRule.Spec.Upstream)
		UpstreamData := string(Upstream)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.AdminProxyRule{
			Id:             &dn,
			ParentLabels:   parentLabels,
			MatchCondition: &MatchConditionData,
			Upstream:       &UpstreamData,
		}
		vAdminProxyRuleList = append(vAdminProxyRuleList, ret)
	}

	log.Debugf("[getApigatewayApiGatewayProxyRulesResolver]Output ProxyRules objects %v", vAdminProxyRuleList)

	return vAdminProxyRuleList, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: Cors Node: ApiGateway PKG: Apigateway
//////////////////////////////////////
func getApigatewayApiGatewayCorsResolver(obj *model.ApigatewayApiGateway, id *string) ([]*model.DomainCORSConfig, error) {
	log.Debugf("[getApigatewayApiGatewayCorsResolver]Parent Object %+v", obj)
	var vDomainCORSConfigList []*model.DomainCORSConfig
	if id != nil && *id != "" {
		log.Debugf("[getApigatewayApiGatewayCorsResolver]Id %q", *id)
		vCORSConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com")).GetCors(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayCorsResolver]Error getting Cors node %q : %s", *id, err)
			return vDomainCORSConfigList, nil
		}
		dn := vCORSConfig.DisplayName()
		parentLabels := map[string]interface{}{"corsconfigs.domain.nexus.vmware.com": dn}
		Origins, _ := json.Marshal(vCORSConfig.Spec.Origins)
		OriginsData := string(Origins)
		Headers, _ := json.Marshal(vCORSConfig.Spec.Headers)
		HeadersData := string(Headers)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.DomainCORSConfig{
			Id:           &dn,
			ParentLabels: parentLabels,
			Origins:      &OriginsData,
			Headers:      &HeadersData,
		}
		vDomainCORSConfigList = append(vDomainCORSConfigList, ret)

		log.Debugf("[getApigatewayApiGatewayCorsResolver]Output Cors objects %v", vDomainCORSConfigList)

		return vDomainCORSConfigList, nil
	}

	log.Debug("[getApigatewayApiGatewayCorsResolver]Id is empty, process all Corss")

	vCORSConfigParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetApiGateway(context.TODO(), getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getApigatewayApiGatewayCorsResolver]Error getting parent node %s", err)
		return vDomainCORSConfigList, nil
	}
	vCORSConfigAllObj, err := vCORSConfigParent.GetAllCors(context.TODO())
	if err != nil {
		log.Errorf("[getApigatewayApiGatewayCorsResolver]Error getting Cors objects %s", err)
		return vDomainCORSConfigList, nil
	}
	for _, i := range vCORSConfigAllObj {
		vCORSConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.vmware.com")).GetCors(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayCorsResolver]Error getting Cors node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vCORSConfig.DisplayName()
		parentLabels := map[string]interface{}{"corsconfigs.domain.nexus.vmware.com": dn}
		Origins, _ := json.Marshal(vCORSConfig.Spec.Origins)
		OriginsData := string(Origins)
		Headers, _ := json.Marshal(vCORSConfig.Spec.Headers)
		HeadersData := string(Headers)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.DomainCORSConfig{
			Id:           &dn,
			ParentLabels: parentLabels,
			Origins:      &OriginsData,
			Headers:      &HeadersData,
		}
		vDomainCORSConfigList = append(vDomainCORSConfigList, ret)
	}

	log.Debugf("[getApigatewayApiGatewayCorsResolver]Output Cors objects %v", vDomainCORSConfigList)

	return vDomainCORSConfigList, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: ApiGateway Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigApiGatewayResolver(obj *model.ConfigConfig, id *string) (*model.ApigatewayApiGateway, error) {
	log.Debugf("[getConfigConfigApiGatewayResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigApiGatewayResolver]Id %q", *id)
		vApiGateway, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetApiGateway(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigApiGatewayResolver]Error getting ApiGateway node %q : %s", *id, err)
			return &model.ApigatewayApiGateway{}, nil
		}
		dn := vApiGateway.DisplayName()
		parentLabels := map[string]interface{}{"apigateways.apigateway.nexus.vmware.com": dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ApigatewayApiGateway{
			Id:           &dn,
			ParentLabels: parentLabels,
		}

		log.Debugf("[getConfigConfigApiGatewayResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getConfigConfigApiGatewayResolver]Id is empty, process all ApiGateways")
	vApiGatewayParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigApiGatewayResolver]Failed to get parent node %s", err)
		return &model.ApigatewayApiGateway{}, nil
	}
	vApiGateway, err := vApiGatewayParent.GetApiGateway(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigApiGatewayResolver]Error getting ApiGateway node %s", err)
		return &model.ApigatewayApiGateway{}, nil
	}
	dn := vApiGateway.DisplayName()
	parentLabels := map[string]interface{}{"apigateways.apigateway.nexus.vmware.com": dn}

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.ApigatewayApiGateway{
		Id:           &dn,
		ParentLabels: parentLabels,
	}

	log.Debugf("[getConfigConfigApiGatewayResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILD RESOLVER (Non Singleton)
// FieldName: Connect Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigConnectResolver(obj *model.ConfigConfig, id *string) (*model.ConnectConnect, error) {
	log.Debugf("[getConfigConfigConnectResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigConnectResolver]Id %q", *id)
		vConnect, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetConnect(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigConnectResolver]Error getting Connect node %q : %s", *id, err)
			return &model.ConnectConnect{}, nil
		}
		dn := vConnect.DisplayName()
		parentLabels := map[string]interface{}{"connects.connect.nexus.vmware.com": dn}

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConnectConnect{
			Id:           &dn,
			ParentLabels: parentLabels,
		}

		log.Debugf("[getConfigConfigConnectResolver]Output object %v", ret)
		return ret, nil
	}
	log.Debug("[getConfigConfigConnectResolver]Id is empty, process all Connects")
	vConnectParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigConnectResolver]Failed to get parent node %s", err)
		return &model.ConnectConnect{}, nil
	}
	vConnect, err := vConnectParent.GetConnect(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigConnectResolver]Error getting Connect node %s", err)
		return &model.ConnectConnect{}, nil
	}
	dn := vConnect.DisplayName()
	parentLabels := map[string]interface{}{"connects.connect.nexus.vmware.com": dn}

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.ConnectConnect{
		Id:           &dn,
		ParentLabels: parentLabels,
	}

	log.Debugf("[getConfigConfigConnectResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: Routes Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigRoutesResolver(obj *model.ConfigConfig, id *string) ([]*model.RouteRoute, error) {
	log.Debugf("[getConfigConfigRoutesResolver]Parent Object %+v", obj)
	var vRouteRouteList []*model.RouteRoute
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigRoutesResolver]Id %q", *id)
		vRoute, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetRoutes(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigRoutesResolver]Error getting Routes node %q : %s", *id, err)
			return vRouteRouteList, nil
		}
		dn := vRoute.DisplayName()
		parentLabels := map[string]interface{}{"routes.route.nexus.vmware.com": dn}
		vUri := string(vRoute.Spec.Uri)
		Service, _ := json.Marshal(vRoute.Spec.Service)
		ServiceData := string(Service)
		Resource, _ := json.Marshal(vRoute.Spec.Resource)
		ResourceData := string(Resource)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.RouteRoute{
			Id:           &dn,
			ParentLabels: parentLabels,
			Uri:          &vUri,
			Service:      &ServiceData,
			Resource:     &ResourceData,
		}
		vRouteRouteList = append(vRouteRouteList, ret)

		log.Debugf("[getConfigConfigRoutesResolver]Output Routes objects %v", vRouteRouteList)

		return vRouteRouteList, nil
	}

	log.Debug("[getConfigConfigRoutesResolver]Id is empty, process all Routess")

	vRouteParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigRoutesResolver]Error getting parent node %s", err)
		return vRouteRouteList, nil
	}
	vRouteAllObj, err := vRouteParent.GetAllRoutes(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigRoutesResolver]Error getting Routes objects %s", err)
		return vRouteRouteList, nil
	}
	for _, i := range vRouteAllObj {
		vRoute, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetRoutes(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConfigConfigRoutesResolver]Error getting Routes node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vRoute.DisplayName()
		parentLabels := map[string]interface{}{"routes.route.nexus.vmware.com": dn}
		vUri := string(vRoute.Spec.Uri)
		Service, _ := json.Marshal(vRoute.Spec.Service)
		ServiceData := string(Service)
		Resource, _ := json.Marshal(vRoute.Spec.Resource)
		ResourceData := string(Resource)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.RouteRoute{
			Id:           &dn,
			ParentLabels: parentLabels,
			Uri:          &vUri,
			Service:      &ServiceData,
			Resource:     &ResourceData,
		}
		vRouteRouteList = append(vRouteRouteList, ret)
	}

	log.Debugf("[getConfigConfigRoutesResolver]Output Routes objects %v", vRouteRouteList)

	return vRouteRouteList, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: Tenant Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigTenantResolver(obj *model.ConfigConfig, id *string) ([]*model.TenantconfigTenant, error) {
	log.Debugf("[getConfigConfigTenantResolver]Parent Object %+v", obj)
	var vTenantconfigTenantList []*model.TenantconfigTenant
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigTenantResolver]Id %q", *id)
		vTenant, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetTenant(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigTenantResolver]Error getting Tenant node %q : %s", *id, err)
			return vTenantconfigTenantList, nil
		}
		dn := vTenant.DisplayName()
		parentLabels := map[string]interface{}{"tenants.tenantconfig.nexus.vmware.com": dn}
		vName := string(vTenant.Spec.Name)
		vDNSSuffix := string(vTenant.Spec.DNSSuffix)
		vSkipSaasTlsVerify := bool(vTenant.Spec.SkipSaasTlsVerify)
		vInstallTenant := bool(vTenant.Spec.InstallTenant)
		vInstallClient := bool(vTenant.Spec.InstallClient)
		vOrderId := string(vTenant.Spec.OrderId)
		Skus, _ := json.Marshal(vTenant.Spec.Skus)
		SkusData := string(Skus)
		FeatureFlags, _ := json.Marshal(vTenant.Spec.FeatureFlags)
		FeatureFlagsData := string(FeatureFlags)
		Labels, _ := json.Marshal(vTenant.Spec.Labels)
		LabelsData := string(Labels)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.TenantconfigTenant{
			Id:                &dn,
			ParentLabels:      parentLabels,
			Name:              &vName,
			DNSSuffix:         &vDNSSuffix,
			SkipSaasTlsVerify: &vSkipSaasTlsVerify,
			InstallTenant:     &vInstallTenant,
			InstallClient:     &vInstallClient,
			OrderId:           &vOrderId,
			Skus:              &SkusData,
			FeatureFlags:      &FeatureFlagsData,
			Labels:            &LabelsData,
		}
		vTenantconfigTenantList = append(vTenantconfigTenantList, ret)

		log.Debugf("[getConfigConfigTenantResolver]Output Tenant objects %v", vTenantconfigTenantList)

		return vTenantconfigTenantList, nil
	}

	log.Debug("[getConfigConfigTenantResolver]Id is empty, process all Tenants")

	vTenantParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigTenantResolver]Error getting parent node %s", err)
		return vTenantconfigTenantList, nil
	}
	vTenantAllObj, err := vTenantParent.GetAllTenant(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigTenantResolver]Error getting Tenant objects %s", err)
		return vTenantconfigTenantList, nil
	}
	for _, i := range vTenantAllObj {
		vTenant, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetTenant(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConfigConfigTenantResolver]Error getting Tenant node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vTenant.DisplayName()
		parentLabels := map[string]interface{}{"tenants.tenantconfig.nexus.vmware.com": dn}
		vName := string(vTenant.Spec.Name)
		vDNSSuffix := string(vTenant.Spec.DNSSuffix)
		vSkipSaasTlsVerify := bool(vTenant.Spec.SkipSaasTlsVerify)
		vInstallTenant := bool(vTenant.Spec.InstallTenant)
		vInstallClient := bool(vTenant.Spec.InstallClient)
		vOrderId := string(vTenant.Spec.OrderId)
		Skus, _ := json.Marshal(vTenant.Spec.Skus)
		SkusData := string(Skus)
		FeatureFlags, _ := json.Marshal(vTenant.Spec.FeatureFlags)
		FeatureFlagsData := string(FeatureFlags)
		Labels, _ := json.Marshal(vTenant.Spec.Labels)
		LabelsData := string(Labels)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.TenantconfigTenant{
			Id:                &dn,
			ParentLabels:      parentLabels,
			Name:              &vName,
			DNSSuffix:         &vDNSSuffix,
			SkipSaasTlsVerify: &vSkipSaasTlsVerify,
			InstallTenant:     &vInstallTenant,
			InstallClient:     &vInstallClient,
			OrderId:           &vOrderId,
			Skus:              &SkusData,
			FeatureFlags:      &FeatureFlagsData,
			Labels:            &LabelsData,
		}
		vTenantconfigTenantList = append(vTenantconfigTenantList, ret)
	}

	log.Debugf("[getConfigConfigTenantResolver]Output Tenant objects %v", vTenantconfigTenantList)

	return vTenantconfigTenantList, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: TenantPolicy Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigTenantPolicyResolver(obj *model.ConfigConfig, id *string) ([]*model.TenantconfigPolicy, error) {
	log.Debugf("[getConfigConfigTenantPolicyResolver]Parent Object %+v", obj)
	var vTenantconfigPolicyList []*model.TenantconfigPolicy
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigTenantPolicyResolver]Id %q", *id)
		vPolicy, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetTenantPolicy(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigTenantPolicyResolver]Error getting TenantPolicy node %q : %s", *id, err)
			return vTenantconfigPolicyList, nil
		}
		dn := vPolicy.DisplayName()
		parentLabels := map[string]interface{}{"policies.tenantconfig.nexus.vmware.com": dn}
		StaticApplications, _ := json.Marshal(vPolicy.Spec.StaticApplications)
		StaticApplicationsData := string(StaticApplications)
		PinApplications, _ := json.Marshal(vPolicy.Spec.PinApplications)
		PinApplicationsData := string(PinApplications)
		vDynamicAppSchedulingDisable := bool(vPolicy.Spec.DynamicAppSchedulingDisable)
		vDisableProvisioning := bool(vPolicy.Spec.DisableProvisioning)
		vDisableAutoScaling := bool(vPolicy.Spec.DisableAutoScaling)
		vDisableAppClusterOnboarding := bool(vPolicy.Spec.DisableAppClusterOnboarding)
		vDisableUpgrade := bool(vPolicy.Spec.DisableUpgrade)
		vOnFailureDowngrade := bool(vPolicy.Spec.OnFailureDowngrade)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.TenantconfigPolicy{
			Id:                          &dn,
			ParentLabels:                parentLabels,
			StaticApplications:          &StaticApplicationsData,
			PinApplications:             &PinApplicationsData,
			DynamicAppSchedulingDisable: &vDynamicAppSchedulingDisable,
			DisableProvisioning:         &vDisableProvisioning,
			DisableAutoScaling:          &vDisableAutoScaling,
			DisableAppClusterOnboarding: &vDisableAppClusterOnboarding,
			DisableUpgrade:              &vDisableUpgrade,
			OnFailureDowngrade:          &vOnFailureDowngrade,
		}
		vTenantconfigPolicyList = append(vTenantconfigPolicyList, ret)

		log.Debugf("[getConfigConfigTenantPolicyResolver]Output TenantPolicy objects %v", vTenantconfigPolicyList)

		return vTenantconfigPolicyList, nil
	}

	log.Debug("[getConfigConfigTenantPolicyResolver]Id is empty, process all TenantPolicys")

	vPolicyParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigTenantPolicyResolver]Error getting parent node %s", err)
		return vTenantconfigPolicyList, nil
	}
	vPolicyAllObj, err := vPolicyParent.GetAllTenantPolicy(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigTenantPolicyResolver]Error getting TenantPolicy objects %s", err)
		return vTenantconfigPolicyList, nil
	}
	for _, i := range vPolicyAllObj {
		vPolicy, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetTenantPolicy(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConfigConfigTenantPolicyResolver]Error getting TenantPolicy node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vPolicy.DisplayName()
		parentLabels := map[string]interface{}{"policies.tenantconfig.nexus.vmware.com": dn}
		StaticApplications, _ := json.Marshal(vPolicy.Spec.StaticApplications)
		StaticApplicationsData := string(StaticApplications)
		PinApplications, _ := json.Marshal(vPolicy.Spec.PinApplications)
		PinApplicationsData := string(PinApplications)
		vDynamicAppSchedulingDisable := bool(vPolicy.Spec.DynamicAppSchedulingDisable)
		vDisableProvisioning := bool(vPolicy.Spec.DisableProvisioning)
		vDisableAutoScaling := bool(vPolicy.Spec.DisableAutoScaling)
		vDisableAppClusterOnboarding := bool(vPolicy.Spec.DisableAppClusterOnboarding)
		vDisableUpgrade := bool(vPolicy.Spec.DisableUpgrade)
		vOnFailureDowngrade := bool(vPolicy.Spec.OnFailureDowngrade)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.TenantconfigPolicy{
			Id:                          &dn,
			ParentLabels:                parentLabels,
			StaticApplications:          &StaticApplicationsData,
			PinApplications:             &PinApplicationsData,
			DynamicAppSchedulingDisable: &vDynamicAppSchedulingDisable,
			DisableProvisioning:         &vDisableProvisioning,
			DisableAutoScaling:          &vDisableAutoScaling,
			DisableAppClusterOnboarding: &vDisableAppClusterOnboarding,
			DisableUpgrade:              &vDisableUpgrade,
			OnFailureDowngrade:          &vOnFailureDowngrade,
		}
		vTenantconfigPolicyList = append(vTenantconfigPolicyList, ret)
	}

	log.Debugf("[getConfigConfigTenantPolicyResolver]Output TenantPolicy objects %v", vTenantconfigPolicyList)

	return vTenantconfigPolicyList, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: User Node: Config PKG: Config
//////////////////////////////////////
func getConfigConfigUserResolver(obj *model.ConfigConfig, id *string) ([]*model.UserUser, error) {
	log.Debugf("[getConfigConfigUserResolver]Parent Object %+v", obj)
	var vUserUserList []*model.UserUser
	if id != nil && *id != "" {
		log.Debugf("[getConfigConfigUserResolver]Id %q", *id)
		vUser, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetUser(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigUserResolver]Error getting User node %q : %s", *id, err)
			return vUserUserList, nil
		}
		dn := vUser.DisplayName()
		parentLabels := map[string]interface{}{"users.user.nexus.vmware.com": dn}
		vUsername := string(vUser.Spec.Username)
		vMail := string(vUser.Spec.Mail)
		vFirstName := string(vUser.Spec.FirstName)
		vLastName := string(vUser.Spec.LastName)
		vPassword := string(vUser.Spec.Password)
		vTenantId := string(vUser.Spec.TenantId)
		vRealm := string(vUser.Spec.Realm)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.UserUser{
			Id:           &dn,
			ParentLabels: parentLabels,
			Username:     &vUsername,
			Mail:         &vMail,
			FirstName:    &vFirstName,
			LastName:     &vLastName,
			Password:     &vPassword,
			TenantId:     &vTenantId,
			Realm:        &vRealm,
		}
		vUserUserList = append(vUserUserList, ret)

		log.Debugf("[getConfigConfigUserResolver]Output User objects %v", vUserUserList)

		return vUserUserList, nil
	}

	log.Debug("[getConfigConfigUserResolver]Id is empty, process all Users")

	vUserParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConfigConfigUserResolver]Error getting parent node %s", err)
		return vUserUserList, nil
	}
	vUserAllObj, err := vUserParent.GetAllUser(context.TODO())
	if err != nil {
		log.Errorf("[getConfigConfigUserResolver]Error getting User objects %s", err)
		return vUserUserList, nil
	}
	for _, i := range vUserAllObj {
		vUser, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetUser(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConfigConfigUserResolver]Error getting User node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vUser.DisplayName()
		parentLabels := map[string]interface{}{"users.user.nexus.vmware.com": dn}
		vUsername := string(vUser.Spec.Username)
		vMail := string(vUser.Spec.Mail)
		vFirstName := string(vUser.Spec.FirstName)
		vLastName := string(vUser.Spec.LastName)
		vPassword := string(vUser.Spec.Password)
		vTenantId := string(vUser.Spec.TenantId)
		vRealm := string(vUser.Spec.Realm)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.UserUser{
			Id:           &dn,
			ParentLabels: parentLabels,
			Username:     &vUsername,
			Mail:         &vMail,
			FirstName:    &vFirstName,
			LastName:     &vLastName,
			Password:     &vPassword,
			TenantId:     &vTenantId,
			Realm:        &vRealm,
		}
		vUserUserList = append(vUserUserList, ret)
	}

	log.Debugf("[getConfigConfigUserResolver]Output User objects %v", vUserUserList)

	return vUserUserList, nil
}

//////////////////////////////////////
// LINK RESOLVER
// FieldName: Tenant Node: User PKG: User
//////////////////////////////////////
func getUserUserTenantResolver(obj *model.UserUser) (*model.TenantconfigTenant, error) {
	log.Debugf("[getUserUserTenantResolver]Parent Object %+v", obj)
	vTenantParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetUser(context.TODO(), getParentName(obj.ParentLabels, "users.user.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getUserUserTenantResolver]Error getting parent node %s", err)
		return &model.TenantconfigTenant{}, nil
	}
	vTenant, err := vTenantParent.GetTenant(context.TODO())
	if err != nil {
		log.Errorf("[getUserUserTenantResolver]Error getting Tenant object %s", err)
		return &model.TenantconfigTenant{}, nil
	}
	dn := vTenant.DisplayName()
	parentLabels := map[string]interface{}{"tenants.tenantconfig.nexus.vmware.com": dn}
	vName := string(vTenant.Spec.Name)
	vDNSSuffix := string(vTenant.Spec.DNSSuffix)
	vSkipSaasTlsVerify := bool(vTenant.Spec.SkipSaasTlsVerify)
	vInstallTenant := bool(vTenant.Spec.InstallTenant)
	vInstallClient := bool(vTenant.Spec.InstallClient)
	vOrderId := string(vTenant.Spec.OrderId)
	Skus, _ := json.Marshal(vTenant.Spec.Skus)
	SkusData := string(Skus)
	FeatureFlags, _ := json.Marshal(vTenant.Spec.FeatureFlags)
	FeatureFlagsData := string(FeatureFlags)
	Labels, _ := json.Marshal(vTenant.Spec.Labels)
	LabelsData := string(Labels)

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.TenantconfigTenant{
		Id:                &dn,
		ParentLabels:      parentLabels,
		Name:              &vName,
		DNSSuffix:         &vDNSSuffix,
		SkipSaasTlsVerify: &vSkipSaasTlsVerify,
		InstallTenant:     &vInstallTenant,
		InstallClient:     &vInstallClient,
		OrderId:           &vOrderId,
		Skus:              &SkusData,
		FeatureFlags:      &FeatureFlagsData,
		Labels:            &LabelsData,
	}
	log.Debugf("[getUserUserTenantResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: Endpoints Node: Connect PKG: Connect
//////////////////////////////////////
func getConnectConnectEndpointsResolver(obj *model.ConnectConnect, id *string) ([]*model.ConnectNexusEndpoint, error) {
	log.Debugf("[getConnectConnectEndpointsResolver]Parent Object %+v", obj)
	var vConnectNexusEndpointList []*model.ConnectNexusEndpoint
	if id != nil && *id != "" {
		log.Debugf("[getConnectConnectEndpointsResolver]Id %q", *id)
		vNexusEndpoint, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.vmware.com")).GetEndpoints(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConnectConnectEndpointsResolver]Error getting Endpoints node %q : %s", *id, err)
			return vConnectNexusEndpointList, nil
		}
		dn := vNexusEndpoint.DisplayName()
		parentLabels := map[string]interface{}{"nexusendpoints.connect.nexus.vmware.com": dn}
		vHost := string(vNexusEndpoint.Spec.Host)
		vPort := string(vNexusEndpoint.Spec.Port)
		vCert := string(vNexusEndpoint.Spec.Cert)
		vPath := string(vNexusEndpoint.Spec.Path)
		Cloud, _ := json.Marshal(vNexusEndpoint.Spec.Cloud)
		CloudData := string(Cloud)
		vServiceAccountName := string(vNexusEndpoint.Spec.ServiceAccountName)
		vClientName := string(vNexusEndpoint.Spec.ClientName)
		vClientRegion := string(vNexusEndpoint.Spec.ClientRegion)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConnectNexusEndpoint{
			Id:                 &dn,
			ParentLabels:       parentLabels,
			Host:               &vHost,
			Port:               &vPort,
			Cert:               &vCert,
			Path:               &vPath,
			Cloud:              &CloudData,
			ServiceAccountName: &vServiceAccountName,
			ClientName:         &vClientName,
			ClientRegion:       &vClientRegion,
		}
		vConnectNexusEndpointList = append(vConnectNexusEndpointList, ret)

		log.Debugf("[getConnectConnectEndpointsResolver]Output Endpoints objects %v", vConnectNexusEndpointList)

		return vConnectNexusEndpointList, nil
	}

	log.Debug("[getConnectConnectEndpointsResolver]Id is empty, process all Endpointss")

	vNexusEndpointParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetConnect(context.TODO(), getParentName(obj.ParentLabels, "connects.connect.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConnectConnectEndpointsResolver]Error getting parent node %s", err)
		return vConnectNexusEndpointList, nil
	}
	vNexusEndpointAllObj, err := vNexusEndpointParent.GetAllEndpoints(context.TODO())
	if err != nil {
		log.Errorf("[getConnectConnectEndpointsResolver]Error getting Endpoints objects %s", err)
		return vConnectNexusEndpointList, nil
	}
	for _, i := range vNexusEndpointAllObj {
		vNexusEndpoint, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.vmware.com")).GetEndpoints(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConnectConnectEndpointsResolver]Error getting Endpoints node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vNexusEndpoint.DisplayName()
		parentLabels := map[string]interface{}{"nexusendpoints.connect.nexus.vmware.com": dn}
		vHost := string(vNexusEndpoint.Spec.Host)
		vPort := string(vNexusEndpoint.Spec.Port)
		vCert := string(vNexusEndpoint.Spec.Cert)
		vPath := string(vNexusEndpoint.Spec.Path)
		Cloud, _ := json.Marshal(vNexusEndpoint.Spec.Cloud)
		CloudData := string(Cloud)
		vServiceAccountName := string(vNexusEndpoint.Spec.ServiceAccountName)
		vClientName := string(vNexusEndpoint.Spec.ClientName)
		vClientRegion := string(vNexusEndpoint.Spec.ClientRegion)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConnectNexusEndpoint{
			Id:                 &dn,
			ParentLabels:       parentLabels,
			Host:               &vHost,
			Port:               &vPort,
			Cert:               &vCert,
			Path:               &vPath,
			Cloud:              &CloudData,
			ServiceAccountName: &vServiceAccountName,
			ClientName:         &vClientName,
			ClientRegion:       &vClientRegion,
		}
		vConnectNexusEndpointList = append(vConnectNexusEndpointList, ret)
	}

	log.Debugf("[getConnectConnectEndpointsResolver]Output Endpoints objects %v", vConnectNexusEndpointList)

	return vConnectNexusEndpointList, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: ReplicationConfig Node: Connect PKG: Connect
//////////////////////////////////////
func getConnectConnectReplicationConfigResolver(obj *model.ConnectConnect, id *string) ([]*model.ConnectReplicationConfig, error) {
	log.Debugf("[getConnectConnectReplicationConfigResolver]Parent Object %+v", obj)
	var vConnectReplicationConfigList []*model.ConnectReplicationConfig
	if id != nil && *id != "" {
		log.Debugf("[getConnectConnectReplicationConfigResolver]Id %q", *id)
		vReplicationConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.vmware.com")).GetReplicationConfig(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConnectConnectReplicationConfigResolver]Error getting ReplicationConfig node %q : %s", *id, err)
			return vConnectReplicationConfigList, nil
		}
		dn := vReplicationConfig.DisplayName()
		parentLabels := map[string]interface{}{"replicationconfigs.connect.nexus.vmware.com": dn}
		vAccessToken := string(vReplicationConfig.Spec.AccessToken)
		Source, _ := json.Marshal(vReplicationConfig.Spec.Source)
		SourceData := string(Source)
		Destination, _ := json.Marshal(vReplicationConfig.Spec.Destination)
		DestinationData := string(Destination)
		StatusEndpoint, _ := json.Marshal(vReplicationConfig.Spec.StatusEndpoint)
		StatusEndpointData := string(StatusEndpoint)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConnectReplicationConfig{
			Id:             &dn,
			ParentLabels:   parentLabels,
			AccessToken:    &vAccessToken,
			Source:         &SourceData,
			Destination:    &DestinationData,
			StatusEndpoint: &StatusEndpointData,
		}
		vConnectReplicationConfigList = append(vConnectReplicationConfigList, ret)

		log.Debugf("[getConnectConnectReplicationConfigResolver]Output ReplicationConfig objects %v", vConnectReplicationConfigList)

		return vConnectReplicationConfigList, nil
	}

	log.Debug("[getConnectConnectReplicationConfigResolver]Id is empty, process all ReplicationConfigs")

	vReplicationConfigParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).GetConnect(context.TODO(), getParentName(obj.ParentLabels, "connects.connect.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConnectConnectReplicationConfigResolver]Error getting parent node %s", err)
		return vConnectReplicationConfigList, nil
	}
	vReplicationConfigAllObj, err := vReplicationConfigParent.GetAllReplicationConfig(context.TODO())
	if err != nil {
		log.Errorf("[getConnectConnectReplicationConfigResolver]Error getting ReplicationConfig objects %s", err)
		return vConnectReplicationConfigList, nil
	}
	for _, i := range vReplicationConfigAllObj {
		vReplicationConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.vmware.com")).GetReplicationConfig(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConnectConnectReplicationConfigResolver]Error getting ReplicationConfig node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vReplicationConfig.DisplayName()
		parentLabels := map[string]interface{}{"replicationconfigs.connect.nexus.vmware.com": dn}
		vAccessToken := string(vReplicationConfig.Spec.AccessToken)
		Source, _ := json.Marshal(vReplicationConfig.Spec.Source)
		SourceData := string(Source)
		Destination, _ := json.Marshal(vReplicationConfig.Spec.Destination)
		DestinationData := string(Destination)
		StatusEndpoint, _ := json.Marshal(vReplicationConfig.Spec.StatusEndpoint)
		StatusEndpointData := string(StatusEndpoint)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.ConnectReplicationConfig{
			Id:             &dn,
			ParentLabels:   parentLabels,
			AccessToken:    &vAccessToken,
			Source:         &SourceData,
			Destination:    &DestinationData,
			StatusEndpoint: &StatusEndpointData,
		}
		vConnectReplicationConfigList = append(vConnectReplicationConfigList, ret)
	}

	log.Debugf("[getConnectConnectReplicationConfigResolver]Output ReplicationConfig objects %v", vConnectReplicationConfigList)

	return vConnectReplicationConfigList, nil
}

//////////////////////////////////////
// LINK RESOLVER
// FieldName: RemoteEndpoint Node: ReplicationConfig PKG: Connect
//////////////////////////////////////
func getConnectReplicationConfigRemoteEndpointResolver(obj *model.ConnectReplicationConfig) (*model.ConnectNexusEndpoint, error) {
	log.Debugf("[getConnectReplicationConfigRemoteEndpointResolver]Parent Object %+v", obj)
	vNexusEndpointParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.vmware.com")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.vmware.com")).GetReplicationConfig(context.TODO(), getParentName(obj.ParentLabels, "replicationconfigs.connect.nexus.vmware.com"))
	if err != nil {
		log.Errorf("[getConnectReplicationConfigRemoteEndpointResolver]Error getting parent node %s", err)
		return &model.ConnectNexusEndpoint{}, nil
	}
	vNexusEndpoint, err := vNexusEndpointParent.GetRemoteEndpoint(context.TODO())
	if err != nil {
		log.Errorf("[getConnectReplicationConfigRemoteEndpointResolver]Error getting RemoteEndpoint object %s", err)
		return &model.ConnectNexusEndpoint{}, nil
	}
	dn := vNexusEndpoint.DisplayName()
	parentLabels := map[string]interface{}{"nexusendpoints.connect.nexus.vmware.com": dn}
	vHost := string(vNexusEndpoint.Spec.Host)
	vPort := string(vNexusEndpoint.Spec.Port)
	vCert := string(vNexusEndpoint.Spec.Cert)
	vPath := string(vNexusEndpoint.Spec.Path)
	Cloud, _ := json.Marshal(vNexusEndpoint.Spec.Cloud)
	CloudData := string(Cloud)
	vServiceAccountName := string(vNexusEndpoint.Spec.ServiceAccountName)
	vClientName := string(vNexusEndpoint.Spec.ClientName)
	vClientRegion := string(vNexusEndpoint.Spec.ClientRegion)

	for k, v := range obj.ParentLabels {
		parentLabels[k] = v
	}
	ret := &model.ConnectNexusEndpoint{
		Id:                 &dn,
		ParentLabels:       parentLabels,
		Host:               &vHost,
		Port:               &vPort,
		Cert:               &vCert,
		Path:               &vPath,
		Cloud:              &CloudData,
		ServiceAccountName: &vServiceAccountName,
		ClientName:         &vClientName,
		ClientRegion:       &vClientRegion,
	}
	log.Debugf("[getConnectReplicationConfigRemoteEndpointResolver]Output object %v", ret)

	return ret, nil
}

//////////////////////////////////////
// CHILDREN RESOLVER
// FieldName: Tenant Node: Runtime PKG: Runtime
//////////////////////////////////////
func getRuntimeRuntimeTenantResolver(obj *model.RuntimeRuntime, id *string) ([]*model.TenantruntimeTenant, error) {
	log.Debugf("[getRuntimeRuntimeTenantResolver]Parent Object %+v", obj)
	var vTenantruntimeTenantList []*model.TenantruntimeTenant
	if id != nil && *id != "" {
		log.Debugf("[getRuntimeRuntimeTenantResolver]Id %q", *id)
		vTenant, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Runtime().GetTenant(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRuntimeRuntimeTenantResolver]Error getting Tenant node %q : %s", *id, err)
			return vTenantruntimeTenantList, nil
		}
		dn := vTenant.DisplayName()
		parentLabels := map[string]interface{}{"tenants.tenantruntime.nexus.vmware.com": dn}
		vNamespace := string(vTenant.Spec.Namespace)
		vTenantName := string(vTenant.Spec.TenantName)
		Attributes, _ := json.Marshal(vTenant.Spec.Attributes)
		AttributesData := string(Attributes)
		vSaasDomainName := string(vTenant.Spec.SaasDomainName)
		vSaasApiDomainName := string(vTenant.Spec.SaasApiDomainName)
		vM7Enabled := string(vTenant.Spec.M7Enabled)
		vLicenseType := string(vTenant.Spec.LicenseType)
		vStreamName := string(vTenant.Spec.StreamName)
		vAwsS3Bucket := string(vTenant.Spec.AwsS3Bucket)
		vAwsKmsKeyId := string(vTenant.Spec.AwsKmsKeyId)
		vM7InstallationScheduled := string(vTenant.Spec.M7InstallationScheduled)
		vReleaseVersion := string(vTenant.Spec.ReleaseVersion)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.TenantruntimeTenant{
			Id:                      &dn,
			ParentLabels:            parentLabels,
			Namespace:               &vNamespace,
			TenantName:              &vTenantName,
			Attributes:              &AttributesData,
			SaasDomainName:          &vSaasDomainName,
			SaasApiDomainName:       &vSaasApiDomainName,
			M7Enabled:               &vM7Enabled,
			LicenseType:             &vLicenseType,
			StreamName:              &vStreamName,
			AwsS3Bucket:             &vAwsS3Bucket,
			AwsKmsKeyId:             &vAwsKmsKeyId,
			M7InstallationScheduled: &vM7InstallationScheduled,
			ReleaseVersion:          &vReleaseVersion,
		}
		vTenantruntimeTenantList = append(vTenantruntimeTenantList, ret)

		log.Debugf("[getRuntimeRuntimeTenantResolver]Output Tenant objects %v", vTenantruntimeTenantList)

		return vTenantruntimeTenantList, nil
	}

	log.Debug("[getRuntimeRuntimeTenantResolver]Id is empty, process all Tenants")

	vTenantParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).GetRuntime(context.TODO())
	if err != nil {
		log.Errorf("[getRuntimeRuntimeTenantResolver]Error getting parent node %s", err)
		return vTenantruntimeTenantList, nil
	}
	vTenantAllObj, err := vTenantParent.GetAllTenant(context.TODO())
	if err != nil {
		log.Errorf("[getRuntimeRuntimeTenantResolver]Error getting Tenant objects %s", err)
		return vTenantruntimeTenantList, nil
	}
	for _, i := range vTenantAllObj {
		vTenant, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.vmware.com")).Runtime().GetTenant(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getRuntimeRuntimeTenantResolver]Error getting Tenant node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vTenant.DisplayName()
		parentLabels := map[string]interface{}{"tenants.tenantruntime.nexus.vmware.com": dn}
		vNamespace := string(vTenant.Spec.Namespace)
		vTenantName := string(vTenant.Spec.TenantName)
		Attributes, _ := json.Marshal(vTenant.Spec.Attributes)
		AttributesData := string(Attributes)
		vSaasDomainName := string(vTenant.Spec.SaasDomainName)
		vSaasApiDomainName := string(vTenant.Spec.SaasApiDomainName)
		vM7Enabled := string(vTenant.Spec.M7Enabled)
		vLicenseType := string(vTenant.Spec.LicenseType)
		vStreamName := string(vTenant.Spec.StreamName)
		vAwsS3Bucket := string(vTenant.Spec.AwsS3Bucket)
		vAwsKmsKeyId := string(vTenant.Spec.AwsKmsKeyId)
		vM7InstallationScheduled := string(vTenant.Spec.M7InstallationScheduled)
		vReleaseVersion := string(vTenant.Spec.ReleaseVersion)

		for k, v := range obj.ParentLabels {
			parentLabels[k] = v
		}
		ret := &model.TenantruntimeTenant{
			Id:                      &dn,
			ParentLabels:            parentLabels,
			Namespace:               &vNamespace,
			TenantName:              &vTenantName,
			Attributes:              &AttributesData,
			SaasDomainName:          &vSaasDomainName,
			SaasApiDomainName:       &vSaasApiDomainName,
			M7Enabled:               &vM7Enabled,
			LicenseType:             &vLicenseType,
			StreamName:              &vStreamName,
			AwsS3Bucket:             &vAwsS3Bucket,
			AwsKmsKeyId:             &vAwsKmsKeyId,
			M7InstallationScheduled: &vM7InstallationScheduled,
			ReleaseVersion:          &vReleaseVersion,
		}
		vTenantruntimeTenantList = append(vTenantruntimeTenantList, ret)
	}

	log.Debugf("[getRuntimeRuntimeTenantResolver]Output Tenant objects %v", vTenantruntimeTenantList)

	return vTenantruntimeTenantList, nil
}
