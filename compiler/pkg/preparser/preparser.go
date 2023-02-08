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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"golang.org/x/tools/imports"
)

var pkgImportToPkg = make(map[string]string, 0)

func Parse(startPath string) map[string][]*parser.Package {
	packages := map[string][]*parser.Package{}
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
					ModPath:  modulePath,
					FullName: pkgImport,
					FileSet:  fileset,
					Pkg:      *v,
				}
				parser.ParseGenDecls(v, &pkg)
				packages[pkg.Name] = append(packages[pkg.Name], &pkg)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to parse DSL: %v", err)
	}

	detectDuplicates(packages)

	return packages
}

func detectDuplicates(packages map[string][]*parser.Package) {
	for _, v := range packages {
		for _, pkg := range v {
			nodes := pkg.GetNodes()
			nexusNodes := pkg.GetNexusNodes()

			nodes = append(nodes, pkg.GetNonStructTypes()...)

			for _, nexusNode := range nexusNodes {
				for _, node := range nodes {
					if node.Name.String() == fmt.Sprintf("%sSpec", nexusNode.Name) ||
						node.Name.String() == fmt.Sprintf("%sList", nexusNode.Name) {
						log.Fatalf(`Duplicated type (%s) found in package %s ("%s" is already used by node "%s")`, node.Name, pkg.Name, node.Name, nexusNode.Name)
					}
				}
			}
		}
	}
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

func Render(dslDir string, packages map[string][]*parser.Package) error {
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

func CopyPkgsToBuild(dslDir string, outputDir string) error {
	dir, err := ioutil.ReadDir(dslDir)
	if err != nil {
		return err
	}

	for _, f := range dir {
		if !f.IsDir() || f.Name() == "build" || f.Name() == "vendor" ||
			f.Name() == "global" || strings.HasPrefix(f.Name(), ".") {
			continue
		}

		opt := cp.Options{
			Skip: func(srcinfo os.FileInfo, src, dest string) (bool, error) {
				return !strings.HasSuffix(src, ".go"), nil
			},
		}

		fmt.Println(outputDir)
		fmt.Println(filepath.Join(outputDir, "model", f.Name()))

		err := cp.Copy(filepath.Join(dslDir, f.Name()), filepath.Join(outputDir, "model", f.Name()), opt)
		if err != nil {
			return err
		}
	}

	return nil
}

var importsTemplate = `{{range .ImportsToRender}}{{.}}{{end}}
`

type ImportsTemplateVars struct {
	ImportsToRender []string
}

func RenderImports(packages map[string][]*parser.Package, outputDir string, modPath string) error {
	importsTemplateVars := ImportsTemplateVars{
		ImportsToRender: []string{},
	}

	for k, pkgs := range packages {
		if k != "global" {
			continue
		}

		for _, v := range pkgs {
			for _, importExpr := range v.GetImportMap() {
				importStr := strings.Trim(importExpr, "\"")
				if strings.HasPrefix(importStr, v.ModPath+"/global") ||
					!strings.HasPrefix(importStr, v.ModPath) {
					continue
				}

				namedImport := "nexus_" + importStr[strings.LastIndex(importStr, "/")+1:]
				importsTemplateVars.ImportsToRender = append(importsTemplateVars.ImportsToRender,
					fmt.Sprintf("%s %s", namedImport, strings.ReplaceAll(importExpr, v.ModPath, modPath)))
			}
		}
	}

	t, err := template.New("tmpl").Parse(importsTemplate)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	if err = t.Execute(&b, importsTemplateVars); err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(outputDir, "model", "nexus-dm-imports"), b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
