package parser_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	examplePath    = "../../example/"
	exampleDSLPath = examplePath + "datamodel"
)

func TestParser(t *testing.T) {
	log.StandardLogger().ExitFunc = nil
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}

func init() {
	conf := &config.Config{
		GroupName:     baseGroupName,
		CrdModulePath: crdModulePath,
	}
	config.ConfigInstance = conf
}
