package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	generated1 "nexustempmodule/nexus-gql/graph/generated"
	model1 "nexustempmodule/nexus-gql/graph/model"
)

// Root is the resolver for the root field.
func (r *queryResolver) Root(ctx context.Context) (*model1.RootRoot, error) {
	return getRootResolver()
}

// QueryExample is the resolver for the QueryExample field.
func (r *config_ConfigResolver) QueryExample(ctx context.Context, obj *model1.ConfigConfig, startTime *string, endTime *string, interval *string, isServiceDeployment *bool, startVal *int) (*model1.NexusGraphqlResponse, error) {
	return getConfigConfigQueryExampleResolver(obj, startTime, endTime, interval, isServiceDeployment, startVal)
}

// ACPPolicies is the resolver for the ACPPolicies field.
func (r *config_ConfigResolver) ACPPolicies(ctx context.Context, obj *model1.ConfigConfig, id *string) ([]*model1.PolicypkgAccessControlPolicy, error) {
	return getConfigConfigACPPoliciesResolver(obj, id)
}

// FooExample is the resolver for the FooExample field.
func (r *config_ConfigResolver) FooExample(ctx context.Context, obj *model1.ConfigConfig, id *string) ([]*model1.ConfigFooTypeABC, error) {
	return getConfigConfigFooExampleResolver(obj, id)
}

// GNS is the resolver for the GNS field.
func (r *config_ConfigResolver) GNS(ctx context.Context, obj *model1.ConfigConfig, id *string) (*model1.GnsGns, error) {
	return getConfigConfigGNSResolver(obj, id)
}

// DNS is the resolver for the DNS field.
func (r *config_ConfigResolver) DNS(ctx context.Context, obj *model1.ConfigConfig) (*model1.GnsDns, error) {
	return getConfigConfigDNSResolver(obj)
}

// VMPPolicies is the resolver for the VMPPolicies field.
func (r *config_ConfigResolver) VMPPolicies(ctx context.Context, obj *model1.ConfigConfig, id *string) (*model1.PolicypkgVMpolicy, error) {
	return getConfigConfigVMPPoliciesResolver(obj, id)
}

// Domain is the resolver for the Domain field.
func (r *config_ConfigResolver) Domain(ctx context.Context, obj *model1.ConfigConfig, id *string) (*model1.ConfigDomain, error) {
	return getConfigConfigDomainResolver(obj, id)
}

// SvcGrpInfo is the resolver for the SvcGrpInfo field.
func (r *config_ConfigResolver) SvcGrpInfo(ctx context.Context, obj *model1.ConfigConfig, id *string) (*model1.ServicegroupSvcGroupLinkInfo, error) {
	return getConfigConfigSvcGrpInfoResolver(obj, id)
}

// QueryGns1 is the resolver for the queryGns1 field.
func (r *gns_GnsResolver) QueryGns1(ctx context.Context, obj *model1.GnsGns, startTime *string, endTime *string, interval *string, isServiceDeployment *bool, startVal *int) (*model1.NexusGraphqlResponse, error) {
	return getGnsGnsqueryGns1Resolver(obj, startTime, endTime, interval, isServiceDeployment, startVal)
}

// QueryGnsQM1 is the resolver for the queryGnsQM1 field.
func (r *gns_GnsResolver) QueryGnsQM1(ctx context.Context, obj *model1.GnsGns) (*model1.TimeSeriesData, error) {
	return getGnsGnsqueryGnsQM1Resolver(obj)
}

// QueryGnsQM is the resolver for the queryGnsQM field.
func (r *gns_GnsResolver) QueryGnsQM(ctx context.Context, obj *model1.GnsGns, startTime *string, endTime *string, timeInterval *string, someUserArg1 *string, someUserArg2 *int, someUserArg3 *bool) (*model1.TimeSeriesData, error) {
	return getGnsGnsqueryGnsQMResolver(obj, startTime, endTime, timeInterval, someUserArg1, someUserArg2, someUserArg3)
}

// GnsAccessControlPolicy is the resolver for the GnsAccessControlPolicy field.
func (r *gns_GnsResolver) GnsAccessControlPolicy(ctx context.Context, obj *model1.GnsGns, id *string) (*model1.PolicypkgAccessControlPolicy, error) {
	return getGnsGnsGnsAccessControlPolicyResolver(obj, id)
}

// FooChild is the resolver for the FooChild field.
func (r *gns_GnsResolver) FooChild(ctx context.Context, obj *model1.GnsGns) (*model1.GnsBarChild, error) {
	return getGnsGnsFooChildResolver(obj)
}

// PolicyConfigs is the resolver for the PolicyConfigs field.
func (r *policypkg_AccessControlPolicyResolver) PolicyConfigs(ctx context.Context, obj *model1.PolicypkgAccessControlPolicy, id *string) ([]*model1.PolicypkgACPConfig, error) {
	return getPolicypkgAccessControlPolicyPolicyConfigsResolver(obj, id)
}

// QueryGns1 is the resolver for the queryGns1 field.
func (r *policypkg_VMpolicyResolver) QueryGns1(ctx context.Context, obj *model1.PolicypkgVMpolicy, startTime *string, endTime *string, interval *string, isServiceDeployment *bool, startVal *int) (*model1.NexusGraphqlResponse, error) {
	return getPolicypkgVMpolicyqueryGns1Resolver(obj, startTime, endTime, interval, isServiceDeployment, startVal)
}

// QueryGnsQM1 is the resolver for the queryGnsQM1 field.
func (r *policypkg_VMpolicyResolver) QueryGnsQM1(ctx context.Context, obj *model1.PolicypkgVMpolicy) (*model1.TimeSeriesData, error) {
	return getPolicypkgVMpolicyqueryGnsQM1Resolver(obj)
}

// QueryGnsQM is the resolver for the queryGnsQM field.
func (r *policypkg_VMpolicyResolver) QueryGnsQM(ctx context.Context, obj *model1.PolicypkgVMpolicy, startTime *string, endTime *string, timeInterval *string, someUserArg1 *string, someUserArg2 *int, someUserArg3 *bool) (*model1.TimeSeriesData, error) {
	return getPolicypkgVMpolicyqueryGnsQMResolver(obj, startTime, endTime, timeInterval, someUserArg1, someUserArg2, someUserArg3)
}

// Config is the resolver for the Config field.
func (r *root_RootResolver) Config(ctx context.Context, obj *model1.RootRoot, id *string) (*model1.ConfigConfig, error) {
	return getRootRootConfigResolver(obj, id)
}

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

// Config_Config returns generated1.Config_ConfigResolver implementation.
func (r *Resolver) Config_Config() generated1.Config_ConfigResolver { return &config_ConfigResolver{r} }

// Gns_Gns returns generated1.Gns_GnsResolver implementation.
func (r *Resolver) Gns_Gns() generated1.Gns_GnsResolver { return &gns_GnsResolver{r} }

// Policypkg_AccessControlPolicy returns generated1.Policypkg_AccessControlPolicyResolver implementation.
func (r *Resolver) Policypkg_AccessControlPolicy() generated1.Policypkg_AccessControlPolicyResolver {
	return &policypkg_AccessControlPolicyResolver{r}
}

// Policypkg_VMpolicy returns generated1.Policypkg_VMpolicyResolver implementation.
func (r *Resolver) Policypkg_VMpolicy() generated1.Policypkg_VMpolicyResolver {
	return &policypkg_VMpolicyResolver{r}
}

// Root_Root returns generated1.Root_RootResolver implementation.
func (r *Resolver) Root_Root() generated1.Root_RootResolver { return &root_RootResolver{r} }

type queryResolver struct{ *Resolver }
type config_ConfigResolver struct{ *Resolver }
type gns_GnsResolver struct{ *Resolver }
type policypkg_AccessControlPolicyResolver struct{ *Resolver }
type policypkg_VMpolicyResolver struct{ *Resolver }
type root_RootResolver struct{ *Resolver }
