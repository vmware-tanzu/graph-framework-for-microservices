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
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/util"
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
	Dir      string

	//new fields
	ModPath  string
	Pkg      ast.Package
	GenDecls []ast.GenDecl
	FileSet  *token.FileSet
}

type FieldAnnotation string

const (
	GRAPHQL_ARGS_ANNOTATION          = FieldAnnotation("nexus-graphql-args")
	GRAPHQL_ALIAS_NAME_ANNOTATION    = FieldAnnotation("nexus-alias-name")
	GRAPHQL_ALIAS_TYPE_ANNOTATION    = FieldAnnotation("nexus-alias-type")
	GRAPHQL_TSM_DIRECTIVE_ANNOTATION = FieldAnnotation("nexus-graphql-tsm-directive")
	GRAPHQL_NULLABLE_ANNOTATION      = FieldAnnotation("nexus-graphql-nullable")
	GRAPHQL_TS_TYPE_ANNOTATION       = FieldAnnotation("nexus-graphql-ts-type")
	GRAPHQL_JSONENCODED_ANNOTATION   = FieldAnnotation("nexus-graphql-jsonencoded")
	GRAPHQL_RELATION_NAME            = FieldAnnotation("nexus-graphql-relation-name")
	GRAPHQL_RELATION_PARAMETERS      = FieldAnnotation("nexus-graphql-relation-parameters")
	GRAPHQL_RELATION_UUIDKEY         = FieldAnnotation("nexus-graphql-relation-uuidkey")
	GRAPHQL_TYPE_NAME                = FieldAnnotation("nexus-graphql-type-name")
	GRAPHQL_PROTOBUF_NAME            = FieldAnnotation("nexus-graphql-protobuf-name")
	GRAPHQL_PROTOBUF_FILE            = FieldAnnotation("nexus-graphql-protobuf-file")
)

// func (p *Package) GetImports() []*ast.ImportSpec
// func (p *Package) GetNodes() []*ast.StructType
// func (p *Package) GetNexusNodes []*ast.StructType
// func (p *Package) GetTypes []*ast.GenDecl
// func (p *Package) GetConsts() []*ast.ValueSpec

// GetSpecFields(n *ast.StructType) []*ast.Field
// GetChildFields(n *ast.StructType) []*ast.Field
// GetLinkFields(n *ast.StructType) []*ast.Field

// IsChildField(f *ast.Field) bool
// IsLinkField(f *ast.Field) bool
// IsMapField(f *ast.Field) bool

// GetNodeFieldName(f *ast.Field) string // Config or GNS
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

func (p *Package) GetImportMap() map[string]string {
	var importMap = make(map[string]string)
	imports := p.GetImports()
	var importKey string
	for _, val := range imports {
		importVal := val.Path.Value
		if val.Name != nil {
			importKey = val.Name.String()
		} else {
			importKey = importVal[strings.LastIndex(importVal, "/")+1 : len(importVal)-1]
		}
		importMap[importKey] = importVal
	}
	return importMap
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

func (p *Package) GetNonStructTypes() []*ast.TypeSpec {
	var nonStructs []*ast.TypeSpec
	for _, genDecl := range p.GenDecls {
		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					nonStructs = append(nonStructs, typeSpec)
				}
			}
		}
	}
	return nonStructs
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

func (p *Package) GetConsts() []*ast.ValueSpec {
	var consts []*ast.ValueSpec
	for _, genDecl := range p.GenDecls {
		if genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					consts = append(consts, valueSpec)
				}
			}
		}
	}
	return consts
}

func (p *Package) IsVarPresent(varName string) bool {
	for _, genDecl := range p.GenDecls {
		if genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					if varName == valueSpec.Names[0].Name {
						return true
					}
				}
			}
		}
	}
	return false
}

func IsNexusNode(n *ast.TypeSpec) bool {
	if n == nil {
		return false
	}

	if val, ok := n.Type.(*ast.StructType); ok {
		for _, field := range val.Fields.List {
			typeString := types.ExprString(field.Type)
			if typeString == "nexus.Node" || typeString == "nexus.SingletonNode" || typeString == "NexusNode" {
				return true
			}
		}
	}

	return false
}

func IsSingletonNode(n *ast.TypeSpec) bool {
	if n == nil {
		return false
	}

	if val, ok := n.Type.(*ast.StructType); ok {
		for _, field := range val.Fields.List {
			typeString := types.ExprString(field.Type)
			if typeString == "nexus.SingletonNode" {
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
	if typeString == "nexus.Node" || typeString == "nexus.SingletonNode" || typeString == "nexus.ID" || typeString == "NexusNode" {
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

func GetNexusFields(n *ast.TypeSpec) []*ast.Field {
	var fields []*ast.Field
	if n == nil {
		return nil
	}
	if val, ok := n.Type.(*ast.StructType); ok {
		for _, f := range val.Fields.List {
			if IsNexusField(f) || IsNexusTypeField(f) {
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
			if strings.ToLower(val.Name) == "child" || strings.ToLower(val.Name) == "children" {
				return true
			}
		}
	}

	return false
}

func IsNamedChildOrLink(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get("nexus"); err == nil {
			if strings.ToLower(val.Name) == "children" || strings.ToLower(val.Name) == "links" {
				return true
			}
		}
	}

	return false
}

func IsOnlyChildrenField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get("nexus"); err == nil {
			if strings.ToLower(val.Name) == "children" {
				return true
			}
		}
	}
	return false
}

func IsOnlyChildField(f *ast.Field) bool {
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

func IsOnlyLinkField(f *ast.Field) bool {
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

func IgnoreField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get("nexus-graphql"); err == nil {
			if strings.ToLower(val.Name) == "ignore:true" {
				return true
			}
		}
	}

	return false
}

func IsJsonStringField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get("nexus-graphql"); err == nil {
			if strings.ToLower(val.Name) == "type:string" {
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
			if strings.ToLower(val.Name) == "link" || strings.ToLower(val.Name) == "links" {
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

func IsPointerToArrayField(f *ast.Field) bool {
	if f == nil {
		return false
	}
	if starExpr, ok := f.Type.(*ast.StarExpr); ok {
		if _, ok := starExpr.X.(*ast.ArrayType); ok {
			return true
		}
	}
	return false
}

func GetNodeFieldName(f *ast.Field) (string, error) {
	if f == nil {
		return "", errors.New("provided field does not exist")
	}

	if len(f.Names) == 0 {
		return "", errors.New("sorry, children, links and node fields without a name are not supported")
	} else if len(f.Names) > 1 {
		return "", errors.New("sorry, only one name of field is supported")
	}

	return f.Names[0].Name, nil
}

func GetFieldName(f *ast.Field) (string, error) {
	if f == nil {
		return "", errors.New("provided field does not exist")
	}
	if len(f.Names) == 0 {
		return "", nil
	}
	return f.Names[0].Name, nil
}

func GetFieldNameJsonTag(f *ast.Field) string {
	tag := GetFieldJsonTag(f)
	if tag == nil {
		return ""
	}
	return tag.Name
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

func FillEmptyTag(ts *structtag.Tags, n, tag string) *structtag.Tags {

	if ts == nil {
		ts = &structtag.Tags{}
	}
	for _, t := range ts.Tags() {
		if t.Key == tag {
			return ts
		}
	}
	if len(n) < 2 {
		return ts
	}
	jt := structtag.Tag{
		Key:     tag,
		Name:    util.GetTag(n),
		Options: nil,
	}
	err := ts.Set(&jt)
	if err != nil {
		log.Fatalf("Failed to set tag: %v", err)
	}

	return ts
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
		sel = types.ExprString(fieldType)
		if expr, ok := fieldType.X.(*ast.SelectorExpr); ok {
			x = types.ExprString(expr.X)
			sel = expr.Sel.String()
		}
	case *ast.Ident:
		sel = fieldType.String()
	}

	fieldType := fmt.Sprintf("%s.%s", x, sel)
	if x == "" {
		fieldType = sel
	}
	if x != "" && star {
		fieldType = fmt.Sprintf("*%s", fieldType)
	}

	return fieldType
}

func IsFieldPointer(f *ast.Field) bool {
	if f == nil {
		return false
	}

	star := false
	switch fieldType := f.Type.(type) {
	case *ast.MapType:
		if _, ok := fieldType.Value.(*ast.StarExpr); ok {
			star = true
		}
	case *ast.StarExpr:
		star = true
	}
	return star
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

func (p *Package) ValueSpecToString(t *ast.ValueSpec) (string, error) {
	if t == nil {
		return "", errors.New("provided valuespec does not exist")
	}

	buf := new(bytes.Buffer)
	err := printer.Fprint(buf, p.FileSet, t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func IsNexusGraphqlNullField(f *ast.Field) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get(string(GRAPHQL_NULLABLE_ANNOTATION)); err == nil {
			if strings.ToLower(val.Name) == "false" {
				return false
			}
		}
	}

	return true
}

func GetFieldAnnotationString(f *ast.Field, annotation FieldAnnotation) string {
	if f == nil {
		return ""
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get(string(annotation)); err == nil {
			return val.String()
		}
	}
	return ""
}

func GetFieldAnnotationVal(f *ast.Field, annotation FieldAnnotation) string {
	if f == nil {
		return ""
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if val, err := tags.Get(string(annotation)); err == nil {
			return val.Value()
		}
	}
	return ""
}

func IsFieldAnnotationPresent(f *ast.Field, annotation FieldAnnotation) bool {
	if f == nil {
		return false
	}

	if f.Tag != nil {
		tags := ParseFieldTags(f.Tag.Value)
		if _, err := tags.Get(string(annotation)); err == nil {
			return true
		}
	}
	return false
}
