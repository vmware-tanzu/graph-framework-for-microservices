package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nexus/compiler/pkg/config"
	crd_generator "gitlab.eng.vmware.com/nexus/compiler/pkg/crd-generator"
	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
)

func main() {
	configFile := flag.String("config-file", "nexus-sdk.yaml", "Config file location.")
	dslDir := flag.String("dsl", "datamodel", "DSL file location.")
	crdDir := flag.String("crd-output", "crds", "CRD file location.")
	logLevel := flag.String("log-level", "ERROR", "Log level")
	flag.Parse()

	lvl, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Failed to configure logging: %v\n", err)
	}
	log.SetLevel(lvl)

	conf, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	pkgs := parser.ParseDSLPkg(*dslDir)
	if err = crd_generator.RenderCRDTemplate(conf.GroupName, conf.CrdModulePath, pkgs, *crdDir); err != nil {
		log.Fatalf("Error rendering crd template: %v", err)
	}
}
