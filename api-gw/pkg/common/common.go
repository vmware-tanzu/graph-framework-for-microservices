package common

import (
	"api-gw/pkg/client"
	"api-gw/pkg/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	userv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/user.nexus.vmware.com/v1"
	"github.com/vmware-tanzu/graph-framework-for-microservices/api/build/common"

	tenant_config_v1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	tenant_runtime_v1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/tenantruntime.nexus.vmware.com/v1"
	nexus_client "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/nexus-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
)

// Add logic for registration retries
// Increasing this for CSP cluster as they uses terraform for provisioning nodes, which would consume lot of time
var REGISTRATION_RETRIES int = 30
var REGISTRATION_WAIT_TIME time.Duration = 10

const DISPLAY_NAME = "nexus/display_name"

// URLs related to CSP interaction
const CSP_GATEWAY_ROOT = "csp/gateway"
const CSP_ORG_API_ROOT = "slc/api/v2/orgs"
const CSP_AUTHORIZE_URL = "am/api/auth/authorize"
const CSP_COMMERCE_API = "commerce/api/v3/subscriptions"
const CSP_ORG_REDIRECT_URL = "/v0/cspauth/discovery"

var NexusApps = []string{
	"nexus-tenant-runtime",
	"tsm-tenant-runtime",
}

var AllowedStates = []string{
	"Synced",
}

type Tenantstatus string

const (
	CREATING        Tenantstatus = "creating"
	CREATED                      = "created"
	DELETING                     = "deleting"
	CREATION_FAILED              = "creation_failed"
)

var CSP_PERMISSION_NAME string
var CSP_SERVICE_ID string
var CSP_SERVICE_NAME string

func GetServiceUnavailableJson(err string, state string, status string, creationDate string) map[string]interface{} {
	return map[string]interface{}{
		"featureFlag": "firstTimeExperience",
		"error":       err,
		"details": map[string]string{
			"state":         state,
			"status":        status,
			"message":       err,
			"creationStart": creationDate,
		},
	}
}

// this should be from the ENV here , which will passed from HELM chart
// currently Global API Gateway is configured with this
//const CSP_SERVICE_NAME = "external/8d190d7b-ebb4-4fc9-b4e9-fb4b14148e50/staging-1e"

func VerifyPermissions(token string, Claims jwt.Claims, permissions map[string]string) (hasAccess bool) {
	perms := Claims.(jwt.MapClaims)["perms"].([]interface{})
	for _, perm := range perms {
		for _, permA := range permissions {
			if perm.(string) == permA {
				return true
			} else {
				if strings.Contains(perm.(string), permA) {
					return true
				}
			}
		}
	}
	return false
}

func ConvertProductIDtoLicense(productid string) string {
	for license, sku := range config.SKUConfig.SKU {
		for _, pid := range sku {
			if pid == productid {
				return license
			}
		}
	}
	return ""
}

func GetCSPServiceOwnerToken() string {
	return os.Getenv("CSP_SERVICE_OWNER_TOKEN")
}

func SetCSPVariables() error {
	CSP_PERMISSION_NAME = os.Getenv(CSPPermissionName)
	fmt.Printf(CSP_PERMISSION_NAME)
	if CSP_PERMISSION_NAME == "" {
		return fmt.Errorf("%s is set to be empty", CSPPermissionName)
	} else {
		CSP_SERVICE_ID = strings.Split(CSP_PERMISSION_NAME, "/")[1]
		CSP_SERVICE_NAME = strings.Split(CSP_PERMISSION_NAME, "/")[2]
	}
	return nil
}

var Permissions map[string]string

func SetCSPPermissionOrg() {
	Permissions = map[string]string{
		"user":      fmt.Sprintf("%s:user", CSP_PERMISSION_NAME),
		"admin":     fmt.Sprintf("%s:admin", CSP_PERMISSION_NAME),
		"svc_admin": fmt.Sprintf("%s:admin", CSP_SERVICE_NAME),
		"svc_user":  fmt.Sprintf("%s:user", CSP_SERVICE_NAME),
		// "csp_billing_owner": "csp:billing_owner",
		// "csp_billin_user":   "csp:billing_user",
	}
}

func GetConfigNode(NexusClient *nexus_client.Clientset, name string) (configNode *nexus_client.ConfigConfig, err error) {
	configNodes, err := NexusClient.Config().ListConfigs(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", DISPLAY_NAME, name)})
	if err != nil {
		return &nexus_client.ConfigConfig{}, err
	}
	for _, configNode := range configNodes {
		return configNode, nil
	}
	return &nexus_client.ConfigConfig{}, fmt.Errorf("config node with name %s not found", name)
}

func GetRuntimeNode(NexusClient *nexus_client.Clientset, name string) (configNode *nexus_client.RuntimeRuntime, err error) {
	runtimeNodes, err := NexusClient.Runtime().ListRuntimes(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", DISPLAY_NAME, name)})
	if err != nil {
		return &nexus_client.RuntimeRuntime{}, err
	}
	for _, runtimeNode := range runtimeNodes {
		return runtimeNode, nil
	}
	return &nexus_client.RuntimeRuntime{}, fmt.Errorf("config node with name %s not found", name)
}

func CheckTenantIfExists(NexusClient *nexus_client.Clientset, tenantName string) (found bool, err error) {
	log.Debugf(fmt.Sprintf("checking if tenant exists for :%s", tenantName))
	configData, err := GetConfigNode(NexusClient, "default")
	if err != nil {
		return false, err
	}
	_, err = configData.GetTenant(context.Background(), tenantName)
	if err != nil {
		if nexus_client.IsChildNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

func CheckTenantRuntimeIfExists(NexusClient *nexus_client.Clientset, tenantName string) (result *nexus_client.TenantruntimeTenant, found bool, err error) {
	runtimeData, err := GetRuntimeNode(NexusClient, "default")
	if err != nil {
		return &nexus_client.TenantruntimeTenant{}, false, err
	}
	result, err = runtimeData.GetTenant(context.Background(), tenantName)
	if err != nil {
		if nexus_client.IsChildNotFound(err) {
			return &nexus_client.TenantruntimeTenant{}, false, nil
		}
		return &nexus_client.TenantruntimeTenant{}, false, err
	}
	return result, true, nil

}

func CreateTenantIfNotExists(NexusClient *nexus_client.Clientset, tenantName string, SKU string) error {
	configData, err := GetConfigNode(NexusClient, "default")
	if err != nil {
		return err
	}
	// if _, ok := AvailableSkus[SKU]; !ok {
	// 	return fmt.Errorf("SKU not supported")
	//}
	found, err := CheckTenantIfExists(NexusClient, tenantName)
	if err != nil {
		return fmt.Errorf("could not get tenant details : %v", err)
	}
	if !found && err == nil {
		tenantObj := tenant_config_v1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Name: tenantName,
			},
			Spec: tenant_config_v1.TenantSpec{
				Name:         tenantName,
				Skus:         []string{SKU},
				FeatureFlags: []string{"Project:Enable"},
			},
		}

		_, err = configData.AddTenant(context.Background(), &tenantObj)
		if err != nil {
			return err
		}
	}
	return nil
}

type TenantState struct {
	Status        Tenantstatus
	Message       string
	CreationStart string
	SKU           string
}

var (
	// UserMap is a map of UserName with CR spec. Eg: "foo" => userv1.UserSpec{}
	UserMap      = make(map[string]userv1.UserSpec)
	userMapMutex = &sync.RWMutex{}
	//Tenantruntime map is to store tenant State information
	TenantRuntimeMap      = make(map[string]TenantState)
	tenantMapMutex        = &sync.RWMutex{}
	tenantDisplayMapMutex = &sync.RWMutex{}
	TenantDisplayNameMap  = make(map[string]string)
)

func CreateUser(NexusClient *nexus_client.Clientset, tenantId string, userObj userv1.User) error {
	ok, err := CheckTenantIfExists(NexusClient, tenantId)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s Tenant not found", tenantId)
	}
	configObj, _ := GetConfigNode(NexusClient, "default")
	_, err = configObj.AddUser(context.Background(), &userObj)
	if err != nil {
		return err
	}
	// Store the users info in a local cache.
	AddUser(userObj.Spec.Username, userObj.Spec)
	return nil
}

func DeleteUserObject(NexusClient *nexus_client.Clientset, username string) error {
	configObj, err := GetConfigNode(NexusClient, "default")
	if err != nil {
		return err
	}
	_, err = configObj.GetUser(context.Background(), username)
	if err != nil {
		if nexus_client.IsChildNotFound(err) {
			DeleteUser(username)
			return nil
		}
		return err
	}
	err = configObj.DeleteUser(context.Background(), username)
	if err != nil {
		return err
	}
	// Store the users info in a local cache.
	DeleteUser(username)
	return nil
}

// AddUser adds username with spec as value to UserMap.
func AddUser(key string, val userv1.UserSpec) {
	userMapMutex.Lock()
	defer userMapMutex.Unlock()
	UserMap[key] = val
}

// DeleteUser deletes the requested entry from the UserMap.
func DeleteUser(key string) {
	userMapMutex.Lock()
	defer userMapMutex.Unlock()
	delete(UserMap, key)
}

// GetUser returns user CR spec for the username.
func GetUser(key string) (userv1.UserSpec, bool) {
	userMapMutex.RLock()
	defer userMapMutex.RUnlock()
	spec, ok := UserMap[key]
	return spec, ok
}

// AddTenantState adds tenantState with spec as value to TenantMap.
func AddTenantState(key string, val TenantState) {
	tenantMapMutex.Lock()
	defer tenantMapMutex.Unlock()
	TenantRuntimeMap[key] = val
}

// DeleteTenantState deletes the requested entry from the TenantMap.
func DeleteTenantState(key string) {
	tenantMapMutex.Lock()
	defer tenantMapMutex.Unlock()
	delete(TenantRuntimeMap, key)
}

// GetTenantState returns tenantState spec for the tenantId.
func GetTenantState(key string) (TenantState, bool) {
	tenantMapMutex.RLock()
	defer tenantMapMutex.RUnlock()
	state, ok := TenantRuntimeMap[key]
	return state, ok
}

// AddTenantState adds tenantState with spec as value to TenantMap.
func AddTenantDisplayName(key string, val string) {
	tenantDisplayMapMutex.Lock()
	defer tenantDisplayMapMutex.Unlock()
	TenantDisplayNameMap[key] = val
}

// DeleteTenantState deletes the requested entry from the TenantMap.
func DeleteTenantDisplayName(key string) {
	tenantDisplayMapMutex.Lock()
	defer tenantDisplayMapMutex.Unlock()
	delete(TenantDisplayNameMap, key)
}

// GetTenantState returns tenantState spec for the tenantId.
func GetTenantDisplayName(key string) (string, bool) {
	tenantDisplayMapMutex.RLock()
	defer tenantDisplayMapMutex.RUnlock()
	state, ok := TenantDisplayNameMap[key]
	return state, ok
}

func DeleteUsersForTenant(tenantName string) error {
	users, err := client.NexusClient.User().ListUsers(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Could not get users due to %s", err)
	}
	configObj, err := GetConfigNode(client.NexusClient, "default")
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Spec.TenantId == tenantName {
			if err := configObj.DeleteUser(context.TODO(), user.Labels[common.DISPLAY_NAME_LABEL]); err != nil {
				return err
			}
		}
		DeleteUser(user.Labels[common.DISPLAY_NAME_LABEL])
	}
	return nil
}

func InitAdminDatamodelCache() error {
	users, err := client.NexusClient.User().ListUsers(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Could not get users due to %s", err)
	}
	for _, user := range users {
		found, err := CheckTenantIfExists(client.NexusClient, user.Spec.TenantId)
		if err != nil {
			return err
		}
		if !found {
			configObj, err := GetConfigNode(client.NexusClient, "default")
			if err != nil {
				return err
			}
			if err := configObj.DeleteUser(context.TODO(), user.Labels[common.DISPLAY_NAME_LABEL]); err != nil {
				return err
			}
		}
		AddUser(user.Spec.Username, user.Spec)
	}
	return nil
}

func GetUserNameFromToken(rawDecodedCred string) string {
	return strings.Split(rawDecodedCred, ":")[0]
}

func CreateCookie(key string, value string, expires time.Time) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode
	if !expires.IsZero() {
		cookie.Expires = expires
	}
	return cookie
}

func GetTenantStatus(state tenant_runtime_v1.TenantStatus) (status Tenantstatus, message string) {
	var appsStarted int = 0
	for _, app := range state.InstalledApplications.NexusApps {
		for _, appV1 := range app.OamApp.Components {
			if slices.Contains(NexusApps, strings.Split(appV1.Name, ".")[1]) {
				if !slices.Contains(AllowedStates, appV1.Sync) {
					return CREATING, fmt.Sprintf("%s still not synced, current state: %s, Health: %s", appV1.Name, appV1.Sync, appV1.Health)
				} else {
					appsStarted = appsStarted + 1
				}
			}
		}
	}
	if appsStarted == len(NexusApps) {
		return CREATED, fmt.Sprintf("%v apps started", NexusApps)
	}
	return CREATING, fmt.Sprintf("Apps not created")
}

func GetServableTenantStatus(tenantName string) (httpStatus int, servableStatus map[string]interface{}) {
	state, ok := GetTenantState(tenantName)
	if !ok {
		return 503, GetServiceUnavailableJson("Tenant not available", "STATE_REGISTRATION", "STATUS_FAILED", "NA")
	}
	if state.Status == CREATING {
		started, _ := time.Parse(time.RFC3339Nano, state.CreationStart)
		currentTime := time.Now()
		if currentTime.Sub(started).Minutes() > float64(30) {
			return 503, GetServiceUnavailableJson("Tenant not started", "STATE_REGISTRATION", "STATUS_FAILED", state.CreationStart)
		}
		return 503, GetServiceUnavailableJson(state.Message, "STATE_REGISTRATION", "STATUS_IN_PROGRESS", state.CreationStart)
	}
	if state.Status == CREATION_FAILED {
		return 503, GetServiceUnavailableJson("Tenant not started", "STATE_REGISTRATION", "STATUS_FAILED", state.CreationStart)
	}
	return 200, map[string]interface{}{}
}

func GenerateServiceDefinitionURL(issuerURL, tenantId string) string {
	cspGwURL := fmt.Sprintf("%s/%s", issuerURL, CSP_GATEWAY_ROOT)
	sdefURL := fmt.Sprintf("%s/%s/services", CSP_ORG_API_ROOT, tenantId)
	return fmt.Sprintf("%s/%s", cspGwURL, sdefURL)
}
