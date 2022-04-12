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
	Name    string
	Parents []string
}

func (node *Node) Walk(fn func(node *Node)) {
	fn(node)

	children := node.SingleChildren
	for k, v := range node.MultipleChildren {
		children[k] = v
	}

	for _, n := range children {
		n.Walk(fn)
	}
}
