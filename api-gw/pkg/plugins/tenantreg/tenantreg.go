package tenantreg

import (
	"fmt"

	tenant_config_v1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var tenantRegistrationLog = ctrl.Log.WithName("tenantreg")

type TenantRegistration interface {
	Name() string
	RegisterTenant(tenant_config_v1.Tenant) error
	UnregisterTenant(string) error
}

var plugins map[string]TenantRegistration

func AddTenantRegPlugin(plugin TenantRegistration) error {

	if plugin.Name() == "" {
		return fmt.Errorf("add tenant reg plugin failed: plugin name cannot be empty")
	}

	if _, ok := plugins[plugin.Name()]; ok {
		return fmt.Errorf("add tenant reg plugin failed: plugin %s already registered", plugin.Name())
	}

	plugins[plugin.Name()] = plugin
	return nil
}

func RemoveTenantRegPlugin(pluginName string) error {
	if pluginName == "" {
		return fmt.Errorf("remove tenant reg plugin failed: plugin name cannot be empty")
	}

	delete(plugins, pluginName)
	return nil
}

func RegisterTenant(tenant tenant_config_v1.Tenant) bool {
	allPluginsSuccess := true
	var err error
	for name, plugin := range plugins {
		if err = plugin.RegisterTenant(tenant); err != nil {
			allPluginsSuccess = false
			tenantRegistrationLog.Error(err, fmt.Sprintf("tenant reg failed for plugin %s", name))
		}
	}
	return allPluginsSuccess
}

func UnregisterTenant(tenantName string) bool {
	allPluginsSuccess := true
	var err error
	for name, plugin := range plugins {
		if err = plugin.UnregisterTenant(tenantName); err != nil {
			allPluginsSuccess = false
			tenantRegistrationLog.Error(err, fmt.Sprintf("tenant unreg failed for plugin %s", name))
		}
	}
	return allPluginsSuccess
}
