package rest

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
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

func GetHttpMethodsResponses(p parser.Package, codes map[string]HTTPCodesResponse) map[string]HTTPMethodsResponses {
	responses := make(map[string]HTTPMethodsResponses)

	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				name := valueSpec.Names[0].Name
				if value, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
					if types.ExprString(value.Type) != "HTTPMethodsResponses" {
						continue
					}

					response := HTTPMethodsResponses{}
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

func extractHttpMethodsKey(key ast.Expr) HTTPMethod {
	switch k := key.(type) {
	case *ast.SelectorExpr:
		return HTTPMethod(httpMethods[k.Sel.String()])
	}

	return ""
}

func extractHttpMethodsValue(value ast.Expr, codes map[string]HTTPCodesResponse) HTTPCodesResponse {
	switch val := value.(type) {
	case *ast.Ident:
		return codes[val.String()]
	}

	return nil
}
