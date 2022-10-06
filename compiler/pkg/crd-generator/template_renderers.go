package crd_generator

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"golang.org/x/tools/imports"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser/rest"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/util"
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

//go:embed template/helper.go.tmpl
var helperTemplateFile []byte

//go:embed template/client.go.tmpl
var clientTemplateFile []byte

func RenderCRDTemplate(baseGroupName, crdModulePath string,
	pkgs parser.Packages, graph map[string]parser.Node, outputDir string,
	httpMethods map[string]nexus.HTTPMethodsResponses, httpCodes map[string]nexus.HTTPCodesResponse) error {
	parentsMap := parser.CreateParentsMap(graph)

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
		err = createFolder(apiFolder)
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
		file, err = RenderTypesTemplate(crdModulePath, pkg)
		if err != nil {
			return err
		}
		log.Debugf("Rendered types template for package %s: %s", pkg.Name, file)
		err = createFile(apiFolder, "types.go", file, true)
		if err != nil {
			return err
		}
		crdFiles, err := RenderCRDBaseTemplate(baseGroupName, pkg, parentsMap, httpMethods, httpCodes)
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

	err := RenderHelper(parentsMap, outputDir, crdModulePath)
	if err != nil {
		return err
	}

	err = RenderClient(baseGroupName, outputDir, crdModulePath, pkgs, parentsMap)
	if err != nil {
		return err
	}
	err = createApiNamesFile(pkgNames, outputDir)
	if err != nil {
		return err
	}

	return nil
}

func RenderHelper(parentsMap map[string]parser.NodeHelper, outputDir string, crdModulePath string) error {
	helperFolder := outputDir + "/helper"
	var err error
	err = createFolder(helperFolder)
	if err != nil {
		return err
	}

	file, err := RenderHelperTemplate(parentsMap, crdModulePath)
	if err != nil {
		return err
	}
	log.Debugf("Rendered helper: %s", file)
	err = createFile(helperFolder, "helper.go", file, true)
	if err != nil {
		return err
	}

	return nil
}

func createFolder(name string) error {
	err := os.MkdirAll(name, os.ModeDir|os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating dir %v failed with an error: %v", name, err)
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
	Imports      string
	CommonImport string
	CRDTypes     string
	Structs      string
	Types        string
	Consts       string
}

func RenderTypesTemplate(crdModulePath string, pkg parser.Package) (*bytes.Buffer, error) {
	aliasNameMap := make(map[string]string)
	var vars typesVars
	vars.Imports = parsePackageImports(pkg, aliasNameMap)
	vars.CRDTypes = parsePackageCRDs(pkg, aliasNameMap)
	vars.Structs = parsePackageStructs(pkg, aliasNameMap)
	vars.Types = parsePackageTypes(pkg)
	vars.Consts = parsePackageConsts(pkg)
	vars.CommonImport = util.GetInternalImport(crdModulePath, "common")

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
	CrdName         string
	GroupName       string
	Singular        string
	Plural          string
	Kind            string
	KindList        string
	ResourceVersion string
	NexusAnnotation string
}

type NexusAnnotation struct {
	Name            string                            `json:"name,omitempty"`
	Hierarchy       []string                          `json:"hierarchy,omitempty"`
	Children        map[string]parser.NodeHelperChild `json:"children,omitempty"`
	Links           map[string]parser.NodeHelperChild `json:"links,omitempty"`
	IsSingleton     bool                              `json:"is_singleton"`
	NexusRestAPIGen nexus.RestAPISpec                 `json:"nexus-rest-api-gen,omitempty"`
	Description     string                            `json:"description,omitempty"`
}

type CrdBaseFile struct {
	Name string
	File *bytes.Buffer
}

func RenderCRDBaseTemplate(baseGroupName string, pkg parser.Package, parentsMap map[string]parser.NodeHelper,
	httpMethods map[string]nexus.HTTPMethodsResponses, httpCodes map[string]nexus.HTTPCodesResponse) ([]CrdBaseFile, error) {
	var crds []CrdBaseFile

	restAPISpecMap := rest.GetRestApiSpecs(pkg, httpMethods, httpCodes, parentsMap)
	for _, node := range pkg.GetNexusNodes() {
		typeName := parser.GetTypeName(node)
		groupName := pkg.Name + "." + baseGroupName
		singular := strings.ToLower(typeName)
		kind := strings.Title(typeName)
		plural := util.ToPlural(singular)
		crdName := fmt.Sprintf("%s.%s", plural, groupName)

		nexusAnnotation := &NexusAnnotation{}
		nexusAnnotation.IsSingleton = parser.IsSingletonNode(node)

		var err error
		parents, ok := parentsMap[crdName]
		if ok {
			nexusAnnotation.Hierarchy = parents.Parents
			nexusAnnotation.Children = parents.Children
			nexusAnnotation.Links = parents.Links
			nexusAnnotation.Name = parents.RestName
		}

		if annotation, ok := parser.GetNexusRestAPIGenAnnotation(pkg, typeName); ok {
			nexusAnnotation.NexusRestAPIGen = restAPISpecMap[annotation]
			rest.ValidateRestApiSpec(restAPISpecMap[annotation], parentsMap, crdName)
		}
		if annotation, ok := parser.GetNexusDescriptionAnnotation(pkg, typeName); ok {
			nexusAnnotation.Description = annotation
		}

		nexusAnnotationStr, err := json.Marshal(nexusAnnotation)
		if err != nil {
			return nil, err
		}

		vars := crdBaseVars{
			CrdName:         crdName,
			GroupName:       groupName,
			Singular:        singular,
			Plural:          plural,
			Kind:            kind,
			KindList:        fmt.Sprintf("%sList", kind),
			NexusAnnotation: string(nexusAnnotationStr),
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

type helperVars struct {
	CrdModulePath      string
	GetCrdParentsMap   string
	GetObjectByCRDName string
}

func RenderHelperTemplate(parentsMap map[string]parser.NodeHelper, crdModulePath string) (*bytes.Buffer, error) {
	keys := make([]string, 0, len(parentsMap))
	for k := range parentsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var vars helperVars
	vars.CrdModulePath = strings.TrimSuffix(crdModulePath, "/")
	vars.GetCrdParentsMap = generateGetCrdParentsMap(keys, parentsMap)
	vars.GetObjectByCRDName = generateGetObjectByCRDName(keys, parentsMap)

	helperTemplate, err := readTemplateFile(helperTemplateFile)
	if err != nil {
		return nil, err
	}

	tmpl, err := renderTemplate(helperTemplate, vars)
	if err != nil {
		return nil, err
	}

	out, err := imports.Process("render.go", tmpl.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(out), nil
}

func RenderClient(baseGroupName, outputDir, crdModulePath string, pkgs parser.Packages, parentsMap map[string]parser.NodeHelper) error {
	clientFolder := outputDir + "/nexus-client"
	file, err := RenderClientTemplate(baseGroupName, crdModulePath, pkgs, parentsMap)
	if err != nil {
		return err
	}
	log.Debugf("Rendered client template: %s", file)
	err = createFile(clientFolder, "client.go", file, true)
	if err != nil {
		return err
	}

	return nil
}

type clientVars struct {
	ApiGroups              []ApiGroupsVars
	CommonImport           string
	HelperImport           string
	BaseClientsetImport    string
	FakeBaseCliensetImport string
	BaseImports            string
	ApiGroupsClient        string
	Nodes                  []apiGroupsClientVars
}

func RenderClientTemplate(baseGroupName, crdModulePath string, pkgs parser.Packages, parentsMap map[string]parser.NodeHelper) (*bytes.Buffer, error) {
	vars, err := generateNexusClientVars(baseGroupName, crdModulePath, pkgs, parentsMap)
	if err != nil {
		return nil, err
	}

	clientTemplate, err := readTemplateFile(clientTemplateFile)
	if err != nil {
		return nil, err
	}

	tmpl, err := renderTemplate(clientTemplate, vars)
	if err != nil {
		return nil, err
	}

	out, err := imports.Process("client.go", tmpl.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(out), nil
}
