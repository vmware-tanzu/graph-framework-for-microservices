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
