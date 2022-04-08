package parser

import (
	"go/doc"
	"regexp"
	"strings"
)

func GetNexusRestAPIGenAnnotation(pkg Package, name string) (string, bool) {
	var annotation string
	annotationPrefix := "nexus-rest-api-gen"
	d := doc.New(&pkg.Pkg, pkg.Name, 4)
	for _, t := range d.Types {
		if t.Name == name {
			if strings.Contains(t.Doc, annotationPrefix) {
				re := regexp.MustCompile(annotationPrefix + ".*")
				annotation = re.FindString(t.Doc)
			}
		}
	}
	if annotation != "" {
		return annotation, true
	}
	return annotation, false
}
