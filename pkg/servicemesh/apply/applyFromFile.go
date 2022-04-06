package apply

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli/pkg/servicemesh/auth"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli/pkg/servicemesh/generated"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli/pkg/utils"
	"gopkg.in/resty.v1"
	yaml "gopkg.in/yaml.v2"
)

var CreateResourceFile string
var RestClient *resty.Client

func ApplyResource(cmd *cobra.Command, args []string) error {
	fmt.Printf("%v\n", args)
	fmt.Printf("%v\n", CreateResourceFile)

	if err := utils.IsFileExist(CreateResourceFile); err == nil {

		yamlFile, err := ioutil.ReadFile(CreateResourceFile)
		if err != nil {
			msg := fmt.Sprintf("error: cannot read file \"%s\"", CreateResourceFile)
			return errors.New(msg)
		}

		yamlDocs := strings.Split(string(yamlFile), "---")
		for _, doc := range yamlDocs {

			// Unmarshall the file to yam	result := yaml.MapSlice{}
			input := yaml.MapSlice{}
			err = yaml.Unmarshal([]byte(doc), &input)
			if err != nil {
				return err
			}

			var metadata yaml.MapSlice
			for _, value := range input {
				if value.Key.(string) == "metadata" {
					metadata = value.Value.(yaml.MapSlice)
					fmt.Printf("%v Type: %T\n", value.Value, value.Value)
				}
			}

			specMap := make(map[string]interface{})
			yaml.Unmarshal([]byte(doc), &specMap)
			fmt.Printf("specMap: %+v Type: %T", specMap["spec"], specMap["spec"])

			serverFQDN, err := auth.ServerName()
			if err != nil {
				fmt.Printf("get serverFQDN name failed with error %v", err)
				return err
			}

			url, body, _ := generated.ApplyDecode(metadata, specMap["spec"])
			fmt.Printf("URL: %v", fmt.Sprintf("https://%s/tsm%s", serverFQDN, url))
			fmt.Printf("\nJsonBody: %v\n", string(body))

			accessToken, err := auth.AccessToken()
			if err != nil {
				fmt.Printf("get acess token failed with error %v", err)
				return err
			}

			RestClient.SetHeaders(map[string]string{
				"csp-auth-token": accessToken,
				"Content-Type":   "application/json",
			})

			var resp *resty.Response
			if strings.Contains(url, "health-checks") {
				resp, err = RestClient.R().SetBody(string(body)).Post(
					fmt.Sprintf("https://%s/tsm%s", serverFQDN, url))
			} else {
				resp, err = RestClient.R().SetBody(string(body)).Put(
					fmt.Sprintf("https://%s/tsm%s", serverFQDN, url))
			}
			if err != nil {
				fmt.Printf("error while put request is sent to saas server %v", err)
				return err
			}
			if resp != nil && resp.RawResponse != nil {
				fmt.Printf("\nResponse Code: %v", resp.RawResponse.Status)
				fmt.Printf("\nResponse Body: %v", string(resp.Body()))
			}
		}

	} else {
		fmt.Printf("File %v read failed with error: %v\n", CreateResourceFile, err)
	}
	return nil
}

func init() {
	RestClient = resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// For non-CSP cluster
	/*
		RestClient = resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).SetQueryParams(map[string]string{
		"token": "y61fPvZB9XXzE5tvYKeTeHT28OnEiKz3sDZvGrSrARq3SpceVYHV8XnpnePpqhEs"})
	*/
}
