package generator

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

func CustomQueryToGraphqlSchema(query nexus.GraphQLQuery) string {
	args := ""
	argsList, ok := query.Args.([]parser.GraphQlArg)
	if ok && len(argsList) > 0 {
		args = "(\n"
		for _, arg := range argsList {
			graphqlType := convertGraphqlStdType(arg.Type)
			// AliasType is to over write arg type with annotation `nexus-alias-type:""`
			if arg.AliasType != "" {
				graphqlType = arg.AliasType
			}
			if graphqlType == "" {
				log.Fatalf("Failed to convert type %s to graphql types, supported types are: "+
					"string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, "+
					"float32, float64", arg.Type)
			}

			args += fmt.Sprintf("        %s: %s\n", arg.Name, graphqlType)
		}
		args += "    )"
	}

	var returnType string
	switch query.ApiType {
	case nexus.GraphQLQueryApi:
		returnType = "NexusGraphqlResponse"
	case nexus.GetMetricsApi:
		returnType = "TimeSeriesData"
	default:
		log.Fatalf("Wrong Api Type of Graphql custom query")
	}

	return fmt.Sprintf("    %s"+args+": "+returnType+"\n", query.Name)
}
