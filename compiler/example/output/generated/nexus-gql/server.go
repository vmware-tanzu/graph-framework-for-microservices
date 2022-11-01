package main

import (
	"net/http"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generatednexus-gql/graph"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/output/crd_generatednexus-gql/graph/generated"

	"github.com/rs/cors"
	"github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen/graphql/handler"
	"github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen/graphql/playground"
)

func StartHttpServer() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	ES := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	Hander_server := handler.NewDefaultServer(ES)
	HttpHandlerFunc := playground.Handler("GraphQL playground", "/apis/graphql/v1/query")
	http.Handle("/", HttpHandlerFunc)
	http.Handle("/query", c.Handler(Hander_server))
}
