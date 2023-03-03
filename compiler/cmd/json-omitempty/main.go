package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	goParser "go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

var pkgImportToPkg = make(map[string]string, 0)

func main(){
	dslDir := flag.String("dsl", "datamodel", "DSL file location.")
	flag.Parse()
	parse(*dslDir)
}

func parse(startPath string){
	fmt.Println("startPath",startPath)
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

	detectJsonTag(packages)
}

func detectJsonTag(packages map[string][]*parser.Package) {
	// Set the path to the package directory
	pkgPath := "/Users/kvishnu/Documents/workspace/common-apis/cosmos-datamodel/org"

	// Parse the package directory to get a list of Go files
	files, err := parsePackageDir(pkgPath)
	if err != nil {
		panic(err)
	}

	// Process each file in the package directory
	for _, file := range files {
		// Parse the file to get its AST
		fset := token.NewFileSet()
		node, err := goParser.ParseFile(fset, file, nil, goParser.ParseComments)
		if err != nil {
			panic(err)
		}

		// Modify the AST to add json tags to each struct field
		ast.Inspect(node, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.TypeSpec:
				// switch tt := t.Type.(type) {
				// case *ast.StructType:
				// 	addJSONTags(tt)
				// }
				for _, nf := range parser.GetSpecFields(t) {
					log.Error("******",t.Name.Name,nf.Names,nf.Tag)
					addJSONTags(nf)
				}
			}
			return true
		})

		// Generate the modified source code
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, node); err != nil {
			panic(err)
		}
		formattedCode := buf.Bytes()

		// Write the modified source code back to the file
		if err := ioutil.WriteFile(file, formattedCode, 0); err != nil {
			panic(err)
		}
	}
	for _, v := range packages {
		for _, pkg := range v {
			nodes := pkg.GetNodes()
			nexusNodes := pkg.GetNexusNodes()

			nodes = append(nodes, pkg.GetNonStructTypes()...)
			for _, nexusNode := range nexusNodes {
				for _, nf := range parser.GetSpecFields(nexusNode) {
					log.Error("******",nexusNode.Name.Name,nf.Names,nf.Tag)
					addJSONTags(nf)
				}
			}
			// for _, nexusNode := range nexusNodes {
			// 	for _, node := range nodes {
			// 		if node.Name.String() == fmt.Sprintf("%sSpec", nexusNode.Name) ||
			// 			node.Name.String() == fmt.Sprintf("%sList", nexusNode.Name) {
			// 			log.Fatalf(`Duplicated type (%s) found in package %s ("%s" is already used by node "%s")`, node.Name, pkg.Name, node.Name, nexusNode.Name)
			// 		}
			// 	}
			// }
		}
	}
}

func parsePackageDir(pkgPath string) ([]string, error) {
	// Get a list of Go files in the package directory
	var files []string
	err := filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func addJSONTags(field *ast.Field) {
	
		if field.Tag == nil {
			field.Tag = &ast.BasicLit{}
		}
		
		if !strings.Contains(field.Tag.Value, "json") {
			tagValue := strings.Trim(field.Tag.Value, "`")
			tagValue += ` json:"` + strcase.ToLowerCamel(field.Names[0].Name) + `,omitempty"`
			field.Tag.Value = "`" + tagValue + "`"
		}
	
}
