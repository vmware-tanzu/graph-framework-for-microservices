package main

import (
	"flag"
	"os"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser/rest"

	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/config"
	crd_generator "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/crd-generator"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
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
		conf.CrdModulePath = "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/"
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
	graph := parser.ParseDSLNodes(*dslDir, conf.GroupName)
	methods, codes := rest.ParseResponses(pkgs)
	if err = crd_generator.RenderCRDTemplate(conf.GroupName, conf.CrdModulePath, pkgs, graph,
		*crdDir, methods, codes); err != nil {
		log.Fatalf("Error rendering crd template: %v", err)
	}
}
