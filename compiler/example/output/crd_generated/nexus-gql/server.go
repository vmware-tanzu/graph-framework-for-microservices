package main

import (
	"net/http"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/nexus-gql/graph"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/nexus-gql/graph/generated"

	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/handler"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/playground"
)

func StartHttpServer() {
	ES := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	Hander_server := handler.NewDefaultServer(ES)
	HttpHandlerFunc := playground.Handler("GraphQL playground", "/query")
	http.Handle("/", HttpHandlerFunc)
	http.Handle("/query", Hander_server)
}
