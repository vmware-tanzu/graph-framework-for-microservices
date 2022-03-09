package parser

import (
	"fmt"
	"go/ast"
	"go/types"
)

// Key format: "???/struct_name", Value: GraphNode
type Graph map[string]GraphNode

func (g *Graph) DumpNodeNames() {
	for _, v := range *g {
		fmt.Println(v.Name)
	}
}

func (g *Graph) DumpComments() {
	for k, v := range *g {
		var comments string
		for _, c := range v.Comments.List {
			comments += "  " + c.Text + "\n"
		}
		fmt.Printf("%s:\n%s", k, comments)
	}
}

func (g *Graph) DumpFields() {
	for _, v := range *g {
		dump := fmt.Sprintf("Fields of node %s:\n", v.Name)
		for _, field := range v.Fields.List {
			var names string
			for _, v := range field.Names {
				names += fmt.Sprintf("%s,", v.Name)
			}
			typeString := types.ExprString(field.Type)
			var tag string
			if field.Tag != nil {
				tag = field.Tag.Value
			}
			dump += fmt.Sprintf("  - type: %s, names:%s tag: %s\n", typeString, names, tag)
		}
		fmt.Println(dump)
	}
}

type GraphNode struct {
	Name             string
	ImportPath       string
	Path             string
	Fields           ast.FieldList     `json:"-"`
	Comments         ast.CommentGroup  `json:"-"`
	Imports          []*ast.ImportSpec `json:"-"`
	Spec             []*ast.Field
	SingleChildren   Graph
	MultipleChildren Graph
	SingleLink       Graph
	MultipleLink     Graph
}

// TODO make sure that it works with named imports
func (g *GraphNode) IsNexusNode() bool {
	for _, field := range g.Fields.List {
		typeString := types.ExprString(field.Type)
		if typeString == "nexus.Node" {
			return true
		}
	}
	return false
}

func (g *GraphNode) Walk(fn func(node GraphNode)) {
	fn(*g)

	childs := g.SingleChildren
	for k, v := range g.MultipleChildren {
		childs[k] = v
	}

	for _, n := range childs {
		n.Walk(fn)
	}
}
