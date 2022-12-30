package parser

import (
	"go/ast"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

func ParseGraphqlQuerySpecs(pkgs Packages) map[string]nexus.GraphQLQuerySpec {
	graphQLQueryMap := make(map[string]nexus.GraphQLQuerySpec)
	for _, pkg := range pkgs {
		GetGraphqlQuerySpecs(graphQLQueryMap, pkg)
	}
	return graphQLQueryMap
}

func GetGraphqlQuerySpecs(queryMap map[string]nexus.GraphQLQuerySpec, p Package) {
	for _, spec := range GetNexusSpecs(p, "nexus.GraphQLQuerySpec") {
		queryMap[p.Name+"."+spec.Name] = parseQuerySpec(spec.Value, p)
	}
}

func parseQuerySpec(v *ast.CompositeLit, p Package) nexus.GraphQLQuerySpec {
	spec := nexus.GraphQLQuerySpec{}
	for _, queryElt := range v.Elts {
		querykv, ok := queryElt.(*ast.KeyValueExpr)
		if !ok {
			log.Fatalf("Wrong format of graphql query spec, please check graphql spec")
		}
		val, ok := querykv.Value.(*ast.CompositeLit)
		if !ok {
			log.Fatalf("Wrong format of graphql query spec value, please check graphql spec")
		}
		for _, querySpecElt := range val.Elts {
			queryComp, ok := querySpecElt.(*ast.CompositeLit)
			if !ok {
				log.Fatalf("Wrong format of graphql query spec field, please check graphql spec")
			}

			newQuery := parseQuery(queryComp, p)
			spec.Queries = append(spec.Queries, newQuery)
		}
	}
	return spec
}

func parseQuery(queryComp *ast.CompositeLit, p Package) (newQuery nexus.GraphQLQuery) {
	newQuery = nexus.GraphQLQuery{}
	for _, queryFieldCompElt := range queryComp.Elts {
		queryFieldExp, ok := queryFieldCompElt.(*ast.KeyValueExpr)
		if !ok {
			log.Fatalf("Wrong format of graphql query spec field element, please check graphql spec")
		}

		queryFieldName := queryFieldExp.Key.(*ast.Ident)
		switch queryFieldName.String() {
		case "Name":
			queryFieldValue := queryFieldExp.Value.(*ast.BasicLit)
			name, err := strconv.Unquote(queryFieldValue.Value)
			if err != nil {
				log.Fatalf("Internal compiler error, failed to unqote name in graphql")
			}
			newQuery.Name = name
		case "ServiceEndpoint":
			queryFieldValue := queryFieldExp.Value.(*ast.CompositeLit)
			for _, serviceEndpointField := range queryFieldValue.Elts {
				serviceEndpointFieldKeyKv := serviceEndpointField.(*ast.KeyValueExpr)
				serviceEndpointFieldKey := serviceEndpointFieldKeyKv.Key.(*ast.Ident)
				if serviceEndpointFieldKey.String() == "Port" {
					serviceEndpointFieldValue := serviceEndpointFieldKeyKv.Value.(*ast.BasicLit)
					port, err := strconv.Atoi(serviceEndpointFieldValue.Value)
					if err != nil {
						log.Fatalf("Internal compiler error, failed to parse port to int in graphql spec")
					}
					newQuery.ServiceEndpoint.Port = port
				}
				if serviceEndpointFieldKey.String() == "Domain" {
					serviceEndpointFieldValue := serviceEndpointFieldKeyKv.Value.(*ast.BasicLit)
					domain, err := strconv.Unquote(serviceEndpointFieldValue.Value)
					if err != nil {
						log.Fatalf("Internal compiler error, failed to unqote domain in graphql")
					}
					newQuery.ServiceEndpoint.Domain = domain
				}
			}
		case "Args":
			queryFieldValue, ok := queryFieldExp.Value.(*ast.CompositeLit)
			if !ok {
				continue
			}
			typ, ok := queryFieldValue.Type.(*ast.Ident)
			if !ok {
				log.Fatalf("Graphql query args must not be imported, wrong type: %v", queryFieldValue.Type)
			}
			// translate args to map[arg.fieldName]arg.type
			args := parseArgs(typ.Name, p)
			newQuery.Args = args
		case "ApiType":
			selExpr, ok := queryFieldExp.Value.(*ast.SelectorExpr)
			if !ok {
				log.Fatalf("Failed to parse ApiType param in graphql custom query")
			}
			sel := selExpr.Sel.String()
			switch sel {
			case "GraphQLQueryApi":
				newQuery.ApiType = nexus.GraphQLQueryApi
			case "GetMetricsApi":
				newQuery.ApiType = nexus.GetMetricsApi
			default:
				newQuery.ApiType = nexus.GraphQLQueryApi
			}
		}
	}
	return
}

type GraphQlArg struct {
	Name      string
	Type      string
	AliasType bool
}

func parseArgs(argsTypeName string, p Package) []GraphQlArg {
	args := make([]GraphQlArg, 0)
	for _, decl := range p.GenDecls {
		for _, spec := range decl.Specs {
			v, ok := spec.(*ast.TypeSpec)
			if !ok || v.Name.Name != argsTypeName {
				continue
			}
			for _, field := range GetSpecFields(v) {
				if len(field.Names) == 0 {
					log.Fatalf("Field in graphql args must be named, args %s", argsTypeName)
				}
				// AliasName Annotation
				var fName, fType string
				var aType bool
				if val := GetFieldAnnotationVal(field, GRAPHQL_ALIAS_NAME_ANNOTATION); val != "" {
					fName = val
				} else {
					fName = field.Names[0].Name
				}
				// AliasType Annotation
				if val := GetFieldAnnotationVal(field, GRAPHQL_ALIAS_TYPE_ANNOTATION); val != "" {
					fType = val
					aType = true
				} else {
					fType = GetFieldType(field)
				}
				args = append(args, GraphQlArg{
					Name:      fName,
					Type:      fType,
					AliasType: aType,
				})
			}
		}
	}
	return args
}
