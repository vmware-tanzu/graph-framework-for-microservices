package parser

import (
	"go/ast"
	"strings"

	log "github.com/sirupsen/logrus"
)

var ReservedTypeNames = []string{"ResourceVersion", "SchemeGroupVersion", "Kind", "Resource", "SchemeBuilder",
	"AddToScheme", "Child", "Link", "NexusStatus", "Id"}

func CheckIfReserved(name string) {
	for _, reservedName := range ReservedTypeNames {
		if strings.EqualFold(name, reservedName) {
			log.Fatalf("Name %s is reserved. Please change type name from %s to something else", reservedName, name)
		}
	}
}

func CheckIfFieldsAreReserved(s *ast.StructType) {
	for _, f := range s.Fields.List {
		if val, ok := f.Type.(*ast.StructType); ok {
			CheckIfFieldsAreReserved(val)
		}
		name, err := GetFieldName(f)
		if err != nil {
			log.Fatalf("failed to determine field name: %v", err)
		}
		CheckIfReserved(name)
	}
}

func CheckIfNameReserved(n *ast.TypeSpec) {
	CheckIfReserved(n.Name.Name)
	if val, ok := n.Type.(*ast.StructType); ok {
		CheckIfFieldsAreReserved(val)
	}
}
