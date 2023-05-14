package login

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
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

/*
	curl -X 'POST' \
	  'https://console-stg.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize' \
	  -H 'accept: application/json' \
	  -H 'Authorization: Basic LU5sA2IOTtVvRtF0NBB2tIl2QxuUB2c0EDXcWC5nRq2OuSk0hZ58zpaO9KsYjGmx' \
	  -H 'Content-Type: application/x-www-form-urlencoded' \
	  -d 'api_token=8R9ZiBscnPQRSx7V9sjoSLy4ktXeH8iRIJKS9rLc5oFE4KMBdw0G2rXFNTvqAl1Y'


	  curl -X 'POST'   'https://console-stg.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize'   -H 'accept: application/json'   -H 'Content-Type: application/x-www-form-urlencoded'   -d 'api_token=8R9ZiBscnPQRSx7V9sjoSLy4ktXeH8iRIJKS9rLc5oFE4KMBdw0G2rXFNTvqAl1Y'

{"id_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpZ25pbmdfMiJ9.eyJzdWIiOiJ2bXdhcmUuY29tOmJkYjk1NTBiLWVkZTItNGMwYy1iMDQ3LTkzNWVlYTgwZTFmZCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJjdXN0b21lcl9udW1iZXIiOiIxNTU1MDM3MDcxIiwiaXNzIjoiaHR0cHM6Ly9nYXotcHJldmlldy5jc3AtdmlkbS1wcm9kLmNvbSIsImdyb3VwX25hbWVzIjpbIlJTQUFyY2hlckhvcml6b25Adm13YXJlLmNvbSIsIlRhYmxlYXUgU2VydmVyIFVzZXJzQHZtd2FyZS5jb20iLCJWTWZvdW5kYXRpb24gR3JvdXAgVVNAU3lzdGVtIERvbWFpbiIsIkFMTCBVU0VSUyBFWENMVURJTkcgRFRTQFN5c3RlbSBEb21haW4iLCJBTEwgVVNFUlMgRVhDTFVESU5HIERJVkVTVElUVVJFU0BTeXN0ZW0gRG9tYWluIiwiZy5keW4udm13YXJlX2FsbF9hY3RpdmVfZW1wbG95ZWVzQHZtd2FyZS5jb20iLCJJTUFQUE9QQmxvY2tAU3lzdGVtIERvbWFpbiIsImcuUm5ELUlULVBsdXJhbHNpZ2h0LXVzZXItZ3JvdXBAdm13YXJlLmNvbSIsIkZURXNAU3lzdGVtIERvbWFpbiIsIlVTIEZURXNAU3lzdGVtIERvbWFpbiIsImcudmNvZGVfYXV0aG9yaXplZF91c2Vyc0B2bXdhcmUuY29tIiwiZy5keW4uaGVscG5vd19hbGxfYWN0aXZlX2VtcGxveWVlc0B2bXdhcmUuY29tIiwiQUxMIFVTRVJTIGV4Y2wuIEJvdHNAU3lzdGVtIERvbWFpbiIsImcuWkVOLUdSQyBDT05UUklCVVRPUkB2bXdhcmUuY29tIiwiZy5keW4uYWxsX3VzZXJzQHZtd2FyZS5jb20iLCJnLmR5bi5hbGxfcGVvcGxlX21hbmFnZXJzQHZtd2FyZS5jb20iLCJnLmR5bi5oZWxwbm93X2FsbF9hY3RpdmVfZW1wbG95ZWVzX2V4Y2x1ZGVfcHJlaGlyZUB2bXdhcmUuY29tIiwiR2l0bGFiX0ZURV9JbnRlcm5fQVdGQFN5c3RlbSBEb21haW4iLCJIdW1hbiBBY2NvdW50c0BTeXN0ZW0gRG9tYWluIl0sImNvbnRleHRfbmFtZSI6IjgzMWEwNmFiLTk3ODEtNDg3Yy1hOWRhLTZjNzM5NzNlNTQwYSIsImdpdmVuX25hbWUiOiJEaW5lc2giLCJhdWQiOlsiY3NwX3N0Z19nYXpfaW50ZXJuYWxfY2xpZW50X2lkIl0sImF1dGhfdGltZSI6MTY0MzMwMjYzMSwiZG9tYWluIjoidm13YXJlLmNvbSIsImdyb3VwX2lkcyI6WyJ2bXdhcmUuY29tOjJhNWQ3ZGQwLTVhNWQtNDJkMy04YzFiLTAwMTdlZTk4MDBjNCIsInZtd2FyZS5jb206YjBlM2VhYzAtZjA4Yy00NWZkLTk2ZmUtOGE1MWEzZTNhZjUwIiwidm13YXJlLmNvbTo1MGQxNmE5NC03YTZhLTRjYmItODhlMC1iNDRiMzM4Y2U5N2QiLCJ2bXdhcmUuY29tOmU0OGFmZjUzLTljZDAtNGM2Mi1hZmE2LWFjMGZjMzQzMTVmMiIsInZtd2FyZS5jb206NDlkZDRmNDgtMWE4Yy00YjVmLWI0NTktZTk0NmM4Y2NjZjgxIiwidm13YXJlLmNvbTo1MjdhYzY2Zi1lNDk4LTQ2MmQtOGY5NC0zYmUyZDYzNGExMDMiLCJ2bXdhcmUuY29tOjgzMzE2N2U0LThhYzEtNGMwOC1iYzdmLWQ4N2JjMzg1MjJmMiIsInZtd2FyZS5jb206YzE2OWU1NTItN2RiNi00NzFkLTkwZmMtMDhjMzBiMzcwYTI5Iiwidm13YXJlLmNvbTo1ZDNiZGM1Ny0yMTZlLTRmODEtYmVmYy0wMjMxNTUwZjYxNzYiLCJ2bXdhcmUuY29tOmJjYTM4NzVlLTExZmItNDVjMi1hNmI4LTY0ZDhhNjEyMDk4OCIsInZtd2FyZS5jb206YTJiMmQ0MmQtZDlmZC00MjAwLWE3M2UtNzQwZGQ5MDc4MWY3Iiwidm13YXJlLmNvbTplMzY0ZDJkYy00OTVlLTQyOWQtOWU0OC1lNmQ3NzdkOTlhZjciLCJ2bXdhcmUuY29tOjBjYzUzMWMyLTA4YmUtNDA3ZC1iOWRhLWNjZDJjNzFkYWY1OCIsInZtd2FyZS5jb206YTE2MmUxNmYtMjM3ZS00NGY0LTk0OTAtZjA5YzhiNTI5YjIxIiwidm13YXJlLmNvbTplYjEwMDBjMy01ZDczLTQwZWUtYjJkZC0yNWIxZDc0NTUzYTciLCJ2bXdhcmUuY29tOmM5OGRhNTZhLTkzZTItNDg4OC1hNDkxLWZkMTRiNmNjZjJmNCIsInZtd2FyZS5jb206ZmM4NTM2MjYtNWU4ZS00YzJmLWJhNzEtZGYwNjczM2I3MzE3Iiwidm13YXJlLmNvbToxZjE4YWI4NC01NjM4LTQ2NTMtOWU5ZS00NmE2YmFhODcyMTIiLCJ2bXdhcmUuY29tOjIyMDgwYzdhLTk2M2ItNDkwOC04ZjE2LWFjNTEyYzE5MmE2YyJdLCJjb250ZXh0IjoiOWVmNWUzMWEtOGQ4Yi00YzA0LTkxNzAtZmE4YWQxNDhlZjkzIiwiZXhwIjoxNjQzMzU1MTgxLCJpYXQiOjE2NDMzNTMzODEsImZhbWlseV9uYW1lIjoiQmFraWFyYWoiLCJqdGkiOiIxZmJlNGI2ZC1iYjE0LTQ5YzktOWZmZi1iMDE0YTIzMGY1YTQiLCJlbWFpbCI6ImRiYWtpYXJhakB2bXdhcmUuY29tIiwiYWNjdCI6ImRiYWtpYXJhakB2bXdhcmUuY29tIiwidXNlcm5hbWUiOiJkYmFraWFyYWoifQ.DvrKW4sfbp9qP18PmWWxPGKpmiQRywoeytk3fbnfs4a7wwCjmBEFQlISU_kWzhRggSdyf2uuaCOIDAfz9NJs9TFM89BjiBFYgNYEQ-xC43UoKf1QEceXyIJ8qmuDbvg59hvHAiWZ2dfcVzyllCRzwUng8HMPoRNhPMPxQ162MNtYw_S_XciH9T-qksXTz2QSXhLzm4t2U6Gxv2pqP9SmxBXoJ2S1ldZsXb-tZgywb-5EcfmrR_RsCAimHR1Tzd4tEmdVxrI9z4O6WKDEiSr_vJh4nD3FH8ySPRjSQzhaP3XtT85l26EDBd-XEb7T88qSqSy32920AlFLlauQoiLTmA","token_type":"bearer","expires_in":1799,"scope":"ALL_PERMISSIONS customer_number openid group_ids group_names","access_token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpZ25pbmdfMiJ9.eyJzdWIiOiJ2bXdhcmUuY29tOmJkYjk1NTBiLWVkZTItNGMwYy1iMDQ3LTkzNWVlYTgwZTFmZCIsImlzcyI6Imh0dHBzOi8vZ2F6LXByZXZpZXcuY3NwLXZpZG0tcHJvZC5jb20iLCJjb250ZXh0X25hbWUiOiI4MzFhMDZhYi05NzgxLTQ4N2MtYTlkYS02YzczOTczZTU0MGEiLCJhenAiOiJjc3Bfc3RnX2dhel9pbnRlcm5hbF9jbGllbnRfaWQiLCJhdXRob3JpemF0aW9uX2RldGFpbHMiOltdLCJkb21haW4iOiJ2bXdhcmUuY29tIiwiY29udGV4dCI6IjllZjVlMzFhLThkOGItNGMwNC05MTcwLWZhOGFkMTQ4ZWY5MyIsInBlcm1zIjpbImNzcDpvcmdfb3duZXIiLCJjc3A6c2VydmljZV9lbmFibGVyIiwiZXh0ZXJuYWwvNGM0MTMyYWMtNWNkOC00MDE0LTliODMtNGU5OWNjMjM3NGZlL3N0YWdpbmctdGVzdC0yOnVzZXIiLCJleHRlcm5hbC82NTRjZWM0MC03YzVhLTQwODgtYmIwYi1hMDY5OWUzZDczZmIvYXBwOmRldiIsImV4dGVybmFsL2E5YjFiOGI1LTg3NTAtNGM5OS04ZTY2LTcyOTcyNTIwZjY3Yy9zdGFnaW5nLTE6dXNlciIsImNzcDpzdXBwb3J0X3VzZXIiLCJleHRlcm5hbC84ZDE5MGQ3Yi1lYmI0LTRmYzktYjRlOS1mYjRiMTQxNDhlNTAvc3RhZ2luZy0xZTp1c2VyIiwiY3NwOnNlcnZpY2Vfb3duZXIiLCJleHRlcm5hbC84ZDE5MGQ3Yi1lYmI0LTRmYzktYjRlOS1mYjRiMTQxNDhlNTAvc3RhZ2luZy0xZTphZG1pbiIsImNzcDpkZXZlbG9wZXIiLCJleHRlcm5hbC82NTRjZWM0MC03YzVhLTQwODgtYmIwYi1hMDY5OWUzZDczZmIvc3RhZ2luZy0yOnVzZXIiLCJleHRlcm5hbC8wODE1YzI3YS1kNDdiLTQyNzAtYTNjZi04MmEwOThiNDQ2Yjkvc3RhZ2luZy0wOnVzZXIiXSwiZXhwIjoxNjQzMzU1MTgxLCJpYXQiOjE2NDMzNTMzODEsImp0aSI6IjFmNjdkZmVhLTFiNmUtNGEzMC05ZTc4LWM5Y2I3NmE4NGY3NyIsImFjY3QiOiJkYmFraWFyYWpAdm13YXJlLmNvbSIsInVzZXJuYW1lIjoiZGJha2lhcmFqIn0.RLeI1XV_-wztP9yJtNpssEliAMlTYSX5nXbLeAtthQoic1nFLxvm8b0JDq3hi1cH9B7yNTyEO438pySjdr1s3f4y-9bomytxRjBaZvtg-xxabUoQYgZjfD_Ealt8FC2wIgZvkUMcYw2UbUL5hdgvPx_R9c3WwhNZ-VxRIBydIVY5mwLif9CIgin7bfxj_t6orvaEENX1w6k93ICH4d0ecRVV4ixfFBhM3HP6KEB0e0gaS-mQ8v7AnKr0D6Gk-APzF1EucMk7Y5ILGGIro0mJj-1eWOn-76x2CTkXD-r1lp7FEZTSIWXyzAxJWO5Qj6lHIUK1yHpk4wD_sQJ5_MapRg","refresh_token":"8R9ZiBscnPQRSx7V9sjoSLy4ktXeH8iRIJKS9rLc5oFE4KMBdw0G2rXFNTvqAl1Y"}

curl -v 'https://staging-1.servicemesh.biz/tsm/v1alpha1/certificates' --header "csp-auth-token: $TOKEN"

Step 1; From API token, get csp auth token

Step 2: store response in a file

Step 3: Read from file and use csp-auth-token
*/
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
