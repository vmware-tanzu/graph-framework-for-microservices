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

var DeleteResourceFile string
var DeleteRestClient *resty.Client

func DeleteResource(cmd *cobra.Command, args []string) error {
	fmt.Printf("%v\n", args)
	fmt.Printf("%v\n", DeleteResourceFile)

	if err := utils.IsFileExist(DeleteResourceFile); err == nil {

		yamlFile, err := ioutil.ReadFile(DeleteResourceFile)
		if err != nil {
			msg := fmt.Sprintf("error: cannot read file \"%s\"", DeleteResourceFile)
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

			serverFQDN, err := auth.ServerName()
			if err != nil {
				fmt.Printf("get serverFQDN name failed with error %v", err)
				return err
			}

			url, _ := generated.DeleteDecode(metadata)
			fmt.Printf("URL: %v", fmt.Sprintf("https://%s/tsm%s", serverFQDN, url))

			accessToken, err := auth.AccessToken()
			if err != nil {
				fmt.Printf("get acess token failed with error %v", err)
				return err
			}

			RestClient.SetHeaders(map[string]string{
				"csp-auth-token": accessToken,
			})

			resp, err := RestClient.R().Delete(
				fmt.Sprintf("https://%s/tsm%s", serverFQDN, url))
			fmt.Printf("\nResponse Error: %v", err)
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
	DeleteRestClient = resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}
