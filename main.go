package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"plugin"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql/handler"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/gqlgen.git/graphql/playground"
)

const defaultPort = "8080"

// Defaulting to build.so for local debug
const defaultPath = "graphql.so"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// PLUGIN_PATH environment variable will be part of the deployment spec: https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-runtime-manifests/-/blob/add_graphql/core/templates/datamodel_installer.yaml#L172
	// PLUGIN will be dynamically unarchieved from datamodel image using a init containerhttps://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-runtime-manifests/-/blob/add_graphql/core/templates/datamodel_installer.yaml#L112
	// defaulting to graphql.so for debugging purposes
	graphqlBuildplugin := os.Getenv("PLUGIN_PATH")
	if graphqlBuildplugin == "" {
		graphqlBuildplugin = "graphql.so"
	}

	if _, err := os.Stat(graphqlBuildplugin); err != nil {
		fmt.Printf("error in checking graphql plugin file %s", graphqlBuildplugin)
		panic(err)
	}
	// Opening graphql plugin file archieved from datamodel image
	pl, err := plugin.Open(graphqlBuildplugin)
	if err != nil {
		fmt.Printf("could not open pluginfile: %s", err)
		panic(err)
	}

	// Lookup resolver object present
	esvar, err := pl.Lookup("ES")
	if err != nil {
		fmt.Printf("could not lookup the graphqlExecutable object: %s", err)
		panic(err)
	}
	// Lookup init method present
	plsm, err := pl.Lookup("NewResolverObject")
	if err != nil {
		fmt.Printf("could not lookup the InitMethod : %s", err)
		panic(err)
	}
	// Execute the init method for initialising resolvers and typecast to expected format
	plsm.(func())()
	esObject := *esvar.(*graphql.ExecutableSchema)

	srv := handler.NewDefaultServer(esObject)

	//Start HTTP router with graphql router and proxy it to resolver
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
