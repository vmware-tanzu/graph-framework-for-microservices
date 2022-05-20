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
	pkgsMap := make(map[string]string)
	err := filepath.Walk(startPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == "build" {
				log.Infof("Ignoring build directory...")
				return nil
			}
			fileset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fileset, path, nil, parser.ParseComments)
			if err != nil {
				log.Fatalf("Failed to parse directory %s: %v", path, err)
			}
			for _, v := range pkgs {
				if _, ok := pkgsMap[v.Name]; ok {
					log.Fatalf("Invalid Package name. Package name <%v> is already defined. Please make sure the package names are not duplicated.", v.Name)
				}
				pkgsMap[v.Name] = v.Name
				pkgImport := strings.TrimSuffix(strings.ReplaceAll(path, startPath, modulePath), "/")
				for _, file := range v.Files {
					for _, decl := range file.Decls {
						genDecl, ok := decl.(*ast.GenDecl)
						if ok {
							for _, spec := range genDecl.Specs {
								if typeSpec, ok := spec.(*ast.TypeSpec); ok {
									if _, ok := typeSpec.Type.(*ast.StructType); ok {
										crdName := util.GetCrdName(typeSpec.Name.Name, v.Name, baseGroupName)
										// Detect root nodes
										if path == startPath {
											rootNodes = append(rootNodes, crdName)
										}

										if IsNexusNode(typeSpec) {
											node := Node{
												Name:             typeSpec.Name.Name,
												PkgName:          v.Name,
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
											nodes[crdName] = node
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
		processNode(&r, nodes, baseGroupName)
		graph[root] = r
	}

	return graph
}

func CreateParentsMap(graph map[string]Node) map[string]NodeHelper {
	parents := make(map[string]NodeHelper)
	for _, root := range graph {
		root.Walk(func(node *Node) {
			children := make(map[string]NodeHelperChild)
			for key, child := range node.SingleChildren {
				children[child.CrdName] = NodeHelperChild{
					IsNamed:      false,
					FieldName:    key,
					FieldNameGvk: util.GetGvkFieldTagName(key),
				}
			}

			for key, child := range node.MultipleChildren {
				children[child.CrdName] = NodeHelperChild{
					IsNamed:      true,
					FieldName:    key,
					FieldNameGvk: util.GetGvkFieldTagName(key),
				}
			}

			parents[node.CrdName] = NodeHelper{
				Name:     node.Name,
				RestName: fmt.Sprintf("%s.%s", node.Name, node.PkgName),
				Parents:  node.Parents,
				Children: children,
			}
		})
	}
	return parents
}

func CreateRestMappings(graph map[string]Node) map[string]string {
	mappings := make(map[string]string)
	for _, root := range graph {
		root.Walk(func(node *Node) {
			mappings[fmt.Sprintf("%s.%s", node.Name, node.PkgName)] = node.CrdName
		})
	}
	return mappings
}

func processNode(node *Node, nodes map[string]Node, baseGroupName string) {
	childFields := GetChildFields(node.TypeSpec)
	linkFields := GetLinkFields(node.TypeSpec)

	processField := func(f *ast.Field, isChild bool, isLink bool) {
		if isChild || isLink {
			if IsArrayField(f) {
				log.Fatalf(`"Invalid Type for %v. Nexus Child or Link can't be an array. Please represent it in the form of a map.`+"\n"+
					`For example: `+
					`myStr []string should be represented in the form of myStr map[string]string`, f.Names)
			}
		}

		if IsFieldPointer(f) {
			log.Fatalf("Pointer type is not allowed. Field <%v> is a pointer. Please make sure nexus child/link types are not pointers.", f.Names)
		}

		isMap := IsMapField(f)
		fieldTypeStr := GetFieldType(f)
		fieldName, _ := GetFieldName(f)

		var key string
		fieldType := strings.Split(fieldTypeStr, ".")
		if len(fieldType) == 1 {
			key = util.GetCrdName(fieldType[0], util.RemoveSpecialChars(node.PkgName), baseGroupName)
		}

		if len(fieldType) == 2 {
			for _, importSpec := range node.Imports {
				importPath, err := strconv.Unquote(importSpec.Path.Value)
				if err != nil {
					log.Fatalf("Failed to parse imports: %v", err)
				}

				// If import is not named then we can build key without looping through nodes
				if importSpec.Name == nil {
					key = util.GetCrdName(fieldType[1], util.RemoveSpecialChars(fieldType[0]), baseGroupName)
				} else {
					for _, n := range nodes {
						if n.FullName == importPath && n.Name == fieldType[1] {
							key = n.CrdName
						}
					}
				}
			}
		}
		if isChild {
			n := nodes[key]
			n.Parents = node.Parents
			n.Parents = append(n.Parents, node.CrdName)
			processNode(&n, nodes, baseGroupName)

			if isMap {
				node.MultipleChildren[fieldName] = n
			} else {
				node.SingleChildren[fieldName] = n
			}
		}

		if isLink {
			if isMap {
				node.MultipleLink[fieldName] = nodes[key]
			} else {
				node.SingleLink[fieldName] = nodes[key]
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
