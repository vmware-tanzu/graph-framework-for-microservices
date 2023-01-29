package main

import (
	"net/http"

	"../../example/test-utils/output-group-name-with-hyphen-datamodel/crd_generated/nexus-gql/graph"
	"../../example/test-utils/output-group-name-with-hyphen-datamodel/crd_generated/nexus-gql/graph/generated"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen/graphql"
	"github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen/graphql/handler"
	"github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen/graphql/playground"
)

func StartHttpServer() {
	initNCErr := graph.InitNexusClientSet()
	if initNCErr != nil {
		log.Errorf("Error initializing nexus client in StartHttpServer: %s", initNCErr)
	}
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
