package validate

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	datamodel "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/nexus/generated/client/clientset/versioned"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/validation.git/pkg/nexus/generated/helper"
)

func Validate(dmClient *datamodel.Clientset, resourceName string, labels map[string]string) (bool, string, error) {
	log.Printf("resourceName: %s, labels: %s", resourceName, labels)
	parents := helper.GetCrdParentsMap()[resourceName]
	log.Printf("parents: %v", parents)

	for _, parent := range parents {
		if label, ok := labels[parent]; ok {
			log.Infof("label %s found, val: %s", parent, label)
			if !helper.CheckIfObjectExist(dmClient, parent, label) {
				log.Warnf("required parent %s with name %s not found", parent, label)
				return false, fmt.Sprintf("required parent %s with name %s not found", parent, label), nil
			}
		} else {
			log.Warnf("label %s not found", parent)
			if !helper.CheckIfObjectExist(dmClient, parent, "default") {
				log.Warnf("required parent %s with name default not found", parent)
				return false, fmt.Sprintf("Required parent %s with name default not found", parent), nil
			}
		}
	}

	return true, "Validation successful", nil
}
