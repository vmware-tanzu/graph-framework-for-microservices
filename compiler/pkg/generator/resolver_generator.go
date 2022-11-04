package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strconv"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/util"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type ReturnStatement struct {
	Alias       string
	ReturnType  string
	FieldCount  int
	CRDName     string
	ChainAPI    string
	IsSingleton bool
}

type FieldProperty struct {
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
	IsArrayStdType          bool
	IsSingleton             bool
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
	FieldCount              int
	CRDName                 string
	ChainAPI                string
	LinkAPI                 string
}

type NodeProperty struct {
	IsParentNode           bool
	HasParent              bool
	IsSingletonNode        bool
	IsNexusNode            bool
	BaseImportPath         string
	CrdName                string
	ResolverCount          int
	PkgName                string
	NodeName               string
	SchemaName             string
	Alias                  string
	ReturnType             string
	GroupResourceNameTitle string
	ChildFields            []FieldProperty
	LinkFields             []FieldProperty
	ChildrenFields         []FieldProperty
	LinksFields            []FieldProperty
	ArrayFields            []FieldProperty
	CustomFields           []FieldProperty
	NonStructFields        []FieldProperty
	GraphqlSchemaFields    []FieldProperty
	ResolverFields         map[string][]FieldProperty
	CustomQueries          []nexus.GraphQLQuery
}

// populateValuesForEachNode populates each node with required resolver properties
func populateValuesForEachNode(nodes []*NodeProperty, linkAPI map[string]string, retMap map[string]ReturnStatement) []NodeProperty {
	var nodeProperties []NodeProperty
	for _, n := range nodes {
		resNodeProp := NodeProperty{}
		resNodeProp.GroupResourceNameTitle = util.GetGroupResourceNameTitle(n.NodeName)
		resNodeProp.Alias = n.Alias
		resNodeProp.ReturnType = n.ReturnType
		resNodeProp.BaseImportPath = n.BaseImportPath
		resNodeProp.GraphqlSchemaFields = n.GraphqlSchemaFields
		resNodeProp.IsSingletonNode = n.IsSingletonNode
		resNodeProp.CustomQueries = n.CustomQueries
		resNodeProp.IsNexusNode = n.IsNexusNode
		resNodeProp.ResolverFields = n.ResolverFields
		resNodeProp.ResolverCount = n.ResolverCount
		resNodeProp.PkgName = n.PkgName
		resNodeProp.NodeName = n.NodeName
		resNodeProp.SchemaName = n.SchemaName

		// populate return values for root of the graph
		if !n.HasParent && n.IsParentNode {
			resNodeProp.Alias = retMap[resNodeProp.PkgName+resNodeProp.NodeName].Alias
			resNodeProp.ReturnType = retMap[resNodeProp.PkgName+resNodeProp.NodeName].ReturnType
			resNodeProp.IsParentNode = true
		}

		// populate return values for child fields
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
			f.LinkAPI = linkAPI[f.PkgName+f.NodeName]
			resNodeProp.ChildFields = append(resNodeProp.ChildFields, f)
			resNodeProp.GraphqlSchemaFields = append(resNodeProp.GraphqlSchemaFields, f)
		}

		// populate return values for multiple child fields
		for _, f := range n.ChildrenFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			f.FieldCount = retMap[f.FieldTypePkgPath].FieldCount
			f.CRDName = retMap[f.FieldTypePkgPath].CRDName
			f.ChainAPI = retMap[f.FieldTypePkgPath].ChainAPI
			f.IsSingleton = retMap[f.FieldTypePkgPath].IsSingleton
			f.LinkAPI = linkAPI[f.PkgName+f.NodeName]
			resNodeProp.ChildrenFields = append(resNodeProp.ChildrenFields, f)
		}

		// populate return values for link fields
		for _, f := range n.LinkFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			f.FieldCount = retMap[f.FieldTypePkgPath].FieldCount
			f.CRDName = retMap[f.FieldTypePkgPath].CRDName
			f.ChainAPI = retMap[f.FieldTypePkgPath].ChainAPI
			f.IsSingleton = retMap[f.FieldTypePkgPath].IsSingleton
			f.LinkAPI = linkAPI[f.PkgName+f.NodeName]
			resNodeProp.LinkFields = append(resNodeProp.LinkFields, f)
		}

		// populate return values for multiple link fields
		for _, f := range n.LinksFields {
			f.ReturnType = retMap[f.FieldTypePkgPath].ReturnType
			f.Alias = retMap[f.FieldTypePkgPath].Alias
			f.FieldCount = retMap[f.FieldTypePkgPath].FieldCount
			f.CRDName = retMap[f.FieldTypePkgPath].CRDName
			f.ChainAPI = retMap[f.FieldTypePkgPath].ChainAPI
			f.IsSingleton = retMap[f.FieldTypePkgPath].IsSingleton
			f.LinkAPI = linkAPI[f.PkgName+f.NodeName]
			resNodeProp.LinksFields = append(resNodeProp.LinksFields, f)
		}

		nodeProperties = append(nodeProperties, resNodeProp)
	}

	return nodeProperties
}

func populateValuesForResolver(nodes []*NodeProperty, parentsMap map[string]parser.NodeHelper,
	crdNameMap, nonStructMap map[string]string) (map[string]string, map[string]ReturnStatement) {
	linkAPI := make(map[string]string)
	retMap := make(map[string]ReturnStatement)

	for _, n := range nodes {
		var (
			retType, aliasVal, listRetVal, ChainAPI string
			fieldCount                              int
			IsSingleton                             bool
		)

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
						ChainAPI += fmt.Sprintf(".%s()", crdNameMap[i])
					} else {
						ChainAPI += fmt.Sprintf(".%s(getParentName(obj.ParentLabels, %q))", crdNameMap[i], i)
					}
				} else {
					if childNode, ok := prevNode.Children[i]; ok {
						if currentNode.IsSingleton {
							ChainAPI += fmt.Sprintf(".%s()", childNode.FieldName)
						} else {
							ChainAPI += fmt.Sprintf(".%s(getParentName(obj.ParentLabels, %q))", childNode.FieldName, i)
						}
					}
				}
				// cache the non-leaf node
				prevNode = currentNode
			}

			// Create linkAPI
			if n.IsSingletonNode {
				IsSingleton = true
				if !n.HasParent && n.IsParentNode {
					linkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO())", ChainAPI, n.PkgName+n.NodeName)
				} else {
					linkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO())", ChainAPI, prevNode.Children[n.CrdName].FieldName)
				}
			} else {
				IsSingleton = false
				if !n.HasParent && n.IsParentNode {
					linkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO(), getParentName(obj.ParentLabels, %q))", ChainAPI, n.PkgName+n.NodeName, n.CrdName)
				} else {
					linkAPI[n.PkgName+n.NodeName] = fmt.Sprintf("%s.Get%s(context.TODO(), getParentName(obj.ParentLabels, %q))", ChainAPI, prevNode.Children[n.CrdName].FieldName, n.CrdName)
				}
			}
		}
		for _, i := range n.ResolverFields[n.PkgName+n.NodeName] {
			if i.IsMapTypeField || i.IsStringType {
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

	return linkAPI, retMap
}

func constructNexusTypeMap(nodes []*NodeProperty) map[string]string {
	crdNameMap := make(map[string]string)
	for _, n := range nodes {
		if n.IsNexusNode || n.IsSingletonNode {
			crdNameMap[n.CrdName] = n.PkgName + n.NodeName
		}
	}
	return crdNameMap
}

// processNonNexusFields process and populates properties for each non nexus fields
// <Domain  string>
func processNonNexusFields(aliasNameMap map[string]string, node *ast.TypeSpec,
	nodeProp *NodeProperty, simpleGroupTypeName string) {
	resField := make(map[string][]FieldProperty)
	for _, f := range parser.GetSpecFields(node) {
		var (
			fieldProp FieldProperty
			err       error
		)
		typeString := ConstructType(aliasNameMap, f)
		// populate each field properties
		if len(f.Names) > 0 {
			fieldProp.FieldName, err = parser.GetNodeFieldName(f)
			if err != nil {
				log.Fatalf("failed to determine field name: %v", err)
			}
			fieldProp.FieldType = typeString
			fieldProp.PkgName = simpleGroupTypeName
			fieldProp.NodeName = node.Name.String()
		}

		if parser.IgnoreField(f) {
			continue
		}

		if parser.IsJsonStringField(f) {
			fieldProp.IsStringType = true
			fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
			resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
		} else {
			stdType := convertGraphqlStdType(typeString)
			// standard type
			if len(stdType) != 0 {
				fieldProp.IsStdTypeField = true
				fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, stdType)
				resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
			} else {
				fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
				fieldProp.IsStringType = true
				resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
			}
		}
		nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
		nodeProp.ResolverFields = resField
	}
}

func findTypeAndPkgForField(ptParts []string, importMap map[string]string, pkgs map[string]parser.Package) (string, *parser.Package) {
	structPkg := ptParts[0]
	structType := ptParts[1]

	pkgPath, ok := importMap[structPkg]
	if !ok {
		log.Errorf("Cannot find the package name %s for the type %s", structPkg, structType)
		return "", nil
	}

	importPath, err := strconv.Unquote(pkgPath)
	if err != nil {
		log.Errorf("Failed to parse package %s for the type %s with error %v", pkgPath, structType, err)
		return "", nil
	}

	p, ok := pkgs[importPath]
	if !ok {
		log.Errorf("Cannot find the package details from the path %s for the type %s", importPath, structType)
		return "", nil
	}

	return structType, &p
}

// processNexusFields process and populates properties for each nexus fields
// <gns.Gns `nexus:"child"`>
func processNexusFields(pkg parser.Package, aliasNameMap map[string]string, node *ast.TypeSpec,
	nodeProp *NodeProperty, simpleGroupTypeName string, pkgs map[string]parser.Package) {
	importMap := pkg.GetImportMap()
	for _, nf := range parser.GetNexusFields(node) {
		var (
			fieldProp FieldProperty
			err       error
		)
		if len(nf.Names) > 0 {
			fieldProp.FieldName, err = parser.GetNodeFieldName(nf)
			if err != nil {
				log.Fatalf("Failed to determine field name: %v", err)
			}
			fieldProp.PkgName = simpleGroupTypeName
			fieldProp.NodeName = node.Name.String()
		}

		// Except for nexus fields (nexus.Node and nexus.SingletonNode),
		// this will check other fields to see whether they have nexus secrets annotated on them
		// If yes, the field is ignored in the response.
		if !parser.IsNexusTypeField(nf) {
			nfType := parser.GetFieldType(nf)
			fieldPkg := &pkg
			structType := nfType

			if ptParts := strings.Split(nfType, "."); len(ptParts) == 2 { //service_group.SvcGroup
				structType, fieldPkg = findTypeAndPkgForField(ptParts, importMap, pkgs)
			}
			if len(structType) != 0 {
				if _, ok := parser.GetNexusSecretSpecAnnotation(*fieldPkg, structType); ok {
					log.Debugf("Ignoring the field %s since the node is annotated as nexus secret", fieldProp.FieldName)
					continue
				}
			}
		}

		// `Ignore:true` annotation used to ignore the specific field `nexus-graphql:"ignore:true"`
		if parser.IgnoreField(nf) {
			continue
		}
		// `type:string` annotation used to consider the type as string `nexus-graphql:"type:string"`
		if parser.IsJsonStringField(nf) {
			fieldProp.IsStringType = true
			fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s", fieldProp.FieldName, "String")
		}

		// denote field is nexus or singletonField type
		if parser.IsNexusTypeField(nf) {
			fieldProp.IsNexusOrSingletonField = true
			// Add Custom Query + ID
			fieldProp.SchemaFieldName = CustomQuerySchema
			for _, customQuery := range nodeProp.CustomQueries {
				fieldProp.SchemaFieldName += CustomQueryToGraphqlSchema(customQuery)
				if unicode.IsUpper(rune(customQuery.Name[0])) {
					var customQueryFieldProp FieldProperty
					customQueryFieldProp.IsResolver = true
					customQueryFieldProp.FieldName = customQuery.Name
					nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, customQueryFieldProp)
				}

			}
		}

		// nexus link field
		typeString := ConstructType(aliasNameMap, nf)
		if parser.IsOnlyLinkField(nf) {
			schemaTypeName, resolverTypeName := ValidateImportPkg(nodeProp.PkgName, typeString, importMap, pkgs)
			// `type:string` annotation used to consider the type as string `nexus-graphql:"type:string"`
			fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, schemaTypeName)
			fieldProp.IsResolver = true
			fieldProp.IsNexusTypeField = true
			fieldProp.FieldType = typeString
			fieldProp.FieldTypePkgPath = resolverTypeName
			fieldProp.SchemaTypeName = schemaTypeName
			fieldProp.BaseTypeName = getBaseNodeType(typeString)
			nodeProp.LinkFields = append(nodeProp.LinkFields, fieldProp)
		}

		// nexus child field
		if parser.IsOnlyChildField(nf) {
			schemaTypeName, resolverTypeName := ValidateImportPkg(nodeProp.PkgName, typeString, importMap, pkgs)
			fieldProp.SchemaFieldName = fmt.Sprintf("%s: %s!", fieldProp.FieldName, schemaTypeName)
			fieldProp.SchemaTypeName = schemaTypeName
			fieldProp.IsNexusTypeField = true
			fieldProp.FieldType = typeString
			fieldProp.FieldTypePkgPath = resolverTypeName
			fieldProp.BaseTypeName = getBaseNodeType(typeString)
			nodeProp.ChildFields = append(nodeProp.ChildFields, fieldProp)
		}

		// nexus children or links field
		if parser.IsNamedChildOrLink(nf) {
			fieldProp.IsChildrenOrLinks = true
			schemaTypeName, resolverTypeName := ValidateImportPkg(nodeProp.PkgName, typeString, importMap, pkgs)
			fieldProp.SchemaFieldName = fmt.Sprintf("%s(Id: ID): [%s!]", fieldProp.FieldName, schemaTypeName)
			fieldProp.IsResolver = true
			fieldProp.IsNexusTypeField = true
			fieldProp.FieldType = typeString
			fieldProp.FieldTypePkgPath = resolverTypeName
			fieldProp.SchemaTypeName = schemaTypeName
			fieldProp.BaseTypeName = getBaseNodeType(typeString)
			if parser.IsOnlyChildrenField(nf) {
				nodeProp.ChildrenFields = append(nodeProp.ChildrenFields, fieldProp)
			} else {
				nodeProp.LinksFields = append(nodeProp.LinksFields, fieldProp)
			}
		}
		// no. of resolver field in a node
		if fieldProp.IsResolver {
			nodeProp.ResolverCount += 1
		}
		if !parser.IsOnlyChildField(nf) {
			nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
		}
	}
}

/*
collect and construct type alias field into map recursively before
populating the nexus node and custom struct type
ex. nonStructMap[pkgName] = nodeType  -->  nonStructMap["root"] = "AliasTypeFoo"
*/
func constructAliasType(sortedPackages []parser.Package) map[string]string {
	nonStructMap := make(map[string]string)
	for _, pkg := range sortedPackages {
		for _, node := range pkg.GetNonStructTypes() {
			pkgName := fmt.Sprintf("%s_%s", pkg.Name, parser.GetTypeName(node))
			// NonStruct Map
			nonStructType := types.ExprString(node.Type)
			nonStructMap[pkgName] = nonStructType
		}
	}

	return nonStructMap
}

func setNexusProperties(nodeHelper parser.NodeHelper, node *ast.TypeSpec, nodeProp *NodeProperty) {
	if len(nodeHelper.Parents) > 0 {
		nodeProp.HasParent = true
	}

	if parser.IsSingletonNode(node) {
		nodeProp.IsSingletonNode = true
	}

	if parser.IsNexusNode(node) {
		nodeProp.IsNexusNode = true
	}
}

// isRootOfGraph intended to allow only one root of the graph,
// if we receive multiple node in such behaviour, then we allow the first node and the rest will be skipped with error
// arg: `parents` indicates the node's parents, `rootOfGraph` indicates if the received node is the root of the graph or not.
func isRootOfGraph(parents []string, rootOfGraph bool) bool {
	if len(parents) == 0 && !rootOfGraph {
		return true
	}

	return rootOfGraph
}

// GenerateGraphqlResolverVars populates the node and its field properties required to generate graphql resolver
func GenerateGraphqlResolverVars(baseGroupName, crdModulePath string, pkgs parser.Packages, parentsMap map[string]parser.NodeHelper) ([]NodeProperty, error) {
	sortedKeys := make([]string, 0, len(pkgs))
	for k := range pkgs {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	sortedPackages := make([]parser.Package, len(pkgs))
	for i, k := range sortedKeys {
		sortedPackages[i] = pkgs[k]
	}

	// Iterate all the struct type and it's fields in the sortedPackages and
	// set the node and field properties accordingly.
	var nodes []*NodeProperty
	aliasNameMap := make(map[string]string)
	rootOfGraph := false
	for _, pkg := range sortedPackages {
		simpleGroupTypeName := util.GetSimpleGroupTypeName(pkg.Name)
		// Iterating struct type
		for _, node := range pkg.GetStructs() {
			// Skip Empty struct type
			if len(parser.GetNexusFields(node)) == 0 && len(parser.GetSpecFields(node)) == 0 {
				continue
			}

			typeName := parser.GetTypeName(node)
			if _, ok := parser.GetNexusSecretSpecAnnotation(pkg, typeName); ok {
				log.Debugf("Ignoring the node %s since the node is annotated as nexus secret", typeName)
				continue
			}

			nodeProp := &NodeProperty{}
			// populate node properties
			nodeProp.PkgName = simpleGroupTypeName
			nodeProp.NodeName = node.Name.String()
			nodeProp.BaseImportPath = crdModulePath
			nodeProp.CrdName = util.GetCrdName(node.Name.String(), pkg.Name, baseGroupName)
			nodeHelper := parentsMap[nodeProp.CrdName]
			nodeProp.IsParentNode = parser.IsNexusNode(node)
			nodeProp.CustomQueries = nodeHelper.GraphqlSpec.Queries

			if parser.IsNexusNode(node) && len(nodeHelper.Parents) == 0 && rootOfGraph {
				log.Errorf("Can't allow multiple root of the graph, skipping Node:%s", nodeProp.NodeName)
				continue
			}

			if parser.IsNexusNode(node) {
				rootOfGraph = isRootOfGraph(nodeHelper.Parents, rootOfGraph)
			}
			setNexusProperties(nodeHelper, node, nodeProp)
			nodeProp.SchemaName = fmt.Sprintf("%s_%s", pkg.Name, parser.GetTypeName(node))

			// Iterate each node's nexus fields and set its properties
			processNexusFields(pkg, aliasNameMap, node, nodeProp, simpleGroupTypeName, pkgs)

			// Iterate each node's non-nexus fields and set its properties
			processNonNexusFields(aliasNameMap, node, nodeProp, simpleGroupTypeName)
			nodes = append(nodes, nodeProp)
		}
	}

	crdNameMap := constructNexusTypeMap(nodes)
	// populate return values of each Node for resolver
	nonStructMap := constructAliasType(sortedPackages)
	linkAPI, retMap := populateValuesForResolver(nodes, parentsMap, crdNameMap, nonStructMap)

	// populate return values of each node
	nodeProperties := populateValuesForEachNode(nodes, linkAPI, retMap)

	return nodeProperties, nil
}
