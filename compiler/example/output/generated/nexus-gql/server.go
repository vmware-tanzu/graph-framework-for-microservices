package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/nexus-gql/graph"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/nexus-gql/graph/generated"

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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	StartHttpServer()
	srv := &http.Server{Addr: fmt.Sprintf(":%s", port)}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Error in starting graphql server")
	}
}
