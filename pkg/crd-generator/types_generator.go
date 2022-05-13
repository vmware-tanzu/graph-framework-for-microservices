package crd_generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"
)

const (
	openapigen  string = "// +k8s:openapi-gen=true"
	deepcopygen string = "// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object"
	clientgen   string = "// +genclient\n// +genclient:noStatus\n// +genclient:nonNamespaced"
)

func parsePackageCRDs(pkg parser.Package, aliasNameMap map[string]string) string {
	var output string
	for _, node := range pkg.GetNexusNodes() {
		output += generateType(pkg, node, aliasNameMap)
	}

	return output
}

func generateType(pkg parser.Package, node *ast.TypeSpec, aliasNameMap map[string]string) string {
	var output string
	output += generateCRDStructType(pkg, node)
	output += generateNodeSpec(node, aliasNameMap)
	output += generateListDef(node)

	return output
}

func generateCRDStructType(pkg parser.Package, node *ast.TypeSpec) string {
	var s struct {
		Name       string
		StatusType string
		CrdName    string
	}
	spec := ""
	if len(parser.GetSpecFields(node)) > 0 ||
		len(parser.GetChildFields(node)) > 0 ||
		len(parser.GetLinkFields(node)) > 0 {
		spec = `Spec {{.Name}}Spec ` + "`" + `json:"spec,omitempty" yaml:"spec,omitempty"` + "`"
	}

	status := ""
	statusField := parser.GetStatusField(node)
	if statusField != nil {
		status = `Status {{.StatusType}}` + "`" + `json:"status,omitempty" yaml:"status,omitempty"` + "`"
		s.StatusType = parser.GetFieldType(statusField)
	}

	s.Name = parser.GetTypeName(node)
	if s.Name == "" {
		log.Fatalf("name of type can't be empty")
	}

	s.CrdName = fmt.Sprintf("%s.%s.%s", util.ToPlural(strings.ToLower(s.Name)), strings.ToLower(pkg.Name), config.ConfigInstance.GroupName)

	var crdTemplate = clientgen + "\n" + deepcopygen + "\n" + openapigen + `
type {{.Name}} struct {
	metav1.TypeMeta    ` + "`" + `json:",inline" yaml:",inline"` + "`" + `
	metav1.ObjectMeta  ` + "`" + `json:"metadata" yaml:"metadata"` + "`" + `
	` + spec + `
	` + status + `
}

func (c *{{.Name}}) CRDName() string {
	return "{{.CrdName}}"
}

func (c *{{.Name}}) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

`

	tmpl, err := template.New("tmpl").Parse(crdTemplate)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	b, err := renderTemplate(tmpl, s)
	if err != nil {
		log.Fatalf("failed to render template: %v", err)
	}
	return b.String()
}

func getTag(f *ast.Field, name string, omitempty bool) string {
	n := util.GetTag(name)
	tag := "json:\"" + n + "\" yaml:\"" + n + "\""
	if omitempty {
		tag = "json:\"" + n + ",omitempty\" yaml:\"" + n + ",omitempty\""
	}

	currentTags := parser.GetFieldTags(f)
	if currentTags != nil && currentTags.Len() > 0 {
		nexusTag, err := currentTags.Get("nexus")
		if currentTags.Len() == 1 && err == nil {
			tag += " " + nexusTag.String()
		} else {
			tag = currentTags.String()
		}
	}

	return fmt.Sprintf("`%s`", tag)
}

func generateNodeSpec(node *ast.TypeSpec, aliasNameMap map[string]string) string {
	var crdTemplate = openapigen + `
type {{.Name}}Spec struct {
{{.Fields}}}

`
	if len(parser.GetSpecFields(node)) == 0 &&
		len(parser.GetChildFields(node)) == 0 &&
		len(parser.GetLinkFields(node)) == 0 {
		return ""
	}
	var specDef struct {
		Name   string
		Fields string
	}
	specDef.Name = parser.GetTypeName(node)

	for _, field := range parser.GetSpecFields(node) {
		name, err := parser.GetFieldName(field)
		if err != nil {
			log.Fatalf("failed to GetFieldName: %v", err)
		}
		specDef.Fields += "\t" + name + " "
		typeString := types.ExprString(field.Type)
		parts := strings.Split(typeString, ".")
		if len(parts) > 1 {
			if val, ok := aliasNameMap[parts[0]]; ok {
				parts[0] = val
			}
			typeString = parts[0] + "." + parts[1]
		}
		specDef.Fields += typeString
		specDef.Fields += " " + getTag(field, name, false) + "\n"
	}

	for _, child := range parser.GetChildFields(node) {
		name, err := parser.GetFieldName(child)
		if err != nil {
			continue
		}
		gvkName := util.GetGvkFieldName(name)
		if parser.IsMapField(child) {
			specDef.Fields += "\t" + gvkName + " map[string]Child"
		} else {
			specDef.Fields += "\t" + gvkName + " *Child"
		}
		specDef.Fields += " " + getTag(child, util.GetGvkFieldTagName(name), true) + "\n"
	}

	for _, link := range parser.GetLinkFields(node) {
		name, err := parser.GetFieldName(link)
		if err != nil {
			log.Fatalf("failed to GetFieldName: %v", err)
		}
		gvkName := util.GetGvkFieldName(name)
		if parser.IsMapField(link) {
			specDef.Fields += "\t" + gvkName + " map[string]Link"
		} else {
			specDef.Fields += "\t" + gvkName + " *Link"
		}
		specDef.Fields += " " + getTag(link, util.GetGvkFieldTagName(name), true) + "\n"
	}

	tmpl, err := template.New("tmpl").Parse(crdTemplate)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	b, err := renderTemplate(tmpl, specDef)
	if err != nil {
		log.Fatalf("failed to render template: %v", err)
	}
	return b.String()

}

func generateListDef(node *ast.TypeSpec) string {
	var listTemplate = deepcopygen + `
type {{.Name}}List struct {
	metav1.TypeMeta   ` + "`" + `json:",inline" yaml:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata" yaml:"metadata"` + "`" + `
	Items []{{.Name}} ` + "`" + `json:"items" yaml:"items"` + "`" + `
}

`
	var s struct {
		Name string
	}
	s.Name = parser.GetTypeName(node)
	if s.Name == "" {
		log.Fatalf("name of type can't be empty")
	}
	tmpl, err := template.New("tmpl").Parse(listTemplate)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	b, err := renderTemplate(tmpl, s)
	if err != nil {
		log.Fatalf("failed to render template: %v", err)
	}
	return b.String()
}

func parsePackageStructs(pkg parser.Package) string {
	var output string
	for _, node := range pkg.GetNodes() {
		t, err := pkg.TypeSpecToString(node)
		if err != nil {
			log.Fatalf("failed to translate type spec to string: %v", err)
		}
		output += openapigen + "\n" + "type " + t + "\n\n"
	}

	return output
}

func parsePackageTypes(pkg parser.Package) string {
	var output string
	for _, node := range pkg.GetTypes() {
		t, err := pkg.GenDeclToString(&node)
		if err != nil {
			log.Fatalf("failed to translate type gen decl to string: %v", err)
		}
		output += t + "\n"
	}

	return output
}

func parsePackageImports(pkg parser.Package, aliasNameMap map[string]string) string {
	var output string
	for _, imp := range GenerateImports(&pkg, aliasNameMap) {
		output += imp + "\n"
	}
	return output
}

func GenerateImports(p *parser.Package, aliasNameMap map[string]string) []string {
	var (
		importList            []string
		aliasName, importPath string
	)
	imports := p.GetImports()
	for _, val := range imports {
		i := val.Path.Value
		if strings.Contains(i, p.ModPath) {
			if val.Name != nil {
				aliasName, importPath = constructImports(val.Name.String(), i)
				aliasNameMap[val.Name.String()] = aliasName
			} else {
				last := i[strings.LastIndex(i, "/")+1 : len(i)-1]
				aliasName, importPath = constructImports(last, i)
				aliasNameMap[last] = aliasName
			}
		} else {
			if val.Name != nil {
				importPath = i
				aliasName = val.Name.String()
			}
		}
		i = fmt.Sprintf("%s %s", aliasName, importPath)
		importList = append(importList, i)
	}
	return importList
}

func constructImports(inputAlias, inputImportPath string) (string, string) {
	re, err := regexp.Compile(`[\_\.]`)
	if err != nil {
		log.Fatalf("failed to construct output import path for import path %v : %v", inputImportPath, err)
	}
	aliasName := fmt.Sprintf("%s%sv1", inputAlias, config.ConfigInstance.GroupName)
	aliasName = re.ReplaceAllString(aliasName, "")

	importPath := fmt.Sprintf("\"%sapis/%s.%s/v1\"", config.ConfigInstance.CrdModulePath, re.ReplaceAllString(inputAlias, ""), config.ConfigInstance.GroupName)
	return aliasName, importPath
}
