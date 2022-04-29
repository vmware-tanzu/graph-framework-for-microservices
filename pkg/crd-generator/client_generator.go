package crd_generator

import (
	"go/ast"
	"sort"
	"strings"
	"text/template"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"

	log "github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

func generateNexusClientVars(baseGroupName, crdModulePath string, pkgs parser.Packages, parentsMap map[string]parser.NodeHelper) (clientVars, error) {
	var vars clientVars

	vars.BaseClientsetImport = `"` + crdModulePath + `client/clientset/versioned"`
	vars.HelperImport = `"` + crdModulePath + `helper"`

	sortedKeys := make([]string, 0, len(pkgs))
	for k := range pkgs {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	sortedPackages := make([]parser.Package, len(pkgs))
	for i, k := range sortedKeys {
		sortedPackages[i] = pkgs[k]
	}

	for _, pkg := range sortedPackages {
		if len(pkg.GetNexusNodes()) > 0 {
			// TODO make version configurable
			version := "v1"

			importPath := util.GetImportPath(pkg.Name, baseGroupName, version)
			baseImportName := util.GetBaseImportName(pkg.Name, baseGroupName, version)

			vars.BaseImports += baseImportName + ` "` + crdModulePath + "apis/" + importPath + `"` +
				"\n" // eg baseroothelloworldv1 "helloworld/nexus/generated/apis/root.helloworld.com/v1"

			groupVarName := util.GetGroupVarName(pkg.Name, baseGroupName, version)
			groupTypeName := util.GetGroupTypeName(pkg.Name, baseGroupName, version)
			vars.ClientsetsApiGroups += groupVarName + " *" + groupTypeName + "\n" // eg rootHelloworldV1 *RootHelloworldV1

			initClient := "client." + groupVarName + " = new" + groupTypeName +
				"(client)\n" // eg client.rootHelloworldV1 = newRootHelloworldV1(client)
			vars.InitApiGroups += initClient

			clientsetMethod := "func (c *Clientset) " + groupTypeName + "() *" + groupTypeName + " {\n" + "return c." +
				groupVarName + "\n}\n" // eg
			// func (c *Clientset) RootHelloworldV1() *RootHelloworldV1 {
			//	return c.rootHelloworldV1
			// }
			vars.ClientsetsApiGroupMethods += clientsetMethod

			var groupVars apiGroupsVars
			groupVars.GroupTypeName = groupTypeName

			for _, node := range pkg.GetNexusNodes() {
				var clientGroupVars apiGroupsClientVars
				clientGroupVars.GroupTypeName = groupTypeName
				clientGroupVars.CrdName = util.GetCrdName(node.Name.String(), pkg.Name, baseGroupName)
				err := resolveNode(baseImportName, pkg, baseGroupName, version, &groupVars, &clientGroupVars, node, parentsMap)
				if err != nil {
					return clientVars{}, err
				}
				vars.Nodes = append(vars.Nodes, clientGroupVars)
			}

			apiGroup, err := renderApiGroup(groupVars)
			if err != nil {
				return clientVars{}, err
			}
			vars.ApiGroups += apiGroup

		}
	}
	return vars, nil
}

func resolveNode(baseImportName string, pkg parser.Package, baseGroupName, version string, groupVars *apiGroupsVars, clientGroupVars *apiGroupsClientVars, node *ast.TypeSpec, parentsMap map[string]parser.NodeHelper) error {
	pkgName := pkg.Name
	baseNodeName := node.Name.Name // eg Root
	groupResourceName := util.GetGroupResourceName(baseNodeName)
	groupResourceNameTitle := util.GetGroupResourceNameTitle(baseNodeName)
	groupResourceType := util.GetGroupResourceType(baseNodeName, pkgName, baseGroupName, version)
	groupVars.GroupResources += groupResourceName + " *" +
		groupResourceType + "\n" // eg roots *rootRootHelloWorld
	groupVars.GroupResourcesInit += groupResourceName + ": &" + groupResourceType +
		"{\n client: client,\n},\n" // eg
	// 		roots: &rootRootHelloWorld{
	//			client: client,
	//		},
	groupVars.GroupResourcesDefs += "type " + groupResourceType + " struct {\n  client *Clientset\n}\n" // eg
	// type rootRootHelloWorld struct {
	//	client *Clientset
	// }
	groupVars.GroupResourcesDefs += "func (obj *" + util.GetGroupTypeName(pkgName, baseGroupName, version) + ") " +
		groupResourceNameTitle + "() *" +
		groupResourceType + " {\n" + "return obj." + groupResourceName + "\n}\n" // eg
	// func (r *RootHelloworldV1) Roots() *rootRootHelloWorld {
	//	return r.roots
	// }

	clientGroupVars.BaseImportName = baseImportName
	clientGroupVars.GroupBaseImport = baseImportName + "." + baseNodeName
	clientGroupVars.GroupResourceType = groupResourceType
	clientGroupVars.GroupResourceNameTitle = groupResourceNameTitle

	// TODO support resolution of links which are not nexus nodes https://jira.eng.vmware.com/browse/NPT-112
	children := parser.GetChildFields(node)
	childrenAndLinks := children
	childrenAndLinks = append(childrenAndLinks, parser.GetLinkFields(node)...)

	if len(children) > 0 {
		clientGroupVars.HasChildren = true
	}
	for _, link := range childrenAndLinks {
		linkInfo := getFieldInfo(pkg, link)
		clientVarsLink := apiGroupsClientVarsLink{
			FieldName:              linkInfo.fieldName,
			FieldNameGvk:           util.GetGvkFieldTagName(linkInfo.fieldName),
			Group:                  util.GetGroupName(linkInfo.pkgName, baseGroupName),
			Kind:                   linkInfo.fieldType,
			GroupBaseImport:        util.GetBaseImportName(linkInfo.pkgName, baseGroupName, version) + "." + linkInfo.fieldType,
			GroupResourceNameTitle: util.GetGroupResourceNameTitle(linkInfo.fieldType),
			GroupTypeName:          util.GetGroupTypeName(linkInfo.pkgName, baseGroupName, version),
		}
		if parser.IsMapField(link) {
			clientVarsLink.IsNamed = true
		} else {
			clientVarsLink.IsNamed = false
		}
		if parser.IsLinkField(link) {
			clientGroupVars.Links = append(clientGroupVars.Links, clientVarsLink)
		} else {
			clientGroupVars.Children = append(clientGroupVars.Children, clientVarsLink)
		}
		clientGroupVars.LinksAndChildren = append(clientGroupVars.LinksAndChildren, clientVarsLink)
	}

	for _, f := range parser.GetSpecFields(node) {
		fieldInfo := getFieldInfo(pkg, f)
		var vars apiGroupsClientVarsLink
		vars.FieldName = fieldInfo.fieldName
		vars.FieldNameTag = util.GetTag(fieldInfo.fieldName)

		clientGroupVars.Fields = append(clientGroupVars.Fields, vars)
	}

	nodeHelper := parentsMap[clientGroupVars.CrdName]
	clientGroupVars.Group = util.GetGroupFromCrdName(clientGroupVars.CrdName)
	clientGroupVars.Kind = nodeHelper.Name
	if len(nodeHelper.Parents) > 0 {
		parentCrdName := nodeHelper.Parents[len(nodeHelper.Parents)-1]
		parentHelper := parentsMap[parentCrdName]

		clientGroupVars.Parent.HasParent = true
		clientGroupVars.Parent.IsNamed = parentHelper.Children[clientGroupVars.CrdName].IsNamed
		clientGroupVars.Parent.CrdName = parentCrdName
		clientGroupVars.Parent.GroupTypeName = util.GetGroupTypeName(
			util.GetPackageNameFromCrdName(parentCrdName), baseGroupName, version)
		clientGroupVars.Parent.GroupResourceNameTitle = util.GetGroupResourceNameTitle(parentHelper.Name)
		clientGroupVars.Parent.GvkFieldName = parentHelper.Children[clientGroupVars.CrdName].FieldNameGvk
	}

	return nil
}

type fieldInfo struct {
	pkgName   string
	fieldName string
	fieldType string
}

func getFieldInfo(pkg parser.Package, f *ast.Field) fieldInfo {
	var info fieldInfo
	var err error
	info.fieldName, err = parser.GetFieldName(f)
	if err != nil {
		log.Fatalf("Failed to get name of field: %v", err)
	}
	currentPkgName := pkg.Name

	chType := parser.GetFieldType(f)
	s := strings.Split(chType, ".")
	info.fieldType = s[len(s)-1]

	split := strings.Split(chType, ".")
	// overwrite pkg name of links or children from different packages
	if len(split) > 1 {
		info.pkgName = split[0]
		// overwrite pkg name for node which uses named import like 'sg "helloworld.com/service-groups"'
		for _, imp := range pkg.GetImports() {
			if imp.Name.String() == info.pkgName {
				s := strings.Split(imp.Path.Value, "/")
				info.pkgName = strings.TrimSuffix(s[len(s)-1], "\"")
			}
		}
		info.pkgName = util.RemoveSpecialChars(info.pkgName)
	} else {
		info.pkgName = util.RemoveSpecialChars(currentPkgName)
	}

	return info
}

var apiGroupTmpl = `
type {{.GroupTypeName}} struct {
	{{.GroupResources}}
}

func new{{.GroupTypeName}}(client *Clientset) *{{.GroupTypeName}} {
	return &{{.GroupTypeName}}{
		{{.GroupResourcesInit}}
	}
}

{{.GroupResourcesDefs}}
`

type apiGroupsVars struct {
	GroupTypeName      string
	GroupResourcesInit string
	GroupResources     string
	GroupResourcesDefs string
}

func renderApiGroup(vars apiGroupsVars) (string, error) {
	tmpl, err := template.New("tmpl").Parse(apiGroupTmpl)
	if err != nil {
		return "", err
	}
	ren, err := renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	return ren.String(), nil
}

type apiGroupsClientVars struct {
	apiGroupsVars
	CrdName                string
	ResolveLinksDelete     string
	HasChildren            bool
	BaseImportName         string
	GroupResourceType      string
	GroupResourceNameTitle string
	GroupBaseImport        string
	Group                  string
	Kind                   string

	Parent struct {
		IsNamed                bool
		HasParent              bool
		CrdName                string
		GvkFieldName           string
		GroupTypeName          string
		GroupResourceNameTitle string
	}
	ForUpdatePatches string

	Links            []apiGroupsClientVarsLink
	Children         []apiGroupsClientVarsLink
	LinksAndChildren []apiGroupsClientVarsLink
	Fields           []apiGroupsClientVarsLink
}

type apiGroupsClientVarsLink struct {
	FieldName              string
	FieldNameTag           string
	FieldNameGvk           string
	Group                  string
	Kind                   string
	GroupBaseImport        string
	IsNamed                bool
	GroupTypeName          string
	GroupResourceNameTitle string
}
