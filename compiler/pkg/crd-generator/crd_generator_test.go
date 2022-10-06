package crd_generator_test

import (
	"go/format"
	"io/ioutil"

	"github.com/vmware-tanzu/graph-framework-for-microservices/common-library/pkg/nexus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	crdgenerator "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/crd-generator"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/parser"
)

const (
	baseGroupName              = "tsm.tanzu.vmware.com"
	crdModulePath              = "nexustempmodule/"
	examplePath                = "../../example/"
	exampleDSLPath             = examplePath + "datamodel"
	exampleCRDOutputPath       = examplePath + "output/_crd_base/"
	exampleCRDApisOutputPath   = exampleCRDOutputPath + "apis"
	gnsExamplePath             = exampleCRDApisOutputPath + "/gns.tsm.tanzu.vmware.com/"
	gnsDocPath                 = gnsExamplePath + "v1/doc.go"
	gnsRegisterGroupPath       = gnsExamplePath + "register.go"
	gnsRegisterCRDPath         = gnsExamplePath + "v1/register.go"
	gnsTypesPath               = gnsExamplePath + "v1/types.go"
	gnsCrdBasePath             = exampleCRDOutputPath + "crds/gns_gns.yaml"
	exampleCRDClientOutputPath = exampleCRDOutputPath + "/nexus-client/client.go"
	exampleCRDHelperOutputPath = exampleCRDOutputPath + "/helper/helper.go"
)

var _ = Describe("Template renderers tests", func() {
	var (
		//err error
		pkgs       parser.Packages
		pkg        parser.Package
		parentsMap map[string]parser.NodeHelper
		ok         bool
		methods    map[string]nexus.HTTPMethodsResponses
		codes      map[string]nexus.HTTPCodesResponse
	)

	BeforeEach(func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"]
		Expect(ok).To(BeTrue())

		graph := parser.ParseDSLNodes(exampleDSLPath, baseGroupName)
		parentsMap = parser.CreateParentsMap(graph)
		Expect(parentsMap).To(HaveLen(8))

		methods, codes = rest.ParseResponses(pkgs)
	})

	It("should parse doc template", func() {
		docBytes, err := crdgenerator.RenderDocTemplate(baseGroupName, pkg)
		Expect(err).NotTo(HaveOccurred())
		formatted, err := format.Source(docBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedDoc, err := ioutil.ReadFile(gnsDocPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedDoc)))
	})

	It("should parse register group template", func() {
		regBytes, err := crdgenerator.RenderRegisterGroupTemplate(baseGroupName, pkg)
		Expect(err).NotTo(HaveOccurred())
		formatted, err := format.Source(regBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedRegisterGroup, err := ioutil.ReadFile(gnsRegisterGroupPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedRegisterGroup)))
	})

	It("should parse register CRD template", func() {
		regBytes, err := crdgenerator.RenderRegisterCRDTemplate(crdModulePath, baseGroupName, pkg)
		Expect(err).NotTo(HaveOccurred())
		formatted, err := format.Source(regBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedRegisterCRD, err := ioutil.ReadFile(gnsRegisterCRDPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedRegisterCRD)))
	})

	It("should parse base crd template", func() {
		files, err := crdgenerator.RenderCRDBaseTemplate(baseGroupName, pkg, parentsMap, methods, codes)
		Expect(err).NotTo(HaveOccurred())
		Expect(files).To(HaveLen(4))

		expectedSdk, err := ioutil.ReadFile(gnsCrdBasePath)
		Expect(err).NotTo(HaveOccurred())

		Expect("gns_gns.yaml").To(Or(Equal(files[0].Name), Equal(files[1].Name), Equal(files[2].Name), Equal(files[3].Name)))
		Expect(string(expectedSdk)).To(Or(Equal(files[0].File.String()), Equal(files[1].File.String()), Equal(files[2].File.String()), Equal(files[3].File.String())))
	})

	It("should parse types template", func() {
		typesBytes, err := crdgenerator.RenderTypesTemplate(crdModulePath, pkg)
		Expect(err).NotTo(HaveOccurred())

		formatted, err := format.Source(typesBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedTypes, err := ioutil.ReadFile(gnsTypesPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedTypes)))
	})

	It("should parse client template", func() {
		clientsBytes, err := crdgenerator.RenderClientTemplate(baseGroupName, crdModulePath, pkgs, parentsMap)
		Expect(err).NotTo(HaveOccurred())

		formatted, err := format.Source(clientsBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedTypes, err := ioutil.ReadFile(exampleCRDClientOutputPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedTypes)))
	})

	It("should parse helper template", func() {
		helperBytes, err := crdgenerator.RenderHelperTemplate(parentsMap, crdModulePath)
		Expect(err).NotTo(HaveOccurred())

		formatted, err := format.Source(helperBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedTypes, err := ioutil.ReadFile(exampleCRDHelperOutputPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedTypes)))
	})
})
