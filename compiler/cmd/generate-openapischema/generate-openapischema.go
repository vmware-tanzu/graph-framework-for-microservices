package main

import (
	"flag"
	"fmt"
	"strings"

	generator "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/openapi_generator/openapi"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kube-openapi.git/pkg/validation/spec"
)

func main() {
	var yamlsPath string
	flag.StringVar(&yamlsPath, "yamls-path", "", "Path to directory containing CRD YAML definitions")
	flag.Parse()
	if yamlsPath == "" {
		panic("yamls-path is empty. Run with -h for help")
	}

	ref := func(pkg string) spec.Ref {
		r, err := spec.NewRef(strings.ToLower(pkg))
		if err != nil {
			panic(err)
		}
		return r
	}
	g, err := generator.NewGenerator(openapi.GetOpenAPIDefinitions(ref))
	if err != nil {
		panic(fmt.Sprintf("Failed creating Generator: %v", err))
	}
	err = g.ResolveRefs()
	if err != nil {
		panic(err)
	}
	if len(g.MissingDefinitions()) > 0 {
		for pkg := range g.MissingDefinitions() {
			fmt.Printf("\n***\nMissing schema for %q\n***\n", pkg)
		}
		readmePath := "https://github.com/vmware-tanzu/graph-framework-for-microservices/compiler/blob/master/" +
			"cmd/generate-openapischema/README.md" +
			"#possible-missing-schema-error-messages-and-how-to-solve-them"
		fmt.Printf("\"openapi-gen\" did not generate all the needed schemas.\n"+
			"Refer to %q for possible causes and solutions\n", readmePath)
		panic("Missing schemas!")
	}
	err = g.UpdateYAMLs(yamlsPath)
	if err != nil {
		panic(err)
	}
}
