package utils

import (
	"time"

	clientv1alpha1 "github.com/vmware-tanzu-private/core/apis/client/v1alpha1"
	csp "github.com/vmware-tanzu-private/core/pkg/v1/auth/csp"
	"github.com/vmware-tanzu-private/core/pkg/v1/client"
)

var (
	// Servers ...
	// tanzu login --endpoint https://console-stg.cloud.vmware.com --server stage0 --staging
	// tanzu login --endpoint https://console.cloud.vmware.com --server prod1
	// must match one of these hardcoded servers
	Servers = map[string]string{
		"stage0":   "https://staging-0.servicemesh.biz",
		"stage2":   "https://staging-0.servicemesh.biz",
		"preprod1": "https://preprod-1.servicemesh.biz",
		"prod1":    "https://prod-1.nsxservicemesh.vmware.com",
		"prod2":    "https://prod-2.nsxservicemesh.vmware.com",
		"local":    "127.0.0.1",
	}

	// Port ... https
	Port string = "443"
	// TestPort ... Used for local host
	TestPort string = "4405"
)

// GetSaasURL ... Gets the SAAS instance url
// Based on the current context, this method will return the
// TSM service endpoint.
func GetSaasURL() (s string) {
	server, err := GetCurrentServer()
	if err != nil {
		return s
	}
	if url, found := Servers[server.Name]; found {
		return url + ":" + Port
	}
	return s
}

// GetTestURL ... Gets the local host
func GetTestURL() (s string) {
	if url, found := Servers["local"]; found {
		return url + ":" + TestPort
	}
	return s
}

// GetCurrentServer ... Gets the current context
func GetCurrentServer() (s *clientv1alpha1.Server, err error) {
	return client.GetCurrentServer()
}

// GetTSMConfig ... Gets the tanzu config
func GetTSMConfig() (cfg *clientv1alpha1.Config, err error) {
	return client.GetConfig()
}

// GetAccessToken ... Gets access token from config
func GetAccessToken() string {
	server, err := GetCurrentServer()
	if err != nil {
		return ""
	}
	return server.GlobalOpts.Auth.AccessToken
}

// IsExpired ... Returns True or False if the token is expired
func IsExpired(tokenExpiry time.Time) bool {
	return csp.IsExpired(tokenExpiry)
}
