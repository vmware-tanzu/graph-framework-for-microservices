package crd_generator

import (
	"fmt"
	"go/types"
	"log"
	"regexp"
	"sort"
	"strings"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"
)

type ReturnStatement struct {
	Alias      string
	ReturnType string
}

type Field_prop struct {
	IsResolver        bool
	IsNexusTypeField  bool
	IsChildOrLink     bool
	IsChildrenOrLinks bool
	IsMapTypeField    bool
	IsArrayTypeField  bool
	IsStdTypeField    bool
	IsCustomTypeField bool
	IsFieldIgnore     bool
	IsStringType      bool
	IsAliasTypeField  bool
	PkgName           string
	NodeName          string
	FieldName         string
	FieldType         string
	FieldTypePkgPath  string
	SchemaFieldName   string
	SchemaTypeName    string
	Alias             string
	ReturnType        string
}

type Node_prop struct {
	IsParent            bool
	IsSingletonNode     bool
	IsNexusNode         bool
	BaseImportPath      string
	ChildLinkFields     []Field_prop
	ChildrenLinksFields []Field_prop
	ArrayFields         []Field_prop
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
			return repName + "_" + part[1], strings.Title(repName) + strings.Title(part[1])
		}
		return pkg.Name + part[1], strings.Title(pkg.Name) + strings.Title(part[1])
	}
	return pkg.Name + "_" + typeString, strings.Title(pkg.Name) + strings.Title(typeString)
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
				// nodeProp.Alias = retMap[nodeProp.PkgName+nodeProp.NodeName].Alias
				// nodeProp.ReturnType = retMap[nodeProp.PkgName+nodeProp.NodeName].ReturnType
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
					// IGNORE FIELD
					if parser.IgnoreField(nf) {
						fieldProp.IsFieldIgnore = true
						continue
					}
					// STRING FIELD
					if parser.IsJsonStringField(nf) {
						fieldProp.IsStringType = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
						continue
					}
					// Nexus Type Fields
					if parser.IsNexusTypeField(nf) {
						fieldProp.IsNexusTypeField = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", "Id", "ID")
					}
					// Child and Link fields
					if parser.IsChildOrLink(nf) {
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.IsChildOrLink = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsResolver = true
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						// fmt.Println("CHILD", fieldProp.FieldName, fieldProp.FieldType)
						// update resolver return
						// fieldProp.Alias = retMap[resolverTypeName].Alias
						// fieldProp.ReturnType = retMap[resolverTypeName].ReturnType
						nodeProp.ChildLinkFields = append(nodeProp.ChildLinkFields, fieldProp)
					}
					// Children and Links fields
					if parser.IsNamedChildOrLink(nf) {
						fieldProp.IsChildrenOrLinks = true
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsResolver = true
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						// update resolver return
						// fieldProp.Alias = retMap[resolverTypeName].Alias
						// fieldProp.ReturnType = retMap[resolverTypeName].ReturnType
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
					// fieldProp.Alias = retMap[nodeProp.PkgName+nodeProp.NodeName].Alias
					// fieldProp.ReturnType = retMap[nodeProp.PkgName+nodeProp.NodeName].ReturnType
					// fmt.Println("NAME:", fieldProp.FieldName, fieldProp.Alias, fieldProp.ReturnType)
					if err != nil {
						log.Fatalf("failed to determine field name: %v", err)
					}
					// IGNORE FIELD
					if parser.IgnoreField(f) {
						fieldProp.IsFieldIgnore = true
						continue
					}
					// STRING FIELD
					if parser.IsJsonStringField(f) {
						fieldProp.IsStringType = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
						continue
					}
					// Check Nexus Field
					if parser.IsNexusTypeField(f) {
						fieldProp.IsNexusTypeField = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", "Id", "ID")
						// MAP
					} else if parser.IsMapField(f) {
						// fmt.Println("    MAP FIELD==>", fieldProp.FieldName, parser.GetFieldType(f))
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
						// fieldProp.IsResolver = true
						// ARRAY
					} else if parser.IsArrayField(f) {
						fieldProp.IsArrayTypeField = true
						fieldProp.IsResolver = true
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
						nodeProp.ArrayFields = append(nodeProp.ArrayFields, fieldProp)
						// CUSTOM FIELDS
					} else {
						stdType := convertGraphqlStdType(typeString)
						if stdType != "" {
							fieldProp.IsStdTypeField = true
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, stdType)
							// fmt.Println("Field-1", fieldProp.FieldName)
							resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)

						} else {
							// check type is present in nonStructMap
							schemaTypeName, _ := validateImportPkg(pkg, typeString, importMap)
							if val, ok := nonStructMap[schemaTypeName]; ok {
								fieldProp.IsAliasTypeField = true
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, convertGraphqlStdType(val))
								fieldProp.SchemaTypeName = schemaTypeName
								resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
								fmt.Println("Field-1", fieldProp.FieldName, nodeProp.PkgName+nodeProp.NodeName)
							} else {
								// CustomType Resolver
								fieldProp.IsCustomTypeField = true
								schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, schemaTypeName)
								fmt.Println("CUSTOM Fields", fieldProp.FieldName)
								fieldProp.IsResolver = true
								fieldProp.FieldType = typeString
								fieldProp.FieldTypePkgPath = resolverTypeName
								fieldProp.SchemaTypeName = schemaTypeName
								// update resolver return
								// fieldProp.Alias = retMap[resolverTypeName].Alias
								// fieldProp.ReturnType = retMap[resolverTypeName].ReturnType
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
	for _, n := range Nodes {
		for _, f := range n.ChildLinkFields {
			fmt.Println(f.FieldName, f.FieldType, f)
		}
	}
	return Nodes, nil
}
