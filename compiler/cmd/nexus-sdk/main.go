package main

import (
	"flag"
	"os"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser/rest"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

func main() {
	configFile := flag.String("config-file", "", "Config file location.")
	dslDir := flag.String("dsl", "datamodel", "DSL file location.")
	crdDir := flag.String("crd-output", "_generated", "CRD file location.")
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

	if conf.CrdModulePath == "" {
		conf.CrdModulePath = "nexustempmodule/"
	}
	// env overwrites config
	envGroup := os.Getenv("GROUP_NAME")
	if envGroup != "" {
		conf.GroupName = envGroup
	}
	if conf.GroupName == "" {
		log.Fatalf("failed to determine CRD group name, please add to config file as" +
			" groupName or as GROUP_NAME enviroment variable")
	}

	config.ConfigInstance = conf
	pkgs := parser.ParseDSLPkg(*dslDir)
	graphlqQueries := parser.ParseGraphqlQuerySpecs(pkgs)
	graph, nonNexusTypes, fileset := parser.ParseDSLNodes(*dslDir, conf.GroupName, pkgs, graphlqQueries)
	methods, codes := rest.ParseResponses(pkgs)
	if err = generator.RenderCRDTemplate(conf.GroupName, conf.CrdModulePath, pkgs, graph,
		*crdDir, methods, codes, nonNexusTypes, fileset); err != nil {
		log.Fatalf("Error rendering crd template: %v", err)
	}
}
