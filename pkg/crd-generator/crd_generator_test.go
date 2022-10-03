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
	baseGroupName                = "tsm.tanzu.vmware.com"
	crdModulePath                = "nexustempmodule/"
	examplePath                  = "../../example/"
	exampleDSLPath               = examplePath + "datamodel"
	exampleCRDOutputPath         = examplePath + "output/_crd_base/"
	exampleCRDApisOutputPath     = exampleCRDOutputPath + "apis"
	gnsExamplePath               = exampleCRDApisOutputPath + "/gns.tsm.tanzu.vmware.com/"
	gnsDocPath                   = gnsExamplePath + "v1/doc.go"
	gnsRegisterGroupPath         = gnsExamplePath + "register.go"
	gnsRegisterCRDPath           = gnsExamplePath + "v1/register.go"
	gnsTypesPath                 = gnsExamplePath + "v1/types.go"
	gnsCrdBasePath               = exampleCRDOutputPath + "crds/gns_gns.yaml"
	exampleCRDClientOutputPath   = exampleCRDOutputPath + "/nexus-client/client.go"
	exampleCRDHelperOutputPath   = exampleCRDOutputPath + "/helper/helper.go"
	gqlgenOutputPath             = exampleCRDOutputPath + "/nexus-gql/gqlgen.yml"
	graphqlServerOutputPath      = exampleCRDOutputPath + "/nexus-gql/server.go"
	gglToolOutputPath            = exampleCRDOutputPath + "/nexus-gql/tools.go"
	gglGraphqlResolverOutputPath = exampleCRDOutputPath + "/nexus-gql/graph/graphqlResolver.go"
	gglResolverOutputPath        = exampleCRDOutputPath + "/nexus-gql/graph/resolver.go"
	gglSchemaGraphqlOutputPath   = exampleCRDOutputPath + "/nexus-gql/graph/schema.graphqls"
	gglSchemaResolverOutputPath  = exampleCRDOutputPath + "/nexus-gql/graph/schema.resolvers.go"
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
		gqlvar     crdgenerator.GraphDetails
	)

	BeforeEach(func() {
		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns"]
		Expect(ok).To(BeTrue())

		graph := parser.ParseDSLNodes(exampleDSLPath, baseGroupName)
		parentsMap = parser.CreateParentsMap(graph)
		Expect(parentsMap).To(HaveLen(10))

		methods, codes = rest.ParseResponses(pkgs)
		gqlvar.BaseImportPath = crdModulePath
		var err error
		gqlvar.Nodes, err = crdgenerator.GenerateGraphqlResolverVars(baseGroupName, crdModulePath, pkgs, parentsMap)
		Expect(err).NotTo(HaveOccurred())
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
		Expect(files).To(HaveLen(5))

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
	// Path:"example/output/_crd_base/nexus-gql/tools.go"
	It("should parse nexus-gql tool template", func() {
		qglBytes, err := crdgenerator.RenderGraphqlToolTemplate()
		Expect(err).NotTo(HaveOccurred())

		formatted, err := format.Source(qglBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedTypes, err := ioutil.ReadFile(gglToolOutputPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedTypes)))
	})
	// Path:"example/output/_crd_base/nexus-gql/graph/schema.resolvers.go"
	It("should parse graph schema resolver template", func() {
		qglBytes, err := crdgenerator.RenderGqlSchemaResolverTemplate()
		Expect(err).NotTo(HaveOccurred())

		formatted, err := format.Source(qglBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedTypes, err := ioutil.ReadFile(gglSchemaResolverOutputPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedTypes)))
	})
	// Path:"example/output/_crd_base/nexus-gql/graph/resolver.go"
	It("should parse graph resolver template", func() {
		qglBytes, err := crdgenerator.RenderGqlResolverTemplate()
		Expect(err).NotTo(HaveOccurred())

		formatted, err := format.Source(qglBytes.Bytes())
		Expect(err).NotTo(HaveOccurred())

		expectedTypes, err := ioutil.ReadFile(gglResolverOutputPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(formatted)).To(Equal(string(expectedTypes)))
	})
	// Path:"example/output/_crd_base/nexus-gql/server.go"
	It("should parse nexus-gql server template", func() {
		var vars crdgenerator.ServerVars
		vars.BaseImportPath = crdModulePath
		_, err := crdgenerator.RenderGqlServerTemplate(vars)
		Expect(err).NotTo(HaveOccurred())
	})
	// Path:"example/output/_crd_base/nexus-gql/graph/schema.graphqls"
	It("should parse graph schema graphql template", func() {
		_, err := crdgenerator.RenderGraphqlSchemaTemplate(gqlvar, crdModulePath)
		Expect(err).NotTo(HaveOccurred())
	})
	// Path:"example/output/_crd_base/nexus-gql/gqlgen.yml"
	It("should parse nexus-gql gqlgen config template", func() {
		_, err := crdgenerator.RenderGQLGenTemplate(gqlvar, crdModulePath)
		Expect(err).NotTo(HaveOccurred())
	})
	// Path:"example/output/_crd_base/nexus-gql/graph/graphqlResolver.go"
	It("should parse graphql resolver template", func() {
		_, err := crdgenerator.RenderGraphqlResolverTemplate(gqlvar, crdModulePath)
		Expect(err).NotTo(HaveOccurred())
	})
})
