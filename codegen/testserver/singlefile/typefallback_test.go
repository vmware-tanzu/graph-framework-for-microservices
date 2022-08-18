package singlefile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/client"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql/handler"
)

func TestTypeFallback(t *testing.T) {
	resolvers := &Stub{}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	resolvers.QueryResolver.Fallback = func(ctx context.Context, arg FallbackToStringEncoding) (FallbackToStringEncoding, error) {
		return arg, nil
	}

	t.Run("fallback to string passthrough", func(t *testing.T) {
		var resp struct {
			Fallback string
		}
		c.MustPost(`query { fallback(arg: A) }`, &resp)
		require.Equal(t, "A", resp.Fallback)
	})
}
