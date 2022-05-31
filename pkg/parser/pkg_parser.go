package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
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

				for _, file := range v.Files {
					for _, decl := range file.Decls {
						if genDecl, ok := decl.(*ast.GenDecl); ok {
							pkg.GenDecls = append(pkg.GenDecls, *genDecl)
						}
					}
				}

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
