package apply

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/auth"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gopkg.in/resty.v1"
	yaml "gopkg.in/yaml.v2"
)

var (
	DeleteResourceFile string
	DeleteRestClient   *resty.Client
)

func DeleteResource(cmd *cobra.Command, args []string) error {
	log.Debugf("Args: %v", args)
	log.Debugf("DeleteResourceFile: %v", DeleteResourceFile)

	err := utils.IsFileExist(DeleteResourceFile)
	if err != nil {
		log.Errorf("File %v read failed with error: %v\n", DeleteResourceFile, err)
		return err
	}

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

		apiVersion := utils.GetAPIVersion(input)
		resourceName := utils.GetResourceName(input)
		objName := utils.GetObjectName(input)
		kindName := utils.GetKindName(input)
		labels := utils.GetLabels(input)
		log.Debugf("Resource Labels: %v", labels)

		if len(apiVersion) == 0 || len(resourceName) == 0 || len(kindName) == 0 || len(objName) == 0 {
			return fmt.Errorf("the given custom resource file is not a standard k8s YAML, apiVersion or kind may be missing ")
		}

		serverInfo, err := auth.ServerInfo()
		if err != nil {
			log.Errorf("Get serverInfo failed with error %v", err)
			return err
		}

		token, err := auth.IdToken()
		if err != nil {
			log.Errorf("Get acess token failed with error %v", err)
			return err
		}

		var (
			resp *resty.Response
			url  string
		)

		//for Non-CSP Cluster
		if !serverInfo.CSPEnabled {
			url = fmt.Sprintf("https://%s/apis/%s/%s/%s?labelSelector=%s&token=%s", serverInfo.Name, apiVersion, resourceName, objName, labels, token)
		} else { // CSP enabled Cluster
			url = fmt.Sprintf("https://%s/apis/%s/%s/%s?labelSelector=%s", serverInfo.Name, apiVersion, resourceName, objName, labels)
			accessToken, err := auth.AccessToken()
			if err != nil {
				log.Errorf("Get acess token failed with error %v", err)
				return err
			}
			RestClient.SetHeaders(map[string]string{
				"csp-auth-token": accessToken,
			})
		}
		log.Debugf(url)
		resp, err = RestClient.R().Delete(url)
		if err != nil {
			log.Errorf("Error while delete request is sent to saas server %v", err)
			return err
		}
		if resp != nil && resp.RawResponse != nil {
			log.Debugf(resp.RawResponse.Status)
			log.Debugf(string(resp.Body()))
		}
		statusOK := (resp.StatusCode() >= 200 && resp.StatusCode() < 300) || resp.StatusCode() == 404
		if !statusOK {
			fmt.Printf("%s/%s deletion failed, HTTP Response Status: %d\n", kindName, objName, resp.StatusCode())
			return nil
		}
		fmt.Printf("%s/%s deleted\n", kindName, objName)
	}
	return nil
}

func init() {
	DeleteRestClient = resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}
