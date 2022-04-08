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

func GetCRDParentsMap() map[string][]string {
	return map[string][]string{
		"accesscontrolpolicies.policy.tsm.tanzu.vmware.com": {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com"},
		"acpconfigs.policy.tsm.tanzu.vmware.com":            {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com", "accesscontrolpolicies.policy.tsm.tanzu.vmware.com"},
		"configs.config.tsm.tanzu.vmware.com":               {"roots.root.tsm.tanzu.vmware.com"},
		"gnses.gns.tsm.tanzu.vmware.com":                    {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com"},
		"roots.root.tsm.tanzu.vmware.com":                   {},
		"svcgroups.service_group.tsm.tanzu.vmware.com":      {"roots.root.tsm.tanzu.vmware.com", "configs.config.tsm.tanzu.vmware.com", "gnses.gns.tsm.tanzu.vmware.com"},
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
	if crdName == "svcgroups.service_group.tsm.tanzu.vmware.com" {
		obj, err := dmClient.Service_groupTsmV1().SvcGroups().Get(context.TODO(), name, metav1.GetOptions{})
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

func GenerateCRDName(crdName string, meta metav1.ObjectMeta) string {
	labels := ParseCRDLabels(crdName, meta.Labels)

	var name string
	for i, key := range labels.Keys() {
		value, _ := labels.Get(key)

		name += fmt.Sprintf("%s:%s", key, value)
		if i < labels.Len()-1 {
			name += "/"
		}
	}

	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}
