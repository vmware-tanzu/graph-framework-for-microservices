package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
)

func main() {
	nodes := parser.ParseDSLNodes("../../example/datamodel/")
	root := nodes["gitlab.eng.vmware.com/nexus/compiler/example/datamodel/Root"]
	root.Walk(func(node *parser.Node) {
		log.Println(node.Name)
		log.Println(node.Parents)
		log.Println("---")
	})
}
