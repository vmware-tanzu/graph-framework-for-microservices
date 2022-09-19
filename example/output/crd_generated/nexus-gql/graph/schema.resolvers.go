package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/nexus-gql/graph/generated"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generated/nexus-gql/graph/model"
)

// Root is the resolver for the root field.
func (r *queryResolver) Root(ctx context.Context) (*model.RootRoot, error) {
	return c.getRootResolver()
}

// QueryServiceTable is the resolver for the queryServiceTable field.
func (r *config_ConfigResolver) QueryServiceTable(ctx context.Context, obj *model.ConfigConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryServiceTableResolver(startTime, endTime, systemServices, showGateways, groupby, noMetrics)
}

// QueryServiceVersionTable is the resolver for the queryServiceVersionTable field.
func (r *config_ConfigResolver) QueryServiceVersionTable(ctx context.Context, obj *model.ConfigConfig, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryServiceVersionTableResolver(startTime, endTime, systemServices, showGateways, noMetrics)
}

// QueryServiceTs is the resolver for the queryServiceTS field.
func (r *config_ConfigResolver) QueryServiceTs(ctx context.Context, obj *model.ConfigConfig, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryServiceTSResolver(svcMetric, startTime, endTime, timeInterval)
}

// QueryIncomingAPIs is the resolver for the queryIncomingAPIs field.
func (r *config_ConfigResolver) QueryIncomingAPIs(ctx context.Context, obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryIncomingAPIsResolver(startTime, endTime, destinationService, destinationServiceVersion, timeInterval, timeZone)
}

// QueryOutgoingAPIs is the resolver for the queryOutgoingAPIs field.
func (r *config_ConfigResolver) QueryOutgoingAPIs(ctx context.Context, obj *model.ConfigConfig, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryOutgoingAPIsResolver(startTime, endTime, timeInterval, timeZone)
}

// QueryIncomingTcp is the resolver for the queryIncomingTCP field.
func (r *config_ConfigResolver) QueryIncomingTcp(ctx context.Context, obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryIncomingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryOutgoingTcp is the resolver for the queryOutgoingTCP field.
func (r *config_ConfigResolver) QueryOutgoingTcp(ctx context.Context, obj *model.ConfigConfig, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryOutgoingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryServiceTopology is the resolver for the queryServiceTopology field.
func (r *config_ConfigResolver) QueryServiceTopology(ctx context.Context, obj *model.ConfigConfig, metricStringArray *string, startTime *string, endTime *string) (*model.TimeSeriesData, error) {
	return c.getConfigConfigqueryServiceTopologyResolver(metricStringArray, startTime, endTime)
}

// Gns is the resolver for the GNS field.
func (r *config_ConfigResolver) Gns(ctx context.Context, obj *model.ConfigConfig) (*model.GnsGns, error) {
	return c.getConfigConfigGNSResolver()
}

// Cluster is the resolver for the Cluster field.
func (r *config_ConfigResolver) Cluster(ctx context.Context, obj *model.ConfigConfig) (*model.ConfigCluster, error) {
	return c.getConfigConfigClusterResolver()
}

// QueryServiceTable is the resolver for the queryServiceTable field.
func (r *gns_BarResolver) QueryServiceTable(ctx context.Context, obj *model.GnsBar, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryServiceTableResolver(startTime, endTime, systemServices, showGateways, groupby, noMetrics)
}

// QueryServiceVersionTable is the resolver for the queryServiceVersionTable field.
func (r *gns_BarResolver) QueryServiceVersionTable(ctx context.Context, obj *model.GnsBar, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryServiceVersionTableResolver(startTime, endTime, systemServices, showGateways, noMetrics)
}

// QueryServiceTs is the resolver for the queryServiceTS field.
func (r *gns_BarResolver) QueryServiceTs(ctx context.Context, obj *model.GnsBar, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryServiceTSResolver(svcMetric, startTime, endTime, timeInterval)
}

// QueryIncomingAPIs is the resolver for the queryIncomingAPIs field.
func (r *gns_BarResolver) QueryIncomingAPIs(ctx context.Context, obj *model.GnsBar, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryIncomingAPIsResolver(startTime, endTime, destinationService, destinationServiceVersion, timeInterval, timeZone)
}

// QueryOutgoingAPIs is the resolver for the queryOutgoingAPIs field.
func (r *gns_BarResolver) QueryOutgoingAPIs(ctx context.Context, obj *model.GnsBar, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryOutgoingAPIsResolver(startTime, endTime, timeInterval, timeZone)
}

// QueryIncomingTcp is the resolver for the queryIncomingTCP field.
func (r *gns_BarResolver) QueryIncomingTcp(ctx context.Context, obj *model.GnsBar, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryIncomingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryOutgoingTcp is the resolver for the queryOutgoingTCP field.
func (r *gns_BarResolver) QueryOutgoingTcp(ctx context.Context, obj *model.GnsBar, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryOutgoingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryServiceTopology is the resolver for the queryServiceTopology field.
func (r *gns_BarResolver) QueryServiceTopology(ctx context.Context, obj *model.GnsBar, metricStringArray *string, startTime *string, endTime *string) (*model.TimeSeriesData, error) {
	return c.getGnsBarqueryServiceTopologyResolver(metricStringArray, startTime, endTime)
}

// QueryServiceTable is the resolver for the queryServiceTable field.
func (r *gns_EmptyDataResolver) QueryServiceTable(ctx context.Context, obj *model.GnsEmptyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryServiceTableResolver(startTime, endTime, systemServices, showGateways, groupby, noMetrics)
}

// QueryServiceVersionTable is the resolver for the queryServiceVersionTable field.
func (r *gns_EmptyDataResolver) QueryServiceVersionTable(ctx context.Context, obj *model.GnsEmptyData, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryServiceVersionTableResolver(startTime, endTime, systemServices, showGateways, noMetrics)
}

// QueryServiceTs is the resolver for the queryServiceTS field.
func (r *gns_EmptyDataResolver) QueryServiceTs(ctx context.Context, obj *model.GnsEmptyData, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryServiceTSResolver(svcMetric, startTime, endTime, timeInterval)
}

// QueryIncomingAPIs is the resolver for the queryIncomingAPIs field.
func (r *gns_EmptyDataResolver) QueryIncomingAPIs(ctx context.Context, obj *model.GnsEmptyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryIncomingAPIsResolver(startTime, endTime, destinationService, destinationServiceVersion, timeInterval, timeZone)
}

// QueryOutgoingAPIs is the resolver for the queryOutgoingAPIs field.
func (r *gns_EmptyDataResolver) QueryOutgoingAPIs(ctx context.Context, obj *model.GnsEmptyData, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryOutgoingAPIsResolver(startTime, endTime, timeInterval, timeZone)
}

// QueryIncomingTcp is the resolver for the queryIncomingTCP field.
func (r *gns_EmptyDataResolver) QueryIncomingTcp(ctx context.Context, obj *model.GnsEmptyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryIncomingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryOutgoingTcp is the resolver for the queryOutgoingTCP field.
func (r *gns_EmptyDataResolver) QueryOutgoingTcp(ctx context.Context, obj *model.GnsEmptyData, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryOutgoingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryServiceTopology is the resolver for the queryServiceTopology field.
func (r *gns_EmptyDataResolver) QueryServiceTopology(ctx context.Context, obj *model.GnsEmptyData, metricStringArray *string, startTime *string, endTime *string) (*model.TimeSeriesData, error) {
	return c.getGnsEmptyDataqueryServiceTopologyResolver(metricStringArray, startTime, endTime)
}

// QueryServiceTable is the resolver for the queryServiceTable field.
func (r *gns_GnsResolver) QueryServiceTable(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryServiceTableResolver(startTime, endTime, systemServices, showGateways, groupby, noMetrics)
}

// QueryServiceVersionTable is the resolver for the queryServiceVersionTable field.
func (r *gns_GnsResolver) QueryServiceVersionTable(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryServiceVersionTableResolver(startTime, endTime, systemServices, showGateways, noMetrics)
}

// QueryServiceTs is the resolver for the queryServiceTS field.
func (r *gns_GnsResolver) QueryServiceTs(ctx context.Context, obj *model.GnsGns, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryServiceTSResolver(svcMetric, startTime, endTime, timeInterval)
}

// QueryIncomingAPIs is the resolver for the queryIncomingAPIs field.
func (r *gns_GnsResolver) QueryIncomingAPIs(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryIncomingAPIsResolver(startTime, endTime, destinationService, destinationServiceVersion, timeInterval, timeZone)
}

// QueryOutgoingAPIs is the resolver for the queryOutgoingAPIs field.
func (r *gns_GnsResolver) QueryOutgoingAPIs(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryOutgoingAPIsResolver(startTime, endTime, timeInterval, timeZone)
}

// QueryIncomingTcp is the resolver for the queryIncomingTCP field.
func (r *gns_GnsResolver) QueryIncomingTcp(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryIncomingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryOutgoingTcp is the resolver for the queryOutgoingTCP field.
func (r *gns_GnsResolver) QueryOutgoingTcp(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryOutgoingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryServiceTopology is the resolver for the queryServiceTopology field.
func (r *gns_GnsResolver) QueryServiceTopology(ctx context.Context, obj *model.GnsGns, metricStringArray *string, startTime *string, endTime *string) (*model.TimeSeriesData, error) {
	return c.getGnsGnsqueryServiceTopologyResolver(metricStringArray, startTime, endTime)
}

// FooLink is the resolver for the FooLink field.
func (r *gns_GnsResolver) FooLink(ctx context.Context, obj *model.GnsGns) (*model.GnsBar, error) {
	return c.getGnsGnsFooLinkResolver()
}

// FooLinks is the resolver for the FooLinks field.
func (r *gns_GnsResolver) FooLinks(ctx context.Context, obj *model.GnsGns, id *string) ([]*model.GnsBar, error) {
	return c.getGnsGnsFooLinksResolver(id)
}

// FooChild is the resolver for the FooChild field.
func (r *gns_GnsResolver) FooChild(ctx context.Context, obj *model.GnsGns) (*model.GnsBar, error) {
	return c.getGnsGnsFooChildResolver()
}

// FooChildren is the resolver for the FooChildren field.
func (r *gns_GnsResolver) FooChildren(ctx context.Context, obj *model.GnsGns, id *string) ([]*model.GnsBar, error) {
	return c.getGnsGnsFooChildrenResolver(id)
}

// Mydesc is the resolver for the Mydesc field.
func (r *gns_GnsResolver) Mydesc(ctx context.Context, obj *model.GnsGns) (*model.GnsDescription, error) {
	return c.getGnsGnsMydescResolver()
}

// HostPort is the resolver for the HostPort field.
func (r *gns_GnsResolver) HostPort(ctx context.Context, obj *model.GnsGns) (*model.GnsHostPort, error) {
	return c.getGnsGnsHostPortResolver()
}

// TestArray is the resolver for the TestArray field.
func (r *gns_GnsResolver) TestArray(ctx context.Context, obj *model.GnsGns) (*model.GnsEmptyData, error) {
	return c.getGnsGnsTestArrayResolver()
}

// QueryServiceTable is the resolver for the queryServiceTable field.
func (r *root_RootResolver) QueryServiceTable(ctx context.Context, obj *model.RootRoot, startTime *string, endTime *string, systemServices *bool, showGateways *bool, groupby *string, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryServiceTableResolver(startTime, endTime, systemServices, showGateways, groupby, noMetrics)
}

// QueryServiceVersionTable is the resolver for the queryServiceVersionTable field.
func (r *root_RootResolver) QueryServiceVersionTable(ctx context.Context, obj *model.RootRoot, startTime *string, endTime *string, systemServices *bool, showGateways *bool, noMetrics *bool) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryServiceVersionTableResolver(startTime, endTime, systemServices, showGateways, noMetrics)
}

// QueryServiceTs is the resolver for the queryServiceTS field.
func (r *root_RootResolver) QueryServiceTs(ctx context.Context, obj *model.RootRoot, svcMetric *string, startTime *string, endTime *string, timeInterval *string) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryServiceTSResolver(svcMetric, startTime, endTime, timeInterval)
}

// QueryIncomingAPIs is the resolver for the queryIncomingAPIs field.
func (r *root_RootResolver) QueryIncomingAPIs(ctx context.Context, obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryIncomingAPIsResolver(startTime, endTime, destinationService, destinationServiceVersion, timeInterval, timeZone)
}

// QueryOutgoingAPIs is the resolver for the queryOutgoingAPIs field.
func (r *root_RootResolver) QueryOutgoingAPIs(ctx context.Context, obj *model.RootRoot, startTime *string, endTime *string, timeInterval *string, timeZone *string) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryOutgoingAPIsResolver(startTime, endTime, timeInterval, timeZone)
}

// QueryIncomingTcp is the resolver for the queryIncomingTCP field.
func (r *root_RootResolver) QueryIncomingTcp(ctx context.Context, obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryIncomingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryOutgoingTcp is the resolver for the queryOutgoingTCP field.
func (r *root_RootResolver) QueryOutgoingTcp(ctx context.Context, obj *model.RootRoot, startTime *string, endTime *string, destinationService *string, destinationServiceVersion *string) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryOutgoingTCPResolver(startTime, endTime, destinationService, destinationServiceVersion)
}

// QueryServiceTopology is the resolver for the queryServiceTopology field.
func (r *root_RootResolver) QueryServiceTopology(ctx context.Context, obj *model.RootRoot, metricStringArray *string, startTime *string, endTime *string) (*model.TimeSeriesData, error) {
	return c.getRootRootqueryServiceTopologyResolver(metricStringArray, startTime, endTime)
}

// Config is the resolver for the Config field.
func (r *root_RootResolver) Config(ctx context.Context, obj *model.RootRoot) (*model.ConfigConfig, error) {
	return c.getRootRootConfigResolver()
}

// CustomBar is the resolver for the CustomBar field.
func (r *root_RootResolver) CustomBar(ctx context.Context, obj *model.RootRoot) (*model.RootBar, error) {
	return c.getRootRootCustomBarResolver()
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Config_Config returns generated.Config_ConfigResolver implementation.
func (r *Resolver) Config_Config() generated.Config_ConfigResolver { return &config_ConfigResolver{r} }

// Gns_Bar returns generated.Gns_BarResolver implementation.
func (r *Resolver) Gns_Bar() generated.Gns_BarResolver { return &gns_BarResolver{r} }

// Gns_EmptyData returns generated.Gns_EmptyDataResolver implementation.
func (r *Resolver) Gns_EmptyData() generated.Gns_EmptyDataResolver { return &gns_EmptyDataResolver{r} }

// Gns_Gns returns generated.Gns_GnsResolver implementation.
func (r *Resolver) Gns_Gns() generated.Gns_GnsResolver { return &gns_GnsResolver{r} }

// Root_Root returns generated.Root_RootResolver implementation.
func (r *Resolver) Root_Root() generated.Root_RootResolver { return &root_RootResolver{r} }

type queryResolver struct{ *Resolver }
type config_ConfigResolver struct{ *Resolver }
type gns_BarResolver struct{ *Resolver }
type gns_EmptyDataResolver struct{ *Resolver }
type gns_GnsResolver struct{ *Resolver }
type root_RootResolver struct{ *Resolver }
