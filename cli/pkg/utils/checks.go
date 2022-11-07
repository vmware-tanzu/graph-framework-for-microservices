package utils

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	common "github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/common"
	"github.com/vmware-tanzu/graph-framework-for-microservices/cli/pkg/servicemesh/version"
	"gopkg.in/yaml.v2"
)

const (
	EnableDebugFlag     = "debug"
	ListPrereqFlag      = "list-prereq"
	SkipPrereqCheckFlag = "skip-prereq-check"
)

func IsDebug(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup(EnableDebugFlag).Changed
}

func VerifyAll(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup("all").Changed
}

func ListPrereq(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup(ListPrereqFlag).Changed
}

func SkipPrereqCheck(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup(SkipPrereqCheckFlag).Changed
}

func GetTagVersion(versionKey, EnvKey string) (string, error) {
	var values version.NexusValues
	resultVersion := os.Getenv(EnvKey)
	if resultVersion == "" {
		yamlFile, err := common.TemplateFs.ReadFile("values.yaml")
		if err != nil {
			return "", fmt.Errorf("error while reading version yamlFile %v", err)
		}

		err = yaml.Unmarshal(yamlFile, &values)
		if err != nil {
			return "", fmt.Errorf("error while unmarshal version yaml data %v", err)
		}
		resultVersion = reflect.ValueOf(values).FieldByName(versionKey).Field(0).String()
	}
	return resultVersion, nil
}
