package controllers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/apps/v1"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

func getEnv(deployList *v1.DeploymentList) string {
	containers := deployList.Items[0].Spec.Template.Spec.Containers
	for _, val := range containers {
		for _, v := range val.Env {
			if v.Name == "REMOTE_ENDPOINT_PORT" {
				return v.Value
			}
		}
	}
	return ""
}
