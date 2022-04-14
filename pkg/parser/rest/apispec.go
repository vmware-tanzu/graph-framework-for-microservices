package rest

import (
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
	"go/ast"
	"go/types"
)

func GetRestApiSpecs(p parser.Package) map[string]nexus.RestAPISpec {
	apiSpecs := make(map[string]nexus.RestAPISpec)

	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				name := valueSpec.Names[0].Name
				if value, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
					if types.ExprString(value.Type) != "RestAPISpec" {
						continue
					}

					log.Println(name)
				}
			}
		}
	}
	//for _, genDecl := range p.GenDecls {
	//	for _, spec := range genDecl.Specs {
	//		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
	//			name := valueSpec.Names[0].Name
	//			if value, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
	//				var restUris []nexus.RestURIs
	//				for _, val := range value.Elts {
	//					uris := val.(*ast.KeyValueExpr)
	//
	//					for _, x := range uris.Value.(*ast.CompositeLit).Elts {
	//						var restSpec nexus.RestURIs
	//						for _, y := range x.(*ast.CompositeLit).Elts {
	//							field := y.(*ast.KeyValueExpr)
	//							key := field.Key.(*ast.Ident)
	//							if key.String() == "Uri" {
	//								restSpec.Uri = field.Value.(*ast.BasicLit).Value
	//								continue
	//							}
	//
	//							//for _, z := range field.Value.(*ast.CompositeLit).Elts {
	//							//	//method := z.(*ast.SelectorExpr).Sel
	//							//	//m := GetMethods(method.Name)
	//							//	//restSpec.Methods = append(restSpec.Methods, m)
	//							//}
	//							restUris = append(restUris, restSpec)
	//						}
	//					}
	//					apiSpecs[name] = nexus.RestAPISpec{
	//						Uris: restUris,
	//					}
	//				}
	//			}
	//		}
	//	}
	//}

	return apiSpecs
}

//func GetMethods(method string) string {
//	if method == "MethodGet" {
//		return "GET"
//	} else if method == "MethodPut" {
//		return "PUT"
//	} else if method == "MethodDelete" {
//		return "DELETE"
//	}
//	return ""
//}
