package graphqlcalls

import (
	"context"
	"log"

	"github.com/Khan/genqlient/graphql"
)

type GraphqlFuncs struct {
	Gclient        graphql.Client
	GraphqlFuncMap map[string]func(context.Context)
}

// function keys
const (
	// graphql function keys
	GET_MANAGERS      = "get_managers"
	GET_EMPLOYEE_ROLE = "get_employee_role"
)

// add the map of function keys to function calls
func (g *GraphqlFuncs) Init() {
	if g.Gclient == nil {
		log.Fatal("Graphql client initialization failed. Gclient is empty")
	}
	g.GraphqlFuncMap = make(map[string]func(context.Context))
	funcMap := g.GraphqlFuncMap
	funcMap[GET_MANAGERS] = g.GetManagers
	funcMap[GET_EMPLOYEE_ROLE] = g.GetEmployeeRole
}
