package helper

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/elliotchance/orderedmap"

	datamodel "vmware/build/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DEFAULT_KEY = "default"
const DISPLAY_NAME_LABEL = "nexus/display_name"
const IS_NAME_HASHED_LABEL = "nexus/is_name_hashed"

func GetCRDParentsMap() map[string][]string {
	return map[string][]string{
		"devs.engineering.vmware.org":         {"roots.orgchart.vmware.org", "leaders.management.vmware.org", "mgrs.management.vmware.org"},
		"employees.role.vmware.org":           {"roots.orgchart.vmware.org"},
		"executives.role.vmware.org":          {"roots.orgchart.vmware.org"},
		"humanresourceses.hr.vmware.org":      {"roots.orgchart.vmware.org", "leaders.management.vmware.org"},
		"leaders.management.vmware.org":       {"roots.orgchart.vmware.org"},
		"mgrs.management.vmware.org":          {"roots.orgchart.vmware.org", "leaders.management.vmware.org"},
		"operationses.engineering.vmware.org": {"roots.orgchart.vmware.org", "leaders.management.vmware.org", "mgrs.management.vmware.org"},
		"roots.orgchart.vmware.org":           {},
	}
}

func GetObjectByCRDName(dmClient *datamodel.Clientset, crdName string, name string) interface{} {
	if crdName == "devs.engineering.vmware.org" {
		obj, err := dmClient.EngineeringVmwareV1().Devs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "employees.role.vmware.org" {
		obj, err := dmClient.RoleVmwareV1().Employees().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "executives.role.vmware.org" {
		obj, err := dmClient.RoleVmwareV1().Executives().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "humanresourceses.hr.vmware.org" {
		obj, err := dmClient.HrVmwareV1().HumanResourceses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "leaders.management.vmware.org" {
		obj, err := dmClient.ManagementVmwareV1().Leaders().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "mgrs.management.vmware.org" {
		obj, err := dmClient.ManagementVmwareV1().Mgrs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "operationses.engineering.vmware.org" {
		obj, err := dmClient.EngineeringVmwareV1().Operationses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "roots.orgchart.vmware.org" {
		obj, err := dmClient.OrgchartVmwareV1().Roots().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}

	return nil
}

func ParseCRDLabels(crdName string, labels map[string]string) *orderedmap.OrderedMap {
	parents := GetCRDParentsMap()[crdName]

	m := orderedmap.NewOrderedMap()
	for _, parent := range parents {
		if label, ok := labels[parent]; ok {
			m.Set(parent, label)
		} else {
			m.Set(parent, DEFAULT_KEY)
		}
	}

	return m
}

func GetHashedName(crdName string, labels map[string]string, name string) string {
	orderedLabels := ParseCRDLabels(crdName, labels)

	var output string
	for i, key := range orderedLabels.Keys() {
		value, _ := orderedLabels.Get(key)

		output += fmt.Sprintf("%s:%s", key, value)
		if i < orderedLabels.Len()-1 {
			output += "/"
		}
	}

	output += fmt.Sprintf("%s:%s", crdName, name)

	h := sha1.New()
	_, _ = h.Write([]byte(output))
	return hex.EncodeToString(h.Sum(nil))
}
