package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AuthToken struct {
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
}

type Server struct {
	Name       string `json:"name"`
	CSPEnabled bool   `json:"csp_enabled"`
	InSecure   bool   `json:"in_secure"`
}

/*
curl -X 'POST' \
  'https://console-stg.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize' \
  -H 'accept: application/json' \
  -H 'Authorization: Basic LU5sA2IOTtVvRtF0NBB2tIl2QxuUB2c0EDXcWC5nRq2OuSk0hZ58zpaO9KsYjGmx' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'api_token=&passcode=&refresh_token=8R9ZiBscnPQRSx7V9sjoSLy4ktXeH8iRIJKS9rLc5oFE4KMBdw0G2rXFNTvqAl1Y'


  curl -X 'POST'   'https://console-stg.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize'   -H 'accept: application/json'   -H 'Content-Type: application/x-www-form-urlencoded'   -d 'api_token=&passcode=&refresh_token=8R9ZiBscnPQRSx7V9sjoSLy4ktXeH8iRIJKS9rLc5oFE4KMBdw0G2rXFNTvqAl1Y'



curl -v 'https://staging-1.servicemesh.biz/tsm/v1alpha1/certificates' --header "csp-auth-token: $TOKEN"

Step 1; From API token, get csp auth token

Step 2: store response in a file

Step 3: Read from file and use csp-auth-token

*/

var (
	tokenFile  string = ".servicemesh.config"
	serverFile string = ".servicemesh.server"
)

func AccessToken() (string, error) {

	homeDir, _ := os.UserHomeDir()
	accessTokenFile := fmt.Sprintf("%s/%s", homeDir, tokenFile)
	data, err := ioutil.ReadFile(accessTokenFile)
	if err != nil {
		fmt.Printf("\nreading access token from file %s failed with error %v", accessTokenFile, err)
		return "", err
	}

	authData := AuthToken{}
	err = json.Unmarshal(data, &authData)
	if err != nil {
		fmt.Printf("\nunmarshaling token data %s failed with error %v", string(data), err)
		return "", err
	}

	return authData.AccessToken, nil
}

func ServerInfo() (*Server, error) {
	homeDir, _ := os.UserHomeDir()
	serverFile := fmt.Sprintf("%s/%s", homeDir, serverFile)
	data, err := ioutil.ReadFile(serverFile)
	if err != nil {
		fmt.Printf("\nreading access token from file %s failed with error %v", serverFile, err)
		return nil, err
	}
	serverInfo := Server{}
	err = json.Unmarshal(data, &serverInfo)
	if err != nil {
		fmt.Printf("\nunmarshaling server name %s failed with error %v", string(data), err)
		return nil, err
	}

	return &serverInfo, nil

}

func IdToken() (string, error) {

	homeDir, _ := os.UserHomeDir()
	accessTokenFile := fmt.Sprintf("%s/%s", homeDir, tokenFile)
	data, err := ioutil.ReadFile(accessTokenFile)
	if err != nil {
		fmt.Printf("\nreading access token from file %s failed with error %v", accessTokenFile, err)
		return "", err
	}

	authData := AuthToken{}
	err = json.Unmarshal(data, &authData)
	if err != nil {
		fmt.Printf("\nunmarshaling token data %s failed with error %v", string(data), err)
		return "", err
	}

	return authData.IdToken, nil
}
