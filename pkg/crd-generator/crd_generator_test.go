package crd_generator_test

import (
	"go/format"
	"io/ioutil"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	crdgenerator "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/crd-generator"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

const (
	baseGroupName            = "tsm.tanzu.vmware.com"
	crdModulePath            = "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/"
	examplePath              = "../../example/"
	exampleDSLPath           = examplePath + "datamodel"
	exampleCRDOutputPath     = examplePath + "output/_crd_base/"
	exampleCRDApisOutputPath = exampleCRDOutputPath + "apis"
	gnsExamplePath           = exampleCRDApisOutputPath + "/gns.tsm.tanzu.vmware.com/"
	gnsDocPath               = gnsExamplePath + "v1/doc.go"
	gnsRegisterGroupPath     = gnsExamplePath + "register.go"
	gnsRegisterCRDPath       = gnsExamplePath + "v1/register.go"
	gnsTypesPath             = gnsExamplePath + "v1/types.go"
	gnsCrdBasePath           = exampleCRDOutputPath + "crds/gns_gns.yaml"
)

var _ = Describe("Template renderers tests", func() {
	var (
		//err error
		pkg        parser.Package
		parentsMap map[string]parser.NodeHelper
		ok         bool
		methods    map[string]nexus.HTTPMethodsResponses
		codes      map[string]nexus.HTTPCodesResponse
	)

	BeforeEach(func() {
		pkgs := parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel//config/gns"]
		Expect(ok).To(BeTrue())

		graph := parser.ParseDSLNodes(exampleDSLPath, baseGroupName)
		parentsMap = parser.CreateParentsMap(graph)
		Expect(parentsMap).To(HaveLen(6))

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
		Expect(files).To(HaveLen(2))

		expectedSdk, err := ioutil.ReadFile(gnsCrdBasePath)
		Expect(err).NotTo(HaveOccurred())

		Expect("gns_gns.yaml").To(Or(Equal(files[0].Name)), Equal(files[1].Name))
		Expect(string(expectedSdk)).To(Or(Equal(files[0].File.String()), Equal(files[1].File.String())))
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
})
