package crd_generator

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
	"text/template"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"

	log "github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

func generateNexusClientVars(baseGroupName, crdModulePath string, pkgs parser.Packages) (clientVars, error) {
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
				err := resolveNode(baseImportName, pkg, baseGroupName, version, &groupVars, &clientGroupVars, node)
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

func resolveNode(baseImportName string, pkg parser.Package, baseGroupName, version string,
	groupVars *apiGroupsVars, clientGroupVars *apiGroupsClientVars, node *ast.TypeSpec) error {

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

	clientGroupVars.GroupBaseImport = baseImportName + "." + baseNodeName
	clientGroupVars.GroupResourceType = groupResourceType
	clientGroupVars.GroupResourceNameTitle = groupResourceNameTitle
	clientGroupVars.ResolveLinksGet = ""
	clientGroupVars.ResolveLinksDelete = ""

	// TODO support resolution of links which are not nexus nodes https://jira.eng.vmware.com/browse/NPT-112
	children := parser.GetChildFields(node)
	childrenAndLinks := children
	childrenAndLinks = append(childrenAndLinks, parser.GetLinkFields(node)...)

	if len(children) > 0 {
		clientGroupVars.HasChildren = true
	}
	for _, link := range childrenAndLinks {
		linkInfo := getFieldInfo(pkg, link)
		var vars resolveLinkVars
		vars.LinkFieldName = linkInfo.fieldName
		vars.LinkGroupTypeName = util.GetGroupTypeName(linkInfo.pkgName, baseGroupName, version)
		vars.LinkGroupResourceNameTitle = util.GetGroupResourceNameTitle(linkInfo.fieldType)

		var resolvedLinksGet, resolvedLinksDelete string
		var err error
		if parser.IsMapField(link) {
			resolvedLinksGet, err = renderLinkResolveTemplate(vars, resolveNamedLinkGetTmpl)
			if err != nil {
				return fmt.Errorf("failed to resolve links or children get client template for link %v: %v",
					linkInfo.fieldName, err)
			}
			if !parser.IsLinkField(link) { // do not resolve softlinks for delete
				resolvedLinksDelete, err = renderLinkResolveTemplate(vars, resolveNamedLinkDeleteTmpl)
				if err != nil {
					return fmt.Errorf("failed to resolve links or children delete client template for link %v: %v",
						linkInfo.fieldName, err)
				}
			}
		} else {
			resolvedLinksGet, err = renderLinkResolveTemplate(vars, resolveLinkGetTmpl)
			if err != nil {
				return fmt.Errorf("failed to resolve links or children get client template for link %v: %v",
					linkInfo.fieldName, err)
			}
			if !parser.IsLinkField(link) { // do not resolve softlinks for delete
				resolvedLinksDelete, err = renderLinkResolveTemplate(vars, resolveLinkDeleteTmpl)
				if err != nil {
					return fmt.Errorf("failed to resolve links or children delete client template for link %v: %v",
						linkInfo.fieldName, err)
				}
			}
		}

		clientGroupVars.ResolveLinksGet += resolvedLinksGet
		clientGroupVars.ResolveLinksDelete += resolvedLinksDelete
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

var resolveLinkGetTmpl = `
	if result.Spec.{{.LinkFieldName}}Gvk != nil {
		field, err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().GetByName(ctx, result.Spec.{{.LinkFieldName}}Gvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.{{.LinkFieldName}} = field
	}
`

var resolveNamedLinkGetTmpl = `
	for k, v := range result.Spec.{{.LinkFieldName}}Gvk {
		obj, err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().GetByName(ctx, v.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.{{.LinkFieldName}}[k] = *obj
	}
`

var resolveLinkDeleteTmpl = `
	if result.Spec.{{.LinkFieldName}}Gvk != nil {
		 err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().DeleteByName(ctx, result.Spec.{{.LinkFieldName}}Gvk.Name)
		if err != nil {
			return err
		}
	}
`

var resolveNamedLinkDeleteTmpl = `
	for _, v := range result.Spec.{{.LinkFieldName}}Gvk {
		err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().DeleteByName(ctx, v.Name)
		if err != nil {
			return err
		}
	}
`

func renderLinkResolveTemplate(vars resolveLinkVars, templateToUse string) (string, error) {
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
func (obj *{{.GroupResourceType}}) Get(ctx context.Context, name string, labels map[string]string) (result *{{.GroupBaseImport}}, err error) {
	hashedName := helper.GetHashedName(name, labels)
	result, err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Get(ctx, hashedName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	{{.ResolveLinksGet}}
	return
}

func (obj *{{.GroupResourceType}}) GetByName(ctx context.Context, name string) (result *{{.GroupBaseImport}}, err error) { 
	result, err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	{{.ResolveLinksGet}}
	return
}

func (obj *{{.GroupResourceType}}) Delete(ctx context.Context, name string, labels map[string]string) (err error) {
	hashedName := helper.GetHashedName(name, labels)
	{{if .HasChildren}}{{.GetForDelete}}
	{{ end }}

	{{.ResolveLinksDelete}}

	err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Delete(ctx, hashedName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return
}

func (obj *{{.GroupResourceType}}) DeleteByName(ctx context.Context, name string) (err error) { 
	{{if .HasChildren}}
{{.GetForDeleteByName}}
	{{ end }}

	{{.ResolveLinksDelete}}

	err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return
}
`

var getForDeleteTmpl = `
	result, err := obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Get(ctx, hashedName, metav1.GetOptions{})
	if err != nil {
		return err
	}
`

var getByNameForDeleteTmpl = `
	result, err := obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
`

type apiGroupsClientVars struct {
	apiGroupsVars
	ResolveLinksGet        string
	ResolveLinksDelete     string
	GetForDeleteByName     string
	GetForDelete           string
	HasChildren            bool
	GroupResourceType      string
	GroupResourceNameTitle string
	GroupBaseImport        string
}

func renderClientApiGroup(vars apiGroupsClientVars) (string, error) {
	tmpl, err := template.New("tmpl").Parse(getForDeleteTmpl)
	if err != nil {
		return "", err
	}
	getBase, err := renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	vars.GetForDelete = getBase.String()
	tmpl, err = template.New("tmpl").Parse(getByNameForDeleteTmpl)
	if err != nil {
		return "", err
	}
	getByNameBase, err := renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	vars.GetForDeleteByName = getByNameBase.String()

	tmpl, err = template.New("tmpl").Parse(apiGroupClientTmpl)
	if err != nil {
		return "", err
	}
	ren, err := renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	return ren.String(), nil
}
