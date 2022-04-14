package rest

import (
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
	"go/ast"
	"go/types"
	"strconv"
)

func GetRestApiSpecs(p parser.Package, httpMethods map[string]nexus.HTTPMethodsResponses, httpCodes map[string]nexus.HTTPCodesResponse) map[string]nexus.RestAPISpec {
	apiSpecs := make(map[string]nexus.RestAPISpec)

	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				name := valueSpec.Names[0].Name
				if value, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
					if types.ExprString(value.Type) != "nexus.RestAPISpec" {
						continue
					}

					apiSpec := nexus.RestAPISpec{}
					for _, elt := range value.Elts {
						uris := elt.(*ast.KeyValueExpr)

						for _, uri := range uris.Value.(*ast.CompositeLit).Elts {
							restUri := extractApiSpecRestURI(uri.(*ast.CompositeLit), httpMethods, httpCodes)
							apiSpec.Uris = append(apiSpec.Uris, restUri)
						}
					}

					apiSpecs[name] = apiSpec
				}
			}
		}
	}

	return apiSpecs
}

func extractApiSpecRestURI(uri *ast.CompositeLit, httpMethods map[string]nexus.HTTPMethodsResponses, httpCodes map[string]nexus.HTTPCodesResponse) nexus.RestURIs {
	restUri := nexus.RestURIs{}
	for _, elt := range uri.Elts {
		kv := elt.(*ast.KeyValueExpr)

		switch types.ExprString(kv.Key) {
		case "Uri":
			key, err := strconv.Unquote(types.ExprString(kv.Value))
			if err != nil {
				log.Errorf("Error %v", err)
			}
			restUri.Uri = key
		case "Methods":
			restUri.Methods = extractApiSpecMethods(kv, httpMethods, httpCodes)
		}
	}

	return restUri
}

func extractApiSpecMethods(methods *ast.KeyValueExpr, httpMethods map[string]nexus.HTTPMethodsResponses, httpCodes map[string]nexus.HTTPCodesResponse) nexus.HTTPMethodsResponses {
	switch val := methods.Value.(type) {
	case *ast.Ident:
		return httpMethods[val.Name]
	case *ast.SelectorExpr:
		return httpMethods[val.Sel.String()]
	case *ast.CompositeLit:
		met := make(nexus.HTTPMethodsResponses)
		for _, elt := range val.Elts {
			kv := elt.(*ast.KeyValueExpr)
			httpKey := extractHttpMethodsKey(kv.Key)
			httpValue := extractHttpMethodsValue(kv.Value, httpCodes)
			met[httpKey] = httpValue
		}
		return met
	}
	return nil
}
