package validate

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nexus/validation/pkg/nexus/generated/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Validate(meta metav1.ObjectMeta, group string, labels map[string]string) (bool, string, error) {
	log.Printf("group: %s, labels: %s", group, labels)
	parents := helper.GetCrdParentsMap()[group]
	log.Printf("parents: %v", parents)

	for _, parent := range parents {
		if label, ok := labels[parent]; ok {
			log.Infof("label %s found, val: %s", parent, label)
		} else {
			log.Warnf("label %s not found", parent)
			return false, fmt.Sprintf("Required label %s not found", parent), nil
		}
	}

	return true, "Validation successful", nil
}
