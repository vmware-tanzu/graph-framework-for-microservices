package graphqlcalls

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

var GraphqlFuncMap map[string]func(context.Context, graphql.Client)

// function keys
const (
	// graphql function keys
	GET_MANAGERS      = "get_managers"
	GET_EMPLOYEE_ROLE = "get_employee_role"
)

// add the map of function keys to function calls
func init() {
	GraphqlFuncMap = make(map[string]func(context.Context, graphql.Client))
	funcMap := GraphqlFuncMap
	funcMap[GET_MANAGERS] = GetManagers
	funcMap[GET_EMPLOYEE_ROLE] = GetEmployeeRole
}
