package crd_generator

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

const (
	//GraphQL standard data types
	Int     string = "Int"
	Float   string = "Float"
	String  string = "String"
	Boolean string = "Boolean"

	// Custom query
	CustomQuerySchema = `Id: ID
	ParentLabels: Map
    queryServiceTable(
        startTime: String
        endTime: String
        SystemServices: Boolean
        ShowGateways: Boolean
        Groupby: String
        noMetrics: Boolean
    ): TimeSeriesData
    queryServiceVersionTable(
        startTime: String
        endTime: String
        SystemServices: Boolean
        ShowGateways: Boolean
        noMetrics: Boolean
    ): TimeSeriesData
    queryServiceTS(
        svcMetric: String
        startTime: String
        endTime: String
        timeInterval: String
    ): TimeSeriesData
    queryIncomingAPIs(
        startTime: String
        endTime: String
        destinationService: String
        destinationServiceVersion: String
        timeInterval: String
        timeZone: String
    ): TimeSeriesData
    queryOutgoingAPIs(
        startTime: String
        endTime: String
		destinationService: String
        destinationServiceVersion: String
        timeInterval: String
        timeZone: String
    ): TimeSeriesData
    queryIncomingTCP(
        startTime: String
        endTime: String
        destinationService: String
        destinationServiceVersion: String
    ): TimeSeriesData
    queryOutgoingTCP(
        startTime: String
        endTime: String
        destinationService: String
        destinationServiceVersion: String
    ): TimeSeriesData
    queryServiceTopology(
        metricStringArray: String
        startTime: String
        endTime: String
    ): TimeSeriesData
`
)

// Convert go standardType to GraphQL standardType
func convertGraphqlStdType(t string) string {
	// remove pointers
	typeWithoutPointer := strings.ReplaceAll(t, "*", "")
	switch typeWithoutPointer {
	case "string":
		return String
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return Int
	case "bool":
		return Boolean
	case "float32", "float64":
		return Float
	default:
		return ""
	}
}

func jsonMarshalResolver(FieldName, PkgName string) string {
	return fmt.Sprintf("%s, _ := json.Marshal(v%s.Spec.%s)\n%sData := string(%s)\n", FieldName, PkgName, FieldName, FieldName, FieldName)
}

func validateImportPkg(pkg parser.Package, typeString string, importMap map[string]string) (string, string) {
	typeWithoutPointers := strings.ReplaceAll(typeString, "*", "")
	if strings.Contains(typeWithoutPointers, ".") {
		part := strings.Split(typeWithoutPointers, ".")
		if val, ok := importMap[part[0]]; ok {
			pkgName := val[strings.LastIndex(val, "/")+1 : len(val)-1]
			repName := strings.ReplaceAll(pkgName, "-", "")
			return repName + "_" + part[1], cases.Title(language.Und, cases.NoLower).String(repName) + cases.Title(language.Und, cases.NoLower).String(part[1])
		}
		return pkg.Name + part[1], cases.Title(language.Und, cases.NoLower).String(pkg.Name) + cases.Title(language.Und, cases.NoLower).String(part[1])
	}
	return pkg.Name + "_" + typeWithoutPointers, cases.Title(language.Und, cases.NoLower).String(pkg.Name) + cases.Title(language.Und, cases.NoLower).String(typeWithoutPointers)
}

func getBaseNodeType(typeString string) string {
	if strings.Contains(typeString, ".") {
		part := strings.Split(typeString, ".")
		return part[1]
	}
	return typeString
}
