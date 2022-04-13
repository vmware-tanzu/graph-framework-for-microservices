package parser

import (
	"go/ast"
)

func (p *Package) GetRestApiSpecVars(restAPIMap map[string]RestAPISpec) {
	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				name := valueSpec.Names[0].Name
				if value, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
					var restUris []RestURIs
					for _, val := range value.Elts {
						uris := val.(*ast.KeyValueExpr)
						for _, x := range uris.Value.(*ast.CompositeLit).Elts {

							var restSpec RestURIs
							for _, y := range x.(*ast.CompositeLit).Elts {
								field := y.(*ast.KeyValueExpr)
								key := field.Key.(*ast.Ident)
								if key.String() == "Uri" {
									restSpec.Uri = field.Value.(*ast.BasicLit).Value
									continue
								}

								//for _, z := range field.Value.(*ast.CompositeLit).Elts {
								//	//method := z.(*ast.SelectorExpr).Sel
								//	//m := GetMethods(method.Name)
								//	//restSpec.Methods = append(restSpec.Methods, m)
								//}
								restUris = append(restUris, restSpec)
							}
						}
						restAPIMap[name] = RestAPISpec{
							Uris: restUris,
						}
					}
				}
			}
		}
	}
}

//
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
