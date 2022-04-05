package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/parser"
)

var _ = Describe("Pkg tests", func() {
	var (
		pkgs   map[string]parser.Package
		pkg    parser.Package
		gnsPkg parser.Package
		ok     bool
	)

	BeforeEach(func() {
		_, err := config.LoadConfig("../../example/nexus-sdk.yaml")
		Expect(err).To(Not(HaveOccurred()))

		pkgs = parser.ParseDSLPkg(exampleDSLPath)
		pkg, ok = pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel"]
		Expect(ok).To(BeTrue())
		gnsPkg, ok = pkgs["gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel//config/gns"]
		Expect(ok).To(BeTrue())
	})

	It("should return import strings", func() {
		imports := pkg.GetImportStrings()

		expectedImports := []string{
			"\"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config\"",
			"\"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus\""}

		Expect(imports).To(Equal(expectedImports))
	})

	It("should check if node is nexus.Node", func() {
		nodes := pkg.GetNexusNodes()
		Expect(parser.IsNexusNode(nodes[0])).To(BeTrue())
	})

	It("should get all structs for gns", func() {
		structs := gnsPkg.GetStructs()
		Expect(structs).To(HaveLen(4))
	})

	It("should get all types for gns", func() {
		types := gnsPkg.GetTypes()
		Expect(types).To(HaveLen(0))
	})

	It("should get child fields", func() {
		nodes := pkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		Expect(childFields).To(HaveLen(1))
	})

	It("should get link fields for gns", func() {
		nodes := gnsPkg.GetNexusNodes()
		linkFields := parser.GetChildFields(nodes[0])
		Expect(linkFields).To(HaveLen(2))
	})

	It("should get spec fields for gns", func() {
		nodes := gnsPkg.GetNexusNodes()
		specFields := parser.GetSpecFields(nodes[0])
		Expect(specFields).To(HaveLen(3))
	})

	It("should get field name", func() {
		nodes := pkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		fieldName, err := parser.GetFieldName(childFields[0])
		Expect(err).NotTo(HaveOccurred())
		Expect(fieldName).To(Equal("Config"))
	})

	It("should get field type", func() {
		nodes := pkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		fieldType := parser.GetFieldType(childFields[0])
		Expect(fieldType).To(Equal("config.Config"))
	})

	It("should check if field is map", func() {
		nodes := gnsPkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		isMap := parser.IsMapField(childFields[0])
		Expect(isMap).To(BeTrue())
	})

	It("should get field type for MapType", func() {
		nodes := gnsPkg.GetNexusNodes()
		childFields := parser.GetChildFields(nodes[0])
		fieldType := parser.GetFieldType(childFields[0])
		Expect(fieldType).To(Equal("service_group.SvcGroup"))
	})

	It("should fail if wrong struct tag is given", func() {
		defer func() { log.StandardLogger().ExitFunc = nil }()

		fail := false
		log.StandardLogger().ExitFunc = func(int) {
			fail = true
		}

		parser.ParseFieldTags("`nexus: \"child\"`")
		Expect(fail).To(BeTrue())
	})
})
