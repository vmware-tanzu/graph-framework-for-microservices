package crd_generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/imports"

	"gitlab.eng.vmware.com/nexus/compiler/pkg/parser"
)

//go:embed template/doc.go.tmpl
var docTemplateFile []byte

//go:embed template/register_group.go.tmpl
var registerGroupTemplateFile []byte

//go:embed template/register_crd.go.tmpl
var registerCRDTemplateFile []byte

//go:embed template/types.go.tmpl
var typesTemplateFile []byte

//go:embed template/crd_base.yaml.tmpl
var crdBaseTemplateFile []byte

func RenderCRDTemplate(baseGroupName, crdModulePath string, pkgs parser.Packages, outputDir string) error {
	pkgNames := make([]string, len(pkgs))
	i := 0
	for _, pkg := range pkgs {
		groupName := pkg.Name + "." + baseGroupName
		pkgNames[i] = groupName + ":v1"
		i++
		groupFolder := outputDir + "/apis/" + groupName + "/"
		crdFolder := outputDir + "/crds"
		apiFolder := groupFolder + "v1"
		var err error
		err = createCRDFolder(apiFolder)
		if err != nil {
			return err
		}
		file, err := RenderDocTemplate(baseGroupName, pkg)
		if err != nil {
			return err
		}
		log.Debugf("Rendered doc template for package %s: %s", pkg.Name, file)
		err = createFile(apiFolder, "doc.go", file, true)
		if err != nil {
			return err
		}
		file, err = RenderRegisterGroupTemplate(baseGroupName, pkg)
		if err != nil {
			return err
		}
		log.Debugf("Rendered register group template for package %s: %s", pkg.Name, file)
		err = createFile(groupFolder, "register.go", file, true)
		if err != nil {
			return err
		}
		file, err = RenderRegisterCRDTemplate(crdModulePath, baseGroupName, pkg)
		if err != nil {
			return err
		}
		log.Debugf("Rendered register CRD template for package %s: %s", pkg.Name, file)
		err = createFile(apiFolder, "register.go", file, true)
		if err != nil {
			return err
		}
		file, err = RenderTypesTemplate(pkg)
		if err != nil {
			return err
		}
		log.Debugf("Rendered types template for package %s: %s", pkg.Name, file)
		err = createFile(apiFolder, "types.go", file, true)
		if err != nil {
			return err
		}
		crdFiles, err := RenderCRDBaseTemplate(baseGroupName, pkg)
		if err != nil {
			return err
		}
		for _, f := range crdFiles {
			log.Debugf("Rendered crd base template for package %s: %s", pkg.Name, f.File)
			err = createFile(crdFolder, f.Name, f.File, false)
			if err != nil {
				return err
			}
		}
	}

	return createApiNamesFile(pkgNames, outputDir)
}

func createCRDFolder(name string) error {
	err := os.MkdirAll(name, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating base-group dir %v failed with an error: %v", name, err)
	}
	return nil
}

func createFile(dirName string, fileName string, file *bytes.Buffer, formatFile bool) error {
	var err error
	formatted := file.Bytes()
	if formatFile {
		formatted, err = format.Source(file.Bytes())
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(dirName+"/"+fileName, formatted, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readTemplateFile(rawTemplate []byte) (*template.Template, error) {
	return template.New("tmpl").Parse(string(rawTemplate))
}

func renderTemplate(tmpl *template.Template, data interface{}) (*bytes.Buffer, error) {
	if tmpl == nil {
		return nil, fmt.Errorf("template can not be nil")
	}
	b := bytes.Buffer{}

	if err := tmpl.Execute(&b, data); err != nil {
		log.Errorf("Template resolve failed with error %v. Vars: %+v", err, data)
		return nil, fmt.Errorf("cannot resolve template: %v", err)
	}

	return &b, nil
}

type docVars struct {
	GroupName   string
	GroupGoName string
	Version     string
}

func RenderDocTemplate(baseGroupName string, pkg parser.Package) (*bytes.Buffer, error) {
	vars := docVars{
		// TODO make it configurable
		Version:     "v1",
		GroupName:   pkg.Name + "." + baseGroupName,
		GroupGoName: strings.Title(strings.ToLower(pkg.Name)) + strings.Title(strings.Split(baseGroupName, ".")[0]),
	}
	if vars.GroupGoName == "" || vars.GroupName == "" {
		return nil, fmt.Errorf("failed to determine group name of package")
	}

	docTemplate, err := readTemplateFile(docTemplateFile)
	if err != nil {
		return nil, err
	}

	return renderTemplate(docTemplate, vars)
}

type registerGroupVars struct {
	GroupName   string
	PackageName string
}

func RenderRegisterGroupTemplate(baseGroupName string, pkg parser.Package) (*bytes.Buffer, error) {
	groupName := pkg.Name + "." + baseGroupName
	vars := registerGroupVars{
		GroupName:   groupName,
		PackageName: strings.Replace(groupName, ".", "_", -1),
	}

	if vars.PackageName == "" || vars.GroupName == "" {
		return nil, fmt.Errorf("failed to determine group name of package")
	}

	registerGroupTemplate, err := readTemplateFile(registerGroupTemplateFile)
	if err != nil {
		return nil, err
	}

	return renderTemplate(registerGroupTemplate, vars)
}

type registerCRDVars struct {
	GroupPackageName   string
	GroupPackageImport string
	ResourceVersion    string
	KnownTypes         string
}

func RenderRegisterCRDTemplate(crdModulePath, baseGroupName string, pkg parser.Package) (*bytes.Buffer, error) {
	var knownTypes string
	for _, node := range pkg.GetNexusNodes() {
		knownTypes += "\t\t&" + parser.GetTypeName(node) + "{},\n\t\t&" + parser.GetTypeName(node) + "List{},\n"
	}

	groupName := pkg.Name + "." + baseGroupName
	vars := registerCRDVars{
		GroupPackageName:   strings.Replace(groupName, ".", "_", -1),
		GroupPackageImport: crdModulePath + "apis/" + groupName,
		// TODO make configurable by some variable in package
		ResourceVersion: "v1",
		KnownTypes:      knownTypes,
	}

	if vars.GroupPackageName == "" ||
		vars.GroupPackageImport == "" ||
		vars.ResourceVersion == "" {
		return nil, fmt.Errorf("failed to determine required registerCRDVars")
	}

	registerCrdTemplate, err := readTemplateFile(registerCRDTemplateFile)
	if err != nil {
		return nil, err
	}
	return renderTemplate(registerCrdTemplate, vars)
}

type typesVars struct {
	Imports  string
	CRDTypes string
	Structs  string
	Types    string
}

func RenderTypesTemplate(pkg parser.Package) (*bytes.Buffer, error) {
	var vars typesVars
	vars.CRDTypes = parsePackageCRDs(pkg)
	vars.Structs = parsePackageStructs(pkg)
	vars.Types = parsePackageTypes(pkg)
	vars.Imports = parsePackageImports(pkg)

	registerCrdTemplate, err := readTemplateFile(typesTemplateFile)
	if err != nil {
		return nil, err
	}

	tmpl, err := renderTemplate(registerCrdTemplate, vars)
	if err != nil {
		return nil, err
	}

	out, err := imports.Process("render.go", tmpl.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(out), nil
}

type crdBaseVars struct {
	GroupName       string
	Singular        string
	Plural          string
	Kind            string
	KindList        string
	ResourceVersion string
}

type CrdBaseFile struct {
	Name string
	File *bytes.Buffer
}

func RenderCRDBaseTemplate(baseGroupName string, pkg parser.Package) ([]CrdBaseFile, error) {
	var crds []CrdBaseFile
	for _, node := range pkg.GetNexusNodes() {
		typeName := parser.GetTypeName(node)
		groupName := pkg.Name + "." + baseGroupName
		singular := strings.ToLower(typeName)
		kind := strings.Title(singular)
		var plural string
		if singular[len(singular)-1:] == "s" {
			plural = fmt.Sprintf("%ses", singular)
		} else if singular[len(singular)-1:] == "y" {
			plural = fmt.Sprintf("%sies", singular[:len(singular)-1])
		} else {
			plural = fmt.Sprintf("%ss", singular)
		}

		vars := crdBaseVars{
			GroupName: groupName,
			Singular:  singular,
			Plural:    plural,
			Kind:      kind,
			KindList:  fmt.Sprintf("%sList", kind),
			// TODO make configurable by some variable in package
			ResourceVersion: "v1",
		}

		if vars.GroupName == "" ||
			vars.Singular == "" ||
			vars.Plural == "" ||
			vars.Kind == "" ||
			vars.KindList == "" ||
			vars.ResourceVersion == "" {
			return nil, fmt.Errorf("failed to determine required registerCRDVars")
		}
		registerCrdTemplate, err := readTemplateFile(crdBaseTemplateFile)
		if err != nil {
			return nil, err
		}
		file, err := renderTemplate(registerCrdTemplate, vars)
		if err != nil {
			return nil, err
		}
		crd := CrdBaseFile{
			Name: pkg.Name + "_" + singular + ".yaml",
			File: file,
		}
		crds = append(crds, crd)
	}
	return crds, nil
}

func createApiNamesFile(apiList []string, outputDir string) error {
	sort.Strings(apiList)
	apiNames := "API_NAMES=\""
	for _, name := range apiList {
		apiNames += name + " "
	}
	apiNames += "\""
	var b bytes.Buffer
	b.WriteString(apiNames)
	return createFile(outputDir, "api_names.sh", &b, false)
}
