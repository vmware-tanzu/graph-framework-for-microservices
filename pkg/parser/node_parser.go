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
				return filepath.SkipDir
			}
			if info.Name() == "vendor" {
				log.Infof("Ignoring vendor directory...")
				return filepath.SkipDir
			}
			fileset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fileset, path, nil, parser.ParseComments)
			if err != nil {
				log.Fatalf("Failed to parse directory %s: %v", path, err)
			}
			for _, v := range pkgs {
				if v.Name == "" {
					log.Fatalf("Failed to get package name for %#v", v)
				}
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
										checkIfReserved(typeSpec.Name.Name)
										crdName := util.GetCrdName(typeSpec.Name.Name, v.Name, baseGroupName)
										if IsNexusNode(typeSpec) {
											// Detect root nodes
											if path == startPath {
												rootNodes = append(rootNodes, crdName)
											}

											node := Node{
												Name:             typeSpec.Name.Name,
												PkgName:          v.Name,
												FullName:         pkgImport,
												CrdName:          crdName,
												IsSingleton:      IsSingletonNode(typeSpec),
												Imports:          file.Imports,
												TypeSpec:         typeSpec,
												Parents:          make([]string, 0),
												SingleChildren:   make(map[string]Node),
												MultipleChildren: make(map[string]Node),
												SingleLink:       make(map[string]Node),
												MultipleLink:     make(map[string]Node),
											}
											if node.CrdName == "" {
												log.Fatalf("Internal compiler failure: Failed to determine crd name of node %v", node.Name)
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

	// TEMP FIX: Make more optimal way to auto discover root nodes.
	// https://jira.eng.vmware.com/browse/NPT-340
	graph := buildGraph(nodes, rootNodes, baseGroupName)
	// Find if any node have root node set as child - if yes remove it from rootNodes
	for _, v := range graph {
		v.Walk(func(node *Node) {
			// iterate from the end because length of the slice may change
			for i := len(rootNodes) - 1; i >= 0; i-- {
				rootNode := rootNodes[i]
				for _, child := range node.SingleChildren {
					// if rootNode is a child then remove it from the slice
					if child.CrdName == rootNode {
						rootNodes = append(rootNodes[:i], rootNodes[i+1:]...)
					}
				}
				for _, child := range node.MultipleChildren {
					// if rootNode is a named child then remove it from the slice
					if child.CrdName == rootNode {
						rootNodes = append(rootNodes[:i], rootNodes[i+1:]...)
					}
				}
			}
		})
	}

	return buildGraph(nodes, rootNodes, baseGroupName)
}

func buildGraph(nodes map[string]Node, rootNodes []string, baseGroupName string) map[string]Node {
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
				if child.CrdName == "" {
					log.Fatalf("Internal compiler failure: Failed to determine crd name of child %v", child)
				}
				children[child.CrdName] = NodeHelperChild{
					IsNamed:      false,
					FieldName:    key,
					FieldNameGvk: util.GetGvkFieldTagName(key),
				}
			}

			for key, child := range node.MultipleChildren {
				if child.CrdName == "" {
					log.Fatalf("Internal compiler failure: Failed to determine crd name of child %v", child)
				}
				children[child.CrdName] = NodeHelperChild{
					IsNamed:      true,
					FieldName:    key,
					FieldNameGvk: util.GetGvkFieldTagName(key),
				}
			}
			links := make(map[string]NodeHelperChild)
			for key, link := range node.SingleLink {
				if link.CrdName == "" {
					log.Fatalf("Internal compiler failure: Failed to determine crd name of link %v", link)
				}
				links[link.CrdName] = NodeHelperChild{
					IsNamed:      false,
					FieldName:    key,
					FieldNameGvk: util.GetGvkFieldTagName(key),
				}
			}

			for key, link := range node.MultipleLink {
				if link.CrdName == "" {
					log.Fatalf("Internal compiler failure: Failed to determine crd name of link %v", link)
				}
				links[link.CrdName] = NodeHelperChild{
					IsNamed:      true,
					FieldName:    key,
					FieldNameGvk: util.GetGvkFieldTagName(key),
				}
			}

			if node.CrdName == "" {
				log.Fatalf("Internal compiler failure: Failed to determine crd name of node %v", node)
			}

			parents[node.CrdName] = NodeHelper{
				Name:        node.Name,
				RestName:    fmt.Sprintf("%s.%s", node.PkgName, node.Name),
				Parents:     node.Parents,
				Children:    children,
				Links:       links,
				IsSingleton: node.IsSingleton,
			}
		})
	}
	return parents
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

		isNamed := IsNamedChildOrLink(f)

		if isNamed {
			childNode := findNodeDefForField(f, node, nodes)
			if childNode == nil {
				log.Fatalf("Internal compiler failure: couldn't determine child node of field %s", f.Names)
			}
			if IsSingletonNode(childNode.TypeSpec) {
				log.Fatalf("Singleton can't be used as a named child, wrong field name %s in node %s",
					f.Names, node.Name)
			}
		}
		fieldName, _ := GetNodeFieldName(f)
		if fieldName == "" {
			log.Fatalf("Internal compiler failure: failed to find field name for field: %v in node %v", f.Names, node.Name)
		}
		key := findFieldKeyForNode(f, node, nodes, baseGroupName)
		if key == "" {
			log.Fatalf("Internal compiler failure: failed to find field key for field %v in node %v", f.Names, node.Name)
		}
		if isChild {
			n, ok := nodes[key]
			if !ok {
				log.Fatalf("Internal compiler failure: couldn't find node for key %v", key)
			}
			n.Parents = node.Parents
			n.Parents = append(n.Parents, node.CrdName)
			processNode(&n, nodes, baseGroupName)

			if isNamed {
				node.MultipleChildren[fieldName] = n
			} else {
				node.SingleChildren[fieldName] = n
			}
		}

		if isLink {
			if isNamed {
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

func findFieldKeyForNode(f *ast.Field, node *Node, nodes map[string]Node, baseGroupName string) (key string) {
	fieldTypeStr := GetFieldType(f)
	fieldType := strings.Split(fieldTypeStr, ".")
	if len(fieldType) == 1 {
		key = util.GetCrdName(fieldType[0], util.RemoveSpecialChars(node.PkgName), baseGroupName)
		return
	} else if len(fieldType) == 2 {
		for _, importSpec := range node.Imports {
			importPath, err := strconv.Unquote(importSpec.Path.Value)
			if err != nil {
				log.Fatalf("Failed to parse imports: %v", err)
			}
			importPathSplit := strings.Split(importPath, ".")
			packageDir := importPathSplit[len(importPathSplit)-1]
			// If import is not named then we can build key without looping through nodes
			if importSpec.Name == nil && packageDir == fieldType[1] {
				key = util.GetCrdName(fieldType[1], util.RemoveSpecialChars(fieldType[0]), baseGroupName)
				return
			} else {
				for _, n := range nodes {
					if n.FullName == importPath && n.Name == fieldType[1] {
						key = n.CrdName
						return
					}
				}
			}
		}
	}
	return
}

func findNodeDefForField(f *ast.Field, baseNode *Node, allNodes map[string]Node) *Node {
	fieldTypeStr := GetFieldType(f)
	fieldType := strings.Split(fieldTypeStr, ".")
	var importPathOfNode string
	var nodeName string
	if len(fieldType) == 1 {
		importPathOfNode = baseNode.FullName
		nodeName = fieldTypeStr
	} else if len(fieldType) == 2 {
		nodeName = fieldType[1]
		importPathOfNode = findMatchingImport(nodeName, baseNode.Imports, allNodes)
	}
	for _, n := range allNodes {
		if n.FullName == importPathOfNode && n.Name == nodeName {
			return &n
		}
	}
	return nil
}

func findMatchingImport(nodeName string, imports []*ast.ImportSpec, allNodes map[string]Node) string {
	for _, importSpec := range imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			log.Fatalf("Failed to parse imports: %v", err)
		}
		importPathSplit := strings.Split(importPath, ".")
		packageDir := importPathSplit[len(importPathSplit)-1]
		if importSpec.Name == nil && packageDir == nodeName {
			return importPath
		} else {
			for _, n := range allNodes {
				if n.FullName == importPath && n.Name == nodeName {
					return importPath
				}
			}
		}
	}
	return ""
}

func checkIfReserved(name string) {
	for _, reservedName := range ReservedTypeNames {
		if name == reservedName {
			log.Fatalf("Name %s is reserved. Please change type name.", reservedName)
		}
	}
}
