package helper

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/elliotchance/orderedmap"

	datamodel "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DEFAULT_KEY = "default"

func GetCRDParentsMap() map[string][]string {
	return map[string][]string{
		"accesscontrolpolicies.policy.tsm.tanzu.vmware.com": {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com", "gnses.gns.tsm.tanzu.vmware.com"},
		"acpconfigs.policy.tsm.tanzu.vmware.com":            {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com", "gnses.gns.tsm.tanzu.vmware.com", "accesscontrolpolicies.policy.tsm.tanzu.vmware.com"},
		"configs.config.tsm.tanzu.vmware.com":               {"roots.root.tsm.tanzu.vmware.com"},
		"gnses.gns.tsm.tanzu.vmware.com":                    {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com"},
		"roots.root.tsm.tanzu.vmware.com":                   {},
		"svcgroups.servicegroup.tsm.tanzu.vmware.com":       {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com", "gnses.gns.tsm.tanzu.vmware.com"},
	}
}

func GetObjectByCRDName(dmClient *datamodel.Clientset, crdName string, name string) interface{} {
	if crdName == "accesscontrolpolicies.policy.tsm.tanzu.vmware.com" {
		obj, err := dmClient.PolicyTsmV1().AccessControlPolicies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "acpconfigs.policy.tsm.tanzu.vmware.com" {
		obj, err := dmClient.PolicyTsmV1().ACPConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "configs.config.tsm.tanzu.vmware.com" {
		obj, err := dmClient.ConfigTsmV1().Configs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gnses.gns.tsm.tanzu.vmware.com" {
		obj, err := dmClient.GnsTsmV1().Gnses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "roots.root.tsm.tanzu.vmware.com" {
		obj, err := dmClient.RootTsmV1().Roots().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "svcgroups.servicegroup.tsm.tanzu.vmware.com" {
		obj, err := dmClient.ServicegroupTsmV1().SvcGroups().Get(context.TODO(), name, metav1.GetOptions{})
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
	h.Write([]byte(output))
	return hex.EncodeToString(h.Sum(nil))
}