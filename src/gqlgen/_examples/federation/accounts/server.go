//go:generate go run ../../../testdata/gqlgen.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/_examples/federation/accounts/graph"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/_examples/federation/accounts/graph/generated"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/handler"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/handler/debug"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/playground"
)

const defaultPort = "4001"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.Use(&debug.Tracer{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
