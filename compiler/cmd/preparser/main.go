package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/preparser"
)

func main() {
	dslDir := flag.String("dsl", "datamodel", "DSL file location.")
	outputDir := flag.String("output", "_generated", "output dir location.")
	modPath := flag.String("modpath", "datamodel", "ModPath for rendered imports")
	logLevel := flag.String("log-level", "ERROR", "Log level")
	flag.Parse()

	lvl, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Failed to configure logging: %v\n", err)
	}
	log.SetLevel(lvl)

	packages := preparser.Parse(*dslDir)
	err = preparser.Render(*dslDir, packages)
	if err != nil {
		log.Fatal(err)
	}

	err = preparser.CopyPkgsToBuild(*dslDir, *outputDir)
	if err != nil {
		log.Fatal(err)
	}

	packages = preparser.Parse(*dslDir)
	err = preparser.RenderImports(packages, *outputDir, *modPath)
	if err != nil {
		log.Fatal(err)
	}
}
