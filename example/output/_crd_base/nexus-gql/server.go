package main

import (
	"nexustempmodule/nexus-gql/graph"
	"nexustempmodule/nexus-gql/graph/generated"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql"
)

var ES graphql.ExecutableSchema

func NewResolverObject() {
	ES = generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
}
