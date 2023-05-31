package preparser

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Preparser tests", func() {
	It("parse datamodel", func() {
		pkgs := Parse(exampleDSLPath)
		Expect(pkgs["global"]).ToNot(BeNil())
	})

	It("should render pkgs", func() {
		temp, err := os.MkdirTemp("", "compiler-tests")
		defer os.RemoveAll(temp)
		Expect(err).ToNot(HaveOccurred())

		cmd := exec.Command("cp", "-r", exampleDSLPath, temp)
		err = cmd.Run()
		Expect(err).ToNot(HaveOccurred())

		dslDir := temp + "/global-package"
		pkgs := Parse(dslDir)
		Expect(pkgs["global"]).ToNot(BeNil())

		err = Render(dslDir, pkgs)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should render imports", func() {
		temp, err := os.MkdirTemp("", "compiler-tests")
		defer os.RemoveAll(temp)
		Expect(err).ToNot(HaveOccurred())

		err = os.Mkdir(temp+"/model", os.ModePerm)
		Expect(err).ToNot(HaveOccurred())

		pkgs := Parse(exampleDSLPath)
		Expect(pkgs["global"]).ToNot(BeNil())

		err = RenderImports(pkgs, temp, "github.com/vmware/graph-framework-for-microservices/test")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should copy pkgs to build", func() {
		temp, err := os.MkdirTemp("", "compiler-tests")
		defer os.RemoveAll(temp)
		Expect(err).ToNot(HaveOccurred())

		err = os.Mkdir(temp+"/model", os.ModePerm)
		Expect(err).ToNot(HaveOccurred())

		pkgs := Parse(exampleDSLPath)
		Expect(pkgs["global"]).ToNot(BeNil())

		err = CopyPkgsToBuild(pkgs, temp)
		Expect(err).ToNot(HaveOccurred())
	})
})
