package parser

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	log "github.com/sirupsen/logrus"
)

type Packages map[string]Package

type PackageNodeRelation struct {
	Name  string
	IsMap bool
	Type  string
}

type Package struct {
	Name     string
	FullName string

	//new fields
	ModPath  string
	Pkg      ast.Package
	GenDecls []ast.GenDecl
	FileSet  *token.FileSet
}

// func (p *Package) GetImports() []*ast.ImportSpec
// func (p *Package) GetNodes() []*ast.StructType
// func (p *Package) GetNexusNodes []*ast.StructType
// func (p *Package) GetTypes []*ast.GenDecl

// GetSpecFields(n *ast.StructType) []*ast.Field
// GetChildFields(n *ast.StructType) []*ast.Field
// GetLinkFields(n *ast.StructType) []*ast.Field

// IsChildField(f *ast.Field) bool
// IsLinkField(f *ast.Field) bool
// IsMapField(f *ast.Field) bool

// GetFieldName(f *ast.Field) string // Config or GNS
// GetFieldType(f *ast.Field) string // -> config.Config or Config

// Go 1.18+
// func (p *Package) ToString[T any](n []T) string

// Go 1.17
// func (p *Package) GenDeclToString(n *ast.GenDecl) string
// func (p *Package) StructToString(n *ast.StructType) string

func (p *Package) GetImports() []*ast.ImportSpec {
	var imports []*ast.ImportSpec
	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.ImportSpec); ok {
				imports = append(imports, typeSpec)
			}
		}
	}
	return imports
}

func (p *Package) GetImportStrings() []string {
	var importList []string
	imports := p.GetImports()
	for _, val := range imports {
		i := val.Path.Value
		if val.Name != nil {
			i = fmt.Sprintf("%s %s", val.Name.String(), i)
		}
		importList = append(importList, i)
	}
	return importList
}

func (p *Package) GetNodes() []*ast.TypeSpec {
	var nodes []*ast.TypeSpec
	structs := p.GetStructs()
	for _, val := range structs {
		if !IsNexusNode(val) {
			nodes = append(nodes, val)
		}
	}
	return nodes
}

func (p *Package) GetNexusNodes() []*ast.TypeSpec {
	var nodes []*ast.TypeSpec
	structs := p.GetStructs()
	for _, val := range structs {
		if IsNexusNode(val) {
			nodes = append(nodes, val)
		}
	}
	return nodes
}

func (p *Package) GetStructs() []*ast.TypeSpec {
	var structs []*ast.TypeSpec
	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if _, ok := typeSpec.Type.(*ast.StructType); ok {
					structs = append(structs, typeSpec)
				}
			}
		}
	}
	return structs
}

func (p *Package) GetTypes() []ast.GenDecl {
	var genDecls []ast.GenDecl
	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					genDecls = append(genDecls, genDecl)
				}
			}
		}
	}
	return genDecls
}

func IsNexusNode(n *ast.TypeSpec) bool {
	if n == nil {
		return false
	}

	if val, ok := n.Type.(*ast.StructType); ok {
		for _, field := range val.Fields.List {
			typeString := types.ExprString(field.Type)
			if typeString == "nexus.Node" {
				return true
			}
		}
	}

	return false
}

func IsNexusField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if _, err := tags.Get("nexus"); err == nil {
			return true
		}
	}

	return false
}

func IsNexusTypeField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	typeString := types.ExprString(f.Type)
	if typeString == "nexus.Node" || typeString == "nexus.ID" {
		return true
	}

	return false
}

func GetTypeName(n *ast.TypeSpec) string {
	return n.Name.Name
}

func GetSpecFields(n *ast.TypeSpec) []*ast.Field {
	var fields []*ast.Field
	if n == nil {
		return nil
	}
	if val, ok := n.Type.(*ast.StructType); ok {
		for _, f := range val.Fields.List {
			if !IsNexusField(f) && !IsNexusTypeField(f) {
				fields = append(fields, f)
			}
		}
	}

	return fields
}

func GetChildFields(n *ast.TypeSpec) []*ast.Field {
	var fields []*ast.Field
	if n == nil {
		return nil
	}

	if val, ok := n.Type.(*ast.StructType); ok {
		for _, f := range val.Fields.List {
			if IsChildField(f) && !IsNexusTypeField(f) {
				fields = append(fields, f)
			}
		}
	}
	return fields
}

func GetLinkFields(n *ast.TypeSpec) []*ast.Field {
	var fields []*ast.Field
	if n == nil {
		return nil
	}

	if val, ok := n.Type.(*ast.StructType); ok {
		for _, f := range val.Fields.List {
			if IsLinkField(f) && !IsNexusTypeField(f) {
				fields = append(fields, f)
			}
		}
	}
	return fields
}

func IsChildField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get("nexus"); err == nil {
			if strings.ToLower(val.Name) == "child" {
				return true
			}
		}
	}

	return false
}

func IsLinkField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get("nexus"); err == nil {
			if strings.ToLower(val.Name) == "link" {
				return true
			}
		}
	}
	return false
}

func GetStatusField(n *ast.TypeSpec) *ast.Field {
	if n == nil {
		return nil
	}

	var field *ast.Field
	if val, ok := n.Type.(*ast.StructType); ok {
		for _, f := range val.Fields.List {
			if IsStatusField(f) && !IsNexusTypeField(f) {
				if field != nil {
					log.Fatalf("Only one field can be a nexus status field")
				}
				field = f
			}
		}
	}
	return field
}

func IsStatusField(f *ast.Field) bool {
	if f == nil {
		return false
	}
	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get("nexus"); err == nil {
			if strings.ToLower(val.Name) == "status" {
				return true
			}
		}
	}
	return false
}

func IsMapField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if _, ok := f.Type.(*ast.MapType); ok {
		return true
	}

	if starExpr, ok := f.Type.(*ast.StarExpr); ok {
		if _, ok := starExpr.X.(*ast.MapType); ok {
			return true
		}
	}
	return false
}

func IsArrayField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if _, ok := f.Type.(*ast.ArrayType); ok {
		return true
	}
	return false
}

func GetFieldName(f *ast.Field) (string, error) {
	if f == nil {
		return "", errors.New("provided field does not exist")
	}

	if len(f.Names) == 0 {
		return "", errors.New("sorry, child and link without a name is not supported")
	} else if len(f.Names) > 1 {
		return "", errors.New("sorry, only one name of field is supported")
	}

	return f.Names[0].Name, nil
}

func GetFieldJsonTag(f *ast.Field) *structtag.Tag {
	if f == nil || f.Tag == nil {
		return nil
	}

	tags := ParseFieldTags(f.Tag.Value)
	if val, err := tags.Get("json"); err == nil {
		return val
	}

	return nil
}

func GetFieldTags(f *ast.Field) *structtag.Tags {
	if f == nil || f.Tag == nil {
		return nil
	}

	return ParseFieldTags(f.Tag.Value)
}

func GetFieldType(f *ast.Field) string {
	if f == nil {
		return ""
	}

	var x, sel string
	star := false
	switch fieldType := f.Type.(type) {
	case *ast.SelectorExpr:
		x = types.ExprString(fieldType.X)
		sel = fieldType.Sel.String()
	case *ast.MapType:
		if selExpr, ok := fieldType.Value.(*ast.SelectorExpr); ok {
			x = types.ExprString(selExpr.X)
			sel = selExpr.Sel.String()
		}

		if ident, ok := fieldType.Value.(*ast.Ident); ok {
			sel = ident.String()
		}
	case *ast.StarExpr:
		star = true
		if expr, ok := fieldType.X.(*ast.SelectorExpr); ok {
			x = types.ExprString(expr.X)
			sel = expr.Sel.String()
		}

		if mapExpr, ok := fieldType.X.(*ast.MapType); ok {
			if selExpr, ok := mapExpr.Value.(*ast.SelectorExpr); ok {
				x = types.ExprString(selExpr.X)
				sel = selExpr.Sel.String()
			}

			if ident, ok := mapExpr.Value.(*ast.Ident); ok {
				sel = ident.String()
			}
		}
	case *ast.Ident:
		sel = fieldType.String()
	}

	fieldType := fmt.Sprintf("%s.%s", x, sel)
	if x == "" {
		fieldType = sel
	}
	if star {
		fieldType = fmt.Sprintf("*%s", fieldType)
	}

	return fieldType
}

func ParseFieldTags(tag string) *structtag.Tags {
	tagsStr, err := strconv.Unquote(tag)
	if err != nil {
		log.Fatalf("Failed to parse field tags: %v", err)
	}
	tags, err := structtag.Parse(tagsStr)
	if err != nil {
		log.Fatalf("Failed to parse field tags: %v, tag: %s", err, tag)
	}

	return tags
}

func (p *Package) TypeSpecToString(t *ast.TypeSpec) (string, error) {
	if t == nil {
		return "", errors.New("provided typespec does not exist")
	}

	buf := new(bytes.Buffer)
	err := printer.Fprint(buf, p.FileSet, t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (p *Package) GenDeclToString(t *ast.GenDecl) (string, error) {
	if t == nil {
		return "", errors.New("provided gendecl does not exist")
	}

	buf := new(bytes.Buffer)
	err := printer.Fprint(buf, p.FileSet, t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
