package parser

import (
	"fmt"
	"go/ast"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

func ParseGraphqlSpecs(pkgs Packages) map[string]nexus.GraphQLSpec {
	graphQLSpecMap := make(map[string]nexus.GraphQLSpec)
	for _, pkg := range pkgs {
		GetGraphqlSpecs(graphQLSpecMap, pkg)
	}
	return graphQLSpecMap
}

func GetGraphqlSpecs(graphQLSpecMap map[string]nexus.GraphQLSpec, p Package) {
	for _, spec := range GetNexusSpecs(p, "nexus.GraphQLSpec") {
		graphQLSpecMap[p.Name+"."+spec.Name] = parseGqlSpec(spec.Value, p)
	}
}

func parseGqlSpec(v *ast.CompositeLit, p Package) nexus.GraphQLSpec {
	spec := nexus.GraphQLSpec{
		IdName:     "",
		IdNullable: true,
	}
	for _, gqlSpecElt := range v.Elts {
		gqlSpecKv, ok := gqlSpecElt.(*ast.KeyValueExpr)
		if !ok {
			log.Fatalf("Wrong format of graphql query spec, please check graphql spec")
		}

		gqlSpecFieldName := gqlSpecKv.Key.(*ast.Ident)
		fmt.Println("1111****++++++>>>>KEY:", gqlSpecFieldName.String())
		switch gqlSpecFieldName.String() {
		case "IdName":
			gqlSpecFieldValue := gqlSpecKv.Value.(*ast.BasicLit)
			name, err := strconv.Unquote(gqlSpecFieldValue.Value)
			if err != nil {
				log.Fatalf("Internal compiler error, failed to unqote name in graphql")
			}
			fmt.Println("*****VAL:", name)
			spec.IdName = name
		case "IdNullable":
			gqlSpecFieldValue := gqlSpecKv.Value.(*ast.Ident)
			val := gqlSpecFieldValue.String()
			fmt.Println("*****VAL:", gqlSpecFieldValue.String())
			if val == "false" {
				spec.IdNullable = false
			} else {
				spec.IdNullable = true
			}
		}
	}
	fmt.Println("SPEC", spec)
	return spec
}
