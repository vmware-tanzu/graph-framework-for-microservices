package helper

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/elliotchance/orderedmap"

	datamodel "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DEFAULT_KEY = "default"

func GetCRDParentsMap() map[string][]string {
	return map[string][]string{
		"apis.apis.nexus.org":             {},
		"configs.config.nexus.org":        {"apis.apis.nexus.org"},
		"extensions.extensions.nexus.org": {"apis.apis.nexus.org", "configs.config.nexus.org"},
		"gateways.gateway.nexus.org":      {"apis.apis.nexus.org", "configs.config.nexus.org"},
		"oidcs.authentication.nexus.org":  {"apis.apis.nexus.org", "configs.config.nexus.org"},
	}
}

func GetObjectByCRDName(dmClient *datamodel.Clientset, crdName string, name string) interface{} {
	if crdName == "apis.apis.nexus.org" {
		obj, err := dmClient.ApisNexusV1().Apis().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "configs.config.nexus.org" {
		obj, err := dmClient.ConfigNexusV1().Configs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "extensions.extensions.nexus.org" {
		obj, err := dmClient.ExtensionsNexusV1().Extensions().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "gateways.gateway.nexus.org" {
		obj, err := dmClient.GatewayNexusV1().Gateways().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "oidcs.authentication.nexus.org" {
		obj, err := dmClient.AuthenticationNexusV1().OIDCs().Get(context.TODO(), name, metav1.GetOptions{})
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

func GetHashedName(crdName string, labels map[string]string) string {
	orderedLabels := ParseCRDLabels(crdName, labels)

	var name string
	for i, key := range orderedLabels.Keys() {
		value, _ := orderedLabels.Get(key)

		name += fmt.Sprintf("%s:%s", key, value)
		if i < orderedLabels.Len()-1 {
			name += "/"
		}
	}

	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}
