package parser

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"

	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/config"
	"golang.org/x/mod/modfile"
)

func GetModulePath(startPath string) string {
	file, err := ioutil.ReadFile(path.Join(startPath, "go.mod"))
	if err != nil {
		log.Fatalf("failed to get module path %v", err)
	}
	return modfile.ModulePath(file)
}

func ConstructImports(inputAlias, inputImportPath string) (string, string) {
	re, err := regexp.Compile(`[\_\.]`)
	if err != nil {
		log.Fatalf("failed to construct output import path for import path %v : %v", inputImportPath, err)
	}
	aliasName := fmt.Sprintf("%s%sv1", inputAlias, config.ConfigInstance.GroupName)
	aliasName = re.ReplaceAllString(aliasName, "")

	importPath := fmt.Sprintf("\"%sapis/%s.%s/v1\"", config.ConfigInstance.CrdModulePath, re.ReplaceAllString(inputAlias, ""), config.ConfigInstance.GroupName)
	return aliasName, importPath
}

func SpecialCharsPresent(name string) bool {
	re, err := regexp.Compile(`[^a-z0-9]`)
	if err != nil {
		log.Fatalf("failed to check for special characters in the package name %v : %v", name, err)
	}
	return re.MatchString(name)
}
