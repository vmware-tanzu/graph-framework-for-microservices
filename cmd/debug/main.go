package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

var root parser.Node

func main() {
	baseGroupName := "tsm.tanzu.vmware.com"
	nodes := parser.ParseDSLNodes("../../example/datamodel/", baseGroupName)

	root = nodes["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/Root"]

	root.Walk(func(node *parser.Node) {
		log.Printf("CRD name: %s\n", node.CrdName)
		log.Printf("Parents: %v\n", node.Parents)
		log.Println("---")
	})

	parents := parser.CreateParentsMap(root)
	log.Println(parents["acpconfigs.tsm.tanzu.vmware.com"])
}
