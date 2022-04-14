package rest

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus.git/nexus"
	"go/ast"
	"go/types"
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
}

func GetHttpMethodsResponses(p parser.Package, codes map[string]nexus.HTTPCodesResponse) map[string]nexus.HTTPMethodsResponses {
	responses := make(map[string]nexus.HTTPMethodsResponses)

	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				name := valueSpec.Names[0].Name
				if value, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
					if types.ExprString(value.Type) != "HTTPMethodsResponses" {
						continue
					}

					response := nexus.HTTPMethodsResponses{}
					for _, elt := range value.Elts {
						kv := elt.(*ast.KeyValueExpr)
						responseKey := extractHttpMethodsKey(kv.Key)
						responseValue := extractHttpMethodsValue(kv.Value, codes)
						response[responseKey] = responseValue
					}
					responses[name] = response
				}
			}
		}
	}

	return responses
}

func extractHttpMethodsKey(key ast.Expr) nexus.HTTPMethod {
	switch k := key.(type) {
	case *ast.SelectorExpr:
		return nexus.HTTPMethod(httpMethods[k.Sel.String()])
	}

	return ""
}

func extractHttpMethodsValue(value ast.Expr, codes map[string]nexus.HTTPCodesResponse) nexus.HTTPCodesResponse {
	switch val := value.(type) {
	case *ast.Ident:
		return codes[val.String()]
	}

	return nil
}
