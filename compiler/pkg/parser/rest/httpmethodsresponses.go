package rest

import (
	"go/ast"
	"go/types"

	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

var httpMethods = map[string]string{
	"MethodGet":     "GET",
	"MethodHead":    "HEAD",
	"MethodPost":    "POST",
	"MethodPut":     "PUT",
	"MethodPatch":   "PATCH", // RFC 5789
	"MethodDelete":  "DELETE",
	"MethodConnect": "CONNECT",
	"MethodOptions": "OPTIONS",
	"MethodTrace":   "TRACE",
	"LIST":          "LIST",
}

func GetHttpMethodsResponses(p parser.Package, responsesMap map[string]nexus.HTTPMethodsResponses, httpCodes map[string]nexus.HTTPCodesResponse) map[string]nexus.HTTPMethodsResponses {
	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				name := valueSpec.Names[0].Name
				if valueSpec.Values == nil {
					continue
				}

				if value, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
					if types.ExprString(value.Type) != "nexus.HTTPMethodsResponses" {
						continue
					}

					response := nexus.HTTPMethodsResponses{}
					for _, elt := range value.Elts {
						kv := elt.(*ast.KeyValueExpr)
						responseKey := extractHttpMethodsKey(kv.Key)
						responseValue := extractHttpMethodsValue(kv.Value, httpCodes)
						response[responseKey] = responseValue
					}
					responsesMap[name] = response
				}
			}
		}
	}

	return responsesMap
}

func extractHttpMethodsKey(key ast.Expr) nexus.HTTPMethod {
	switch k := key.(type) {
	case *ast.SelectorExpr:
		return nexus.HTTPMethod(httpMethods[k.Sel.String()])
	}

	return ""
}

func extractHttpMethodsValue(value ast.Expr, httpCodes map[string]nexus.HTTPCodesResponse) nexus.HTTPCodesResponse {
	switch val := value.(type) {
	case *ast.Ident:
		return httpCodes[val.String()]
	case *ast.SelectorExpr:
		return httpCodes[val.Sel.String()]
	case *ast.CompositeLit:
		return extractHttpCodesResponse(val)
	}

	return nil
}
