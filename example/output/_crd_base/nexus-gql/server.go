package main

import (
	"nexustempmodule/nexus-gql/graph"
	"nexustempmodule/nexus-gql/graph/generated"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql/handler"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql/playground"
	"net/http"
)


func StartHttpServer() {
	ES := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	Hander_server := handler.NewDefaultServer(ES)
	HttpHandlerFunc = playground.Handler("GraphQL playground", "/query")
	http.Handle("/", HttpHandlerFunc)
	http.Handle("/query", Hander_server)
}


