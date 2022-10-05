package crd_generator

import (
	"fmt"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

func CustomQueryToGraphqlSchema(query nexus.GraphQLQuery) string {
	args := ""
	argsList, ok := query.Args.([]parser.GraphQlArg)
	if ok && len(argsList) > 0 {
		args = "(\n"
		for _, arg := range argsList {
			args += fmt.Sprintf("        %s: %s\n", arg.Name, convertGraphqlStdType(arg.Type))
		}
		args += "    )"
	}
	return fmt.Sprintf("    query%s"+args+": NexusGraphqlResponse\n", query.Name)
}
