package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"

	log "github.com/sirupsen/logrus"
)

// ParseDSLNodes walks recursively through given path and looks for structs types definitions to add them to graph
func ParseDSLNodes(startPath string, baseGroupName string) map[string]Node {
	modulePath := GetModulePath(startPath)

	rootNodes := make([]string, 0)
	nodes := make(map[string]Node)
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
								if typeSpec, ok := spec.(*ast.TypeSpec); ok {
									if _, ok := typeSpec.Type.(*ast.StructType); ok {
										// Detect root nodes
										if path == startPath {
											rootNodes = append(rootNodes, fmt.Sprintf("%s/%s", pkgImport, typeSpec.Name.Name))
										}

										plural := util.ToPlural(strings.ToLower(typeSpec.Name.Name))
										crdName := fmt.Sprintf("%s.%s", plural, baseGroupName)

										if IsNexusNode(typeSpec) {
											node := Node{
												Name:             typeSpec.Name.Name,
												FullName:         pkgImport,
												CrdName:          crdName,
												Imports:          file.Imports,
												TypeSpec:         typeSpec,
												Parents:          make([]string, 0),
												SingleChildren:   make(map[string]Node),
												MultipleChildren: make(map[string]Node),
												SingleLink:       make(map[string]Node),
												MultipleLink:     make(map[string]Node),
											}

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

	graph := make(map[string]Node)
	for _, root := range rootNodes {
		r := nodes[root]
		processNode(&r, nodes)
		graph[root] = r
	}

	return graph
}

func CreateParentsMap(root Node) map[string][]string {
	parents := make(map[string][]string)
	root.Walk(func(node *Node) {
		parents[node.CrdName] = node.Parents
	})
	return parents
}

func processNode(node *Node, nodes map[string]Node) {
	childFields := GetChildFields(node.TypeSpec)
	linkFields := GetLinkFields(node.TypeSpec)

	processField := func(f *ast.Field, isChild bool, isLink bool) {
		isMap := IsMapField(f)
		fieldTypeStr := GetFieldType(f)

		var key string
		fieldType := strings.Split(fieldTypeStr, ".")
		if len(fieldType) == 1 {
			key = fmt.Sprintf("%s/%s", node.FullName, fieldType[0])
		}

		if len(fieldType) == 2 {
			for _, importSpec := range node.Imports {
				importPath, err := strconv.Unquote(importSpec.Path.Value)
				if err != nil {
					log.Fatalf("Failed to parse imports: %v", err)
				}

				if strings.HasSuffix(importPath, fieldType[0]) || importSpec.Name.String() == fieldType[0] {
					key = fmt.Sprintf("%s/%s", importPath, fieldType[1])
				}
			}
		}

		if isChild {
			n := nodes[key]
			n.Parents = node.Parents
			n.Parents = append(n.Parents, node.CrdName)
			processNode(&n, nodes)

			if isMap {
				node.MultipleChildren[fieldTypeStr] = n
			} else {
				node.SingleChildren[fieldTypeStr] = n
			}
		}

		if isLink {
			if isMap {
				node.MultipleLink[fieldTypeStr] = nodes[key]
			} else {
				node.SingleLink[fieldTypeStr] = nodes[key]
			}
		}

	}

	for _, child := range childFields {
		processField(child, true, false)
	}

	for _, link := range linkFields {
		processField(link, false, true)
	}
}
