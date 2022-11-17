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
	}
	k8sApiConfig := getK8sAPIEndpointConfig()
	nexusClient, err := nexus_client.NewForConfig(k8sApiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client config: %s", err)
	}
	nc = nexusClient
	var vNexusList []*model.ApiNexus
	if id != nil && *id != "" {
		log.Debugf("[getRootResolver]Id: %q", *id)
		vNexus, err := nc.GetApiNexus(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getRootResolver]Error getting Nexus node %q: %s", *id, err)
			return nil, nil
		}
		dn := vNexus.DisplayName()
		parentLabels := map[string]interface{}{"nexuses.api.nexus.org": dn}

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
		parentLabels := map[string]interface{}{"nexuses.api.nexus.org": dn}

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
		vConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).GetConfig(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApiNexusConfigResolver]Error getting Config node %q : %s", *id, err)
			return &model.ConfigConfig{}, nil
		}
		dn := vConfig.DisplayName()
		parentLabels := map[string]interface{}{"configs.config.nexus.org": dn}

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
	vConfigParent, err := nc.GetApiNexus(context.TODO(), getParentName(obj.ParentLabels, "nexuses.api.nexus.org"))
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
	parentLabels := map[string]interface{}{"configs.config.nexus.org": dn}

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
// CHILD RESOLVER (Non Singleton)
// FieldName: Authn Node: ApiGateway PKG: Apigateway
//////////////////////////////////////
func getApigatewayApiGatewayAuthnResolver(obj *model.ApigatewayApiGateway, id *string) (*model.AuthenticationOIDC, error) {
	log.Debugf("[getApigatewayApiGatewayAuthnResolver]Parent Object %+v", obj)
	if id != nil && *id != "" {
		log.Debugf("[getApigatewayApiGatewayAuthnResolver]Id %q", *id)
		vOIDC, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org")).GetAuthn(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayAuthnResolver]Error getting Authn node %q : %s", *id, err)
			return &model.AuthenticationOIDC{}, nil
		}
		dn := vOIDC.DisplayName()
		parentLabels := map[string]interface{}{"oidcs.authentication.nexus.org": dn}
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
	vOIDCParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetApiGateway(context.TODO(), getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org"))
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
	parentLabels := map[string]interface{}{"oidcs.authentication.nexus.org": dn}
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
		vProxyRule, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org")).GetProxyRules(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayProxyRulesResolver]Error getting ProxyRules node %q : %s", *id, err)
			return vAdminProxyRuleList, nil
		}
		dn := vProxyRule.DisplayName()
		parentLabels := map[string]interface{}{"proxyrules.admin.nexus.org": dn}
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

	vProxyRuleParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetApiGateway(context.TODO(), getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org"))
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
		vProxyRule, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org")).GetProxyRules(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayProxyRulesResolver]Error getting ProxyRules node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vProxyRule.DisplayName()
		parentLabels := map[string]interface{}{"proxyrules.admin.nexus.org": dn}
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
		vCORSConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org")).GetCors(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayCorsResolver]Error getting Cors node %q : %s", *id, err)
			return vDomainCORSConfigList, nil
		}
		dn := vCORSConfig.DisplayName()
		parentLabels := map[string]interface{}{"corsconfigs.domain.nexus.org": dn}
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

	vCORSConfigParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetApiGateway(context.TODO(), getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org"))
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
		vCORSConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).ApiGateway(getParentName(obj.ParentLabels, "apigateways.apigateway.nexus.org")).GetCors(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getApigatewayApiGatewayCorsResolver]Error getting Cors node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vCORSConfig.DisplayName()
		parentLabels := map[string]interface{}{"corsconfigs.domain.nexus.org": dn}
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
		vApiGateway, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetApiGateway(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigApiGatewayResolver]Error getting ApiGateway node %q : %s", *id, err)
			return &model.ApigatewayApiGateway{}, nil
		}
		dn := vApiGateway.DisplayName()
		parentLabels := map[string]interface{}{"apigateways.apigateway.nexus.org": dn}

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
	vApiGatewayParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.org"))
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
	parentLabels := map[string]interface{}{"apigateways.apigateway.nexus.org": dn}

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
		vConnect, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetConnect(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigConnectResolver]Error getting Connect node %q : %s", *id, err)
			return &model.ConnectConnect{}, nil
		}
		dn := vConnect.DisplayName()
		parentLabels := map[string]interface{}{"connects.connect.nexus.org": dn}

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
	vConnectParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.org"))
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
	parentLabels := map[string]interface{}{"connects.connect.nexus.org": dn}

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
		vRoute, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetRoutes(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConfigConfigRoutesResolver]Error getting Routes node %q : %s", *id, err)
			return vRouteRouteList, nil
		}
		dn := vRoute.DisplayName()
		parentLabels := map[string]interface{}{"routes.route.nexus.org": dn}
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

	vRouteParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).GetConfig(context.TODO(), getParentName(obj.ParentLabels, "configs.config.nexus.org"))
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
		vRoute, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetRoutes(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConfigConfigRoutesResolver]Error getting Routes node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vRoute.DisplayName()
		parentLabels := map[string]interface{}{"routes.route.nexus.org": dn}
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
// FieldName: Endpoints Node: Connect PKG: Connect
//////////////////////////////////////
func getConnectConnectEndpointsResolver(obj *model.ConnectConnect, id *string) ([]*model.ConnectNexusEndpoint, error) {
	log.Debugf("[getConnectConnectEndpointsResolver]Parent Object %+v", obj)
	var vConnectNexusEndpointList []*model.ConnectNexusEndpoint
	if id != nil && *id != "" {
		log.Debugf("[getConnectConnectEndpointsResolver]Id %q", *id)
		vNexusEndpoint, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.org")).GetEndpoints(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConnectConnectEndpointsResolver]Error getting Endpoints node %q : %s", *id, err)
			return vConnectNexusEndpointList, nil
		}
		dn := vNexusEndpoint.DisplayName()
		parentLabels := map[string]interface{}{"nexusendpoints.connect.nexus.org": dn}
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

	vNexusEndpointParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetConnect(context.TODO(), getParentName(obj.ParentLabels, "connects.connect.nexus.org"))
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
		vNexusEndpoint, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.org")).GetEndpoints(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConnectConnectEndpointsResolver]Error getting Endpoints node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vNexusEndpoint.DisplayName()
		parentLabels := map[string]interface{}{"nexusendpoints.connect.nexus.org": dn}
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
		vReplicationConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.org")).GetReplicationConfig(context.TODO(), *id)
		if err != nil {
			log.Errorf("[getConnectConnectReplicationConfigResolver]Error getting ReplicationConfig node %q : %s", *id, err)
			return vConnectReplicationConfigList, nil
		}
		dn := vReplicationConfig.DisplayName()
		parentLabels := map[string]interface{}{"replicationconfigs.connect.nexus.org": dn}
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

	vReplicationConfigParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).GetConnect(context.TODO(), getParentName(obj.ParentLabels, "connects.connect.nexus.org"))
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
		vReplicationConfig, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.org")).GetReplicationConfig(context.TODO(), i.DisplayName())
		if err != nil {
			log.Errorf("[getConnectConnectReplicationConfigResolver]Error getting ReplicationConfig node %q : %s", i.DisplayName(), err)
			continue
		}
		dn := vReplicationConfig.DisplayName()
		parentLabels := map[string]interface{}{"replicationconfigs.connect.nexus.org": dn}
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
	vNexusEndpointParent, err := nc.ApiNexus(getParentName(obj.ParentLabels, "nexuses.api.nexus.org")).Config(getParentName(obj.ParentLabels, "configs.config.nexus.org")).Connect(getParentName(obj.ParentLabels, "connects.connect.nexus.org")).GetReplicationConfig(context.TODO(), getParentName(obj.ParentLabels, "replicationconfigs.connect.nexus.org"))
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
	parentLabels := map[string]interface{}{"nexusendpoints.connect.nexus.org": dn}
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
