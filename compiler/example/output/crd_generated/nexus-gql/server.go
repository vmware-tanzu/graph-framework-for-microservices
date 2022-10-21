package main

import (
	"net/http"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/nexus-gql/graph"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/crd_generated/nexus-gql/graph/generated"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql/handler"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql/playground"
)

func StartHttpServer() {
	ES := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	Hander_server := handler.NewDefaultServer(ES)
	HttpHandlerFunc := playground.Handler("GraphQL playground", "/query")
	http.Handle("/", HttpHandlerFunc)
	http.Handle("/query", Hander_server)
}
