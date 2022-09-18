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

// Gns is the resolver for the GNS field.
func (r *config_ConfigResolver) Gns(ctx context.Context, obj *model.ConfigConfig) (*model.GnsGns, error) {
	return c.getConfigConfigGNSResolver()
}

// Cluster is the resolver for the Cluster field.
func (r *config_ConfigResolver) Cluster(ctx context.Context, obj *model.ConfigConfig) (*model.ConfigCluster, error) {
	return c.getConfigConfigClusterResolver()
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

// Gns_Gns returns generated.Gns_GnsResolver implementation.
func (r *Resolver) Gns_Gns() generated.Gns_GnsResolver { return &gns_GnsResolver{r} }

// Root_Root returns generated.Root_RootResolver implementation.
func (r *Resolver) Root_Root() generated.Root_RootResolver { return &root_RootResolver{r} }

type queryResolver struct{ *Resolver }
type config_ConfigResolver struct{ *Resolver }
type gns_GnsResolver struct{ *Resolver }
type root_RootResolver struct{ *Resolver }
