package parser

import (
	"go/ast"
	"strconv"

	log "github.com/sirupsen/logrus"
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

var GraphQLQueryMap = make(map[string]nexus.GraphQLQuerySpec)

func ParseGraphqlQuerySpecs(pkgs Packages) {
	for _, pkg := range pkgs {
		GetGraphqlQuerySpecs(GraphQLQueryMap, pkg)
	}
}

func GetGraphqlQuerySpecs(queryMap map[string]nexus.GraphQLQuerySpec, p Package) {
	for _, spec := range GetNexusSpecs(p, "nexus.GraphQLQuerySpec") {
		queryMap[p.Name+"."+spec.Name] = parseQuerySpec(spec.Value)
	}
}

func parseQuerySpec(v *ast.CompositeLit) nexus.GraphQLQuerySpec {
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

			newQuery := parseQuery(queryComp)
			spec.Queries = append(spec.Queries, newQuery)
		}
	}
	return spec
}

func parseQuery(queryComp *ast.CompositeLit) (newQuery nexus.GraphQLQuery) {
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
			if ok {
				newQuery.Args = queryFieldValue.Type
			}
		}
	}
	return
}
