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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ReturnStatement struct {
	Alias      string
	ReturnType string
}

type Field_prop struct {
	IsResolver              bool
	IsNexusTypeField        bool
	IsNexusOrSingletonField bool
	IsChildOrLink           bool
	IsChildrenOrLinks       bool
	IsMapTypeField          bool
	IsArrayTypeField        bool
	IsStdTypeField          bool
	IsCustomTypeField       bool
	IsFieldIgnore           bool
	IsStringType            bool
	IsAliasTypeField        bool
	PkgName                 string
	NodeName                string
	FieldName               string
	FieldType               string
	FieldTypePkgPath        string
	ModelType               string
	SchemaFieldName         string
	SchemaTypeName          string
	BaseTypeName            string
	Alias                   string
	ReturnType              string
}

type Node_prop struct {
	IsParentNode        bool
	HasParent           bool
	IsSingletonNode     bool
	IsNexusNode         bool
	BaseImportPath      string
	CrdName             string
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

func getArraySchema(FieldName, schemaTypeName string, nonStructMap map[string]string) string {
	if val, ok := nonStructMap[schemaTypeName]; ok {
		if convertGraphqlStdType(val) != "" {
			return fmt.Sprintf("%s: [%s]!", FieldName, convertGraphqlStdType(val))
		} else {
			if strings.HasPrefix(val, "[]") {
				arrStd := regexp.MustCompile(`^(\[])`).ReplaceAllString(val, "")
				return fmt.Sprintf("%s: [%s]!", FieldName, convertGraphqlStdType(arrStd))
			}
			return fmt.Sprintf("%s: [%s]!", FieldName, val)
		}
	} else {
		return fmt.Sprintf("%s: [%s]!", FieldName, schemaTypeName)
	}
}

func jsonMarshalResolver(FieldName, PkgName string) string {
	return fmt.Sprintf("%s, _ := json.Marshal(v%s.Spec.%s)\n%sData := string(%s)\n", FieldName, PkgName, FieldName, FieldName, FieldName)
}
func jsonMarshalCustomResolver(FieldName, PkgName, FieldType string) string {
	return fmt.Sprintf("%s, _ := json.Marshal(v%s.Spec.%s.%s)\n%sData := string(%s)\n", FieldName, PkgName, FieldType, FieldName, FieldName, FieldName)
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

func getBaseNodeType(pkg parser.Package, typeString string, importMap map[string]string) string {
	if strings.Contains(typeString, ".") {
		part := strings.Split(typeString, ".")
		return part[1]
	}
	return typeString
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
	// Iterate All Non Struct Node from all pkg and store the details in nonStructMap
	for _, pkg := range sortedPackages {
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
			// NonStruct Map
			nonStructType := types.ExprString(node.Type)
			nonStructMap[pkgName] = nonStructType
		}
	}
	for _, pkg := range sortedPackages {
		simpleGroupTypeName := util.GetSimpleGroupTypeName(pkg.Name)
		for _, node := range pkg.GetStructs() {
			var nodeProp Node_prop
			resField := make(map[string][]Field_prop)
			// GRAPHQL SCHEMA NAME
			nodeProp.PkgName = simpleGroupTypeName
			nodeProp.NodeName = node.Name.String()
			nodeProp.BaseImportPath = crdModulePath
			nodeProp.CrdName = util.GetCrdName(node.Name.String(), pkg.Name, baseGroupName)
			nodeHelper := parentsMap[nodeProp.CrdName]
			nodeProp.IsParentNode = parser.IsSingletonNode(node)
			if len(nodeHelper.Parents) > 0 {
				nodeProp.HasParent = true
			}
			if parser.IsSingletonNode(node) {
				nodeProp.IsSingletonNode = true
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
					}
					fieldProp.PkgName = simpleGroupTypeName
					fieldProp.NodeName = node.Name.String()
					// IGNORE FIELD using Annotation `nexus-graphql:"ignore:true"`
					if parser.IgnoreField(nf) {
						fieldProp.IsFieldIgnore = true
						continue
					}
					// Convert to String type using annotation `nexus-graphql:"type:string"`
					if parser.IsJsonStringField(nf) {
						fieldProp.IsStringType = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
						continue
					}
					// Nexus Type Fields
					if parser.IsNexusTypeField(nf) {
						fieldProp.IsNexusOrSingletonField = true
						//fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", "Id", "ID")
						// Custom Query + ID
						fieldProp.SchemaFieldName = CustomQuerySchema
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
						fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
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
						fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
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
						resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
					} else if parser.IsMapField(f) {
						fieldProp.IsMapTypeField = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
						resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
						// ARRAY FIELD
					} else if parser.IsArrayField(f) {
						fieldProp.IsArrayTypeField = true
						fieldProp.IsResolver = true
						fmt.Println("TYPE ARRAY:", typeString, strings.ReplaceAll(typeString, "[]", ""))
						schemaTypeName, _ := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, schemaTypeName)

						fieldProp.SchemaTypeName = schemaTypeName
						fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
						fmt.Println("Array Fields", f.Names)
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
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: [%s]!", fieldProp.FieldName, stdType)
							fieldProp.SchemaTypeName = schemaTypeName
							fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
							fieldProp.FieldType = strings.ReplaceAll(typeString, "[]", "")
							fieldProp.ModelType = convertGoStdType(strings.ReplaceAll(typeString, "[]", ""))
							fieldProp.FieldTypePkgPath = ""
							fmt.Println("AAA:", stdType, fieldProp.PkgName, fieldProp.SchemaTypeName, fieldProp.BaseTypeName, fieldProp.FieldTypePkgPath)
						} else {
							schemaTypeName, resolverTypeName := validateImportPkg(pkg, arr, importMap)
							fieldProp.SchemaFieldName = getArraySchema(fieldProp.FieldName, schemaTypeName, nonStructMap)
							fieldProp.ModelType = "model." + fieldProp.PkgName + strings.ReplaceAll(typeString, "[]", "")
							fieldProp.FieldTypePkgPath = resolverTypeName
							if val, ok := nonStructMap[schemaTypeName]; ok {
								fieldProp.IsAliasTypeField = true
								fieldProp.ModelType = convertGoStdType(val)
								fieldProp.FieldTypePkgPath = ""
							}
							fieldProp.SchemaTypeName = schemaTypeName
							fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
							fieldProp.FieldType = fieldProp.PkgName + strings.ReplaceAll(typeString, "[]", "")
							fmt.Println("BBB:", stdType, fieldProp.PkgName, fieldProp.SchemaTypeName, fieldProp.BaseTypeName, fieldProp.FieldTypePkgPath)
						}
						nodeProp.ArrayFields = append(nodeProp.ArrayFields, fieldProp)
						resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
						// CUSTOM FIELDS
					} else {
						stdType := convertGraphqlStdType(typeString)
						if stdType != "" {
							fieldProp.IsStdTypeField = true
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, stdType)
							resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)

						} else {
							// check type is present in nonStructMap
							schemaTypeName, _ := validateImportPkg(pkg, typeString, importMap)
							if val, ok := nonStructMap[schemaTypeName]; ok {
								fieldProp.IsAliasTypeField = true
								if convertGoStdType(val) != "" {
									fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, convertGraphqlStdType(val))
								} else if strings.HasPrefix(val, "map") {
									fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
								} else if strings.HasPrefix(val, "[]") {
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
										fieldProp.SchemaFieldName = fmt.Sprintf("%s: [%s!]", fieldProp.FieldName, stdType)
									} else {
										schemaTypeName, _ := validateImportPkg(pkg, arr, importMap)
										fieldProp.SchemaFieldName = getArraySchema(fieldProp.FieldName, schemaTypeName, nonStructMap)
										fieldProp.SchemaTypeName = schemaTypeName
										fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
									}
								} else {
									fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, val)
								}
								fieldProp.SchemaTypeName = schemaTypeName
								fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
								resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
							} else {
								// CustomType Resolver

								fieldProp.IsCustomTypeField = true
								schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
								fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, schemaTypeName)
								fieldProp.IsResolver = true
								fieldProp.FieldType = typeString
								fieldProp.FieldTypePkgPath = resolverTypeName
								fieldProp.SchemaTypeName = schemaTypeName
								fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
								nodeProp.CustomFields = append(nodeProp.CustomFields, fieldProp)
							}
						}
					}
				} else {
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
	// Calculate customType
	customApi := make(map[string]string)
	for _, n := range Nodes {
		for _, f := range n.CustomFields {
			customApi[f.BaseTypeName] = f.FieldName
		}
	}
	// calculate Returnvalue for resolver
	retMap := make(map[string]ReturnStatement)
	ListRetMap := make(map[string]ReturnStatement)
	for _, n := range Nodes {
		var retType string
		var aliasVal string
		var listRetVal string

		// if n.IsNexusNode || n.IsSingletonNode {
		// 	retType += fmt.Sprintf("\t%s: &%s,\n", "Id", "Id")
		// }
		// Check Return field len
		if len(n.ResolverFields[n.PkgName+n.NodeName]) > 0 {
			retType += fmt.Sprintf("ret := &model.%s%s {\n", n.PkgName, n.NodeName)
		}
		for _, i := range n.ResolverFields[n.PkgName+n.NodeName] {
			if i.IsAliasTypeField {
				if val, ok := nonStructMap[i.SchemaTypeName]; ok {
					if strings.HasPrefix(val, "map") {
						retType += fmt.Sprintf("\t%s: &%sData,\n", i.FieldName, i.FieldName)
						if !n.IsNexusNode {
							aliasVal += jsonMarshalCustomResolver(i.FieldName, n.PkgName, i.NodeName)
						} else {
							aliasVal += jsonMarshalResolver(i.FieldName, n.NodeName)
						}
					} else if strings.HasPrefix(val, "[]") {
						fmt.Println("ARRAY ALIAS", i.FieldName, i.FieldType)
					} else if strings.HasPrefix(val, "*") {
						fmt.Println("Alias-Pointer:", i.NodeName, i.FieldName, i.FieldType)
					} else {
						fmt.Println("1111######", i.FieldName, i.FieldType, val, i.IsArrayTypeField)
						if convertGoStdType(val) != "" && !i.IsArrayTypeField {
							retType += fmt.Sprintf("\t%s: &v%s,\n", i.FieldName, i.FieldName)
							listRetVal += fmt.Sprintf("v%s := %s(&i.%s)\n", i.FieldName, convertGoStdType(val), i.FieldName)
							if !n.IsNexusNode {
								fmt.Println("RETURN STATEMENT", i.FieldName, i.FieldType, val)
								aliasVal += fmt.Sprintf("v%s := %s(v%s.Spec.%s.%s)\n", i.FieldName, convertGoStdType(val), i.PkgName, customApi[i.NodeName], i.FieldName)
							} else {
								aliasVal += fmt.Sprintf("v%s := %s(v%s.Spec.%s)\n", i.FieldName, convertGoStdType(val), i.NodeName, i.FieldName)
							}
						} else {
							fmt.Println("Not found")
						}
					}
				}
			} else if i.IsArrayTypeField {
				fmt.Println("*******")
			} else if i.IsMapTypeField {
				retType += fmt.Sprintf("\t%s: &%sData,\n", i.FieldName, i.FieldName)
				if !n.IsNexusNode {
					aliasVal += jsonMarshalCustomResolver(i.FieldName, n.PkgName, i.NodeName)
				} else {
					aliasVal += jsonMarshalResolver(i.FieldName, n.NodeName)
				}
			} else if i.IsStringType {
				retType += fmt.Sprintf("\t%s: &%sData,\n", i.FieldName, i.FieldName)
				if !n.IsNexusNode {
					aliasVal += jsonMarshalCustomResolver(i.FieldName, n.PkgName, i.NodeName)
				} else {
					aliasVal += jsonMarshalResolver(i.FieldName, n.NodeName)
				}
			} else if i.IsStdTypeField {
				fmt.Println("2222######", i.FieldName)
				if convertGoStdType(i.FieldType) != "" {
					fmt.Println("RETURN STATEMENT-2", n.NodeName, i.FieldName, i.FieldType)
					retType += fmt.Sprintf("\t%s: &v%s,\n", i.FieldName, i.FieldName)
					listRetVal += fmt.Sprintf("v%s := %s(&i.%s)\n", i.FieldName, convertGoStdType(i.FieldType), i.FieldName)
					if !n.IsNexusNode {
						aliasVal += fmt.Sprintf("v%s := %s(v%s.Spec.%s.%s)\n", i.FieldName, convertGoStdType(i.FieldType), i.PkgName, customApi[i.NodeName], i.FieldName)
					} else {
						aliasVal += fmt.Sprintf("v%s := %s(v%s.Spec.%s)\n", i.FieldName, convertGoStdType(i.FieldType), i.NodeName, i.FieldName)
					}
				}
			}
		}
		if len(n.ResolverFields[n.PkgName+n.NodeName]) > 0 {
			retType += "\t}"
		}

		retMap[n.PkgName+n.NodeName] = ReturnStatement{
			Alias:      aliasVal,
			ReturnType: retType,
		}
		ListRetMap[n.PkgName+n.NodeName] = ReturnStatement{
			Alias:      listRetVal,
			ReturnType: retType,
		}
		fmt.Println("re", n.PkgName+n.NodeName)
	}
	var ResNodes []Node_prop
	for _, n := range Nodes {
		var resNodeProp Node_prop
		resNodeProp.Alias = n.Alias
		resNodeProp.ReturnType = n.ReturnType
		resNodeProp.BaseImportPath = n.BaseImportPath
		resNodeProp.GraphqlSchemaFields = n.GraphqlSchemaFields
		resNodeProp.IsSingletonNode = n.IsSingletonNode
		resNodeProp.IsNexusNode = n.IsNexusNode
		resNodeProp.ResolverFields = n.ResolverFields
		resNodeProp.ResolverCount = n.ResolverCount
		resNodeProp.PkgName = n.PkgName
		resNodeProp.NodeName = n.NodeName
		resNodeProp.SchemaName = n.SchemaName
		if !n.HasParent && n.IsParentNode {
			resNodeProp.Alias = retMap[resNodeProp.PkgName+resNodeProp.NodeName].Alias
			resNodeProp.ReturnType = retMap[resNodeProp.PkgName+resNodeProp.NodeName].ReturnType
			resNodeProp.IsParentNode = true
		}
		for _, f := range n.ChildLinkFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			resNodeProp.ChildLinkFields = append(resNodeProp.ChildLinkFields, f)
		}
		for _, f := range n.ChildrenLinksFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			resNodeProp.ChildrenLinksFields = append(resNodeProp.ChildrenLinksFields, f)
		}
		for _, f := range n.CustomFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			resNodeProp.CustomFields = append(resNodeProp.CustomFields, f)
		}
		for _, f := range n.NonStructFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			resNodeProp.NonStructFields = append(resNodeProp.NonStructFields, f)
		}
		for _, f := range n.ArrayFields {
			f.ReturnType = ListRetMap[f.FieldTypePkgPath].ReturnType
			f.Alias = ListRetMap[f.FieldTypePkgPath].Alias
			fmt.Println("AAARRR", f.FieldName, f.FieldTypePkgPath, ListRetMap[f.FieldTypePkgPath].ReturnType)
			resNodeProp.ArrayFields = append(resNodeProp.ArrayFields, f)
		}
		ResNodes = append(ResNodes, resNodeProp)
	}
	return ResNodes, nil
}
