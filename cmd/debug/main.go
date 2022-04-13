package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser/rest"
)

var root parser.Node

func main() {
	// baseGroupName := "tsm.tanzu.vmware.com"
	// nodes := parser.ParseDSLNodes("../../example/datamodel/", baseGroupName)

	// root = nodes["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/Root"]

	// root.Walk(func(node *parser.Node) {
	// 	log.Printf("CRD name: %s\n", node.CrdName)
	// 	log.Printf("Parents: %v\n", node.Parents)
	// 	log.Println("---")
	// })

	pkgs := parser.ParseDSLPkg("../../example/datamodel/")
	nexus := pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"]

	responseCodes := rest.GetHttpCodesResponses(nexus)
	log.Println(responseCodes)
	//gns := pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns"]

	//mp := make(map[string]parser.RestAPISpec)
	//gns.GetRestApiSpecVars(mp)
	//
	//annotation, ok := parser.GetNexusRestAPIGenAnnotation(gns, "Gns")
	//if ok {
	//	log.Println(annotation)
	//	a := strings.Split(annotation, ":")
	//	name := a[1]
	//	value := mp[name]
	//	log.Println(value)
	//}
	//
	//annotation, ok = parser.GetNexusRestAPIGenAnnotation(gns, "Dns")
	//if ok {
	//	log.Println(annotation)
	//	a := strings.Split(annotation, ":")
	//	name := a[1]
	//	value := mp[name]
	//	log.Println(value)
	//}
}
