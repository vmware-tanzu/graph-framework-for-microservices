package crd_generator

import (
	"sort"
	"strings"
	"text/template"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"
	"k8s.io/gengo/namer"
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

			groupName := pkg.Name + "." + baseGroupName                    // eg root.helloworld.com
			importPath := groupName + "/" + version                        // eg root.helloworld.com/v1
			baseImportName := `base` + util.RemoveSpecialChars(importPath) // eg baseroothelloworldv1
			vars.BaseImports += baseImportName + ` "` + crdModulePath + "apis/" + importPath + `"` +
				"\n" // eg baseroothelloworldv1 "helloworld/nexus/generated/apis/root.helloworld.com/v1"

			groupGoName := namer.IC(strings.Split(baseGroupName, ".")[0])
			//baseGroupNameNoSpecials := util.RemoveSpecialChars(baseGroupName)
			groupVarName := pkg.Name + groupGoName +
				strings.Title(version) // eg rootHelloworldV1
			groupTypeName := strings.Title(pkg.Name) + groupGoName +
				strings.Title(version) // eg RootHelloworldV1
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
				baseNodeName := node.Name.Name                                       // eg Root
				groupResourceName := strings.ToLower(util.ToPlural(baseNodeName))    // eg roots
				groupResourceNameTitle := strings.Title(util.ToPlural(baseNodeName)) // eg Roots
				groupResourceType := strings.ToLower(baseNodeName) + groupTypeName   // eg rootRootHelloworld
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
				groupVars.GroupResourcesDefs += "func (obj *" + groupTypeName + ") " + groupResourceNameTitle + "() *" +
					groupResourceType + " {\n" + "return obj." + groupResourceName + "\n}\n" // eg
				// func (r *RootHelloworldV1) Roots() *rootRootHelloWorld {
				//	return r.roots
				// }

				clientGroupVars.GroupBaseImport = baseImportName + "." + baseNodeName
				clientGroupVars.GroupResourceType = groupResourceType
				clientGroupVars.GroupResourceNameTitle = groupResourceNameTitle

			}

			apiGroup, err := renderApiGroup(groupVars)
			if err != nil {
				return clientVars{}, err
			}
			vars.ApiGroups += apiGroup

			apiGroupClient, err := renderClientApiGroup(clientGroupVars)
			if err != nil {
				return clientVars{}, err
			}
			vars.ApiGroupsClient += apiGroupClient
		}
	}
	return vars, nil
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

	// resolve children
    // TODO
	// resolve links
    // TODO
	return
}
`

type apiGroupsClientVars struct {
	apiGroupsVars
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
