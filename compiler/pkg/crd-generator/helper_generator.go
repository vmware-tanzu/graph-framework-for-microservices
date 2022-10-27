package crd_generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	log "github.com/sirupsen/logrus"
)

func generateGetCrdParentsMap(keys []string, parentsMap map[string]parser.NodeHelper) string {
	var output string
	output += "map[string][]string{\n"
	for _, k := range keys {
		v := parentsMap[k]
		kv := fmt.Sprintf(`	"%s": {`, k)
		for _, p := range v.Parents {
			kv += fmt.Sprintf(`"%s",`, p)
		}
		kv += "},\n"
		output += kv
	}
	output += "}"
	return output
}

func generateGetObjectByCRDName(keys []string, parentsMap map[string]parser.NodeHelper) string {
	var output string
	var ifTemplate = `if crdName == "{{.CrdName}}" {
		obj, err := dmClient.{{.Method}}().{{.Plural}}().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}`

	for _, k := range keys {
		v := parentsMap[k]
		var s struct {
			CrdName string
			Method  string
			Plural  string
		}

		tmpl, err := template.New("tmpl").Parse(ifTemplate)
		if err != nil {
			log.Fatalf("failed to parse template: %v", err)
		}

		s.CrdName = k
		s.Method = getMethodName(k)
		s.Plural = util.ToPlural(v.Name)

		b, err := renderTemplate(tmpl, s)
		if err != nil {
			log.Fatalf("failed to render template: %v", err)
		}
		output += fmt.Sprintf("%s\n", b.String())
	}

	return output
}

func getMethodName(crdName string) string {
	crdName = strings.ReplaceAll(crdName, "-", ".")
	parts := strings.Split(crdName, ".")
	return fmt.Sprintf("%s%sV1", cases.Title(language.Und, cases.NoLower).String(parts[1]), cases.Title(language.Und, cases.NoLower).String(parts[2]))
}
