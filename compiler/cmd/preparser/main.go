package main

import (
	"flag"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/preparser"
)

func main() {
	configFile := flag.String("config-file", "", "Config file location.")
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

	conf := &config.Config{}
	if *configFile != "" {
		conf, err = config.LoadConfig(*configFile)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
	}
	config.ConfigInstance = conf
	packages := preparser.Parse(*dslDir)
	err = preparser.Render(*dslDir, packages)
	if err != nil {
		log.Fatal(err)
	}

	packages = preparser.Parse(*dslDir)
	err = preparser.RenderImports(packages, *outputDir, *modPath)
	if err != nil {
		log.Fatal(err)
	}

	err = preparser.CopyPkgsToBuild(packages, *outputDir)
	if err != nil {
		log.Fatal(err)
	}
}
