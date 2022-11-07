package apply

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"

	gyaml "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/log"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/auth"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/utils"
	"gopkg.in/resty.v1"
	yaml "gopkg.in/yaml.v2"
)

var (
	GetResourceFile    string
	Labels             string
	DefaultGetHelpFunc func(*cobra.Command, []string)
)

func GetResource(cmd *cobra.Command, args []string) error {
	log.Debugf("Args: %v GetResourceFile: %v Labels: %v", args, GetResourceFile, Labels)

	serverInfo, err := auth.ServerInfo()
	if err != nil {
		log.Errorf("Get serverInfo failed with error %v", err)
		return err
	}

	token, err := auth.IdToken()
	if err != nil {
		log.Errorf("Get Id token failed with error %v", err)
		return err
	}

	var resourceName, apiVersion, objName, kindName string
	// GET Request with Positional Args
	if len(args) != 0 {
		metadata := strings.SplitN(args[0], ".", 2)

		if len(metadata) > 1 {
			resourceName = metadata[0]
			apiVersion = metadata[1]
		}

		if len(args) > 1 {
			objName = args[1]
		}

		shortApiVersion, shortResourceName, _ := GetShortName(args[0], token, serverInfo)
		if shortApiVersion != "" && shortResourceName != "" {
			return GetRequest(token, shortApiVersion, shortResourceName, objName, kindName, Labels, serverInfo)
		}
		return GetRequest(token, apiVersion, resourceName, objName, kindName, Labels, serverInfo)
	}

	// GET Request with Declarative Config Yaml
	err = utils.IsFileExist(GetResourceFile)
	if err != nil {
		log.Errorf("File %v read failed with error: %v\n", GetResourceFile, err)
		return err
	}
	yamlFile, err := ioutil.ReadFile(GetResourceFile)
	if err != nil {
		msg := fmt.Sprintf("error: cannot read file \"%s\"", GetResourceFile)
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

		apiVersion = utils.GetAPIVersion(input)
		resourceName = utils.GetResourceName(input)
		objName = utils.GetObjectName(input)
		kindName = utils.GetKindName(input)
		labels := utils.GetLabels(input)
		log.Debugf("Resource Labels: %v", labels)

		if len(apiVersion) == 0 || len(resourceName) == 0 || len(kindName) == 0 || len(objName) == 0 {
			return fmt.Errorf("the given custom resource file is not a standard k8s YAML, apiVersion or kind may be missing ")
		}
		GetRequest(token, apiVersion, resourceName, objName, kindName, labels, serverInfo)
	}
	return nil
}

func GetRequest(token, apiVersion, resourceName, objName, kindName, labels string, serverInfo *auth.Server) error {
	var (
		resp *resty.Response
		url  string
	)

	apiVersionSuffix := "/v1"

	if !strings.HasSuffix(apiVersion, apiVersionSuffix) {
		apiVersion = apiVersion + apiVersionSuffix
	}
	//for Non-CSP Cluster
	if !serverInfo.CSPEnabled && !serverInfo.InSecure {
		url = fmt.Sprintf("https://%s/apis/%s/%s/%s?labelSelector=%s&token=%s", serverInfo.Name, apiVersion, resourceName, objName, labels, token)
		RestClient.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
	} else if serverInfo.InSecure { // Local/Kind Clusters
		url = fmt.Sprintf("http://%s/apis/%s/%s/%s?labelSelector=%s", serverInfo.Name, apiVersion, resourceName, objName, labels)
		RestClient.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
	} else { // CSP enabled Cluster
		url = fmt.Sprintf("https://%s/apis/%s/%s/%s?labelSelector=%s", serverInfo.Name, apiVersion, resourceName, objName, labels)
		accessToken, err := auth.AccessToken()
		if err != nil {
			log.Errorf("Get access token failed with error %v", err)
			return err
		}
		RestClient.SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"csp-auth-token": accessToken,
		})
	}

	log.Debugf(url)
	resp, err := RestClient.R().Get(url)
	if err != nil {
		log.Errorf("Error while get request is sent to saas server %v", err)
		return err
	}

	if len(kindName) != 0 {
		fmt.Printf("%s/%s \n", kindName, objName)
	}

	var data interface{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return fmt.Errorf("error while unmarshaling the response body: %v", err)
	}
	log.Debugf("Response Body: %v", data)

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		log.Debugf("No of Elements: %v\n", s.Len())
		for i := 0; i < s.Len(); i++ {
			newVal, _ := json.Marshal(s.Index(i).Interface())
			y, _ := gyaml.JSONToYAML(newVal)
			fmt.Printf("%s", string(y))

		}
	case reflect.Map:
		s := reflect.ValueOf(data)
		newVal, _ := json.Marshal(s.Interface())
		y, _ := gyaml.JSONToYAML(newVal)
		fmt.Printf("%s", string(y))
	}
	return nil
}

func GetSpecRequest(token, crd string, serverInfo *auth.Server) error {
	var (
		resp *resty.Response
		url  string
	)
	//for Non-CSP Cluster
	if !serverInfo.CSPEnabled && !serverInfo.InSecure {
		url = fmt.Sprintf("https://%s/declarative/apis?crd=%s&token=%s", serverInfo.Name, crd, token)
		RestClient.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
	} else if serverInfo.InSecure { // Local/Kind Clusters
		url = fmt.Sprintf("http://%s/declarative/apis?crd=%s&token=%s", serverInfo.Name, crd, token)
		RestClient.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
	} else { // CSP enabled Cluster
		url = fmt.Sprintf("https://%s/declarative/apis?crd=%s", serverInfo.Name, crd)
		accessToken, err := auth.AccessToken()
		if err != nil {
			log.Errorf("Get access token failed with error %v", err)
			return err
		}
		RestClient.SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"csp-auth-token": accessToken,
		})
	}

	log.Debugf(url)
	resp, err := RestClient.R().Get(url)
	if err != nil {
		log.Errorf("Error while get request is sent to saas server %v", err)
		return err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("%s spec not found", crd)
	}

	fmt.Print(string(resp.Body()))

	return nil
}

func GetApisRequest(token string, serverInfo *auth.Server) (map[string]interface{}, error) {
	var resp *resty.Response
	var url string

	//for Non-CSP Cluster
	if !serverInfo.CSPEnabled && !serverInfo.InSecure {
		url = fmt.Sprintf("https://%s/declarative/apis?token=%s", serverInfo.Name, token)
		RestClient.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
	} else if serverInfo.InSecure { // Local/Kind Clusters
		url = fmt.Sprintf("http://%s/declarative/apis", serverInfo.Name)
		RestClient.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
	} else { // CSP enabled Cluster
		url = fmt.Sprintf("https://%s/declarative/apis", serverInfo.Name)
		accessToken, err := auth.AccessToken()
		if err != nil {
			log.Errorf("get access token failed with error %v", err)
			return nil, err
		}
		RestClient.SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"csp-auth-token": accessToken,
		})
	}

	log.Debugf(url)
	resp, err := RestClient.R().Get(url)
	if err != nil {
		log.Errorf("error while get request is sent to saas server %v", err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshaling the response body: %v", err)
	}
	//log.Debugf("Response Body: %v", data)
	return data, nil
}

func GetApisList(token string, serverInfo *auth.Server) error {
	data, err := GetApisRequest(token, serverInfo)
	if err != nil {
		return fmt.Errorf("error while executing get apis request %v", err)
	}

	var apis []string
	apisMap := make(map[string]string)
	for _, val := range data {
		methods := val.(map[string]interface{})
		var group, resource, short string
		for k, methodInfo := range methods {
			if info, ok := methodInfo.(map[string]interface{}); ok {
				if v, ok := info["kind"]; ok {
					resource = strings.ToLower(utils.ToPlural(v.(string)))
				}
				if v, ok := info["group"]; ok {
					group = v.(string)
				}
			}

			if k == "short" {
				short = methodInfo.(map[string]interface{})["name"].(string)
			}
		}

		api := fmt.Sprintf("%s.%s", resource, group)
		if _, ok := apisMap[api]; !ok {
			apis = append(apis, api)
			apisMap[api] = short
		}
	}

	fmt.Printf("\nAvailable APIs:\n")
	sort.Strings(apis)
	for _, api := range apis {
		out := fmt.Sprintf("· %s", api)
		if apisMap[api] != "" {
			out += fmt.Sprintf(" ( %s )", apisMap[api])
		}
		fmt.Println(out)
	}

	return nil
}

func GetShortName(shortName, token string, serverInfo *auth.Server) (string, string, error) {
	data, err := GetApisRequest(token, serverInfo)
	if err != nil {
		return "", "", fmt.Errorf("error while executing get apis request %v", err)
	}

	for _, val := range data {
		methods := val.(map[string]interface{})
		if short, ok := methods["short"].(map[string]interface{}); ok {
			if short["name"] != shortName {
				continue
			}
		} else {
			continue
		}

		var apiVersion, resourceName string
		for _, info := range methods {
			if info, ok := info.(map[string]interface{}); ok {
				if v, ok := info["kind"]; ok {
					resourceName = strings.ToLower(utils.ToPlural(v.(string)))
				}
				if v, ok := info["group"]; ok {
					apiVersion = v.(string)
				}
			}
		}
		return apiVersion, resourceName, nil
	}

	return "", "", nil
}

func GetHelp(cmd *cobra.Command, args []string) {
	DefaultGetHelpFunc(cmd, args)

	serverInfo, err := auth.ServerInfo()
	if err != nil {
		log.Errorf("Get serverInfo failed with error %v", err)
		return
	}

	token, err := auth.IdToken()
	if err != nil {
		log.Errorf("Get Id token failed with error %v", err)
		return
	}

	err = GetApisList(token, serverInfo)
	if err != nil {
		log.Errorf("Get apis request failed with error %v", err)
	}
}

func GetSpec(cmd *cobra.Command, args []string) error {
	log.Debugf("Args: %v", args)

	serverInfo, err := auth.ServerInfo()
	if err != nil {
		log.Errorf("Get serverInfo failed with error %v", err)
		return err
	}

	token, err := auth.IdToken()
	if err != nil {
		log.Errorf("Get Id token failed with error %v", err)
		return err
	}

	apiVer, resName, _ := GetShortName(args[0], token, serverInfo)
	if apiVer != "" && resName != "" {
		return GetSpecRequest(token, fmt.Sprintf("%s.%s", resName, apiVer), serverInfo)
	}
	return GetSpecRequest(token, args[0], serverInfo)
}

func init() {
	RestClient = resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}
