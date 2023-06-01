package api_test

import (
	"testing"

	log "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	Uri         = "/v1alpha1/project/{projectId}/global-namespaces"
	ResourceUri = "/v1alpha1/project/{projectId}/global-namespaces/{id}"
	ListUri     = "/v1alpha1/global-namespaces/test"
)

func TestApi(t *testing.T) {
	log.StandardLogger().ExitFunc = nil
	RegisterFailHandler(Fail)
	RunSpecs(t, "Declarative Suite")
}
