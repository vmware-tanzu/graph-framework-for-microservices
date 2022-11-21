package openapi_generator_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi/pkg/validation/spec"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generator Suite")
}

func createFileWithEmptyYAMLDefinitions(tmpDir string, names []string) string {
	content := ""
	for _, name := range names {
		content += getEmptyYAMLDefinition(name)
	}
	tmpFile := fmt.Sprintf("%s/%s.yaml", tmpDir, strings.Join(names, "_"))
	err := ioutil.WriteFile(tmpFile, []byte(content), 0665)
	Expect(err).NotTo(HaveOccurred())
	return tmpFile
}

func getEmptyYAMLDefinition(name string) string {
	format := fmt.Sprintf(`---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    nexus: |
      {"name":%q,"is_singleton":false,"nexus-rest-api-gen":{"uris":null}}
  creationTimestamp: null
  name: NAMEs.test.it
spec:
  conversion:
    strategy: None
  group: test.it
  names:
    kind: CAPITAL_NAME
    listKind: CAPITAL_NAMEList
    plural: NAMEs
    shortNames:
      - NAME
    singular: NAME
  scope: Namespaced
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: null
  storedVersions:
  - v1
`, name)
	capitalName := strings.ToUpper(name[:1]) + name[1:]
	format = strings.ReplaceAll(format, "CAPITAL_NAME", capitalName)
	return strings.ReplaceAll(format, "NAME", name)
}

func getSchemaName(name string) string {
	return fmt.Sprintf("nexustempmodule/apis/test.it/v1.%s", name)
}

func compareTmpFileWithExpectedFile(actualFile, expectedFile string) {
	actual, err := ioutil.ReadFile(actualFile)
	Expect(err).NotTo(HaveOccurred())

	expected, err := ioutil.ReadFile(expectedFile)
	Expect(err).NotTo(HaveOccurred())

	Expect(actual).To(MatchYAML(string(expected)))
}

func fooDefinition() common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"fizz": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"string"},
						},
					},
					"buzz": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"integer"},
						},
					},
				},
			},
		},
	}
}
