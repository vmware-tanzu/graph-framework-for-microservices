package parser

import (
	"go/doc"
	"regexp"
	"strings"
)

const (
	NexusRestApiGenAnnotation  = "nexus-rest-api-gen"
	NexusDescriptionAnnotation = "nexus-description"
	NexusGraphqlAnnotation     = "nexus-graphql-query"
	NexusSecretSpecAnnotation  = "nexus-secret-spec"
)

func GetNexusSecretSpecAnnotation(pkg Package, name string) (string, bool) {
	return getNexusAnnotation(pkg, name, NexusSecretSpecAnnotation)
}

func GetNexusRestAPIGenAnnotation(pkg Package, name string) (string, bool) {
	return getNexusAnnotation(pkg, name, NexusRestApiGenAnnotation)
}

func GetNexusDescriptionAnnotation(pkg Package, name string) (string, bool) {
	return getNexusAnnotation(pkg, name, NexusDescriptionAnnotation)
}

func GetNexusGraphqlAnnotation(pkg Package, name string) (string, bool) {
	return getNexusAnnotation(pkg, name, NexusGraphqlAnnotation)
}

func getNexusAnnotation(pkg Package, name string, annotationName string) (string, bool) {
	var annotationValue string

	d := doc.New(&pkg.Pkg, pkg.Name, 4)
	for _, t := range d.Types {
		if t.Name == name {
			if strings.Contains(t.Doc, annotationName) {
				re := regexp.MustCompile(annotationName + ".*")
				annotationValue = re.FindString(t.Doc)
			}
		}
	}

	if annotationValue != "" {
		val := strings.Split(annotationValue, ":")
		if len(val) == 2 {
			return strings.TrimSpace(val[1]), true
		}

		return annotationValue, false
	}

	return annotationValue, false
}
