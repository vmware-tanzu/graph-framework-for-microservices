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
const DISPLAY_NAME_LABEL = "nexus/display_name"
const IS_NAME_HASHED_LABEL = "nexus/is_name_hashed"

func GetCRDParentsMap() map[string][]string {
	return map[string][]string{
		"apigateways.apigateway.nexus.org":     {"nexuses.api.nexus.org", "configs.config.nexus.org"},
		"configs.config.nexus.org":             {"nexuses.api.nexus.org"},
		"connects.connect.nexus.org":           {"nexuses.api.nexus.org", "configs.config.nexus.org"},
		"nexusendpoints.connect.nexus.org":     {"nexuses.api.nexus.org", "configs.config.nexus.org", "connects.connect.nexus.org"},
		"nexuses.api.nexus.org":                {},
		"oidcs.authentication.nexus.org":       {"nexuses.api.nexus.org", "configs.config.nexus.org", "apigateways.apigateway.nexus.org"},
		"replicationconfigs.connect.nexus.org": {"nexuses.api.nexus.org", "configs.config.nexus.org", "connects.connect.nexus.org"},
		"routes.route.nexus.org":               {"nexuses.api.nexus.org", "configs.config.nexus.org"},
	}
}

func GetObjectByCRDName(dmClient *datamodel.Clientset, crdName string, name string) interface{} {
	if crdName == "apigateways.apigateway.nexus.org" {
		obj, err := dmClient.ApigatewayNexusV1().ApiGateways().Get(context.TODO(), name, metav1.GetOptions{})
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
	if crdName == "connects.connect.nexus.org" {
		obj, err := dmClient.ConnectNexusV1().Connects().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nexusendpoints.connect.nexus.org" {
		obj, err := dmClient.ConnectNexusV1().NexusEndpoints().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nexuses.api.nexus.org" {
		obj, err := dmClient.ApiNexusV1().Nexuses().Get(context.TODO(), name, metav1.GetOptions{})
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
	if crdName == "replicationconfigs.connect.nexus.org" {
		obj, err := dmClient.ConnectNexusV1().ReplicationConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "routes.route.nexus.org" {
		obj, err := dmClient.RouteNexusV1().Routes().Get(context.TODO(), name, metav1.GetOptions{})
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
