package crd_generator

import (
	"fmt"
	"go/types"
	"log"
	"sort"
	"strings"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ReturnStatement struct {
	Alias       string
	ReturnType  string
	FieldCount  int
	CRDName     string
	ChainAPI    string
	IsSingleton bool
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
	IsPointerTypeField      bool
	IsStringType            bool
	IsAliasTypeField        bool
	PkgName                 string
	NodeName                string
	FieldName               string
	FieldType               string
	FieldTypePkgPath        string
	ModelType               string
	IsArrayStdType          bool
	SchemaFieldName         string
	SchemaTypeName          string
	BaseTypeName            string
	Alias                   string
	ReturnType              string
	FieldCount              int
	CRDName                 string
	ChainAPI                string
	IsSingleton             bool
	LinkAPI                 string
}

type Node_prop struct {
	IsParentNode           bool
	HasParent              bool
	IsSingletonNode        bool
	IsNexusNode            bool
	BaseImportPath         string
	CrdName                string
	ChildFields            []Field_prop
	LinkFields             []Field_prop
	ChildrenFields         []Field_prop
	LinksFields            []Field_prop
	ArrayFields            []Field_prop
	CustomFields           []Field_prop
	NonStructFields        []Field_prop
	GraphqlSchemaFields    []Field_prop
	ResolverFields         map[string][]Field_prop
	ResolverCount          int
	PkgName                string
	NodeName               string
	SchemaName             string
	Alias                  string
	ReturnType             string
	GroupResourceNameTitle string
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
	// Iterate over all non struct type from sortedPackages and store the details in nonStructMap
	// nonStructMap[pkgName] = nodeType  --> ex: nonStructMap["root"] = "AliasTypeFoo"
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
	// Iterate All Nodes in the sortedPackages
	for _, pkg := range sortedPackages {
		simpleGroupTypeName := util.GetSimpleGroupTypeName(pkg.Name)
		for _, node := range pkg.GetStructs() {
			var nodeProp Node_prop
			resField := make(map[string][]Field_prop)
			// Fill Node Properties
			nodeProp.PkgName = simpleGroupTypeName
			nodeProp.NodeName = node.Name.String()
			nodeProp.BaseImportPath = crdModulePath
			nodeProp.CrdName = util.GetCrdName(node.Name.String(), pkg.Name, baseGroupName)
			nodeHelper := parentsMap[nodeProp.CrdName]
			nodeProp.IsParentNode = parser.IsNexusNode(node)
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
						continue
					}
					// Convert to String type using annotation `nexus-graphql:"type:string"`
					if parser.IsJsonStringField(nf) {
						fieldProp.IsStringType = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
					}
					// NexusOrSingletonField Type Fields
					if parser.IsNexusTypeField(nf) {
						fieldProp.IsNexusOrSingletonField = true
						// Add Custom Query + ID
						fieldProp.SchemaFieldName = CustomQuerySchema
					}
					// Nexus Child and Link fields
					if parser.IsOnlyLinkField(nf) {
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsResolver = true
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
						nodeProp.LinkFields = append(nodeProp.LinkFields, fieldProp)
					}
					if parser.IsOnlyChildField(nf) {
						schemaTypeName, resolverTypeName := validateImportPkg(pkg, typeString, importMap)
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, schemaTypeName)
						fieldProp.IsNexusTypeField = true
						fieldProp.FieldType = typeString
						fieldProp.FieldTypePkgPath = resolverTypeName
						fieldProp.SchemaTypeName = schemaTypeName
						fieldProp.BaseTypeName = getBaseNodeType(pkg, typeString, importMap)
						nodeProp.ChildFields = append(nodeProp.ChildFields, fieldProp)
					}
					// Nexus Children and Links fields details
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
						if parser.IsOnlyChildrenField(nf) {
							nodeProp.ChildrenFields = append(nodeProp.ChildrenFields, fieldProp)
						} else {
							nodeProp.LinksFields = append(nodeProp.LinksFields, fieldProp)
						}
					}
				}
				// No of resolver fields in a nodes
				if fieldProp.IsResolver {
					nodeProp.ResolverCount += 1
				}
				if !parser.IsOnlyChildField(nf) {
					nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
				}
			}
			// Iterate Non Nexus Fields
			for _, f := range parser.GetSpecFields(node) {
				var fieldProp Field_prop
				var err error
				typeString := ConstructType(aliasNameMap, f)
				if f != nil {
					// Fill all field information in fieldProp
					fieldProp.FieldName, err = parser.GetNodeFieldName(f)
					fieldProp.FieldType = typeString
					fieldProp.PkgName = simpleGroupTypeName
					fieldProp.NodeName = node.Name.String()
					if err != nil {
						log.Fatalf("failed to determine field name: %v", err)
					}
					// IGNORE FIELD
					if parser.IgnoreField(f) {
						continue
					}
					// STRING FIELD
					if parser.IsJsonStringField(f) {
						fieldProp.IsStringType = true
						fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
						resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
					} else {
						stdType := convertGraphqlStdType(typeString)
						// STANDARD TYPE CHECK
						if len(stdType) != 0 {
							fieldProp.IsStdTypeField = true
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, stdType)
							resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
						} else {
							// JSON MARSHAL
							fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
							fieldProp.IsStringType = true
							resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
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
	// Create CustomType map
	customApi := make(map[string]string)
	CRDNameMap := make(map[string]string)
	for _, n := range Nodes {
		if n.IsNexusNode || n.IsSingletonNode {
			CRDNameMap[n.CrdName] = n.PkgName + n.NodeName
		}
		for _, f := range n.CustomFields {
			customApi[f.BaseTypeName] = f.FieldName
		}
	}
	// Collect ReturnValue of each Node for resolver
	LinkAPI := make(map[string]string)
	retMap := make(map[string]ReturnStatement)
	for _, n := range Nodes {
		var retType string
		var aliasVal string
		var listRetVal string
		var fieldCount int
		var IsSingleton bool
		var ChainAPI string

		retType += fmt.Sprintf("ret := &model.%s%s {\n", n.PkgName, n.NodeName)
		if n.IsNexusNode || n.IsSingletonNode {
			retType += fmt.Sprintf("\t%s: &%s,\n", "Id", "dn")
			aliasVal += fmt.Sprintf("%s := v%s.DisplayName()\n", "dn", n.NodeName)

			retType += fmt.Sprintf("\t%s: %s,\n", "ParentLabels", "parentLabels")
			aliasVal += fmt.Sprintf("%s := map[string]interface{}{%q:%s}\n", "parentLabels", n.CrdName, "dn")

			ChainAPI += "nc"
			var prevNode parser.NodeHelper
			for _, i := range parentsMap[n.CrdName].Parents {
				currentNode := parentsMap[i]

				if len(currentNode.Parents) == 0 {
					// root of the graph
					if currentNode.IsSingleton {
						ChainAPI += fmt.Sprintf(".%s()", CRDNameMap[i])
					} else {
						ChainAPI += fmt.Sprintf(".%s(obj.ParentLabels[\"%s\"].(string))", CRDNameMap[i], i)
					}
				} else {
					if childNode, ok := prevNode.Children[i]; ok {
						if currentNode.IsSingleton {
							ChainAPI += fmt.Sprintf(".%s()", childNode.FieldName)
						} else {
							ChainAPI += fmt.Sprintf(".%s(obj.ParentLabels[\"%s\"].(string))", childNode.FieldName, i)
						}
					} else {
						panic(fmt.Sprintf("unable to find child %s in parent node of %s", i, prevNode.RestName))
					}

				}

				// cache the non-leaf node
				prevNode = currentNode

				//				if parentsMap[i].IsSingleton {
				//					ChainAPI += fmt.Sprintf(".%s()", CRDNameMap[i])
				//				} else {
				//					ChainAPI += fmt.Sprintf(".%s(obj.ParentLabels[\"%s\"].(string))", CRDNameMap[i], i)
				//				}
			}
			// Create LinkAPI
			if n.IsSingletonNode {
				IsSingleton = true
				if !n.HasParent && n.IsParentNode {
					LinkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO())", ChainAPI, n.PkgName+n.NodeName)
				} else {
					LinkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO())", ChainAPI, prevNode.Children[n.CrdName].FieldName)
				}
			} else {
				IsSingleton = false
				if !n.HasParent && n.IsParentNode {
					LinkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO(), obj.ParentLabels[\"%s\"].(string))", ChainAPI, n.PkgName+n.NodeName, n.CrdName)
				} else {
					LinkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO(), obj.ParentLabels[\"%s\"].(string))", ChainAPI, prevNode.Children[n.CrdName].FieldName, n.CrdName)
				}
			}
		}
		for _, i := range n.ResolverFields[n.PkgName+n.NodeName] {
			if i.IsAliasTypeField {
				if val, ok := nonStructMap[i.SchemaTypeName]; ok {
					if strings.HasPrefix(val, "map") {
						fieldCount += 1
						retType += fmt.Sprintf("\t%s: &%sData,\n", i.FieldName, i.FieldName)
						aliasVal += jsonMarshalResolver(i.FieldName, n.NodeName)
					} else {
						if len(convertGoStdType(val)) != 0 && !i.IsArrayTypeField {
							fieldCount += 1
							retType += fmt.Sprintf("\t%s: &v%s,\n", i.FieldName, i.FieldName)
							aliasVal += fmt.Sprintf("v%s := %s(v%s.Spec.%s)\n", i.FieldName, convertGoStdType(val), i.NodeName, i.FieldName)
						}
					}
				}
			} else if i.IsMapTypeField || i.IsStringType {
				fieldCount += 1
				retType += fmt.Sprintf("\t%s: &%sData,\n", i.FieldName, i.FieldName)
				aliasVal += jsonMarshalResolver(i.FieldName, n.NodeName)
			} else if i.IsStdTypeField {
				if len(convertGoStdType(i.FieldType)) != 0 {
					fieldCount += 1
					retType += fmt.Sprintf("\t%s: &v%s,\n", i.FieldName, i.FieldName)
					listRetVal += fmt.Sprintf("v%s := %s(i.%s)\n", i.FieldName, convertGoStdType(i.FieldType), i.FieldName)
					aliasVal += fmt.Sprintf("v%s := %s(v%s.Spec.%s)\n", i.FieldName, convertGoStdType(i.FieldType), i.NodeName, i.FieldName)
				}
			}
		}
		retType += "\t}"
		retMap[n.PkgName+n.NodeName] = ReturnStatement{
			Alias:       aliasVal,
			ReturnType:  retType,
			FieldCount:  fieldCount,
			CRDName:     n.CrdName,
			ChainAPI:    ChainAPI,
			IsSingleton: IsSingleton,
		}
	}
	// set return value to each node
	var ResNodes []Node_prop
	for _, n := range Nodes {
		var resNodeProp Node_prop
		resNodeProp.GroupResourceNameTitle = util.GetGroupResourceNameTitle(n.NodeName)
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
		for _, f := range n.ChildFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			f.FieldCount = retMap[f.FieldTypePkgPath].FieldCount
			f.CRDName = retMap[f.FieldTypePkgPath].CRDName
			f.ChainAPI = retMap[f.FieldTypePkgPath].ChainAPI
			f.IsSingleton = retMap[f.FieldTypePkgPath].IsSingleton
			f.IsResolver = true
			resNodeProp.ResolverCount = 1
			if f.IsSingleton {
				f.SchemaFieldName = fmt.Sprintf("%s: %s!", f.FieldName, f.SchemaTypeName)
			} else {
				f.SchemaFieldName = fmt.Sprintf("%s(Id: ID): %s!", f.FieldName, f.SchemaTypeName)
			}
			f.LinkAPI = LinkAPI[f.PkgName+f.NodeName]
			resNodeProp.ChildFields = append(resNodeProp.ChildFields, f)
			resNodeProp.GraphqlSchemaFields = append(resNodeProp.GraphqlSchemaFields, f)
		}
		for _, f := range n.ChildrenFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			f.FieldCount = retMap[f.FieldTypePkgPath].FieldCount
			f.CRDName = retMap[f.FieldTypePkgPath].CRDName
			f.ChainAPI = retMap[f.FieldTypePkgPath].ChainAPI
			f.IsSingleton = retMap[f.FieldTypePkgPath].IsSingleton
			f.LinkAPI = LinkAPI[f.PkgName+f.NodeName]
			resNodeProp.ChildrenFields = append(resNodeProp.ChildrenFields, f)
		}
		for _, f := range n.LinkFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			f.FieldCount = retMap[f.FieldTypePkgPath].FieldCount
			f.CRDName = retMap[f.FieldTypePkgPath].CRDName
			f.ChainAPI = retMap[f.FieldTypePkgPath].ChainAPI
			f.IsSingleton = retMap[f.FieldTypePkgPath].IsSingleton
			f.LinkAPI = LinkAPI[f.PkgName+f.NodeName]
			resNodeProp.LinkFields = append(resNodeProp.LinkFields, f)
		}
		for _, f := range n.LinksFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			f.FieldCount = retMap[f.FieldTypePkgPath].FieldCount
			f.CRDName = retMap[f.FieldTypePkgPath].CRDName
			f.ChainAPI = retMap[f.FieldTypePkgPath].ChainAPI
			f.IsSingleton = retMap[f.FieldTypePkgPath].IsSingleton
			f.LinkAPI = LinkAPI[f.PkgName+f.NodeName]
			resNodeProp.LinksFields = append(resNodeProp.LinksFields, f)
		}
		ResNodes = append(ResNodes, resNodeProp)
	}
	return ResNodes, nil
}
