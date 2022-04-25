package parser

import (
	"go/doc"
	"regexp"
	"strings"
)

const NexusRestApiGenAnnotation = "nexus-rest-api-gen"

func GetNexusRestAPIGenAnnotation(pkg Package, name string) (string, bool) {
	var annotation string

	d := doc.New(&pkg.Pkg, pkg.Name, 4)
	for _, t := range d.Types {
		if t.Name == name {
			if strings.Contains(t.Doc, NexusRestApiGenAnnotation) {
				re := regexp.MustCompile(NexusRestApiGenAnnotation + ".*")
				annotation = re.FindString(t.Doc)
			}
		}
	}

	if annotation != "" {
		val := strings.Split(annotation, ":")
		if len(val) == 2 {
			return strings.TrimSpace(val[1]), true
		}

		return annotation, false
	}

	return annotation, false
}
