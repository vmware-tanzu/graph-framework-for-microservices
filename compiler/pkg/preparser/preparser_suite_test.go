package preparser

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/pkg/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	exampleDSLPath = "../../example/test-utils/global-package"
	baseGroupName  = "tsm.tanzu.vmware.com"
	crdModulePath  = "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/output/generated/"
)

func TestPreparser(t *testing.T) {
	log.StandardLogger().ExitFunc = nil
	RegisterFailHandler(Fail)
	RunSpecs(t, "Preparser Suite")
}

func init() {
	conf := &config.Config{
		GroupName:     baseGroupName,
		CrdModulePath: crdModulePath,
	}
	config.ConfigInstance = conf
}
