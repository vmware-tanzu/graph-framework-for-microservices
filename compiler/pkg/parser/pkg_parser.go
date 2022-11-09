package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"
)

// ParseDSLPkg walks recursively through given path and looks for structs types definitions to add them to a Package map
func ParseDSLPkg(startPath string) Packages {
	modulePath := GetModulePath(startPath)

	packages := make(Packages)
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

			for _, f := range config.ConfigInstance.IgnoredDirs {
				if info.Name() == f {
					log.Infof(fmt.Sprintf("Ignoring %v directory from config", f))
					return filepath.SkipDir
				}
			}

			fileset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fileset, path, nil, parser.ParseComments)
			if err != nil {
				log.Fatalf("failed to parse directory %s: %v", path, err)
			}
			for _, v := range pkgs {
				if v.Name == "nexus" {
					log.Infof("Ignoring nexus package...")
					continue
				}

				if SpecialCharsPresent(v.Name) {
					log.Fatalf("Invalid package-name <%v>, special characters are not allowed. Please use only lowercase alphanumeric characters.", v.Name)
				}
				pkgImport := strings.TrimSuffix(strings.ReplaceAll(path, startPath, modulePath), "/")
				pkg := Package{
					Name:     v.Name,
					FullName: pkgImport,
					ModPath:  modulePath,
					FileSet:  fileset,
					Pkg:      *v,
				}
				parseGenDecls(v, &pkg)
				packages[pkgImport] = pkg
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to parse DSL: %v", err)
	}

	return packages
}

func parseGenDecls(v *ast.Package, pkg *Package) {
	sortedKeys := make([]string, 0, len(v.Files))
	for k := range v.Files {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		file := v.Files[k]
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				pkg.GenDecls = append(pkg.GenDecls, *genDecl)
			}
		}
	}
}
