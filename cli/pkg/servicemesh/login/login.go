package login

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
	"gopkg.in/resty.v1"
)

type AccessToken struct {
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
}

type ServerInfo struct {
	Name       string `json:"name"`
	CSPEnabled bool   `json:"csp_enabled"`
	InSecure   bool   `json:"in_secure"`
}

var (
	tokenFile     string = ".servicemesh.config"
	ApiToken      string
	serverFile    string = ".servicemesh.server"
	Server        string
	IsPrivateSaas bool
	IsInSecure    bool
)

func Login(cmd *cobra.Command, args []string) error {
	log.Debugf("Args: %v ApiToken: %v, Server: %v, IsPrivateSaas%v", args, ApiToken, Server, Server)

	if IsInSecure {
		serverInfo := ServerInfo{
			Name:     Server,
			InSecure: true,
		}
		homeDir, _ := os.UserHomeDir()
		file, _ := json.MarshalIndent(serverInfo, "", " ")
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", homeDir, serverFile), file, 0644)
		if err != nil {
			fmt.Printf("\nwriting server info to file %s failed with error %v", serverFile, err)
			return err
		}
		return nil
	}

	if IsPrivateSaas {
		serverInfo := ServerInfo{
			Name:       Server,
			CSPEnabled: false,
		}
		homeDir, _ := os.UserHomeDir()
		file, _ := json.MarshalIndent(serverInfo, "", " ")
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", homeDir, serverFile), file, 0644)
		if err != nil {
			fmt.Printf("\nwriting server info to file %s failed with error %v", serverFile, err)
			return err
		}

		var accessToken AccessToken
		accessToken.IdToken = ApiToken
		homeDir, _ = os.UserHomeDir()
		file, _ = json.MarshalIndent(accessToken, "", " ")
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", homeDir, tokenFile), file, 0644)
		if err != nil {
			fmt.Printf("\nwriting access token to file %s failed with error %v", tokenFile, err)
			return err
		}
		return nil
	}

	cspClient := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"accept":       "application/json",
	})
	body := fmt.Sprintf("api_token=%s", ApiToken)

	cspUrl := "https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize"
	if strings.HasPrefix(Server, "stag") {
		cspUrl = "https://console-stg.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize"
	}

	resp, err := cspClient.R().SetBody(body).Post(cspUrl)
	fmt.Printf("\nResponse Error: %v", err)
	if resp != nil && resp.RawResponse != nil {
		fmt.Printf("\nResponse Code: %v", resp.RawResponse.Status)
		fmt.Printf("\nResponse Body: %v", string(resp.Body()))

		var accessToken AccessToken
		if err := json.Unmarshal(resp.Body(), &accessToken); err != nil {
			fmt.Printf("\nJson unmarshal of %s failed with error %v", string(resp.Body()), err)
			return err
		}

		homeDir, _ := os.UserHomeDir()
		file, _ := json.MarshalIndent(accessToken, "", " ")
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", homeDir, tokenFile), file, 0644)
		if err != nil {
			fmt.Printf("\nwriting access token to file %s failed with error %v", tokenFile, err)
			return err
		}

		serverInfo := ServerInfo{
			Name:       Server,
			CSPEnabled: true,
		}

		homeDir, _ = os.UserHomeDir()
		file, _ = json.MarshalIndent(serverInfo, "", " ")
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", homeDir, serverFile), file, 0644)
		if err != nil {
			fmt.Printf("\nwriting server info to file %s failed with error %v", serverFile, err)
			return err
		}
		return nil
	}
	return err
}
