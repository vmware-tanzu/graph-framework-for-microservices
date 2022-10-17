package main

import (
	"log"
	"net/http"

	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/_examples/selection"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/handler"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/playground"
)

func main() {
	http.Handle("/", playground.Handler("Selection Demo", "/query"))
	http.Handle("/query", handler.NewDefaultServer(selection.NewExecutableSchema(selection.Config{Resolvers: &selection.Resolver{}})))
	log.Fatal(http.ListenAndServe(":8086", nil))
}
