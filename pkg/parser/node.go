package parser

import (
	"go/ast"
)

type Node struct {
	Name             string
	PkgName          string
	FullName         string
	CrdName          string
	Imports          []*ast.ImportSpec
	TypeSpec         *ast.TypeSpec
	Parents          []string
	SingleChildren   map[string]Node
	MultipleChildren map[string]Node
	SingleLink       map[string]Node
	MultipleLink     map[string]Node
}

type NodeHelper struct {
	Name     string
	Parents  []string
	Children map[string]NodeHelperChild // CRD Name => NodeHelperChild
}

type NodeHelperChild struct {
	FieldName    string `json:"fieldName"`
	FieldNameGvk string `json:"fieldNameGvk"`
	IsNamed      bool   `json:"isNamed"`
}

func (node *Node) Walk(fn func(node *Node)) {
	fn(node)

	for _, n := range node.MultipleChildren {
		n.Walk(fn)
	}

	for _, n := range node.SingleChildren {
		n.Walk(fn)
	}
}
