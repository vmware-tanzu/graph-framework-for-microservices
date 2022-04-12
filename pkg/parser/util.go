package parser

import (
	"io/ioutil"
	"path"
	"regexp"

	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/modfile"
)

func GetModulePath(startPath string) string {
	file, err := ioutil.ReadFile(path.Join(startPath, "go.mod"))
	if err != nil {
		log.Fatalf("failed to get module path %v", err)
	}
	return modfile.ModulePath(file)
}

func SpecialCharsPresent(name string) bool {
	re, err := regexp.Compile(`[^a-z0-9]`)
	if err != nil {
		log.Fatalf("failed to check for special characters in the package name %v : %v", name, err)
	}
	return re.MatchString(name)
}
