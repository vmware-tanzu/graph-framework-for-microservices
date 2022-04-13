package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

//TODO: Move this to COMMON nexus repo

// HTTPMethod type.
type HTTPMethod string

// ResponseCode type.
type ResponseCode int

// HTTPResponse type.
type HTTPResponse struct {
	Description string `json:"description"`
}

// HTTPCodesResponse code to response type.
type HTTPCodesResponse map[ResponseCode]HTTPResponse

// HTTPMethodsResponses to response mapping.
type HTTPMethodsResponses map[HTTPMethod]HTTPCodesResponse

// RestURIs and associated data.
type RestURIs struct {
	Uri     string               `json:"uri"`
	Methods HTTPMethodsResponses `json:"methods"`
}

type RestAPISpec struct {
	Uris []RestURIs `json:"uris"`
}

const DefaultHTTPErrorCode ResponseCode = http.StatusNotImplemented

var DefaultHTTPError = HTTPResponse{Description: http.StatusText(http.StatusNotImplemented)}
