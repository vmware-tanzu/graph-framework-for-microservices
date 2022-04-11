package crd_generator

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

func generateNexusClientVars(baseGroupName, crdModulePath string, pkgs parser.Packages) (clientVars, error) {
	var vars clientVars

	vars.BaseClientsetImport = `"` + crdModulePath + `client/clientset/versioned"`

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

			importPath := getImportPath(pkg.Name, baseGroupName, version)
			baseImportName := getBaseImportName(pkg.Name, baseGroupName, version)

			vars.BaseImports += baseImportName + ` "` + crdModulePath + "apis/" + importPath + `"` +
				"\n" // eg baseroothelloworldv1 "helloworld/nexus/generated/apis/root.helloworld.com/v1"

			groupVarName := getGroupVarName(pkg.Name, baseGroupName, version)
			groupTypeName := getGroupTypeName(pkg.Name, baseGroupName, version)
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

			var clientGroupVars apiGroupsClientVars
			clientGroupVars.GroupTypeName = groupTypeName

			for _, node := range pkg.GetNexusNodes() {
				err := resolveNode(baseImportName, pkg.Name, baseGroupName, version, &groupVars, &clientGroupVars, node)
				if err != nil {
					return clientVars{}, err
				}
				apiGroupClient, err := renderClientApiGroup(clientGroupVars)
				if err != nil {
					return clientVars{}, err
				}
				vars.ApiGroupsClient += apiGroupClient
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

func resolveNode(baseImportName, pkgName, baseGroupName, version string,
	groupVars *apiGroupsVars, clientGroupVars *apiGroupsClientVars, node *ast.TypeSpec) error {

	baseNodeName := node.Name.Name // eg Root
	groupResourceName := getGroupResourceName(baseNodeName)
	groupResourceNameTitle := getGroupResourceNameTitle(baseNodeName)
	groupResourceType := getGroupResourceType(baseNodeName, pkgName, baseGroupName, version)
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
	groupVars.GroupResourcesDefs += "func (obj *" + getGroupTypeName(pkgName, baseGroupName, version) + ") " +
		groupResourceNameTitle + "() *" +
		groupResourceType + " {\n" + "return obj." + groupResourceName + "\n}\n" // eg
	// func (r *RootHelloworldV1) Roots() *rootRootHelloWorld {
	//	return r.roots
	// }

	clientGroupVars.GroupBaseImport = baseImportName + "." + baseNodeName
	clientGroupVars.GroupResourceType = groupResourceType
	clientGroupVars.GroupResourceNameTitle = groupResourceNameTitle
	clientGroupVars.ResolveLinks = ""

	childrenAndLinks := parser.GetChildFields(node)
	childrenAndLinks = append(childrenAndLinks, parser.GetLinkFields(node)...)

	for _, link := range childrenAndLinks {
		linkInfo := getFieldInfo(pkgName, link)
		var vars resolveLinkVars
		vars.LinkFieldName = linkInfo.fieldName
		vars.LinkGroupTypeName = getGroupTypeName(linkInfo.pkgName, baseGroupName, version)
		vars.LinkGroupResourceNameTitle = getGroupResourceNameTitle(linkInfo.fieldType)

		var resolvedLinks string
		var err error
		if parser.IsMapField(link) {
			resolvedLinks, err = renderLinkResolveTemplate(vars, true)
		} else {
			resolvedLinks, err = renderLinkResolveTemplate(vars, false)
		}
		if err != nil {
			return fmt.Errorf("failed to resolve links or children client template for link %v: %v", linkInfo.fieldName, err)
		}
		clientGroupVars.ResolveLinks += resolvedLinks

	}
	return nil
}

type fieldInfo struct {
	pkgName   string
	fieldName string
	fieldType string
}

func getFieldInfo(currentPkgName string, f *ast.Field) fieldInfo {
	var info fieldInfo
	var err error
	info.fieldName, err = parser.GetFieldName(f)
	if err != nil {
		log.Fatalf("Failed to get name of field: %v", err)
	}

	chType := parser.GetFieldType(f)
	s := strings.Split(chType, ".")
	info.fieldType = s[len(s)-1]

	split := strings.Split(chType, ".")
	if len(split) > 1 {
		info.pkgName = split[0]
	} else {
		info.pkgName = currentPkgName
	}

	return info

}

type resolveLinkVars struct {
	LinkFieldName              string
	LinkGroupTypeName          string
	LinkGroupResourceNameTitle string
}

var resolveLinkTmpl = `
	if result.Spec.{{.LinkFieldName}}Gvk.Name != "" {
		field, err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().Get(ctx, result.Spec.{{.LinkFieldName}}Gvk.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.{{.LinkFieldName}} = *field
	}
`

var resolveNamedLinkTmpl = `
	for k, v := range result.Spec.{{.LinkFieldName}}Gvk {
		obj, err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().Get(ctx, v.Name, options)
		if err != nil {
			return nil, err
		}
		result.Spec.{{.LinkFieldName}}[k] = *obj
	}
`

func renderLinkResolveTemplate(vars resolveLinkVars, named bool) (string, error) {
	var templateToUse string
	if named {
		templateToUse = resolveNamedLinkTmpl
	} else {
		templateToUse = resolveLinkTmpl
	}

	tmpl, err := template.New("tmpl").Parse(templateToUse)
	if err != nil {
		return "", err
	}
	ren, err := renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	return ren.String(), nil
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

var apiGroupClientTmpl = `
func (obj *{{.GroupResourceType}}) Get(ctx context.Context, name string, options metav1.GetOptions) (result *{{.GroupBaseImport}}, err error) {
	result, err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	{{.ResolveLinks}}
	return
}
`

type apiGroupsClientVars struct {
	apiGroupsVars
	ResolveLinks           string
	GroupResourceType      string
	GroupResourceNameTitle string
	GroupBaseImport        string
}

func renderClientApiGroup(vars apiGroupsClientVars) (string, error) {
	tmpl, err := template.New("tmpl").Parse(apiGroupClientTmpl)
	if err != nil {
		return "", err
	}
	ren, err := renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	return ren.String(), nil
}
