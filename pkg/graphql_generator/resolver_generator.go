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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ReturnStatement struct {
	Alias      string
	ReturnType string
}

type Field_prop struct {
	IsResolver       bool
	IsNexusTypeField bool
	PkgName          string
	NodeName         string
	FieldName        string
	FieldType        string
	FieldTypePkgPath string
	SchemaFieldName  string
	SchemaTypeName   string
	Alias            string
	ReturnType       string
}

type Node_prop struct {
	IsSingletonNode     bool
	IsNexusNode         bool
	BaseImportPath      string
	ChildLinkFields     []Field_prop
	ChildrenLinksFields []Field_prop
	MapFields           []Field_prop
	CustomFields        []Field_prop
	NonStructFields     []Field_prop
	GraphqlSchemaFields []Field_prop
	ResolverFields      map[string][]Field_prop
	ResolverCount       int
	PkgName             string
	NodeName            string
	SchemaName          string
	Alias               string
	ReturnType          string
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

func validateImportPkg(pkg parser.Package, typeString string, importMap map[string]string) (string, string) {
	if strings.Contains(typeString, ".") {
		part := strings.Split(typeString, ".")
		if val, ok := importMap[part[0]]; ok {
			pkgName := val[strings.LastIndex(val, "/")+1 : len(val)-1]
			repName := strings.ReplaceAll(pkgName, "-", "")
			return repName + "_" + part[1], cases.Title(language.Und, cases.NoLower).String(repName) + cases.Title(language.Und, cases.NoLower).String(part[1])
		}
		return pkg.Name + part[1], cases.Title(language.Und, cases.NoLower).String(pkg.Name) + cases.Title(language.Und, cases.NoLower).String(part[1])
	}
	return pkg.Name + "_" + typeString, cases.Title(language.Und, cases.NoLower).String(pkg.Name) + cases.Title(language.Und, cases.NoLower).String(typeString)
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

func GetResolverReturnVal(sortedPackages []parser.Package, nonStructMap map[string]string) map[string]ReturnStatement {
	var Nodes []Node_prop
	aliasNameMap := make(map[string]string)
	for _, pkg := range sortedPackages {
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
			resField := make(map[string][]Field_prop)
			// GRAPHQL SCHEMA NAME
			if parser.IsSingletonNode(node) {
				nodeProp.IsSingletonNode = true
			} else if parser.IsNexusNode(node) {
				nodeProp.IsNexusNode = true
			}
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
			nodeProp.NodeName = node.Name.String()
			for _, nf := range parser.GetNexusFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, nf)
				if nf != nil {
					if len(nf.Names) > 0 {
						fieldProp.FieldName, err = parser.GetNodeFieldName(nf)
						if err != nil {
							log.Fatalf("failed to determine field name: %v", err)
						}
					}
					fieldProp.PkgName = simpleGroupTypeName
					fieldProp.NodeName = node.Name.String()
					// Child and Link fields
					if parser.IsChildOrLink(nf) {
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsResolver = true
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						nodeProp.ChildLinkFields = append(nodeProp.ChildLinkFields, fieldProp)
					}
					// Children and Links fields
					if parser.IsNamedChildOrLink(nf) {
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsResolver = true
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						nodeProp.ChildrenLinksFields = append(nodeProp.ChildrenLinksFields, fieldProp)
					}
				}
				if fieldProp.IsResolver {
					nodeProp.ResolverCount += 1
				}
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
			}
			for _, f := range parser.GetSpecFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, f)
				if f != nil {
					// GET FIELD NAME
					fieldProp.FieldName, err = parser.GetNodeFieldName(f)
					fieldProp.FieldType = typeString
					if err != nil {
						log.Fatalf("failed to determine field name: %v", err)
					}
					// Check Nexus Field
					if parser.IsNexusTypeField(f) {
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", "Id", "ID")
						fieldProp.IsNexusTypeField = true
						// MAP
					} else if parser.IsMapField(f) {
						// fmt.Println("    MAP FIELD==>", fieldProp.FieldName, parser.GetFieldType(f))
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
							schemaTypeName, resolverTypeName := validateImportPkg(pkg, arr, importMap)
							fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, schemaTypeName)
							fieldProp.SchemaTypeName = schemaTypeName
							fieldProp.FieldType = typeString
							fieldProp.FieldTypePkgPath = resolverTypeName
						}
						fieldProp.IsResolver = true
						// CUSTOM FIELDS
					} else {
						stdType := convertGraphqlStdType(typeString)
						if stdType != "" {
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, stdType)
							// fmt.Println("Res Field-1", fieldProp.FieldName)
							resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)

						} else {
							// check type is present in nonStructMap
							schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
							if val, ok := nonStructMap[schemaTypeName]; ok {
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, convertGraphqlStdType(val))
								fieldProp.FieldType = typeString
								fieldProp.FieldTypePkgPath = resolverTypeName
								fieldProp.SchemaTypeName = schemaTypeName
								resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
							} else {
								schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, schemaTypeName)
								fmt.Println("3)schemaTypeName", schemaTypeName)
								fieldProp.IsResolver = true
								fieldProp.FieldType = typeString
								fieldProp.FieldTypePkgPath = resolverTypeName
								fieldProp.SchemaTypeName = schemaTypeName
								nodeProp.CustomFields = append(nodeProp.CustomFields, fieldProp)
							}

						}
					}
				} else {
					// fmt.Println("#### NO FIELDS")
					continue
				}
				if fieldProp.IsResolver {
					nodeProp.ResolverCount += 1
				}
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
				nodeProp.ResolverFields = resField
			}
			Nodes = append(Nodes, nodeProp)
		}
	}
	// calculate Returnvalue for resolver
	retMap := make(map[string]ReturnStatement)
	customApi := make(map[string]string)
	for _, n := range Nodes {
		var retType string
		var aliasVal string
		retType += fmt.Sprintf("ret := &model.%s%s {\n", n.PkgName, n.NodeName)
		if n.IsNexusNode || n.IsSingletonNode {
			retType += fmt.Sprintf("\t\t%s: &%s,\n", "Id", "Id")
		}
		for _, i := range n.ResolverFields[n.PkgName+n.NodeName] {
			customApi[i.FieldType] = i.FieldName
			if _, ok := nonStructMap[i.SchemaTypeName]; ok {
				retType += fmt.Sprintf("\t\t%s: &v%s,\n", i.FieldName, i.FieldType)
				aliasVal += fmt.Sprintf("v%s := %s(v%s.Spec.%s)\n\t", i.FieldType, nonStructMap[i.SchemaTypeName], n.NodeName, i.FieldType)
			} else {
				if i.IsNexusTypeField {
					retType += fmt.Sprintf("\t\t%s: &v%s.Spec.%s,\n", i.FieldName, n.NodeName, i.FieldName)
				} else {
					retType += fmt.Sprintf("\t\t%s: &v%s.Spec.%s.%s,\n", i.FieldName, n.PkgName, n.NodeName, i.FieldName)
				}
			}
		}
		retType += "\t}"

		retMap[n.PkgName+n.NodeName] = ReturnStatement{
			Alias:      aliasVal,
			ReturnType: retType,
		}
	}
	return retMap
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
	// **********
	retMap := GetResolverReturnVal(sortedPackages, nonStructMap)
	for _, pkg := range sortedPackages {
		// if i > 0 {
		// 	break
		// }
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
			resField := make(map[string][]Field_prop)
			// GRAPHQL SCHEMA NAME
			nodeProp.PkgName = simpleGroupTypeName
			nodeProp.NodeName = node.Name.String()
			nodeProp.BaseImportPath = crdModulePath
			if parser.IsSingletonNode(node) {
				nodeProp.IsSingletonNode = true
				nodeProp.Alias = retMap[nodeProp.PkgName+nodeProp.NodeName].Alias
				nodeProp.ReturnType = retMap[nodeProp.PkgName+nodeProp.NodeName].ReturnType
				fmt.Println("IsSingletonNode", node.Name.Name, nodeProp.PkgName+nodeProp.NodeName)
			}
			if parser.IsNexusNode(node) {
				nodeProp.IsNexusNode = true
			}
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

			for _, nf := range parser.GetNexusFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, nf)
				if nf != nil {
					if len(nf.Names) > 0 {
						fieldProp.FieldName, err = parser.GetNodeFieldName(nf)
						if err != nil {
							log.Fatalf("failed to determine field name: %v", err)
						}
						// fmt.Println("Nexus-FIELDS", nf.Names[0].Name)
					}
					fieldProp.PkgName = simpleGroupTypeName
					fieldProp.NodeName = node.Name.String()

					// Nexus Type Fields
					if parser.IsNexusTypeField(nf) {
						// fmt.Println("ID")
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", "Id", "ID")
					}
					// Child and Link fields
					if parser.IsChildOrLink(nf) {
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsResolver = true
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						// fmt.Println("CHILD", fieldProp.FieldName, fieldProp.FieldType)
						// update resolver return
						fieldProp.Alias = retMap[resolverTypeName].Alias
						fieldProp.ReturnType = retMap[resolverTypeName].ReturnType
						nodeProp.ChildLinkFields = append(nodeProp.ChildLinkFields, fieldProp)
					}
					// Children and Links fields
					if parser.IsNamedChildOrLink(nf) {
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsResolver = true
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						// update resolver return
						fieldProp.Alias = retMap[resolverTypeName].Alias
						fieldProp.ReturnType = retMap[resolverTypeName].ReturnType
						nodeProp.ChildrenLinksFields = append(nodeProp.ChildrenLinksFields, fieldProp)
					}
				}
				if fieldProp.IsResolver {
					nodeProp.ResolverCount += 1
				}
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
			}
			for _, f := range parser.GetSpecFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, f)
				if f != nil {
					// GET FIELD NAME
					fieldProp.FieldName, err = parser.GetNodeFieldName(f)
					fieldProp.FieldType = typeString
					fieldProp.PkgName = simpleGroupTypeName
					fieldProp.NodeName = node.Name.String()
					// update resolver return
					fieldProp.Alias = retMap[nodeProp.PkgName+nodeProp.NodeName].Alias
					fieldProp.ReturnType = retMap[nodeProp.PkgName+nodeProp.NodeName].ReturnType
					// fmt.Println("NAME:", fieldProp.FieldName, fieldProp.Alias, fieldProp.ReturnType)
					if err != nil {
						log.Fatalf("failed to determine field name: %v", err)
					}
					// Check Nexus Field
					if parser.IsNexusTypeField(f) {
						fieldProp.IsNexusTypeField = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", "Id", "ID")
						// MAP
					} else if parser.IsMapField(f) {
						// fmt.Println("    MAP FIELD==>", fieldProp.FieldName, parser.GetFieldType(f))
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
							schemaTypeName, _ := validateImportPkg(pkg, arr, importMap)
							fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, schemaTypeName)
							fieldProp.SchemaTypeName = schemaTypeName
						}
						fieldProp.IsResolver = true
						// CUSTOM FIELDS
					} else {
						stdType := convertGraphqlStdType(typeString)
						if stdType != "" {
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, stdType)
							// fmt.Println("Field-1", fieldProp.FieldName)
							resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)

						} else {
							// check type is present in nonStructMap
							schemaTypeName, _ := validateImportPkg(pkg, typeString, importMap)
							if val, ok := nonStructMap[schemaTypeName]; ok {
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, convertGraphqlStdType(val))
								fieldProp.SchemaTypeName = schemaTypeName
								resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
								fmt.Println("Field-1", fieldProp.FieldName, nodeProp.PkgName+nodeProp.NodeName)
							} else {
								// CustomType Resolver
								schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, schemaTypeName)
								fmt.Println("CUSTOM Fields", fieldProp.FieldName)
								fieldProp.IsResolver = true
								fieldProp.FieldType = typeString
								fieldProp.FieldTypePkgPath = resolverTypeName
								fieldProp.SchemaTypeName = schemaTypeName
								// update resolver return
								fieldProp.Alias = retMap[resolverTypeName].Alias
								fieldProp.ReturnType = retMap[resolverTypeName].ReturnType
								nodeProp.CustomFields = append(nodeProp.CustomFields, fieldProp)
							}
						}
					}
				} else {
					// fmt.Println("#### NO FIELDS")
					continue
				}
				if fieldProp.IsResolver {
					nodeProp.ResolverCount += 1
				}
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
				nodeProp.ResolverFields = resField
			}

			// for _, cf := range parser.GetChildFields(node) {
			// 	var fieldProp Field_prop
			// 	var err error
			// 	typeString := ConstructType(aliasNameMap, cf)
			// 	if cf != nil {
			// 		// GET FIELD NAME
			// 		fieldProp.FieldName, err = parser.GetNodeFieldName(cf)
			// 		fieldProp.IsResolver = true
			// 		if err != nil {
			// 			log.Fatalf("failed to determine field name: %v", err)
			// 		}
			// 		if parser.IsOnlyChildField(cf) {
			// 			// fmt.Println("    Child FIELDS:-->", fieldProp.FieldName, "type:", typeString)
			// 			pkgName := validateImportPkg(pkg, typeString, importMap)
			// 			fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, pkgName)
			// 		} else {
			// 			// fmt.Println("    Children FIELDS:-->", fieldProp.FieldName, "type:", typeString)
			// 			pkgName := validateImportPkg(pkg, typeString, importMap)
			// 			fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, pkgName)
			// 		}

			// 	}
			// 	nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
			// 	nodeProp.ResolverCount += 1
			// }
			// for _, lf := range parser.GetLinkFields(node) {
			// 	var fieldProp Field_prop
			// 	var err error
			// 	typeString := ConstructType(aliasNameMap, lf)
			// 	if lf != nil {
			// 		// GET FIELD NAME
			// 		fieldProp.FieldName, err = parser.GetNodeFieldName(lf)
			// 		fieldProp.IsResolver = true
			// 		if err != nil {
			// 			log.Fatalf("failed to determine field name: %v", err)
			// 		}
			// 		if parser.IsOnlyLinkField(lf) {
			// 			// fmt.Println("    Link FIELDS:-->", fieldProp.FieldName, "type:", typeString)
			// 			pkgName := validateImportPkg(pkg, typeString, importMap)
			// 			fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, pkgName)
			// 		} else {
			// 			// fmt.Println("    Links FIELDS:-->", fieldProp.FieldName, "type:", typeString)
			// 			pkgName := validateImportPkg(pkg, typeString, importMap)
			// 			fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, pkgName)
			// 		}

			// 	}
			// 	nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
			// 	nodeProp.ResolverCount += 1
			// }
			Nodes = append(Nodes, nodeProp)
		}
	}
	// fmt.Println("***************** GRAPHQL SCHEMA ******************")
	// for _, n := range Nodes {
	// 	fmt.Println("type ", n.SchemaName, "{")
	// 	for _, f := range n.GraphqlSchemaFields {
	// 		fmt.Println("\t", f.SchemaFieldName, "    -->", f.IsResolver)
	// 	}
	// 	fmt.Println("}")
	// }
	// fmt.Println("***************** GRAPHQL CONFIG ******************")
	// for _, n := range Nodes {
	// 	if n.ResolverCount > 0 {
	// 		fmt.Println(n.SchemaName, ":")
	// 		fmt.Println("  fields:", n.ResolverCount)
	// 		for _, f := range n.GraphqlSchemaFields {
	// 			if f.IsResolver {
	// 				fmt.Println("    ", f.FieldName, ":\n      resolver: true")
	// 			}
	// 		}
	// 	}
	// }
	// for _, n := range Nodes {
	// 	for _, j := range n.ReturnStatementMap {
	// 		if _, ok := j["RootRoot"]; ok {
	// 			fmt.Println("11111", j["RootRoot"].Alias)
	// 			fmt.Println("22222", j["RootRoot"].ReturnType)
	// 		}
	// 	}
	// }

	return Nodes, nil
}

// ret := &model.RootRoot{
// 		DisplayName: &vRoot.Spec.DisplayName,
// }
