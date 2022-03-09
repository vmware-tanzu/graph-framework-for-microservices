package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
)

func main() {
	// Package parser
	pkgs := parser.ParseDSLPkg("../../example/datamodel/")

	for key, pkg := range pkgs {
		if pkg.Name != "root" {
			continue
		}

		fmt.Println(key)

		for _, node := range pkg.GetNexusNodes() {
			cfg, err := parser.GetNexusNodeConfig(pkg, node.Name.String())
			if err != nil {
				panic(err)
			}
			spew.Dump(cfg)

			childFields := parser.GetChildFields(node)

			fmt.Printf("Child fields:\n")
			for _, f := range childFields {
				name, err := parser.GetFieldName(f)
				if err != nil {
					panic(err)
				}
				t := parser.GetFieldType(f)
				isMap := parser.IsMapField(f)
				fmt.Printf(" - %s:%s isMap:%t\n", name, t, isMap)
			}

		}

	}
}
