package version

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gopkg.in/yaml.v2"
)

type NexusValues struct {
	NexusCli                versionFields `yaml:"nexusCli"`
	NexusCompiler           versionFields `yaml:"nexusCompiler"`
	NexusAppTemplates       versionFields `yaml:"nexusAppTemplates"`
	NexusDatamodelTemplates versionFields `yaml:"nexusDatamodelTemplates"`
	NexusRuntime            versionFields `yaml:"nexusRuntime"`
	TSMCli                  versionFields `yaml:"tsmCli"`
}

type versionFields struct {
	Version string `yaml:"version"`
}

func Version(cmd *cobra.Command, args []string) error {
	var values NexusValues

	if err := GetNexusValues(&values); err != nil {
		return err
	}

	fmt.Printf("NexusCli: %s\n", values.NexusCli.Version)
	fmt.Printf("NexusCompiler: %s\n", values.NexusCompiler.Version)
	fmt.Printf("NexusAppTemplates: %s\n", values.NexusAppTemplates.Version)
	fmt.Printf("NexusDatamodelTemplates: %s\n", values.NexusDatamodelTemplates.Version)
	fmt.Printf("NexusRuntimeManifets: %s\n", values.NexusRuntime.Version)
	return nil

}

func TSMVersion(cmd *cobra.Command, args []string) error {
	var values NexusValues

	if err := GetNexusValues(&values); err != nil {
		return err
	}

	fmt.Printf("TSMCli: %s\n", values.TSMCli.Version)
	return nil

}

func GetNexusValues(values *NexusValues) error {
	yamlFile, err := common.TemplateFs.ReadFile("values.yaml")
	if err != nil {
		return fmt.Errorf("error while reading version yamlFile %v", err)

	}

	err = yaml.Unmarshal(yamlFile, &values)
	if err != nil {
		return fmt.Errorf("error while unmarshal version yaml data %v", err)
	}

	return nil
}

func GetLatestNexusVersion() (string, error) {
	const gcrImages = "https://gcr.io/v2/nsx-sm/nexus/nexus-cli/tags/list"
	output, err := exec.Command("curl", gcrImages).Output()
	if err != nil {
		errMsg := fmt.Sprintf("Nexus CLI Upgrade check: failed to fetch latest tag from nexus-cli image registry. Please ensure you are able to access: `%s`", gcrImages)
		fmt.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}
	if len(output) == 0 {
		return "", fmt.Errorf("no tags found")
	}

	var curlResp map[string]interface{}

	err = json.Unmarshal(output, &curlResp)
	if err != nil {
		return "", fmt.Errorf("while parsing the image tags, unmarshal failed: %v", err)
	}

	for k, v := range curlResp {
		if k == "manifest" {
			for _, v1 := range v.(map[string]interface{}) {
				for k2, v2 := range v1.(map[string]interface{}) {
					if k2 == "tag" {
						tags := v2.([]interface{})
						if len(tags) == 2 {
							for _, tag := range tags {
								if tag == "latest" {
									continue
								}
								return tag.(string), nil
							}
						}
					}
				}
			}
		}
	}
	return "", fmt.Errorf("could not get latest version")
}

func format(versionString string) string {
	if strings.HasPrefix(versionString, "v") {
		versionString = strings.TrimPrefix(versionString, "v")
	}
	if strings.HasSuffix(versionString, "\n") {
		versionString = strings.TrimSuffix(versionString, "\n")
	}
	return versionString
}
