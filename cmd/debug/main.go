package main

import (
	"encoding/json"
	"fmt"
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
	methods, codes := rest.ParseResponses(pkgs)
	log.Println(methods)

	gns := pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns"]

	apiSpecs := rest.GetRestApiSpecs(gns, methods, codes)

	b, err := json.MarshalIndent(apiSpecs, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
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
