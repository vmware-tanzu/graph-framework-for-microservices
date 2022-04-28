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
		var vars resolveLinkVars
		vars.LinkFieldName = linkInfo.fieldName
		vars.LinkFieldType = linkInfo.fieldType
		vars.LinkFieldNameTag = strings.ToLower(linkInfo.fieldName)
		vars.LinkGroupTypeName = util.GetGroupTypeName(linkInfo.pkgName, baseGroupName, version)
		vars.LinkGroupResourceNameTitle = util.GetGroupResourceNameTitle(linkInfo.fieldType)
		vars.LinkBaseImport = util.GetBaseImportName(linkInfo.pkgName, baseGroupName, version)

		var resolvedLinksGet, resolvedLinksDelete, resolvedLinksCreate string
		var err error
		if parser.IsMapField(link) {
			resolvedLinksGet, err = renderLinkResolveTemplate(vars, resolveNamedLinkGetTmpl)
			if err != nil {
				return fmt.Errorf("failed to resolve links or children get client template for link %v: %v",
					linkInfo.fieldName, err)
			}
			if !parser.IsLinkField(link) { // do not resolve softlinks for delete/create
				resolvedLinksDelete, err = renderLinkResolveTemplate(vars, resolveNamedLinkDeleteTmpl)
				if err != nil {
					return fmt.Errorf("failed to resolve links or children delete client template for link %v: %v",
						linkInfo.fieldName, err)
				}
				resolvedLinksCreate, err = renderLinkResolveTemplate(vars, resolveLinkCreateTmpl)
				if err != nil {
					return fmt.Errorf("failed to resolve links or children create client template for link %v: %v",
						linkInfo.fieldName, err)
				}
			} else {
				clientVarsLink := apiGroupsClientVarsLink{
					FieldName:       linkInfo.fieldName,
					FieldNameGvk:    util.GetGvkFieldTagName(linkInfo.fieldName),
					Group:           util.GetGroupName(linkInfo.pkgName, baseGroupName),
					Kind:            linkInfo.fieldType,
					GroupBaseImport: util.GetBaseImportName(linkInfo.pkgName, baseGroupName, version) + "." + linkInfo.fieldType,
					IsNamed:         true,
				}
				clientGroupVars.Links = append(clientGroupVars.Links, clientVarsLink)
			}
		} else {
			resolvedLinksGet, err = renderLinkResolveTemplate(vars, resolveLinkGetTmpl)
			if err != nil {
				return fmt.Errorf("failed to resolve links or children get client template for link %v: %v",
					linkInfo.fieldName, err)
			}

			if !parser.IsLinkField(link) { // do not resolve softlinks for delete/create
				resolvedLinksDelete, err = renderLinkResolveTemplate(vars, resolveLinkDeleteTmpl)
				if err != nil {
					return fmt.Errorf("failed to resolve links or children delete client template for link %v: %v",
						linkInfo.fieldName, err)
				}
				resolvedLinksCreate, err = renderLinkResolveTemplate(vars, resolveLinkCreateTmpl)
				if err != nil {
					return fmt.Errorf("failed to resolve links or children create client template for link %v: %v",
						linkInfo.fieldName, err)
				}
			} else {
				clientVarsLink := apiGroupsClientVarsLink{
					FieldName:       linkInfo.fieldName,
					FieldNameGvk:    util.GetGvkFieldTagName(linkInfo.fieldName),
					Group:           util.GetGroupName(linkInfo.pkgName, baseGroupName),
					Kind:            linkInfo.fieldType,
					GroupBaseImport: util.GetBaseImportName(linkInfo.pkgName, baseGroupName, version) + "." + linkInfo.fieldType,
					IsNamed:         false,
				}
				clientGroupVars.Links = append(clientGroupVars.Links, clientVarsLink)
			}
		}

		clientGroupVars.ResolveLinksGet += resolvedLinksGet
		clientGroupVars.ResolveLinksDelete += resolvedLinksDelete
		clientGroupVars.ResolveLinksCreate += resolvedLinksCreate
	}

	for _, f := range parser.GetSpecFields(node) {
		fieldInfo := getFieldInfo(pkg, f)
		var vars resolveLinkVars
		vars.LinkFieldName = fieldInfo.fieldName
		vars.LinkFieldNameTag = util.GetTag(fieldInfo.fieldName)
		resolvedPatches, err := renderLinkResolveTemplate(vars, patchForUpdateTmpl)
		if err != nil {
			return fmt.Errorf("failed to parse ")
		}
		clientGroupVars.ForUpdatePatches += resolvedPatches
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

type resolveLinkVars struct {
	LinkFieldName              string
	LinkFieldType              string
	LinkFieldNameTag           string
	LinkGroupTypeName          string
	LinkGroupResourceNameTitle string
	LinkBaseImport             string
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
	result.Spec.{{.LinkFieldName}} = make(map[string]{{.LinkBaseImport}}.{{.LinkFieldType}}, len(result.Spec.{{.LinkFieldName}}Gvk))
	for _, v := range result.Spec.{{.LinkFieldName}}Gvk {
		field, err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().GetByName(ctx, v.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.{{.LinkFieldName}}[field.GetLabels()["nexus/display_name"]] = *field
	}
`

var resolveLinkDeleteTmpl = `
	if result.Spec.{{.LinkFieldName}}Gvk != nil {
		 err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().DeleteByName(ctx, result.Spec.{{.LinkFieldName}}Gvk.Name, parents)
		if err != nil {
			return err
		}
	}
`

var resolveNamedLinkDeleteTmpl = `
	for _, v := range result.Spec.{{.LinkFieldName}}Gvk {
		err := obj.client.{{.LinkGroupTypeName}}().{{.LinkGroupResourceNameTitle}}().DeleteByName(ctx, v.Name, parents)
		if err != nil {
			return err
		}
	}
`

var resolveLinkCreateTmpl = `
	objToCreate.Spec.{{.LinkFieldName}} = nil
	objToCreate.Spec.{{.LinkFieldName}}Gvk = nil
`

var patchForUpdateTmpl = `
	patchValue{{.LinkFieldName}} := objToUpdate.Spec.{{.LinkFieldName}}
	patchOp{{.LinkFieldName}} := PatchOp{
		Op:    "replace",
		Path:  "/spec/{{.LinkFieldNameTag}}",
		Value: patchValue{{.LinkFieldName}},
	}
	patch = append(patch, patchOp{{.LinkFieldName}})
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
// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *{{.GroupResourceType}}) Get(ctx context.Context, name string, parents map[string]string) (result *{{.GroupBaseImport}}, err error) {
	hashedName := helper.GetHashedName("{{.CrdName}}", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name 
func (obj *{{.GroupResourceType}}) GetByName(ctx context.Context, name string) (result *{{.GroupBaseImport}}, err error) { 
	result, err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return obj.resolveLinks(ctx, result)
}

func (obj *{{.GroupResourceType}}) resolveLinks(ctx context.Context, raw *{{.GroupBaseImport}}) (result *{{.GroupBaseImport}}, err error) {
	result = raw
	{{.ResolveLinksGet}}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *{{.GroupResourceType}}) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents["nexus/is_name_hashed"] = "true"
	hashedName := helper.GetHashedName("{{.CrdName}}", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *{{.GroupResourceType}}) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) { 
	{{if .HasChildren}}
{{.GetForDeleteByName}}
	{{ end }}

	{{.ResolveLinksDelete}}

	err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	{{if .Parent.HasParent}}
{{.Parent.UpdateParentForDelete}}
	{{ end }}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *{{.GroupResourceType}}) Create(ctx context.Context, objToCreate *{{.GroupBaseImport}}, parents map[string]string) (result *{{.GroupBaseImport}}, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels["nexus/is_name_hashed"] != "true" {
		objToCreate.Labels["nexus/display_name"] = objToCreate.GetName()
		objToCreate.Labels["nexus/is_name_hashed"] = "true"
		hashedName := helper.GetHashedName("{{.CrdName}}", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *{{.GroupResourceType}}) CreateByName(ctx context.Context, objToCreate *{{.GroupBaseImport}}, parents map[string]string) (result *{{.GroupBaseImport}}, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels["nexus/display_name"]; !ok {
		objToCreate.Labels["nexus/display_name"] = objToCreate.GetName()
	}

	{{.ResolveLinksCreate}}

	result, err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	{{if .Parent.HasParent}}
{{.Parent.UpdateParentForCreate}}
	{{ end }}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *{{.GroupResourceType}}) Update(ctx context.Context, objToUpdate *{{.GroupBaseImport}}, parents map[string]string) (result *{{.GroupBaseImport}}, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels["nexus/is_name_hashed"] != "true" {
		objToUpdate.Labels["nexus/display_name"] = objToUpdate.GetName()
		objToUpdate.Labels["nexus/is_name_hashed"] = "true"
		hashedName := helper.GetHashedName("{{.CrdName}}", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *{{.GroupResourceType}}) UpdateByName(ctx context.Context, objToUpdate *{{.GroupBaseImport}}) (result *{{.GroupBaseImport}}, err error) {
	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)
	{{.ForUpdatePatches}}
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.{{.GroupTypeName}}().{{.GroupResourceNameTitle}}().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return obj.resolveLinks(ctx, result)
}

{{ range $key, $link := .Links }}
// Add{{$link.FieldName}} updates srcObj with linkToAdd object
func (obj *{{$.GroupResourceType}}) Add{{$link.FieldName}}(ctx context.Context, srcObj *{{$.GroupBaseImport}}, linkToAdd *{{$link.GroupBaseImport}}) (result *{{$.GroupBaseImport}}, err error) {
	{{ if $link.IsNamed }}
	payload := "{\"spec\": {\"{{$link.FieldNameGvk}}\": {\"" + linkToAdd.Name + "\": {\"name\": \"" + linkToAdd.Name + "\",\"kind\": \"{{$link.Kind}}\", \"group\": \"{{$link.Group}}\"}}}}"
	result, err = obj.client.baseClient.{{$.GroupTypeName}}().{{$.GroupResourceNameTitle}}().Patch(ctx, srcObj.Name, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	{{ else }}
	var patch Patch
	patchOp := PatchOp{
		Op:   "replace",
		Path: "/spec/{{$link.FieldNameGvk}}",
		Value: {{$.BaseImportName}}.Child{
			Group: "{{$link.Group}}",
			Kind:  "{{$link.Kind}}",
			Name:  linkToAdd.Name,
		},
	}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.{{$.GroupTypeName}}().{{$.GroupResourceNameTitle}}().Patch(ctx, srcObj.Name, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	{{ end }}

	return obj.resolveLinks(ctx, result)
}

// Remove{{$link.FieldName}} removes linkToRemove object from srcObj
func (obj *{{$.GroupResourceType}}) Remove{{$link.FieldName}}(ctx context.Context, srcObj *{{$.GroupBaseImport}}, linkToRemove *{{$link.GroupBaseImport}}) (result *{{$.GroupBaseImport}}, err error) {
	var patch Patch
	{{if $link.IsNamed}}
	patchOp := PatchOp{
		Op:    "remove",
		Path:  "/spec/{{$link.FieldNameGvk}}/" + linkToRemove.Name,
	}
	{{ else }}
	patchOp := PatchOp{
		Op:    "remove",
		Path:  "/spec/{{$link.FieldNameGvk}}",
	}
	{{ end }}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.{{$.GroupTypeName}}().{{$.GroupResourceNameTitle}}().Patch(ctx, srcObj.Name, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return obj.resolveLinks(ctx, result)
}
{{ end }}
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
	if parents == nil {
		parents = make(map[string]string, 1)
	}

	if _, ok := result.GetLabels()["nexus/display_name"]; ok {
		parents["{{.CrdName}}"] = result.GetLabels()["nexus/display_name"]
	} else {
		parents["{{.CrdName}}"] = name
	}
`

var updateParentForCreate = `
	parentName, ok := parents["{{.Parent.CrdName}}"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if objToCreate.Labels["nexus/is_name_hashed"] == "true" {
		parentName = helper.GetHashedName("{{.Parent.CrdName}}", parents, parentName)
	}
	{{if .Parent.IsNamed}}
	payload := "{\"spec\": {\"{{.Parent.GvkFieldName}}\": {\"" + objToCreate.Name + "\": {\"name\": \"" + objToCreate.Name + "\",\"kind\": \"{{.Kind}}\", \"group\": \"{{.Group}}\"}}}}"
	_, err = obj.client.baseClient.{{.Parent.GroupTypeName}}().{{.Parent.GroupResourceNameTitle}}().Patch(ctx, parentName, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	{{ else }}
	var patch Patch
	patchOp := PatchOp{
		Op:    "replace",
		Path:  "/spec/{{.Parent.GvkFieldName}}",
		Value: {{.BaseImportName}}.Child{
			Group: "{{.Group}}",
			Kind:  "{{.Kind}}",
			Name:  objToCreate.Name,
		},
	}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	_, err = obj.client.baseClient.{{.Parent.GroupTypeName}}().{{.Parent.GroupResourceNameTitle}}().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	{{ end }}
`

var updateParentForDelete = `
	var patch Patch
	{{if .Parent.IsNamed}}
	patchOp := PatchOp{
		Op:    "remove",
		Path:  "/spec/{{.Parent.GvkFieldName}}/" + name,
	}
	{{ else }}
	patchOp := PatchOp{
		Op:    "remove",
		Path:  "/spec/{{.Parent.GvkFieldName}}",
	}
	{{ end }}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return err
	}
	parentName, ok := parents["{{.Parent.CrdName}}"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if parents["nexus/is_name_hashed"] == "true" {
		parentName = helper.GetHashedName("{{.Parent.CrdName}}", parents, parentName)
	}
	_, err = obj.client.baseClient.{{.Parent.GroupTypeName}}().{{.Parent.GroupResourceNameTitle}}().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return err
	}
`

type apiGroupsClientVars struct {
	apiGroupsVars
	CrdName                string
	ResolveLinksGet        string
	ResolveLinksDelete     string
	ResolveLinksCreate     string
	GetForDeleteByName     string
	GetForDelete           string
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
		UpdateParentForCreate  string
		UpdateParentForDelete  string
		GroupTypeName          string
		GroupResourceNameTitle string
	}
	ForUpdatePatches string

	Links []apiGroupsClientVarsLink
}

type apiGroupsClientVarsLink struct {
	FieldName       string
	FieldNameGvk    string
	Group           string
	Kind            string
	GroupBaseImport string
	IsNamed         bool
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
	// Parent
	tmpl, err = template.New("tmpl").Parse(updateParentForCreate)
	if err != nil {
		return "", err
	}
	updateParentBase, err := renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	vars.Parent.UpdateParentForCreate = updateParentBase.String()
	tmpl, err = template.New("tmpl").Parse(updateParentForDelete)
	if err != nil {
		return "", err
	}
	updateParentBase, err = renderTemplate(tmpl, vars)
	if err != nil {
		return "", err
	}
	vars.Parent.UpdateParentForDelete = updateParentBase.String()

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
