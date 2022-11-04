package generator

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
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

func getPkgName(pkgs parser.Packages, pkgPath string) string {
	importPath, err := strconv.Unquote(pkgPath)
	if err != nil {
		log.Errorf("Failed to parse the package path : %s: %v", pkgPath, err)
	}
	return pkgs[importPath].Name
}

func genSchemaResolverName(fn1, fn2 string) (string, string) {
	return fmt.Sprintf("%s_%s", strings.ToLower(fn1), fn2), cases.Title(language.Und, cases.NoLower).String(fn1) + cases.Title(language.Und, cases.NoLower).String(fn2)
}

func ValidateImportPkg(pkgName, typeString string, importMap map[string]string, pkgs parser.Packages) (string, string) {
	typeWithoutPointers := strings.ReplaceAll(typeString, "*", "")
	if strings.Contains(typeWithoutPointers, ".") {
		part := strings.Split(typeWithoutPointers, ".")
		if val, ok := importMap[part[0]]; ok {
			pkgName := getPkgName(pkgs, val)
			repName := strings.ReplaceAll(pkgName, "-", "")
			return genSchemaResolverName(repName, part[1])
		}
		for _, v := range importMap {
			pkgName := getPkgName(pkgs, v)
			if pkgName == part[0] {
				repName := strings.ReplaceAll(pkgName, "-", "")
				return genSchemaResolverName(repName, part[1])
			}
		}
		return genSchemaResolverName(pkgName, part[1])
	}
	return genSchemaResolverName(pkgName, typeWithoutPointers)
}

func getBaseNodeType(typeString string) string {
	if strings.Contains(typeString, ".") {
		part := strings.Split(typeString, ".")
		return part[1]
	}
	return typeString
}
