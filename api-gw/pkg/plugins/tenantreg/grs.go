//go:build grs

package tenantreg

import (
	"api-gw/pkg/common"
	"api-gw/pkg/model"
	"context"
	"fmt"
	"os"

	tenant_config_v1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	reg_svc "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/registration-service/global"
	ctrl "sigs.k8s.io/controller-runtime"
)

var grsRegistrationServiceLog = ctrl.Log.WithName("grs")

var SKUstoIntMap = map[string]int32{
	"LICENSE_ADVANCE":          int32(reg_svc.TenantRequest_LICENSE_ADVANCE),
	"LICENSE_ENTERPRISE":       int32(reg_svc.TenantRequest_LICENSE_ENTERPRISE),
	"LICENSE_ADVANCE_SCALE":    int32(reg_svc.TenantRequest_LICENSE_ADVANCE_SCALE),
	"LICENSE_ENTERPRISE_SCALE": int32(reg_svc.TenantRequest_LICENSE_ENTERPRISE_SCALE),
}

type GrsTenantRegistration struct {
	conn               *model.ConnectorObject
	reg_service_client reg_svc.GlobalRegistrationClient
}

func (g GrsTenantRegistration) RegisterTenant(tenant tenant_config_v1.Tenant) error {

	if g.conn == nil || g.reg_service_client == nil {
		return fmt.Errorf("%s tenant reg plugin not ready; not connected", g.Name())
	}

	regs, err := g.reg_service_client.RegisterTenant(context.Background(), &reg_svc.TenantRequest{
		Name:    tenant.Labels[common.DISPLAY_NAME],
		License: reg_svc.TenantRequest_License(AvailableSkus[tenant.Spec.Skus[0]]),
	})

	if err != nil {
		return fmt.Errorf("tenant register failed due to error %s", err)
	}

	if regs.Code != 0 {
		return fmt.Errorf("tenant register failed as server returned code %+v", regs.Code)
	}
	return nil
}

func (g GrsTenantRegistration) UnregisterTenant(tenantName string) error {
	if g.conn == nil || g.reg_service_client == nil {
		return fmt.Errorf("%s tenant unreg plugin not ready; not connected", g.Name())
	}

	regs, err := g.reg_service_client.UnregisterTenant(context.Background(), &reg_svc.TenantRequest{Name: tenantName})
	if err != nil {
		return fmt.Errorf("tenant unregister failed due to error %s", err)
	}

	if regs.Code != 0 {
		return fmt.Errorf("tenant unregister failed as server returned code %+v", regs.Code)
	}
	return nil
}

func (g GrsTenantRegistration) Name() string {
	return "grs"
}

func (g *GrsTenantRegistration) connect() {

	grpcConnector := model.ConnectorObject{
		Service:  "global-registration-service:30031",
		Protocol: "grpc",
	}

	err := grpcConnector.InitConnection()
	if err != nil {
		grsRegistrationServiceLog.Error(err, "unable to reconcile TenantConfig")
		os.Exit(1)
	}

	g.reg_service_client = reg_svc.NewGlobalRegistrationClient(grpcConnector.Connection)
	g.conn = &grpcConnector
}

func init() {

	var grs GrsTenantRegistration
	// Register GRS Tenant Registration as one of the tenant registration plugins.
	AddTenantRegPlugin(grs)

	// Initiate connection to GRS.
	go grs.connect()
}
