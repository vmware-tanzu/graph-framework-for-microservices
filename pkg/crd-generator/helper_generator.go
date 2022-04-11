package crd_generator

import (
	"fmt"
	"strings"
	"text/template"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"

	log "github.com/sirupsen/logrus"
)

func generateGetCrdParentsMap(parentsMap map[string]parser.NodeHelper) string {
	var output string
	output += "map[string][]string{\n"
	for k, v := range parentsMap {
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

func generateGetObjectByCRDName(parentsMap map[string]parser.NodeHelper) string {
	var output string
	var ifTemplate = `if crdName == "{{.CrdName}}" {
		obj, err := dmClient.{{.Method}}().{{.Plural}}().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}`

	for k, v := range parentsMap {
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

		parts := strings.Split(k, ".")
		s.Method = fmt.Sprintf("%s%sV1", strings.Title(parts[1]), strings.Title(parts[2]))
		s.Plural = util.ToPlural(v.Name)

		b, err := renderTemplate(tmpl, s)
		if err != nil {
			log.Fatalf("failed to render template: %v", err)
		}
		output += fmt.Sprintf("%s\n", b.String())
	}

	return output
}
