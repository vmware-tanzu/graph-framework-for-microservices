package helper

import (
	"context"
	log "github.com/sirupsen/logrus"
	datamodel "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/nexus/generated/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetCrdParentsMap() map[string][]string {
	return map[string][]string{
		"roots.root.helloworld.com":     {""},
		"configs.config.helloworld.com": {"roots.root.helloworld.com"},
	}
}

func CheckIfObjectExist(dmClient *datamodel.Clientset, crdName string, name string) bool {
	if crdName == "roots.root.helloworld.com" {
		_, err := dmClient.RootHelloworldV1().Roots("default").Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			log.Printf("roots not found %v", err)
			return false
		}
		return true
	}

	if crdName == "configs.config.helloworld.com" {
		_, err := dmClient.ConfigHelloworldV1().Configs("default").Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	}

	return false
}
