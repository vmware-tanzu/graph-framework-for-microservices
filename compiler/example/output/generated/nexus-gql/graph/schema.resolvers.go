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
func (r *config_ConfigResolver) ACPPolicies(ctx context.Context, obj *model.ConfigConfig, id *string) ([]*model.PolicyAccessControlPolicy, error) {
	return getConfigConfigACPPoliciesResolver(obj, id)
}

// FooExample is the resolver for the FooExample field.
func (r *config_ConfigResolver) FooExample(ctx context.Context, obj *model.ConfigConfig, id *string) ([]*model.ConfigFooType, error) {
	return getConfigConfigFooExampleResolver(obj, id)
}

// Gns is the resolver for the GNS field.
func (r *config_ConfigResolver) Gns(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.GnsGns, error) {
	return getConfigConfigGNSResolver(obj, id)
}

// Dns is the resolver for the DNS field.
func (r *config_ConfigResolver) Dns(ctx context.Context, obj *model.ConfigConfig) (*model.GnsDns, error) {
	return getConfigConfigDNSResolver(obj)
}

// VMPPolicies is the resolver for the VMPPolicies field.
func (r *config_ConfigResolver) VMPPolicies(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.PolicyVMpolicy, error) {
	return getConfigConfigVMPPoliciesResolver(obj, id)
}

// Domain is the resolver for the Domain field.
func (r *config_ConfigResolver) Domain(ctx context.Context, obj *model.ConfigConfig, id *string) (*model.ConfigDomain, error) {
	return getConfigConfigDomainResolver(obj, id)
}

// QueryGns1 is the resolver for the queryGns1 field.
func (r *gns_GnsResolver) QueryGns1(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, interval *string, isServiceDeployment *bool, startVal *int) (*model.NexusGraphqlResponse, error) {
	return getGnsGnsqueryGns1Resolver(obj, startTime, endTime, interval, isServiceDeployment, startVal)
}

// QueryGnsQm1 is the resolver for the queryGnsQM1 field.
func (r *gns_GnsResolver) QueryGnsQm1(ctx context.Context, obj *model.GnsGns) (*model.TimeSeriesData, error) {
	return getGnsGnsqueryGnsQM1Resolver(obj)
}

// QueryGnsQm is the resolver for the queryGnsQM field.
func (r *gns_GnsResolver) QueryGnsQm(ctx context.Context, obj *model.GnsGns, startTime *string, endTime *string, interval *string) (*model.TimeSeriesData, error) {
	return getGnsGnsqueryGnsQMResolver(obj, startTime, endTime, interval)
}

// GnsServiceGroups is the resolver for the GnsServiceGroups field.
func (r *gns_GnsResolver) GnsServiceGroups(ctx context.Context, obj *model.GnsGns, id *string) ([]*model.ServicegroupSvcGroup, error) {
	return getGnsGnsGnsServiceGroupsResolver(obj, id)
}

// GnsAccessControlPolicy is the resolver for the GnsAccessControlPolicy field.
func (r *gns_GnsResolver) GnsAccessControlPolicy(ctx context.Context, obj *model.GnsGns, id *string) (*model.PolicyAccessControlPolicy, error) {
	return getGnsGnsGnsAccessControlPolicyResolver(obj, id)
}

// FooChild is the resolver for the FooChild field.
func (r *gns_GnsResolver) FooChild(ctx context.Context, obj *model.GnsGns) (*model.GnsBarChild, error) {
	return getGnsGnsFooChildResolver(obj)
}

// Foo is the resolver for the Foo field.
func (r *gns_GnsResolver) Foo(ctx context.Context, obj *model.GnsGns, id *string) (*model.GnsFoo, error) {
	return getGnsGnsFooResolver(obj, id)
}

// DestSvcGroups is the resolver for the DestSvcGroups field.
func (r *policy_ACPConfigResolver) DestSvcGroups(ctx context.Context, obj *model.PolicyACPConfig, id *string) ([]*model.ServicegroupSvcGroup, error) {
	return getPolicyACPConfigDestSvcGroupsResolver(obj, id)
}

// SourceSvcGroups is the resolver for the SourceSvcGroups field.
func (r *policy_ACPConfigResolver) SourceSvcGroups(ctx context.Context, obj *model.PolicyACPConfig, id *string) ([]*model.ServicegroupSvcGroup, error) {
	return getPolicyACPConfigSourceSvcGroupsResolver(obj, id)
}

// PolicyConfigs is the resolver for the PolicyConfigs field.
func (r *policy_AccessControlPolicyResolver) PolicyConfigs(ctx context.Context, obj *model.PolicyAccessControlPolicy, id *string) ([]*model.PolicyACPConfig, error) {
	return getPolicyAccessControlPolicyPolicyConfigsResolver(obj, id)
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

// Policy_ACPConfig returns generated.Policy_ACPConfigResolver implementation.
func (r *Resolver) Policy_ACPConfig() generated.Policy_ACPConfigResolver {
	return &policy_ACPConfigResolver{r}
}

// Policy_AccessControlPolicy returns generated.Policy_AccessControlPolicyResolver implementation.
func (r *Resolver) Policy_AccessControlPolicy() generated.Policy_AccessControlPolicyResolver {
	return &policy_AccessControlPolicyResolver{r}
}

// Root_Root returns generated.Root_RootResolver implementation.
func (r *Resolver) Root_Root() generated.Root_RootResolver { return &root_RootResolver{r} }

type queryResolver struct{ *Resolver }
type config_ConfigResolver struct{ *Resolver }
type gns_GnsResolver struct{ *Resolver }
type policy_ACPConfigResolver struct{ *Resolver }
type policy_AccessControlPolicyResolver struct{ *Resolver }
type root_RootResolver struct{ *Resolver }
