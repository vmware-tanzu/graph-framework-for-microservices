package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-gql/graph/generated"
	"golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-gql/graph/model"
)

// Root is the resolver for the root field.
func (r *queryResolver) Root(ctx context.Context, id *string) ([]*model.ApiNexus, error) {
	return getRootResolver(id)
}

// Config is the resolver for the Config field.
func (r *api_NexusResolver) Config(ctx context.Context, obj *model.ApiNexus, id *string) (*model.ConfigConfig, error) {
	return getApiNexusConfigResolver(obj, id)
}

// ProxyRules is the resolver for the ProxyRules field.
func (r *apigateway_ApiGatewayResolver) ProxyRules(ctx context.Context, obj *model.ApigatewayApiGateway, id *string) ([]*model.AdminProxyRule, error) {
	return getApigatewayApiGatewayProxyRulesResolver(obj, id)
}

// Cors is the resolver for the Cors field.
func (r *apigateway_ApiGatewayResolver) Cors(ctx context.Context, obj *model.ApigatewayApiGateway, id *string) ([]*model.DomainCORSConfig, error) {
	return getApigatewayApiGatewayCorsResolver(obj, id)
}

// Authn is the resolver for the Authn field.
func (r *apigateway_ApiGatewayResolver) Authn(ctx context.Context, obj *model.ApigatewayApiGateway, id *string) (*model.AuthenticationOIDC, error) {
	return getApigatewayApiGatewayAuthnResolver(obj, id)
}

// Routes is the resolver for the Routes field.
func (r *config_ConfigResolver) Routes(ctx context.Context, obj *model.ConfigConfig, id *string) ([]*model.RouteRoute, error) {
	return getConfigConfigRoutesResolver(obj, id)
}

// ApiGateway is the resolver for the ApiGateway field.
func (r *config_ConfigResolver) ApiGateway(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.ApigatewayApiGateway, error) {
	return getConfigConfigApiGatewayResolver(obj, id)
}

// Connect is the resolver for the Connect field.
func (r *config_ConfigResolver) Connect(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.ConnectConnect, error) {
	return getConfigConfigConnectResolver(obj, id)
}

// Endpoints is the resolver for the Endpoints field.
func (r *connect_ConnectResolver) Endpoints(ctx context.Context, obj *model.ConnectConnect, id *string) ([]*model.ConnectNexusEndpoint, error) {
	return getConnectConnectEndpointsResolver(obj, id)
}

// ReplicationConfig is the resolver for the ReplicationConfig field.
func (r *connect_ConnectResolver) ReplicationConfig(ctx context.Context, obj *model.ConnectConnect, id *string) ([]*model.ConnectReplicationConfig, error) {
	return getConnectConnectReplicationConfigResolver(obj, id)
}

// RemoteEndpoint is the resolver for the RemoteEndpoint field.
func (r *connect_ReplicationConfigResolver) RemoteEndpoint(ctx context.Context, obj *model.ConnectReplicationConfig) (*model.ConnectNexusEndpoint, error) {
	return getConnectReplicationConfigRemoteEndpointResolver(obj)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Api_Nexus returns generated.Api_NexusResolver implementation.
func (r *Resolver) Api_Nexus() generated.Api_NexusResolver { return &api_NexusResolver{r} }

// Apigateway_ApiGateway returns generated.Apigateway_ApiGatewayResolver implementation.
func (r *Resolver) Apigateway_ApiGateway() generated.Apigateway_ApiGatewayResolver {
	return &apigateway_ApiGatewayResolver{r}
}

// Config_Config returns generated.Config_ConfigResolver implementation.
func (r *Resolver) Config_Config() generated.Config_ConfigResolver { return &config_ConfigResolver{r} }

// Connect_Connect returns generated.Connect_ConnectResolver implementation.
func (r *Resolver) Connect_Connect() generated.Connect_ConnectResolver {
	return &connect_ConnectResolver{r}
}

// Connect_ReplicationConfig returns generated.Connect_ReplicationConfigResolver implementation.
func (r *Resolver) Connect_ReplicationConfig() generated.Connect_ReplicationConfigResolver {
	return &connect_ReplicationConfigResolver{r}
}

type queryResolver struct{ *Resolver }
type api_NexusResolver struct{ *Resolver }
type apigateway_ApiGatewayResolver struct{ *Resolver }
type config_ConfigResolver struct{ *Resolver }
type connect_ConnectResolver struct{ *Resolver }
type connect_ReplicationConfigResolver struct{ *Resolver }
