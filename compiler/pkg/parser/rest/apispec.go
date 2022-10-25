package rest

import (
	"go/ast"
	"go/types"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

var uris = map[string]string{}

func GetRestApiSpecs(p parser.Package, httpMethods map[string]nexus.HTTPMethodsResponses,
	httpCodes map[string]nexus.HTTPCodesResponse, parentsMap map[string]parser.NodeHelper) map[string]nexus.RestAPISpec {
	apiSpecs := make(map[string]nexus.RestAPISpec)

	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				name := valueSpec.Names[0].Name
				if valueSpec.Values == nil {
					continue
				}
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
		case "QueryParams":
			restUri.QueryParams = extractApiSpecQueryParams(kv)
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

func extractApiSpecQueryParams(kv *ast.KeyValueExpr) []string {
	var params []string
	switch val := kv.Value.(type) {
	case *ast.CompositeLit:
		for _, v := range val.Elts {
			lit := v.(*ast.BasicLit)

			param, err := strconv.Unquote(lit.Value)
			if err != nil {
				log.Errorf("Error %v", err)
			}
			params = append(params, param)
		}
	}
	return params
}

func ValidateRestApiSpec(apiSpec nexus.RestAPISpec, parentsMap map[string]parser.NodeHelper, crdName string) {
	r := regexp.MustCompile(`{([^{}]+)}`)
	crdHelper := parentsMap[crdName]

	for _, uri := range apiSpec.Uris {
		uriRegex, _ := regexp.Compile("{.*?}")
		redactedUri := uriRegex.ReplaceAllString(uri.Uri, "{param}")

		if u, ok := uris[redactedUri]; ok {
			log.Fatalf("RestApiSpec: Duplicate found: %s and %s", u, uri.Uri)
		}

		uriParams := r.FindAllStringSubmatch(uri.Uri, -1)
		if _, ok := uri.Methods["LIST"]; ok {
			if nodeExist(crdHelper.RestName, uriParams) || queryParamExist(crdHelper.RestName, uri.QueryParams) {
				log.Fatalf("RestApiSpec: Provided node name (%s) cannot be applied as a param because endpoint is a list. URI: %s", crdHelper.RestName, uri.Uri)
			}
		}

		// Check if node name is in both URI and Query param
		// We are ignoring checking for node in URL because endpoint can be a list, and then we don't need this param
		if nodeExist(crdHelper.RestName, uriParams) && queryParamExist(crdHelper.RestName, uri.QueryParams) {
			log.Fatalf("RestApiSpec: Provided node name (%s) cannot be applied to both URI Param and Query Param. URI: %s", crdHelper.RestName, uri.Uri)
		}

		for _, parentCrd := range crdHelper.Parents {
			parentCrdHelper := parentsMap[parentCrd]
			parentName := parentCrdHelper.RestName

			if parentCrdHelper.IsSingleton {
				continue
			}

			if nodeExist(parentName, uriParams) && queryParamExist(parentName, uri.QueryParams) {
				log.Fatalf("RestApiSpec: Provided node name (%s) cannot be applied to both URI Param and Query Param. URI: %s", parentName, uri.Uri)
			}

			if !nodeExist(parentName, uriParams) && !queryParamExist(parentName, uri.QueryParams) {
				log.Fatalf("RestApiSpec: Provided node name (%s) not found for uri: %s", parentName, uri.Uri)
			}
		}

		uris[redactedUri] = uri.Uri
	}
}

func nodeExist(name string, params [][]string) bool {
	for _, p := range params {
		if p[1] == name {
			return true
		}
	}

	return false
}

func queryParamExist(name string, params []string) bool {
	for _, p := range params {
		if p == name {
			return true
		}
	}

	return false
}
