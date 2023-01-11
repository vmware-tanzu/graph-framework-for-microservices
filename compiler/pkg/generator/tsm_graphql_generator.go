package generator

import (
	"fmt"
	"go/ast"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/util"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type GraphQLSchemaType int

const (
	Standard GraphQLSchemaType = iota
	Array
	JsonMarshal
	Child
	Link
	NamedChild
	NamedLink
)

// tsmPopulateValuesForEachNode populates each node with required resolver properties
func tsmPopulateValuesForEachNode(nodes []*NodeProperty, linkAPI map[string]string, retMap map[string]ReturnStatement) []NodeProperty {
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
			f.LinkAPI = linkAPI[f.PkgName+f.NodeName]
			resNodeProp.ChildFields = append(resNodeProp.ChildFields, f)
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

func tsmPopulateValuesForResolver(nodes []*NodeProperty, parentsMap map[string]parser.NodeHelper,
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

// tsmProcessNonNexusFields process and populates properties for each non nexus fields
// <Domain  string>
func tsmProcessNonNexusFields(pkg parser.Package, aliasNameMap map[string]string, node *ast.TypeSpec,
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
			fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(JsonMarshal, fieldProp.FieldName, "String", "id: ID", f, pkg)
			resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
		} else if parser.IsFieldAnnotationPresent(f, parser.GRAPHQL_ALIAS_TYPE_ANNOTATION) || parser.IsFieldAnnotationPresent(f, parser.GRAPHQL_ALIAS_NAME_ANNOTATION) {
			fieldProp.SchemaFieldName = GetGraphQLAliasValue(fieldProp.FieldName, f)
			resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
		} else if parser.IsArrayField(f) {
			arr := regexp.MustCompile(`^(\[])`).ReplaceAllString(typeString, "")
			if !strings.Contains(arr, ".") {
				stdType := convertGraphqlStdType(arr)
				fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(Array, fieldProp.FieldName, stdType, "id: ID", f, pkg)
				resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
			}
		} else {
			stdType := convertGraphqlStdType(typeString)
			// standard type
			if len(stdType) != 0 {
				fieldProp.IsStdTypeField = true
				fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(Standard, fieldProp.FieldName, stdType, "id: ID", f, pkg)
				resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
			} else {
				fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(Standard, fieldProp.FieldName, "String", "id: ID", f, pkg)
				fieldProp.IsStringType = true
				resField[nodeProp.PkgName+nodeProp.NodeName] = append(resField[nodeProp.PkgName+nodeProp.NodeName], fieldProp)
			}
		}
		nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
		nodeProp.ResolverFields = resField
	}
}

// tsmProcessNexusFields process and populates properties for each nexus fields
// <gns.Gns `nexus:"child"`>
func tsmProcessNexusFields(pkg parser.Package, aliasNameMap map[string]string, node *ast.TypeSpec,
	nodeProp *NodeProperty, simpleGroupTypeName string, pkgs map[string]parser.Package,
	gqlSpecMap map[string]nexus.GraphQLSpec) {
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
			fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(JsonMarshal, fieldProp.FieldName, "String", "id: ID", nf, pkg)
		}

		// denote field is nexus or singletonField type
		if parser.IsNexusTypeField(nf) {
			fieldProp.IsNexusOrSingletonField = true
			// get nexus schemaFieldName from GraphQlSpec "IdName" & "IdNullable"
			fieldProp.SchemaFieldName = GetNexusSchemaFieldName(nodeProp.GraphQlSpec)
			for _, customQuery := range nodeProp.CustomQueries {
				cq := CustomQueryToGraphqlSchema(customQuery)
				// In TSM DM "@timeseriesAPI" directives is need to added along with returnType "TimeSeriesData"
				fieldProp.SchemaFieldName += "\n" + strings.ReplaceAll(cq, "TimeSeriesData", fmt.Sprintf("TimeSeriesData @timeseriesAPI(file: \"../../tsquery/timeSeriesQuery\", handler: \"%s\")", customQuery.Name))
				var customQueryFieldProp FieldProperty
				customQueryFieldProp.IsResolver = true
				customQueryFieldProp.FieldName = customQuery.Name
				nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, customQueryFieldProp)
				nodeProp.ResolverCount += 1
			}
		}

		// nexus link field
		typeString := ConstructType(aliasNameMap, nf)
		if parser.IsOnlyLinkField(nf) {
			schemaTypeName, resolverTypeName := ValidateImportPkg(nodeProp.PkgName, typeString, importMap, pkgs)
			fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(Link, fieldProp.FieldName, schemaTypeName, "id: ID", nf, pkg)
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
			fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(Child, fieldProp.FieldName, schemaTypeName, "id: ID", nf, pkg)
			fieldProp.SchemaTypeName = schemaTypeName
			fieldProp.IsResolver = true
			fieldProp.IsNexusTypeField = true
			fieldProp.FieldType = typeString
			fieldProp.FieldTypePkgPath = resolverTypeName
			fieldProp.BaseTypeName = getBaseNodeType(typeString)
			nodeProp.ChildFields = append(nodeProp.ChildFields, fieldProp)
		}

		// nexus children or links field
		if parser.IsNamedChildOrLink(nf) {
			var listArg string
			fieldProp.IsChildrenOrLinks = true
			schemaTypeName, resolverTypeName := ValidateImportPkg(nodeProp.PkgName, typeString, importMap, pkgs)
			// Annotation `nexus-graphql-args:"name: String"` use to specify graphql arguments
			AnnotatedGqlArgs := parser.GetFieldAnnotationVal(nf, parser.GRAPHQL_ARGS_ANNOTATION)
			if AnnotatedGqlArgs != "" {
				listArg = AnnotatedGqlArgs
			} else {
				if val, ok := parser.GetNexusGraphqlSpecAnnotation(pkg, typeString); ok {
					gqlspec := gqlSpecMap[fmt.Sprintf("%s.%s", pkg.Name, val)]
					listArg = GetNexusSchemaFieldName(gqlspec)
				} else {
					listArg = GetNodeDetails(nodeProp.PkgName, typeString, importMap, pkgs, gqlSpecMap)
				}
			}

			sType := NamedChild
			if parser.IsLinkField(nf) {
				sType = NamedLink
			}
			fieldProp.SchemaFieldName = getTsmGraphqlSchemaFieldName(sType, fieldProp.FieldName, schemaTypeName, listArg, nf, pkg)
			fieldProp.IsResolver = true
			fieldProp.IsNexusTypeField = true
			fieldProp.FieldType = typeString
			fieldProp.FieldTypePkgPath = resolverTypeName
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
		nodeProp.GraphqlSchemaFields = append(nodeProp.GraphqlSchemaFields, fieldProp)
	}
}

// GenerateTsmGraphqlSchemaVars populates the node and its field properties required to generate graphql resolver
func GenerateTsmGraphqlSchemaVars(baseGroupName, crdModulePath string, pkgs parser.Packages, parentsMap map[string]parser.NodeHelper) ([]NodeProperty, error) {
	sortedKeys := make([]string, 0, len(pkgs))
	gqlSpecMap := parser.ParseGraphqlSpecs(pkgs)
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
	//rootOfGraph := false
	for _, pkg := range sortedPackages {
		simpleGroupTypeName := util.GetSimpleGroupTypeName(pkg.Name)
		// Iterating struct type
		for _, node := range pkg.GetStructs() {
			// Nexus GraphQlSpec by default "IdNullable" value is true
			gqlspec := nexus.GraphQLSpec{
				IdName:     "",
				IdNullable: true,
			}
			// Skip Empty struct type
			if len(parser.GetNexusFields(node)) == 0 && len(parser.GetSpecFields(node)) == 0 {
				continue
			}

			typeName := parser.GetTypeName(node)
			if _, ok := parser.GetNexusSecretSpecAnnotation(pkg, typeName); ok {
				log.Debugf("Ignoring the node %s since the node is annotated as nexus secret", typeName)
				continue
			}

			if val, ok := parser.GetNexusGraphqlSpecAnnotation(pkg, typeName); ok {
				gqlspec = gqlSpecMap[fmt.Sprintf("%s.%s", pkg.Name, val)]
			}

			nodeProp := &NodeProperty{}
			// populate node properties
			nodeProp.PkgName = simpleGroupTypeName
			nodeProp.NodeName = node.Name.String()
			nodeProp.BaseImportPath = crdModulePath
			nodeProp.CrdName = util.GetCrdName(node.Name.String(), pkg.Name, baseGroupName)
			nodeHelper := parentsMap[nodeProp.CrdName]
			nodeProp.IsParentNode = parser.IsNexusNode(node)
			nodeProp.CustomQueries = nodeHelper.GraphqlQuerySpec.Queries
			nodeProp.GraphQlSpec = gqlspec

			setNexusProperties(nodeHelper, node, nodeProp)
			if pkg.Name == "global" {
				nodeProp.SchemaName = parser.GetTypeName(node)
			} else if pkg.Name == "tsm" {
				continue
			} else {
				nodeProp.SchemaName = fmt.Sprintf("%s_%s", pkg.Name, parser.GetTypeName(node))
			}
			// Iterate each node's nexus fields and set its properties
			tsmProcessNexusFields(pkg, aliasNameMap, node, nodeProp, simpleGroupTypeName, pkgs, gqlSpecMap)

			// Iterate each node's non-nexus fields and set its properties
			tsmProcessNonNexusFields(pkg, aliasNameMap, node, nodeProp, simpleGroupTypeName)
			nodes = append(nodes, nodeProp)
		}
	}

	crdNameMap := constructNexusTypeMap(nodes)
	// populate return values of each Node for resolver
	nonStructMap := constructAliasType(sortedPackages)
	linkAPI, retMap := tsmPopulateValuesForResolver(nodes, parentsMap, crdNameMap, nonStructMap)

	// populate return values of each node
	nodeProperties := tsmPopulateValuesForEachNode(nodes, linkAPI, retMap)

	return nodeProperties, nil
}
