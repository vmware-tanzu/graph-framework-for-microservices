package crd_generator

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

func CustomQueryToGraphqlSchema(query nexus.GraphQLQuery) string {
	args := ""
	argsList, ok := query.Args.([]parser.GraphQlArg)
	if ok && len(argsList) > 0 {
		args = "(\n"
		for _, arg := range argsList {
			graphqlType := convertGraphqlStdType(arg.Type)
			if graphqlType == "" {
				log.Fatalf("Failed to convert type %s to graphql types, supported types are: "+
					"string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, "+
					"float32, float64", arg.Type)
			}

			args += fmt.Sprintf("        %s: %s\n", arg.Name, graphqlType)
		}
		args += "    )"
	}
	return fmt.Sprintf("    %s"+args+": NexusGraphqlResponse\n", query.Name)
}
