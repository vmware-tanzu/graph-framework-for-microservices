package main

import (
	"net/http"

	"nexustempmodule/nexus-gql/graph"
	"nexustempmodule/nexus-gql/graph/generated"

	"github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen/graphql/handler"
	"github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen/graphql/playground"
)

func StartHttpServer() {
	ES := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	Hander_server := handler.NewDefaultServer(ES)
	HttpHandlerFunc := playground.Handler("GraphQL playground", "/query")
	http.Handle("/", HttpHandlerFunc)
	http.Handle("/query", Hander_server)
}
