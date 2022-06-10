package nexus_client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	baseClientset "golang-appnet.eng.vmware.com/nexus-sdk/api/build/client/clientset/versioned"
	"golang-appnet.eng.vmware.com/nexus-sdk/api/build/helper"

	baseapisnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/apis.nexus.org/v1"
	baseauthenticationnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.org/v1"
	baseconfignexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.org/v1"
	baseextensionsnexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/extensions.nexus.org/v1"
	basegatewaynexusorgv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/gateway.nexus.org/v1"
)

type Clientset struct {
	baseClient            *baseClientset.Clientset
	apisNexusV1           *ApisNexusV1
	extensionsNexusV1     *ExtensionsNexusV1
	authenticationNexusV1 *AuthenticationNexusV1
	configNexusV1         *ConfigNexusV1
	gatewayNexusV1        *GatewayNexusV1
}

func NewForConfig(config *rest.Config) (*Clientset, error) {
	baseClient, err := baseClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client := &Clientset{}
	client.baseClient = baseClient
	client.apisNexusV1 = newApisNexusV1(client)
	client.extensionsNexusV1 = newExtensionsNexusV1(client)
	client.authenticationNexusV1 = newAuthenticationNexusV1(client)
	client.configNexusV1 = newConfigNexusV1(client)
	client.gatewayNexusV1 = newGatewayNexusV1(client)

	return client, nil
}

func (c *Clientset) ApisNexusV1() *ApisNexusV1 {
	return c.apisNexusV1
}
func (c *Clientset) ExtensionsNexusV1() *ExtensionsNexusV1 {
	return c.extensionsNexusV1
}
func (c *Clientset) AuthenticationNexusV1() *AuthenticationNexusV1 {
	return c.authenticationNexusV1
}
func (c *Clientset) ConfigNexusV1() *ConfigNexusV1 {
	return c.configNexusV1
}
func (c *Clientset) GatewayNexusV1() *GatewayNexusV1 {
	return c.gatewayNexusV1
}

type ApisNexusV1 struct {
	apis *apiApisNexusV1
}

func newApisNexusV1(client *Clientset) *ApisNexusV1 {
	return &ApisNexusV1{
		apis: &apiApisNexusV1{
			client: client,
		},
	}
}

type apiApisNexusV1 struct {
	client *Clientset
}

func (obj *ApisNexusV1) Apis() *apiApisNexusV1 {
	return obj.apis
}

type ExtensionsNexusV1 struct {
	extensions *extensionExtensionsNexusV1
}

func newExtensionsNexusV1(client *Clientset) *ExtensionsNexusV1 {
	return &ExtensionsNexusV1{
		extensions: &extensionExtensionsNexusV1{
			client: client,
		},
	}
}

type extensionExtensionsNexusV1 struct {
	client *Clientset
}

func (obj *ExtensionsNexusV1) Extensions() *extensionExtensionsNexusV1 {
	return obj.extensions
}

type AuthenticationNexusV1 struct {
	oidcs *oidcAuthenticationNexusV1
}

func newAuthenticationNexusV1(client *Clientset) *AuthenticationNexusV1 {
	return &AuthenticationNexusV1{
		oidcs: &oidcAuthenticationNexusV1{
			client: client,
		},
	}
}

type oidcAuthenticationNexusV1 struct {
	client *Clientset
}

func (obj *AuthenticationNexusV1) OIDCs() *oidcAuthenticationNexusV1 {
	return obj.oidcs
}

type ConfigNexusV1 struct {
	configs *configConfigNexusV1
}

func newConfigNexusV1(client *Clientset) *ConfigNexusV1 {
	return &ConfigNexusV1{
		configs: &configConfigNexusV1{
			client: client,
		},
	}
}

type configConfigNexusV1 struct {
	client *Clientset
}

func (obj *ConfigNexusV1) Configs() *configConfigNexusV1 {
	return obj.configs
}

type GatewayNexusV1 struct {
	gateways *gatewayGatewayNexusV1
}

func newGatewayNexusV1(client *Clientset) *GatewayNexusV1 {
	return &GatewayNexusV1{
		gateways: &gatewayGatewayNexusV1{
			client: client,
		},
	}
}

type gatewayGatewayNexusV1 struct {
	client *Clientset
}

func (obj *GatewayNexusV1) Gateways() *gatewayGatewayNexusV1 {
	return obj.gateways
}

func (obj *apiApisNexusV1) Get(ctx context.Context, name string, labels map[string]string) (result *baseapisnexusorgv1.Api, err error) {
	hashedName := helper.GetHashedName(name, labels)
	result, err = obj.client.baseClient.ApisNexusV1().Apis().Get(ctx, hashedName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if result.Spec.ConfigGvk.Name != "" {
		field, err := obj.client.ConfigNexusV1().Configs().GetByName(ctx, result.Spec.ConfigGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.Config = *field
	}

	return
}

func (obj *apiApisNexusV1) GetByName(ctx context.Context, name string) (result *baseapisnexusorgv1.Api, err error) {
	result, err = obj.client.baseClient.ApisNexusV1().Apis().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if result.Spec.ConfigGvk.Name != "" {
		field, err := obj.client.ConfigNexusV1().Configs().GetByName(ctx, result.Spec.ConfigGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.Config = *field
	}

	return
}

func (obj *extensionExtensionsNexusV1) Get(ctx context.Context, name string, labels map[string]string) (result *baseextensionsnexusorgv1.Extension, err error) {
	hashedName := helper.GetHashedName(name, labels)
	result, err = obj.client.baseClient.ExtensionsNexusV1().Extensions().Get(ctx, hashedName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return
}

func (obj *extensionExtensionsNexusV1) GetByName(ctx context.Context, name string) (result *baseextensionsnexusorgv1.Extension, err error) {
	result, err = obj.client.baseClient.ExtensionsNexusV1().Extensions().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return
}

func (obj *oidcAuthenticationNexusV1) Get(ctx context.Context, name string, labels map[string]string) (result *baseauthenticationnexusorgv1.OIDC, err error) {
	hashedName := helper.GetHashedName(name, labels)
	result, err = obj.client.baseClient.AuthenticationNexusV1().OIDCs().Get(ctx, hashedName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return
}

func (obj *oidcAuthenticationNexusV1) GetByName(ctx context.Context, name string) (result *baseauthenticationnexusorgv1.OIDC, err error) {
	result, err = obj.client.baseClient.AuthenticationNexusV1().OIDCs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return
}

func (obj *configConfigNexusV1) Get(ctx context.Context, name string, labels map[string]string) (result *baseconfignexusorgv1.Config, err error) {
	hashedName := helper.GetHashedName(name, labels)
	result, err = obj.client.baseClient.ConfigNexusV1().Configs().Get(ctx, hashedName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if result.Spec.GatewayGvk.Name != "" {
		field, err := obj.client.GatewayNexusV1().Gateways().GetByName(ctx, result.Spec.GatewayGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.Gateway = *field
	}

	if result.Spec.ApiExtensionsGvk.Name != "" {
		field, err := obj.client.ApiextensionsNexusV1().Extensions().GetByName(ctx, result.Spec.ApiExtensionsGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.ApiExtensions = *field
	}

	for k, v := range result.Spec.AuthNGvk {
		obj, err := obj.client.AuthnNexusV1().OIDCs().GetByName(ctx, v.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.AuthN[k] = *obj
	}

	return
}

func (obj *configConfigNexusV1) GetByName(ctx context.Context, name string) (result *baseconfignexusorgv1.Config, err error) {
	result, err = obj.client.baseClient.ConfigNexusV1().Configs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if result.Spec.GatewayGvk.Name != "" {
		field, err := obj.client.GatewayNexusV1().Gateways().GetByName(ctx, result.Spec.GatewayGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.Gateway = *field
	}

	if result.Spec.ApiExtensionsGvk.Name != "" {
		field, err := obj.client.ApiextensionsNexusV1().Extensions().GetByName(ctx, result.Spec.ApiExtensionsGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.ApiExtensions = *field
	}

	for k, v := range result.Spec.AuthNGvk {
		obj, err := obj.client.AuthnNexusV1().OIDCs().GetByName(ctx, v.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.AuthN[k] = *obj
	}

	return
}

func (obj *gatewayGatewayNexusV1) Get(ctx context.Context, name string, labels map[string]string) (result *basegatewaynexusorgv1.Gateway, err error) {
	hashedName := helper.GetHashedName(name, labels)
	result, err = obj.client.baseClient.GatewayNexusV1().Gateways().Get(ctx, hashedName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if result.Spec.AuthnGvk.Name != "" {
		field, err := obj.client.AuthnNexusV1().OIDCs().GetByName(ctx, result.Spec.AuthnGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.Authn = *field
	}

	return
}

func (obj *gatewayGatewayNexusV1) GetByName(ctx context.Context, name string) (result *basegatewaynexusorgv1.Gateway, err error) {
	result, err = obj.client.baseClient.GatewayNexusV1().Gateways().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if result.Spec.AuthnGvk.Name != "" {
		field, err := obj.client.AuthnNexusV1().OIDCs().GetByName(ctx, result.Spec.AuthnGvk.Name)
		if err != nil {
			return nil, err
		}
		result.Spec.Authn = *field
	}

	return
}
