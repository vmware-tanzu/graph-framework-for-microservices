package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ParseDSLNodes walks recursively through given path and looks for structs types definitions to add them to graph
func ParseDSLNodes(startPath string) Graph {
	modulePath := GetModulePath(startPath)

	rootNodes := make([]string, 0)
	nodes := make(Graph)
	err := filepath.Walk(startPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			fileset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fileset, path, nil, parser.ParseComments)
			if err != nil {
				log.Fatalf("Failed to parse directory %s: %v", path, err)
			}
			for _, v := range pkgs {
				pkgImport := strings.TrimSuffix(strings.ReplaceAll(path, startPath, fmt.Sprintf("%s/", modulePath)), "/")
				for _, file := range v.Files {
					for _, decl := range file.Decls {
						genDecl, ok := decl.(*ast.GenDecl)
						if ok {
							for _, spec := range genDecl.Specs {
								typeSpec, ok := spec.(*ast.TypeSpec)
								if ok {
									val, ok := typeSpec.Type.(*ast.StructType)
									if ok {
										// Detect root nodes
										if path == startPath {
											rootNodes = append(rootNodes, fmt.Sprintf("%s/%s", pkgImport, typeSpec.Name.Name))
										}

										// Create node
										node := GraphNode{
											Name:             typeSpec.Name.Name,
											ImportPath:       pkgImport,
											Path:             path,
											Imports:          file.Imports,
											Spec:             make([]*ast.Field, 0),
											SingleChildren:   make(Graph),
											MultipleChildren: make(Graph),
											SingleLink:       make(Graph),
											MultipleLink:     make(Graph),
										}
										if val.Fields != nil {
											node.Fields = *val.Fields
										}
										if genDecl.Doc != nil {
											node.Comments = *genDecl.Doc
										}

										// Add nexus.Node to nodes map
										if node.IsNexusNode() {
											nodes[fmt.Sprintf("%s/%s", pkgImport, node.Name)] = node
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to ParseDSLNodes %v", err)
	}

	graph := make(Graph)
	for _, root := range rootNodes {
		r := nodes[root]
		ProcessNode(&r, nodes)
		graph[root] = r
	}

	return graph
}

func ProcessNode(node *GraphNode, nodes Graph) {
	for _, field := range node.Fields.List {
		if field.Tag != nil {
			// Get type definition from a field
			isMap := false
			var x, sel string
			switch fieldType := field.Type.(type) {
			case *ast.SelectorExpr:
				x = types.ExprString(fieldType.X)
				sel = fieldType.Sel.String()
			case *ast.MapType:
				isMap = true
				switch mapType := fieldType.Value.(type) {
				case *ast.SelectorExpr:
					x = types.ExprString(mapType.X)
					sel = mapType.Sel.String()
				case *ast.Ident:
					sel = mapType.String()
				}
			}

			// Parse field tags
			tagsStr, err := strconv.Unquote(field.Tag.Value)
			if err != nil {
				log.Fatalf("Failed to parse field tags: %v", err)
			}
			tags := reflect.StructTag(tagsStr)

			// Prepare correct key to get node from nodes map
			if val, ok := tags.Lookup("nexus"); ok {
				var key string
				if x == "" {
					key = fmt.Sprintf("%s/%s", node.ImportPath, sel)
				} else {
					for _, importSpec := range node.Imports {
						importPath, err := strconv.Unquote(importSpec.Path.Value)
						if err != nil {
							log.Fatalf("Failed to parse imports: %v", err)
						}
						if strings.HasSuffix(importPath, x) || importSpec.Name.String() == x {
							key = fmt.Sprintf("%s/%s", importPath, sel)
						}
					}
				}

				relationKey := fmt.Sprintf("%s.%s", x, sel)
				if len(field.Names) == 0 {
					log.Fatalf("Sorry, child and link without a name is not supported, node: %s", node.Name)
				} else if len(field.Names) == 1 {
					relationKey = field.Names[0].Name
				} else if len(field.Names) > 1 {
					log.Fatalf("Sorry, only one name of field is supported, node: %s, field: %s",
						node.Name, field.Names)
				}

				// Add node to Children or Link map
				if val == "child" {
					n := nodes[key]
					ProcessNode(&n, nodes)

					if isMap {
						node.MultipleChildren[relationKey] = n
					} else {
						node.SingleChildren[relationKey] = n
					}
				} else if val == "link" {
					if isMap {
						node.MultipleLink[relationKey] = nodes[key]
					} else {
						node.SingleLink[relationKey] = nodes[key]
					}
				} else {
					node.Spec = append(node.Spec, field)
				}
			}
		} else {
			typeString := types.ExprString(field.Type)
			if typeString != "nexus.Node" {
				node.Spec = append(node.Spec, field)
			}
		}
	}
}
