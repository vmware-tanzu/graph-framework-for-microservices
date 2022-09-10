package graphql_generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"regexp"
	"sort"
	"strings"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"
)

type Field_prop struct {
	IsResolver      bool
	PkgName         string
	NodeName        string
	FieldName       string
	FieldType       string
	SchemaFieldName string
}

type Node_prop struct {
	IsSingletonNode     bool
	Child               []Field_prop
	Children            []Field_prop
	Link                []Field_prop
	Links               []Field_prop
	MapFields           []Field_prop
	CustomFields        []Field_prop
	NonStructFields     []Field_prop
	GraphqlSchemaFields []Field_prop
	ResolverCount       int
	PkgName             string
	NodeName            string
	SchemaName          string
}

const (
	Int     string = "Int"
	Float   string = "Float"
	String  string = "String"
	Boolean string = "Boolean"
	ID      string = "String"
)

// Convert go standardType to GraphQL standardType
func convertGraphqlStdType(t string) string {
	switch t {
	case "string":
		return String
	case "int", "int32", "int64", "uint16":
		return Int
	case "bool":
		return Boolean
	case "float", "float32", "float64":
		return Float
	default:
		return ""
	}
}

func validateImportPkg(pkg parser.Package, typeString string, importMap map[string]string) string {
	if strings.Contains(typeString, ".") {
		part := strings.Split(typeString, ".")
		if val, ok := importMap[part[0]]; ok {
			pkgName := val[strings.LastIndex(val, "/")+1 : len(val)-1]
			repName := strings.ReplaceAll(pkgName, "-", "")
			return repName + "_" + part[1]
		}
		return "_" + part[1]
	}
	return pkg.Name + "_" + typeString
}

// TODO: https://jira.eng.vmware.com/browse/NPT-296
// Support cross-package imports for the following additional types:
// 1. map[gns.MyStr][]gns.MyStr
// 2. map[string]map[string]gns.MyStr
// 3. []map[string]gns.MyStr
// 4. **gns.MyStr
func ConstructType(aliasNameMap map[string]string, field *ast.Field) string {
	typeString := types.ExprString(field.Type)

	// Check if the field is imported from a different package.
	if !strings.Contains(typeString, ".") {
		return typeString
	}

	switch {
	case parser.IsMapField(field):
		// TODO: Check if the function GetFieldType(field) can be reused for cases other than:
		// map[string]gns.MyStr
		// https://jira.eng.vmware.com/browse/NPT-296
		mapParts := regexp.MustCompile(`^(map\[)`).ReplaceAllString(typeString, "")
		mapStr := regexp.MustCompile(`\]`).Split(mapParts, -1)
		var types []string
		for _, val := range mapStr {
			parts := strings.Split(val, ".")
			if len(parts) > 1 {
				parts = ConstructTypeParts(aliasNameMap, parts)
				val = parts[0] + "." + parts[1]
			}
			types = append(types, val)
		}
		typeString = fmt.Sprintf("map[%s]%s", types[0], types[1])
	case parser.IsArrayField(field):
		arr := regexp.MustCompile(`^(\[])`).ReplaceAllString(typeString, "")
		parts := strings.Split(arr, ".")
		if len(parts) > 1 {
			parts = ConstructTypeParts(aliasNameMap, parts)
			typeString = fmt.Sprintf("[]%s.%s", parts[0], parts[1])
		}
	default:
		parts := strings.Split(typeString, ".")
		if len(parts) > 1 {
			parts = ConstructTypeParts(aliasNameMap, parts)
			typeString = fmt.Sprintf("%s.%s", parts[0], parts[1])
		}
	}
	return typeString
}

func ConstructTypeParts(aliasNameMap map[string]string, parts []string) []string {
	if strings.Contains(parts[0], "*") {
		if val, ok := aliasNameMap[strings.TrimLeft(parts[0], "*")]; ok {
			parts[0] = "*" + val
		}
	} else {
		if val, ok := aliasNameMap[parts[0]]; ok {
			parts[0] = val
		}
	}
	return parts
}

func GenerateGraphqlResolverVars(baseGroupName, crdModulePath string, pkgs parser.Packages, parentsMap map[string]parser.NodeHelper) ([]Node_prop, error) {
	var Nodes []Node_prop
	aliasNameMap := make(map[string]string)
	sortedKeys := make([]string, 0, len(pkgs))
	nonStructMap := make(map[string]string)
	for k := range pkgs {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	sortedPackages := make([]parser.Package, len(pkgs))
	for i, k := range sortedKeys {
		sortedPackages[i] = pkgs[k]
	}
	for i, pkg := range sortedPackages {
		if i > 0 {
			break
		}
		// Get all non struct type for enum
		for _, node := range pkg.GetNonStructTypes() {
			var pkgName string
			if pkg.FullName == pkg.ModPath {
				pkgName = pkg.Name + "_" + parser.GetTypeName(node)
			} else if pkg.Name != "" {
				pkgName = pkg.Name + "_" + parser.GetTypeName(node)
			} else {
				pkgPath := pkg.FullName
				libName := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
				specTypePrefix := strings.ReplaceAll(libName, "-", "")
				pkgName = specTypePrefix + "_" + parser.GetTypeName(node)
			}
			nonStructMap[pkgName] = types.ExprString(node.Type)
		}
		simpleGroupTypeName := util.GetSimpleGroupTypeName(pkg.Name)
		for _, node := range pkg.GetStructs() {
			var nodeProp Node_prop
			// GRAPHQL SCHEMA NAME
			if pkg.FullName == pkg.ModPath {
				nodeProp.SchemaName = pkg.Name + "_" + parser.GetTypeName(node)
			} else if pkg.Name != "" {
				nodeProp.SchemaName = pkg.Name + "_" + parser.GetTypeName(node)
			} else {
				pkgPath := pkg.FullName
				libName := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
				specTypePrefix := strings.ReplaceAll(libName, "-", "")
				nodeProp.SchemaName = specTypePrefix + "_" + parser.GetTypeName(node)
			}
			importMap := pkg.GetImportMap()
			nodeProp.PkgName = simpleGroupTypeName
			// fmt.Println("  NODE:-->", node.Name.String(), nodeProp.SchemaName)
			nodeProp.NodeName = node.Name.String()
			for _, f := range parser.GetSpecFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, f)
				if f != nil {
					// GET FIELD NAME
					fieldProp.FieldName, err = parser.GetNodeFieldName(f)
					if err != nil {
						log.Fatalf("failed to determine field name: %v", err)
					}
					// Check Nexus Field
					if parser.IsNexusTypeField(f) {
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", "Id", "ID")
						// MAP
					} else if parser.IsMapField(f) {
						fmt.Println("    MAP FIELD==>", fieldProp.FieldName, parser.GetFieldType(f))
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
						fieldProp.IsResolver = true
						// ARRAY
					} else if parser.IsArrayField(f) {
						var stdType string
						arr := regexp.MustCompile(`^(\[])`).ReplaceAllString(typeString, "")
						parts := strings.Split(arr, ".")
						if len(parts) > 1 {
							parts = ConstructTypeParts(aliasNameMap, parts)
							stdType = convertGraphqlStdType(parts[1])
						} else {
							stdType = convertGraphqlStdType(parts[0])
						}
						if stdType != "" {
							fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, stdType)
						} else {
							fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, validateImportPkg(pkg, arr, importMap))
						}
						fieldProp.IsResolver = true
						// CUSTOM FIELDS
					} else {
						stdType := convertGraphqlStdType(typeString)
						if stdType != "" {
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, stdType)
						} else {
							// check type is present in nonStructMap
							if val, ok := nonStructMap[validateImportPkg(pkg, typeString, importMap)]; ok {
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, convertGraphqlStdType(val))
							} else {
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, validateImportPkg(pkg, typeString, importMap))
								fieldProp.IsResolver = true
							}
						}
					}
				} else {
					fmt.Println("#### NO FIELDS")
					continue
				}
				if fieldProp.IsResolver {
					nodeProp.ResolverCount += 1
				}
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
			}
			for _, cf := range parser.GetChildFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, cf)
				if cf != nil {
					// GET FIELD NAME
					fieldProp.FieldName, err = parser.GetNodeFieldName(cf)
					fieldProp.IsResolver = true
					if err != nil {
						log.Fatalf("failed to determine field name: %v", err)
					}
					if parser.IsOnlyChildField(cf) {
						fmt.Println("    Child FIELDS:-->", fieldProp.FieldName, "type:", typeString)
						pkgName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, pkgName)
					} else {
						fmt.Println("    Children FIELDS:-->", fieldProp.FieldName, "type:", typeString)
						pkgName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, pkgName)
					}

				}
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
				nodeProp.ResolverCount += 1
			}
			for _, lf := range parser.GetLinkFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, lf)
				if lf != nil {
					// GET FIELD NAME
					fieldProp.FieldName, err = parser.GetNodeFieldName(lf)
					fieldProp.IsResolver = true
					if err != nil {
						log.Fatalf("failed to determine field name: %v", err)
					}
					if parser.IsOnlyLinkField(lf) {
						fmt.Println("    Link FIELDS:-->", fieldProp.FieldName, "type:", typeString)
						pkgName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, pkgName)
					} else {
						fmt.Println("    Links FIELDS:-->", fieldProp.FieldName, "type:", typeString)
						pkgName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, pkgName)
					}

				}
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
				nodeProp.ResolverCount += 1

			}
			Nodes = append(Nodes, nodeProp)
		}
	}
	fmt.Println("***************** GRAPHQL SCHEMA ******************")
	for _, n := range Nodes {
		fmt.Println("type ", n.SchemaName, "{")
		for _, f := range n.GraphqlSchemaFields {
			fmt.Println("\t", f.SchemaFieldName, "    -->", f.IsResolver)
		}
		fmt.Println("}")
	}
	fmt.Println("***************** GRAPHQL CONFIG ******************")
	for _, n := range Nodes {
		if n.ResolverCount > 0 {
			fmt.Println(n.SchemaName, ":")
			fmt.Println("  fields:", n.ResolverCount)
			for _, f := range n.GraphqlSchemaFields {
				if f.IsResolver {
					fmt.Println("    ", f.FieldName, ":\n      resolver: true")
				}
			}
		}
	}
	return Nodes, nil
}
