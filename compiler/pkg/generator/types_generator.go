package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/util"
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
		Name                 string
		StatusType           string
		StatusName           string
		StatusNameFirstLower string
		CrdName              string
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
		status = `{{.StatusName}} {{.StatusType}}` + "`" + `json:"{{.StatusNameFirstLower}},omitempty" yaml:"{{.StatusNameFirstLower}},omitempty"` + "`"
		s.StatusType = parser.GetFieldType(statusField)
		statusName, err := parser.GetFieldName(statusField)
		if err != nil {
			log.Fatalf("failed to determine field name: %v", err)
		}
		if statusName == "" {
			log.Fatalf("name of the user defined status field in nexus node can't be empty")
		}
		s.StatusName = statusName
		s.StatusNameFirstLower = util.GetTag(statusName)
	}

	s.Name = parser.GetTypeName(node)
	if s.Name == "" {
		log.Fatalf("name of type can't be empty")
	}

	s.CrdName = fmt.Sprintf("%s.%s.%s", strings.ToLower(util.ToPlural(s.Name)), strings.ToLower(pkg.Name), config.ConfigInstance.GroupName)

	var crdTemplate = clientgen + "\n" + deepcopygen + "\n" + openapigen + `
type {{.Name}} struct {
	metav1.TypeMeta    ` + "`" + `json:",inline" yaml:",inline"` + "`" + `
	metav1.ObjectMeta  ` + "`" + `json:"metadata" yaml:"metadata"` + "`" + `
	` + spec + `
	Status {{.Name}}NexusStatus` + "`" + `json:"status,omitempty" yaml:"status,omitempty"` + "`" + `
}

` + openapigen + `
type {{.Name}}NexusStatus struct{
	` + status + `
	Nexus NexusStatus` + "`" + `json:"nexus,omitempty" yaml:"nexus,omitempty"` + "`" + `
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
		if err == nil {
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
		var name string
		var err error
		if parser.IsChildField(field) || parser.IsLinkField(field) {
			name, err = parser.GetNodeFieldName(field)
		} else {
			name, err = parser.GetFieldName(field)
		}
		if err != nil {
			log.Fatalf("failed to determine field name: %v", err)
		}
		if comments := getNexusValidationComments(field); comments != nil {
			for _, comment := range comments {
				specDef.Fields += comment + "\n"
			}
		}
		specDef.Fields += "\t" + name + " "
		typeString := ConstructType(aliasNameMap, field)
		// Type is set to "any" for field with annotation "nexus-graphql-jsonencoded"
		if parser.IsFieldAnnotationPresent(field, parser.GRAPHQL_JSONENCODED_ANNOTATION) {
			specDef.Fields += "nexus.NexusGenericObject"
		} else {
			specDef.Fields += typeString
		}
		specDef.Fields += " " + getTag(field, name, false) + "\n"

	}

	for _, child := range parser.GetChildFields(node) {
		name, err := parser.GetNodeFieldName(child)
		if err != nil {
			continue
		}
		gvkName := util.GetGvkFieldName(name)
		if parser.IsNamedChildOrLink(child) {
			specDef.Fields += "\t" + gvkName + " map[string]Child"
		} else {
			specDef.Fields += "\t" + gvkName + " *Child"
		}
		specDef.Fields += " " + getTag(child, util.GetGvkFieldTagName(name), true) + "\n"
	}

	for _, link := range parser.GetLinkFields(node) {
		name, err := parser.GetNodeFieldName(link)
		if err != nil {
			log.Fatalf("failed to GetFieldName: %v", err)
		}
		gvkName := util.GetGvkFieldName(name)
		if parser.IsNamedChildOrLink(link) {
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

func parsePackageStructs(pkg parser.Package, aliasNameMap map[string]string) string {
	var output string
	for _, node := range pkg.GetNodes() {
		output += generateNonNexusTypes(node, aliasNameMap)
	}
	return output
}

func generateNonNexusTypes(node *ast.TypeSpec, aliasNameMap map[string]string) string {
	var crdTemplate = openapigen + `
type {{.Name}} struct {
{{.Fields}}}

`

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
		if comments := getNexusValidationComments(field); comments != nil {
			for _, comment := range comments {
				specDef.Fields += comment + "\n"
			}
		}
		specDef.Fields += "\t" + name + " "
		typeString := ConstructType(aliasNameMap, field)

		currentTags := parser.GetFieldTags(field)
		currentTags = parser.FillEmptyTag(currentTags, name, "json")
		currentTags = parser.FillEmptyTag(currentTags, name, "yaml")

		if currentTags != nil && currentTags.Len() > 0 {
			specDef.Fields += typeString
			specDef.Fields += " " + fmt.Sprintf("`%s`", currentTags.String()) + "\n"
		} else {
			specDef.Fields += typeString + "\n"
		}
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

func parsePackageConsts(pkg parser.Package) string {
	var c struct {
		Consts string
	}
	var count int = 0
	for _, varConst := range pkg.GetConsts() {
		count++
		if count > 1 {
			c.Consts += "\n"
		}
		t, err := pkg.ValueSpecToString(varConst)
		if err != nil {
			log.Fatalf("failed to translate type gen decl to string: %v", err)
		}
		c.Consts += "\t" + t
	}
	var constTemplate string = ""
	switch count {
	case 0:
		return constTemplate
	case 1:
		constTemplate = `const {{.Consts}}`
	default:
		constTemplate = `
const (
{{.Consts}}
)`
	}

	tmpl, err := template.New("tmpl").Parse(constTemplate)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	b, err := renderTemplate(tmpl, c)
	if err != nil {
		log.Fatalf("failed to render template: %v", err)
	}
	return b.String()
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
		iUnquoted, _ := strconv.Unquote(i)
		if strings.HasPrefix(iUnquoted, p.ModPath) {
			if val.Name != nil {
				aliasName, importPath = constructImports(val.Name.String(), i)
				aliasNameMap[val.Name.String()] = aliasName
			} else {
				last := i[strings.LastIndex(i, "/")+1 : len(i)-1]
				aliasName, importPath = constructImports(last, i)
				aliasNameMap[last] = aliasName
			}
		} else {
			aliasName = ""
			importPath = i
			if val.Name != nil {
				aliasName = val.Name.String()
			}
		}
		i = fmt.Sprintf("%s %s", aliasName, importPath)

		// Ignore metav1 external import because it's already present in types.go template
		if strings.HasPrefix(i, "metav1") {
			continue
		}

		importList = append(importList, i)
	}
	return importList
}

func constructImports(inputAlias, inputImportPath string) (string, string) {
	re, err := regexp.Compile(`[\_\.-]`)
	if err != nil {
		log.Fatalf("failed to construct output import path for import path %v : %v", inputImportPath, err)
	}
	aliasName := fmt.Sprintf("%s%sv1", inputAlias, config.ConfigInstance.GroupName)
	aliasName = re.ReplaceAllString(aliasName, "")

	importPath := fmt.Sprintf("\"%sapis/%s.%s/v1\"", config.ConfigInstance.CrdModulePath, re.ReplaceAllString(inputAlias, ""), config.ConfigInstance.GroupName)
	return aliasName, importPath
}

// TODO: https://jira.eng.vmware.com/browse/NPT-296
// Support cross-package imports for the following additional types:
// 1. map[gns.MyStr][]gns.MyStr
// 2. map[string]map[string]gns.MyStr
// 3. []map[string]gns.MyStr
// 4. **gns.MyStr
func ConstructType(aliasNameMap map[string]string, field *ast.Field) string {
	typeString := types.ExprString(field.Type)

	// Check if the field is imported from a different package.
	if !strings.Contains(typeString, ".") {
		return typeString
	}

	switch {
	case parser.IsMapField(field):
		// TODO: Check if the function GetFieldType(field) can be reused for cases other than:
		// map[string]gns.MyStr
		// https://jira.eng.vmware.com/browse/NPT-296
		mapParts := regexp.MustCompile(`^(map\[)`).ReplaceAllString(typeString, "")
		mapStr := regexp.MustCompile(`\]`).Split(mapParts, -1)
		var types []string
		for _, val := range mapStr {
			parts := strings.Split(val, ".")
			if len(parts) > 1 {
				parts = ConstructTypeParts(aliasNameMap, parts)
				val = parts[0] + "." + parts[1]
			}
			types = append(types, val)
		}
		typeString = fmt.Sprintf("map[%s]%s", types[0], types[1])
	case parser.IsArrayField(field):
		arr := regexp.MustCompile(`^(\[])`).ReplaceAllString(typeString, "")
		parts := strings.Split(arr, ".")
		if len(parts) > 1 {
			parts = ConstructTypeParts(aliasNameMap, parts)
			typeString = fmt.Sprintf("[]%s.%s", parts[0], parts[1])
		}
	default:
		parts := strings.Split(typeString, ".")
		if len(parts) > 1 {
			parts = ConstructTypeParts(aliasNameMap, parts)
			typeString = fmt.Sprintf("%s.%s", parts[0], parts[1])
		}
	}
	return typeString
}

func ConstructTypeParts(aliasNameMap map[string]string, parts []string) []string {
	if strings.Contains(parts[0], "*") {
		if val, ok := aliasNameMap[strings.TrimLeft(parts[0], "*")]; ok {
			parts[0] = "*" + val
		}
	} else {
		if val, ok := aliasNameMap[parts[0]]; ok {
			parts[0] = val
		}
	}
	return parts
}

func getNexusValidationComments(field *ast.Field) []string {
	var comments []string
	if field.Doc != nil {
		for _, val := range field.Doc.List {
			if strings.Contains(val.Text, "nexus-validation:") {
				comments = append(comments, val.Text)
			}

		}
	}
	return comments
}

func convertGoStdType(t string) string {
	switch t {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "int"
	case "bool":
		return "bool"
	case "float32", "float64":
		return "float64"
	default:
		return ""
	}
}
