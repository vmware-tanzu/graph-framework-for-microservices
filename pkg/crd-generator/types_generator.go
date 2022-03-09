package crd_generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"text/template"
	"unicode"

	log "github.com/sirupsen/logrus"

	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
)

const (
	openapigen  string = "// +k8s:openapi-gen=true"
	deepcopygen string = "// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object"
	clientgen   string = "// +genclient\n// +genclient:noStatus"
)

func parsePackageCRDs(pkg parser.Package) string {
	var output string
	for _, node := range pkg.GetNexusNodes() {
		output += generateType(node)
	}

	return output
}

func generateType(node *ast.TypeSpec) string {
	var output string
	output += generateCRDStructType(node)
	output += generateNodeSpec(node)
	output += generateListDef(node)

	return output
}

func generateCRDStructType(node *ast.TypeSpec) string {
	var s struct {
		Name       string
		StatusType string
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

	var crdTemplate = clientgen + "\n" + deepcopygen + "\n" + openapigen + `
type {{.Name}} struct {
	metav1.TypeMeta    ` + "`" + `json:",inline" yaml:",inline"` + "`" + `
	metav1.ObjectMeta  ` + "`" + `json:"metadata" yaml:"metadata"` + "`" + `
	` + spec + `
	` + status + `
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
	n := string(unicode.ToLower(rune(name[0]))) + name[1:]
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

func generateNodeSpec(node *ast.TypeSpec) string {
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
		specDef.Fields += typeString
		specDef.Fields += " " + getTag(field, name, false) + "\n"
	}

	for _, child := range parser.GetChildFields(node) {
		name, err := parser.GetFieldName(child)
		if err != nil {
			continue
		}
		if parser.IsMapField(child) {
			specDef.Fields += "\t" + name + " map[string]Child"
		} else {
			specDef.Fields += "\t" + name + " Child"
		}
		specDef.Fields += " " + getTag(child, name, true) + "\n"
	}

	for _, link := range parser.GetLinkFields(node) {
		name, err := parser.GetFieldName(link)
		if err != nil {
			log.Fatalf("failed to GetFieldName: %v", err)
		}
		if parser.IsMapField(link) {
			specDef.Fields += "\t" + name + " map[string]Link"
		} else {
			specDef.Fields += "\t" + name + " Link"
		}
		specDef.Fields += " " + getTag(link, name, true) + "\n"
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

func parsePackageImports(pkg parser.Package) string {
	var output string
	for _, imp := range pkg.GetImportStrings() {
		output += imp + "\n"
	}
	return output
}
