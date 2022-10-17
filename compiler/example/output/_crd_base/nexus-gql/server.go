package main

import (
	"net/http"

	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/handler"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/playground"
	"nexustempmodule/nexus-gql/graph"
	"nexustempmodule/nexus-gql/graph/generated"
)

func StartHttpServer() {
	ES := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	Hander_server := handler.NewDefaultServer(ES)
	HttpHandlerFunc := playground.Handler("GraphQL playground", "/query")
	http.Handle("/", HttpHandlerFunc)
	http.Handle("/query", Hander_server)
}
