package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/nexus-gql/graph/generated"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/nexus-gql/graph/model"
)

// Root is the resolver for the root field.
func (r *queryResolver) Root(ctx context.Context) (*model.RootRoot, error) {
	return getRootResolver()
}

// QueryExample is the resolver for the QueryExample field.
func (r *config_ConfigResolver) QueryExample(ctx context.Context, obj *model.ConfigConfig, startTime *string, endTime *string, interval *string, isServiceDeployment *bool, startVal *int) (*model.NexusGraphqlResponse, error) {
	return getConfigConfigQueryExampleResolver(obj, startTime, endTime, interval, isServiceDeployment, startVal)
}

// ACPPolicies is the resolver for the ACPPolicies field.
func (r *config_ConfigResolver) ACPPolicies(ctx context.Context, obj *model.ConfigConfig, id *string) ([]*model.PolicypkgAccessControlPolicy, error) {
	return getConfigConfigACPPoliciesResolver(obj, id)
}

// FooExample is the resolver for the FooExample field.
func (r *config_ConfigResolver) FooExample(ctx context.Context, obj *model.ConfigConfig, id *string) ([]*model.ConfigFooTypeABC, error) {
	return getConfigConfigFooExampleResolver(obj, id)
}

// GNS is the resolver for the GNS field.
func (r *config_ConfigResolver) GNS(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.GnsGns, error) {
	return getConfigConfigGNSResolver(obj, id)
}

// DNS is the resolver for the DNS field.
func (r *config_ConfigResolver) DNS(ctx context.Context, obj *model.ConfigConfig) (*model.GnsDns, error) {
	return getConfigConfigDNSResolver(obj)
}

// VMPPolicies is the resolver for the VMPPolicies field.
func (r *config_ConfigResolver) VMPPolicies(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.PolicypkgVMpolicy, error) {
	return getConfigConfigVMPPoliciesResolver(obj, id)
}

// Domain is the resolver for the Domain field.
func (r *config_ConfigResolver) Domain(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.ConfigDomain, error) {
	return getConfigConfigDomainResolver(obj, id)
}

// SvcGrpInfo is the resolver for the SvcGrpInfo field.
func (r *config_ConfigResolver) SvcGrpInfo(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.ServicegroupSvcGroupLinkInfo, error) {
	return getConfigConfigSvcGrpInfoResolver(obj, id)
}

// QueryGns1 is the resolver for the queryGns1 field.
func (r *gns_GnsResolver) QueryGns1(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, interval *string, isServiceDeployment *bool, startVal *int) (*model.NexusGraphqlResponse, error) {
	return getGnsGnsqueryGns1Resolver(obj, startTime, endTime, interval, isServiceDeployment, startVal)
}

// QueryGnsQM1 is the resolver for the queryGnsQM1 field.
func (r *gns_GnsResolver) QueryGnsQM1(ctx context.Context, obj *model.GnsGns) (*model.TimeSeriesData, error) {
	return getGnsGnsqueryGnsQM1Resolver(obj)
}

// QueryGnsQM is the resolver for the queryGnsQM field.
func (r *gns_GnsResolver) QueryGnsQM(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, timeInterval *string, someUserArg1 *string, someUserArg2 *int, someUserArg3 *bool) (*model.TimeSeriesData, error) {
	return getGnsGnsqueryGnsQMResolver(obj, startTime, endTime, timeInterval, someUserArg1, someUserArg2, someUserArg3)
}

// GnsAccessControlPolicy is the resolver for the GnsAccessControlPolicy field.
func (r *gns_GnsResolver) GnsAccessControlPolicy(ctx context.Context, obj *model.GnsGns, id *string) (*model.PolicypkgAccessControlPolicy, error) {
	return getGnsGnsGnsAccessControlPolicyResolver(obj, id)
}

// FooChild is the resolver for the FooChild field.
func (r *gns_GnsResolver) FooChild(ctx context.Context, obj *model.GnsGns) (*model.GnsBarChild, error) {
	return getGnsGnsFooChildResolver(obj)
}

// PolicyConfigs is the resolver for the PolicyConfigs field.
func (r *policypkg_AccessControlPolicyResolver) PolicyConfigs(ctx context.Context, obj *model.PolicypkgAccessControlPolicy, id *string) ([]*model.PolicypkgACPConfig, error) {
	return getPolicypkgAccessControlPolicyPolicyConfigsResolver(obj, id)
}

// QueryGns1 is the resolver for the queryGns1 field.
func (r *policypkg_VMpolicyResolver) QueryGns1(ctx context.Context, obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, interval *string, isServiceDeployment *bool, startVal *int) (*model.NexusGraphqlResponse, error) {
	return getPolicypkgVMpolicyqueryGns1Resolver(obj, startTime, endTime, interval, isServiceDeployment, startVal)
}

// QueryGnsQM1 is the resolver for the queryGnsQM1 field.
func (r *policypkg_VMpolicyResolver) QueryGnsQM1(ctx context.Context, obj *model.PolicypkgVMpolicy) (*model.TimeSeriesData, error) {
	return getPolicypkgVMpolicyqueryGnsQM1Resolver(obj)
}

// QueryGnsQM is the resolver for the queryGnsQM field.
func (r *policypkg_VMpolicyResolver) QueryGnsQM(ctx context.Context, obj *model.PolicypkgVMpolicy, startTime *string, endTime *string, timeInterval *string, someUserArg1 *string, someUserArg2 *int, someUserArg3 *bool) (*model.TimeSeriesData, error) {
	return getPolicypkgVMpolicyqueryGnsQMResolver(obj, startTime, endTime, timeInterval, someUserArg1, someUserArg2, someUserArg3)
}

// Config is the resolver for the Config field.
func (r *root_RootResolver) Config(ctx context.Context, obj *model.RootRoot, id *string) (*model.ConfigConfig, error) {
	return getRootRootConfigResolver(obj, id)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Config_Config returns generated.Config_ConfigResolver implementation.
func (r *Resolver) Config_Config() generated.Config_ConfigResolver { return &config_ConfigResolver{r} }

// Gns_Gns returns generated.Gns_GnsResolver implementation.
func (r *Resolver) Gns_Gns() generated.Gns_GnsResolver { return &gns_GnsResolver{r} }

// Policypkg_AccessControlPolicy returns generated.Policypkg_AccessControlPolicyResolver implementation.
func (r *Resolver) Policypkg_AccessControlPolicy() generated.Policypkg_AccessControlPolicyResolver {
	return &policypkg_AccessControlPolicyResolver{r}
}

// Policypkg_VMpolicy returns generated.Policypkg_VMpolicyResolver implementation.
func (r *Resolver) Policypkg_VMpolicy() generated.Policypkg_VMpolicyResolver {
	return &policypkg_VMpolicyResolver{r}
}

// Root_Root returns generated.Root_RootResolver implementation.
func (r *Resolver) Root_Root() generated.Root_RootResolver { return &root_RootResolver{r} }

type queryResolver struct{ *Resolver }
type config_ConfigResolver struct{ *Resolver }
type gns_GnsResolver struct{ *Resolver }
type policypkg_AccessControlPolicyResolver struct{ *Resolver }
type policypkg_VMpolicyResolver struct{ *Resolver }
type root_RootResolver struct{ *Resolver }
