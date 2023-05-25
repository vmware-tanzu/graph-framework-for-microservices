package parser

import (
	"go/ast"
	"go/types"
	"os"
	"path"
	"regexp"

	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/modfile"
)

func GetModulePath(startPath string) string {
	file, err := os.ReadFile(path.Join(startPath, "go.mod"))
	if err != nil {
		log.Fatalf("failed to get module path %v", err)
	}
	return modfile.ModulePath(file)
}

func SpecialCharsPresent(name string) bool {
	re, err := regexp.Compile(`[^a-z0-9]`)
	if err != nil {
		log.Fatalf("failed to check for special characters in the package name %v : %v", name, err)
	}
	return re.MatchString(name)
}

//TODO: Move this to COMMON nexus repo

// HTTPMethod type.
type HTTPMethod string

// ResponseCode type.
type ResponseCode int

// HTTPResponse type.
type HTTPResponse struct {
	Description string `json:"description"`
}

// HTTPCodesResponse code to response type.
type HTTPCodesResponse map[ResponseCode]HTTPResponse

// HTTPMethodsResponses to response mapping.
type HTTPMethodsResponses map[HTTPMethod]HTTPCodesResponse

// RestURIs and associated data.
type RestURIs struct {
	Uri     string               `json:"uri"`
	Methods HTTPMethodsResponses `json:"methods"`
}

type RestAPISpec struct {
	Uris []RestURIs `json:"uris"`
}

func CheckValueSpec(valueSpec *ast.ValueSpec) bool {
	if len(valueSpec.Names) == 0 || valueSpec.Values == nil {
		return false
	}
	return true
}

type NexusSpec struct {
	Name  string
	Value *ast.CompositeLit
}

func GetNexusSpecs(p Package, nexusType string) (specs []NexusSpec) {
	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || !CheckValueSpec(valueSpec) {
				continue
			}
			name := valueSpec.Names[0].Name
			value, ok := valueSpec.Values[0].(*ast.CompositeLit)
			if !ok || types.ExprString(value.Type) != nexusType {
				continue
			}
			specs = append(specs, NexusSpec{
				Name:  name,
				Value: value,
			})
		}

	}
	return
}

var basicTypes = []string{
	"bool",
	"byte",
	"complex64",
	"complex128",
	"error",
	"float32",
	"float64",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"rune",
	"string",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",
}

func isBasicType(fieldType ast.Expr) bool {
	switch typ := fieldType.(type) {
	case *ast.Ident:
		return isBasicIdent(typ.Name)
	case *ast.ArrayType:
		eltType, ok := typ.Elt.(*ast.Ident)
		if !ok {
			return false
		}
		return isBasicIdent(eltType.Name)
	case *ast.MapType:
		return false // Assuming maps are not considered basic types
	case *ast.StructType:
		return false // Assuming structs are not considered basic types
	}

	return false
}

func isBasicIdent(name string) bool {
	for _, t := range basicTypes {
		if t == name {
			return true
		}
	}

	return false
}
