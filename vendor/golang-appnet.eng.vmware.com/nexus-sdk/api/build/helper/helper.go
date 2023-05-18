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
		"apigateways.apigateway.nexus.vmware.com":     {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com"},
		"configs.config.nexus.vmware.com":             {"nexuses.api.nexus.vmware.com"},
		"connects.connect.nexus.vmware.com":           {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com"},
		"corsconfigs.domain.nexus.vmware.com":         {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com", "apigateways.apigateway.nexus.vmware.com"},
		"nexusendpoints.connect.nexus.vmware.com":     {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com", "connects.connect.nexus.vmware.com"},
		"nexuses.api.nexus.vmware.com":                {},
		"oidcs.authentication.nexus.vmware.com":       {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com", "apigateways.apigateway.nexus.vmware.com"},
		"policies.tenantconfig.nexus.vmware.com":      {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com"},
		"proxyrules.admin.nexus.vmware.com":           {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com", "apigateways.apigateway.nexus.vmware.com"},
		"replicationconfigs.connect.nexus.vmware.com": {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com", "connects.connect.nexus.vmware.com"},
		"routes.route.nexus.vmware.com":               {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com"},
		"runtimes.runtime.nexus.vmware.com":           {"nexuses.api.nexus.vmware.com"},
		"tenants.tenantconfig.nexus.vmware.com":       {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com"},
		"tenants.tenantruntime.nexus.vmware.com":      {"nexuses.api.nexus.vmware.com", "runtimes.runtime.nexus.vmware.com"},
		"users.user.nexus.vmware.com":                 {"nexuses.api.nexus.vmware.com", "configs.config.nexus.vmware.com"},
	}
}

func GetObjectByCRDName(dmClient *datamodel.Clientset, crdName string, name string) interface{} {
	if crdName == "apigateways.apigateway.nexus.vmware.com" {
		obj, err := dmClient.ApigatewayNexusV1().ApiGateways().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "configs.config.nexus.vmware.com" {
		obj, err := dmClient.ConfigNexusV1().Configs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "connects.connect.nexus.vmware.com" {
		obj, err := dmClient.ConnectNexusV1().Connects().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "corsconfigs.domain.nexus.vmware.com" {
		obj, err := dmClient.DomainNexusV1().CORSConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nexusendpoints.connect.nexus.vmware.com" {
		obj, err := dmClient.ConnectNexusV1().NexusEndpoints().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "nexuses.api.nexus.vmware.com" {
		obj, err := dmClient.ApiNexusV1().Nexuses().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "oidcs.authentication.nexus.vmware.com" {
		obj, err := dmClient.AuthenticationNexusV1().OIDCs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "policies.tenantconfig.nexus.vmware.com" {
		obj, err := dmClient.TenantconfigNexusV1().Policies().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "proxyrules.admin.nexus.vmware.com" {
		obj, err := dmClient.AdminNexusV1().ProxyRules().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "replicationconfigs.connect.nexus.vmware.com" {
		obj, err := dmClient.ConnectNexusV1().ReplicationConfigs().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "routes.route.nexus.vmware.com" {
		obj, err := dmClient.RouteNexusV1().Routes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "runtimes.runtime.nexus.vmware.com" {
		obj, err := dmClient.RuntimeNexusV1().Runtimes().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "tenants.tenantconfig.nexus.vmware.com" {
		obj, err := dmClient.TenantconfigNexusV1().Tenants().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "tenants.tenantruntime.nexus.vmware.com" {
		obj, err := dmClient.TenantruntimeNexusV1().Tenants().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return obj
	}
	if crdName == "users.user.nexus.vmware.com" {
		obj, err := dmClient.UserNexusV1().Users().Get(context.TODO(), name, metav1.GetOptions{})
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
