package preparser

import (
	"bytes"
	"fmt"
	"go/ast"
	goParser "go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"golang.org/x/tools/imports"
)

var pkgImportToPkg = make(map[string]string, 0)

func Parse(startPath string) map[string][]parser.Package {
	packages := map[string][]parser.Package{}
	modulePath := parser.GetModulePath(startPath)
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
			pkgs, err := goParser.ParseDir(fileset, path, nil, goParser.ParseComments)
			if err != nil {
				log.Fatalf("failed to parse directory %s: %v", path, err)
			}
			pkgImport := strings.TrimSuffix(strings.ReplaceAll(path, startPath, modulePath), "/")
			for _, v := range pkgs {
				if v.Name == "nexus" {
					log.Infof("Ignoring nexus package...")
					continue
				}

				if parser.SpecialCharsPresent(v.Name) {
					log.Fatalf("Invalid package-name <%v>, special characters are not allowed. Please use only lowercase alphanumeric characters.", v.Name)
				}
				pkgImportToPkg[pkgImport] = v.Name

				pkg := parser.Package{
					Name:     v.Name,
					FullName: pkgImport,
					FileSet:  fileset,
					Pkg:      *v,
				}
				packages[pkg.Name] = append(packages[pkg.Name], pkg)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to parse DSL: %v", err)
	}

	return packages
}

func removeImportIdentifierFromFields(file *ast.File, pkg string) *ast.File {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, field := range structType.Fields.List {
				selectorExpr, ok := field.Type.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				if types.ExprString(selectorExpr.X) == "nexus" {
					continue
				}

				var modify bool
				for _, imp := range file.Imports {
					pkgImport, _ := strconv.Unquote(imp.Path.Value)
					if types.ExprString(selectorExpr.X) == imp.Name.String() {
						if val := pkgImportToPkg[pkgImport]; val == pkg {
							modify = true
						}
					}
				}

				if types.ExprString(selectorExpr.X) == pkg {
					modify = true
				}

				if modify {
					field.Type = selectorExpr.Sel
				}

			}
		}
	}
	return file
}

func Render(dslDir string, packages map[string][]parser.Package) error {
	for k, pkgs := range packages {
		if len(pkgs) == 1 {
			continue
		}
		pkgDir := filepath.Join(dslDir, k)
		//TODO: create unique directory (e.g. global may already exist)
		err := os.MkdirAll(pkgDir, os.ModePerm)
		if err != nil {
			return err
		}

		created := map[string]int{}
		for _, pkg := range pkgs {
			for _, file := range pkg.Pkg.Files {
				file = removeImportIdentifierFromFields(file, pkg.Name)

				// remove src file
				srcFile := pkg.FileSet.Position(file.Package).Filename
				err := os.Remove(srcFile)
				if err != nil {
					return err
				}

				filename := filepath.Base(srcFile)
				if _, ok := created[filename]; ok {
					created[filename]++
					filename = fmt.Sprintf("%d_%s", created[filename], filename)
				}

				// render AST to buffer
				var buf bytes.Buffer
				err = printer.Fprint(&buf, pkg.FileSet, file)
				if err != nil {
					return err
				}

				// format file & organize imports using imports package
				out, err := imports.Process(filename, buf.Bytes(), nil)
				if err != nil {
					return err
				}

				// create file and write output
				f, err := os.Create(filepath.Join(pkgDir, filename))
				if err != nil {
					return err
				}
				_, err = f.Write(out)
				if err != nil {
					return err
				}

				err = f.Close()
				if err != nil {
					return err
				}

				created[filename] = 0
			}
		}

	}

	return nil
}
