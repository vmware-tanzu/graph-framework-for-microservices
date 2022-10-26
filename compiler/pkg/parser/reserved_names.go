package parser

import (
	"fmt"
	"go/ast"
	"strings"

	log "github.com/sirupsen/logrus"
)

var ReservedTypeNames = []string{"ResourceVersion", "SchemeGroupVersion", "Kind", "Resource", "SchemeBuilder",
	"AddToScheme", "Child", "Link", "NexusStatus", "Id"}

func CheckIfReserved(name string) error {
	for _, reservedName := range ReservedTypeNames {
		if strings.EqualFold(name, reservedName) {
			return fmt.Errorf("name %s is reserved", reservedName)
		}
	}
	return nil
}

func CheckIfFieldsAreReserved(s *ast.StructType) {
	for _, f := range s.Fields.List {
		name, err := GetFieldName(f)
		if err != nil {
			log.Fatalf("failed to determine field name: %v", err)
		}
		err = CheckIfReserved(name)
		if err != nil {
			log.Fatalf("Please change the field name %s as %s ", name, err.Error())
		}
	}
}

func CheckIfNameReserved(n *ast.TypeSpec) {
	err := CheckIfReserved(n.Name.Name)
	if err != nil {
		log.Fatalf("Please change the type name %s as %s ", n.Name.Name, err.Error())
	}
	if val, ok := n.Type.(*ast.StructType); ok {
		CheckIfFieldsAreReserved(val)
	}
}
