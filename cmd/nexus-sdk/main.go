package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nexus/compiler/pkg/config"
	crd_generator "gitlab.eng.vmware.com/nexus/compiler/pkg/crd-generator"
	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
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

	// env overwrites config
	envPath := os.Getenv("CRD_MODULE_PATH")
	if envPath != "" {
		conf.CrdModulePath = envPath
	}
	if conf.CrdModulePath == "" {
		log.Fatalf("failed to determine Crd module path, please add to config file as" +
			" crdModulePath or as CRD_MODULE_PATH enviroment variable")
	}
	envGroup := os.Getenv("GROUP_NAME")
	if envGroup != "" {
		conf.GroupName = envGroup
	}
	if conf.GroupName == "" {
		log.Fatalf("failed to determine CRD group name, please add to config file as" +
			" groupName or as GROUP_NAME enviroment variable")
	}

	pkgs := parser.ParseDSLPkg(*dslDir)
	if err = crd_generator.RenderCRDTemplate(conf.GroupName, conf.CrdModulePath, pkgs, *crdDir); err != nil {
		log.Fatalf("Error rendering crd template: %v", err)
	}
}
