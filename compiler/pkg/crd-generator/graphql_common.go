package crd_generator

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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

func validateImportPkg(pkgName, typeString string, importMap map[string]string) (string, string) {
	typeWithoutPointers := strings.ReplaceAll(typeString, "*", "")
	if strings.Contains(typeWithoutPointers, ".") {
		part := strings.Split(typeWithoutPointers, ".")
		if val, ok := importMap[part[0]]; ok {
			pkgName := val[strings.LastIndex(val, "/")+1 : len(val)-1]
			repName := strings.ReplaceAll(pkgName, "-", "")
			return repName + "_" + part[1], cases.Title(language.Und, cases.NoLower).String(repName) + cases.Title(language.Und, cases.NoLower).String(part[1])
		}
		return strings.ToLower(pkgName) + part[1], cases.Title(language.Und, cases.NoLower).String(pkgName) + cases.Title(language.Und, cases.NoLower).String(part[1])
	}
	return strings.ToLower(pkgName) + "_" + typeWithoutPointers, cases.Title(language.Und, cases.NoLower).String(pkgName) + cases.Title(language.Und, cases.NoLower).String(typeWithoutPointers)
}

func getBaseNodeType(typeString string) string {
	if strings.Contains(typeString, ".") {
		part := strings.Split(typeString, ".")
		return part[1]
	}
	return typeString
}

func getPkgName(pkg parser.Package) string {
	if pkg.FullName == pkg.ModPath {
		return pkg.Name
	}
	pkgPath := pkg.FullName
	libName := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
	return strings.ReplaceAll(libName, "-", "")
}
