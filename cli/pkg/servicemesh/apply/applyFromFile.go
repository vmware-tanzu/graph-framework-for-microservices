package apply

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	yamltojson "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/log"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/servicemesh/auth"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/utils"
	"gopkg.in/resty.v1"
	yaml "gopkg.in/yaml.v2"
)

var (
	CreateResourceFile string
	RestClient         *resty.Client
	IsTSMCli           bool
)

func ApplyResource(cmd *cobra.Command, args []string) error {
	log.Debugf("Args: %v", args)
	log.Debugf("CreateResourceFile: %v", CreateResourceFile)

	err := utils.IsFileExist(CreateResourceFile)
	if err != nil {
		log.Errorf("File %v read failed with error: %v\n", CreateResourceFile, err)
		return err
	}

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

		var url string
		apiVersion := utils.GetAPIVersion(input)
		resourceName := utils.GetResourceName(input)
		objName := utils.GetObjectName(input)
		kindName := utils.GetKindName(input)

		if len(apiVersion) == 0 || len(resourceName) == 0 || len(kindName) == 0 || len(objName) == 0 {
			return fmt.Errorf("the given custom resource file is not a standard k8s YAML, apiVersion or kind may be missing ")
		}

		IsTSMCli = strings.Contains(apiVersion, common.NexusGroupSuffix)
		log.Debugf("Checking the request is for TSM CLI: %s:%s:%b", apiVersion, common.NexusGroupSuffix, IsTSMCli)

		specMap := make(map[string]interface{})
		err = yaml.Unmarshal([]byte(doc), &specMap)
		if err != nil {
			log.Errorf("Error while Unmarshal the spec %v", err)
			return err
		}
		log.Debugf("Spec: %v", specMap["spec"])

		if IsTSMCli {
			specMap["apiVersion"] = ""
			specMap["kind"] = ""
		}

		b, err := yaml.Marshal(specMap)
		if err != nil {
			panic(err)
		}

		serverInfo, err := auth.ServerInfo()
		if err != nil {
			log.Errorf("Get serverInfo failed with error %v", err)
			return err
		}
		body, err := yamltojson.YAMLToJSON(b)
		if err != nil {
			log.Errorf("Yaml to Json conversion failed with error %v", err)
			return err
		}

		var resp *resty.Response
		//for Non-CSP Cluster
		if !serverInfo.CSPEnabled && !serverInfo.InSecure {
			token, err := auth.IdToken()
			if err != nil {
				log.Errorf("Get acess token failed with error %v", err)
				return err
			}
			url = fmt.Sprintf("https://%s/apis/%s/%s?token=%s", serverInfo.Name, apiVersion, resourceName, token)
			RestClient.SetHeaders(map[string]string{
				"Content-Type": "application/json",
			})
		} else if serverInfo.InSecure { // Local/Kind Clusters
			url = fmt.Sprintf("http://%s/apis/%s/%s", serverInfo.Name, apiVersion, resourceName)
			RestClient.SetHeaders(map[string]string{
				"Content-Type": "application/json",
			})
		} else { // CSP enabled Cluster
			url = fmt.Sprintf("https://%s/apis/%s/%s", serverInfo.Name, apiVersion, resourceName)
			accessToken, err := auth.AccessToken()
			if err != nil {
				log.Errorf("Get acess token failed with error %v", err)
				return err
			}
			RestClient.SetHeaders(map[string]string{
				"Content-Type":   "application/json",
				"csp-auth-token": accessToken,
			})
		}

		log.Debugf("Resource URL: %v", url)
		log.Debugf("Resource Body: %s", string(body))
		if IsTSMCli {
			resp, err = RestClient.R().SetBody(string(body)).Put(url)
		} else {
			resp, err = RestClient.R().SetBody(string(body)).Post(url)
		}
		if err != nil {
			log.Errorf("Error while put request is sent to saas server %v", err)
			return err
		}

		if resp != nil && resp.RawResponse != nil {
			log.Debugf(resp.RawResponse.Status)
			log.Debugf(string(resp.Body()))
		}
		statusOK := resp.StatusCode() >= 200 && resp.StatusCode() < 300
		if !statusOK {
			fmt.Printf("%s/%s creation failed, HTTP Response Status: %d\n", kindName, objName, resp.StatusCode())
			return nil
		}
		fmt.Printf("%s/%s created\n", kindName, objName)
	}
	return nil
}

func init() {
	RestClient = resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}
