package parser

import (
	"go/ast"
)

type Node struct {
	Name             string
	FullName         string
	Imports          []*ast.ImportSpec
	TypeSpec         *ast.TypeSpec
	Parents          []string
	SingleChildren   map[string]Node
	MultipleChildren map[string]Node
	SingleLink       map[string]Node
	MultipleLink     map[string]Node
}

func (n *Node) Walk(fn func(node *Node)) {
	fn(n)

	childs := n.SingleChildren
	for k, v := range n.MultipleChildren {
		childs[k] = v
	}

	for _, n := range childs {
		n.Walk(fn)
	}
}
