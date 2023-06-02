package registration_test

import (
	"api-gw/internal/tenant/registration"
	"api-gw/pkg/client"
	"api-gw/pkg/common"
	"api-gw/pkg/config"
	"api-gw/pkg/envoy"
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	reg_svc_mock "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/mocks/pkg/registration-service/global"
	reg_svc "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/registration-service/global"
	apinexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/api.nexus.vmware.com/v1"
	confignexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.vmware.com/v1"
	runtimenexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/runtime.nexus.vmware.com/v1"
	tenant_config_v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantruntime.nexus.vmware.com/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Registration Function tests", func() {

	AfterSuite(func() {
		envoy.XDSListener.Close()
	})

	It("should check registration happens correctly", func() {
		client.NexusClient = nexus_client.NewFakeClient()
		_, err := client.NexusClient.Api().CreateNexusByName(context.TODO(), &apinexusv1.Nexus{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = common.GetConfigNode(client.NexusClient, "default")
		Expect(err).NotTo(BeNil())

		_, err = client.NexusClient.Config().CreateConfigByName(context.TODO(), &confignexusv1.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "943ea6107388dc0d02a4c4d861295cd2ce24d551",
				Labels: map[string]string{
					common.DISPLAY_NAME: "default",
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		runtime, err := client.NexusClient.Runtime().CreateRuntimeByName(context.TODO(), &runtimenexusv1.Runtime{
			ObjectMeta: metav1.ObjectMeta{
				Name: "e817339e4e7bf29fa47ca62dd272b44282d271b8",
				Labels: map[string]string{
					common.DISPLAY_NAME: "default",
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		config.GlobalStaticRouteConfig = &config.GlobalStaticRoutes{
			Suffix: []string{"js", "css", "png"},
			Prefix: []string{"/home", "/allspark-static"},
		}

		envoy.Init(nil, nil, nil, logrus.Level(log.Level()))
		snap, err := envoy.GenerateNewSnapshot(nil, nil, nil, nil)
		Expect(snap).NotTo(BeNil())
		Expect(err).To(BeNil())

		ctrl := gomock.NewController(GinkgoT())
		regClient := reg_svc_mock.NewMockGlobalRegistrationClient(ctrl)

		gomock.InOrder(
			regClient.EXPECT().RegisterTenant(gomock.Any(), gomock.Any()).Return(&reg_svc.TenantResponse{
				Code: 0,
			}, nil),
		)

		gomock.InOrder(
			regClient.EXPECT().UnregisterTenant(gomock.Any(), gomock.Any()).Return(&reg_svc.TenantResponse{
				Code: 0,
			}, nil),
		)

		err = registration.AddTenantToSystem(tenant_config_v1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Name: "8088123",
				Labels: map[string]string{
					common.DISPLAY_NAME: "test",
				},
			},
			Spec: tenant_config_v1.TenantSpec{
				Name: "test",
				Skus: []string{"advance"},
			},
		}, regClient)
		Expect(err).NotTo(HaveOccurred())

		gomock.InOrder(
			regClient.EXPECT().RegisterTenant(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("could not create tenant")),
		)

		err = common.RegisterTenant(regClient, "test", reg_svc.TenantRequest_License(common.AvailableSkus["advance"]))
		Expect(err).To(HaveOccurred())

		err = common.UnregisterTenant(regClient, "test", reg_svc.TenantRequest_License(common.AvailableSkus["advance"]))
		Expect(err).NotTo(HaveOccurred())

		gomock.InOrder(
			regClient.EXPECT().UnregisterTenant(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("could not create tenant")),
		)

		err = common.UnregisterTenant(regClient, "test", reg_svc.TenantRequest_License(common.AvailableSkus["advance"]))
		Expect(err).To(HaveOccurred())

		// Test InitRuntime

		gomock.InOrder(
			regClient.EXPECT().UnregisterTenant(gomock.Any(), gomock.Any()).Return(&reg_svc.TenantResponse{
				Code: 0,
			}, nil),
		)

		runtime.AddTenant(context.Background(), &v1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
				Labels: map[string]string{
					common.DISPLAY_NAME: "test",
				},
			},
			Spec: v1.TenantSpec{
				Attributes: v1.Attributes{
					Skus: []string{"LICENSE_ADVANCE"},
				},
			},
		})

		err = registration.InitTenantRuntimeCache(regClient)
		Expect(err).To(BeNil())

		err = common.CreateTenantIfNotExists(client.NexusClient, "testing", "advance")
		Expect(err).To(BeNil())

		gomock.InOrder(
			regClient.EXPECT().RegisterTenant(gomock.Any(), gomock.Any()).Return(&reg_svc.TenantResponse{
				Code: 0,
			}, nil),
		)

		err = registration.InitTenantConfig(regClient)
		Expect(err).To(BeNil())

		_, ok := common.GetTenantState("testing")
		Expect(ok).To(BeTrue())
		runtime.AddTenant(context.Background(), &v1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testing",
				Labels: map[string]string{
					common.DISPLAY_NAME: "testing",
				},
			},
			Spec: v1.TenantSpec{
				Attributes: v1.Attributes{
					Skus: []string{"LICENSE_ADVANCE"},
				},
			},
		})

		gomock.InOrder(
			regClient.EXPECT().UnregisterTenant(gomock.Any(), gomock.Any()).Return(&reg_svc.TenantResponse{
				Code: 0,
			}, nil),
		)
		err = registration.InitTenantRuntimeCache(regClient)
		Expect(err).To(BeNil())

	})
})
