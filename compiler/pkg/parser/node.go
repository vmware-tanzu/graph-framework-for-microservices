package parser

import (
	"go/ast"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Node struct {
	Name             string
	PkgName          string
	FullName         string
	CrdName          string
	IsSingleton      bool
	Imports          []*ast.ImportSpec
	TypeSpec         *ast.TypeSpec
	Parents          []string
	SingleChildren   map[string]Node
	MultipleChildren map[string]Node
	SingleLink       map[string]Node
	MultipleLink     map[string]Node
	GraphqlQuerySpec nexus.GraphQLQuerySpec
	GraphqlSpec      nexus.GraphQLSpec
}

type NodeHelper struct {
	Name             string
	RestName         string
	Parents          []string
	Children         map[string]NodeHelperChild // CRD Name => NodeHelperChild
	Links            map[string]NodeHelperChild // FieldName => NodeHelperChild
	RestMappings     map[string]string
	IsSingleton      bool
	GraphqlQuerySpec nexus.GraphQLQuerySpec
	GraphqlSpec      nexus.GraphQLSpec
}

type NodeHelperChild struct {
	FieldName    string `json:"fieldName"`
	FieldNameGvk string `json:"fieldNameGvk"`
	IsNamed      bool   `json:"isNamed"`
}

type NonNexusTypes struct {
	Types  map[string]ast.Decl
	Values []string
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
