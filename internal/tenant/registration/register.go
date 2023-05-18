package registration

import (
	"api-gw/pkg/client"
	"api-gw/pkg/common"
	"api-gw/pkg/envoy"
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"

	reg_svc "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/registration-service/global"
	tenant_config_v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AddTenantToSystem(tenantconfig tenant_config_v1.Tenant, regclient reg_svc.GlobalRegistrationClient) error {
	tenantName := tenantconfig.Labels[common.DISPLAY_NAME]
	envoy.AddTenantConfig(&envoy.TenantConfig{
		Name: tenantName,
	})
	log.Info("Registering tenant in registration service")

	result, exists, err := common.CheckTenantRuntimeIfExists(client.NexusClient, tenantName)
	if err != nil {
		log.Errorf("Could not get the tenant runtime for tenant %s due to %s", tenantName, err.Error())
		return err
	}
	//Adding tenant Displayname to map , this is due to the reconciler will not provide name when delete event is triggered
	common.AddTenantDisplayName(tenantconfig.Name, tenantName)
	if !exists {
		// Add retries for creating tenant to combat grs registration
		//Adding tenant State to creation in progress , this is to make sure the status calls to get tenant does not fail till tenantruntime CR is created
		common.AddTenantState(tenantName, common.TenantState{
			Status:        common.CREATING,
			Message:       "Tenant Provisioning in progress",
			CreationStart: tenantconfig.ObjectMeta.CreationTimestamp.Format(time.RFC3339Nano),
			SKU:           tenantconfig.Spec.Skus[0],
		})
		var registration_retry int = 0
		for registration_retry < common.REGISTRATION_RETRIES {
			registration_retry = registration_retry + 1
			err := common.RegisterTenant(regclient, tenantName, reg_svc.TenantRequest_License(common.AvailableSkus[tenantconfig.Spec.Skus[0]]))
			if err != nil {
				log.Errorf("RegisterTenant Failed: exceeded maximum retries %d", common.REGISTRATION_RETRIES)
				if registration_retry == common.REGISTRATION_RETRIES {
					common.AddTenantState(tenantName, common.TenantState{
						Status:        common.CREATION_FAILED,
						Message:       "Tenant Provisioning Failure due to could not register tenant",
						CreationStart: tenantconfig.ObjectMeta.CreationTimestamp.Format(time.RFC3339Nano),
						SKU:           tenantconfig.Spec.Skus[0],
					})
					return err
				} else {
					log.Debugf("RegisterTenant Failed : continue to retry for %d time", registration_retry)
					time.Sleep(common.REGISTRATION_WAIT_TIME)
					continue
				}
			}
			break
		}
	} else {
		// Adding status of apps as the tenantruntime exists already
		status, message := common.GetTenantStatus(result.Status.AppStatus)
		if status == common.CREATED {
			common.AddTenantState(tenantName, common.TenantState{
				Status:        common.CREATED,
				Message:       message,
				CreationStart: tenantconfig.ObjectMeta.CreationTimestamp.Format(time.RFC3339Nano),
				SKU:           tenantconfig.Spec.Skus[0],
			})
		}
	}
	return nil
}

func InitTenantConfig(reg_client reg_svc.GlobalRegistrationClient) error {
	tenantconfigs, err := client.NexusClient.Tenantconfig().ListTenants(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Could not get tenant config due to %s", err)
	}
	for _, tenantconfig := range tenantconfigs {
		// Not returning error till the GRS issue is fixed for provisioning multiple tenants at once
		AddTenantToSystem(*tenantconfig.Tenant, reg_client)
	}
	return nil
}

func InitTenantRuntimeCache(reg_client reg_svc.GlobalRegistrationClient) error {
	tenantruntimes, err := client.NexusClient.Tenantruntime().ListTenants(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Could not get tenant runtime due to %s", err)
	}
	for _, tenantruntime := range tenantruntimes {
		tenantName := tenantruntime.Labels[common.DISPLAY_NAME]
		log.Debugf(fmt.Sprintf("Tenant name %s", tenantName))
		found, err := common.CheckTenantIfExists(client.NexusClient, tenantName)
		if err != nil {
			return fmt.Errorf("Could not get tenant config for tenant %s due to %s", tenantName, err)
		}
		if !found {
			log.Infof("Calling Unregister tenant as the tenant config CR is not present: %s", tenantName)
			sku := reg_svc.TenantRequest_License(common.SKUstoIntMap[tenantruntime.Spec.Attributes.Skus[0]])
			common.DeleteTenantState(tenantName)
			// Ignoring error in unregister cause it could be due to other blocking calls till mutiple tenant provisioning is supported
			common.UnregisterTenant(reg_client, tenantName, sku)

		} else {
			status, message := common.GetTenantStatus(tenantruntime.Status.AppStatus)
			log.Debugf("Tenant Runtime object for tenant: %s, object: %v, status: %v, message: %s", tenantName, tenantruntime.Status.AppStatus, status, message)
			tenantState, ok := common.GetTenantState(tenantName)
			log.Debugf("Existing Tenant status object: %v", tenantState)
			if ok {
				state := common.TenantState{
					Status:        status,
					Message:       message,
					CreationStart: tenantState.CreationStart,
					SKU:           tenantState.SKU,
				}
				common.AddTenantState(tenantName, state)
			}
		}
	}
	return nil
}
