package datamodel

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/pkg/common"
	"gopkg.in/yaml.v2"
)

// Currently when we create a datamodel we have a nexus.yaml Properties file , which has the groupName as of today there is no other way for user to update the properties File
// Adding this config with --set flag allows us to update multiple Properties at once, Example: nexus datamodel config --set dockerRepo=test,groupName=test.com,..
var data map[string]string
var get bool
var ConfigureCmd = &cobra.Command{
	Use:   "config",
	Short: "configure a Nexus Datamodel",
	//Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: RunSetGet,
}

func RunSetGet(cmd *cobra.Command, args []string) error {
	_, err := os.Stat(common.NexusDMPropertiesFile)
	if err != nil {
		return fmt.Errorf("could not find %s", common.NexusDMPropertiesFile)
	}
	inData, err := ioutil.ReadFile(common.NexusDMPropertiesFile)
	if err != nil {
		return fmt.Errorf("Could not read yamlfile ")
	}
	mapData := make(map[string]string)
	err = yaml.Unmarshal(inData, &mapData)
	if err != nil {
		return fmt.Errorf("Could not unmarshal yamlfile")
	}
	if data != nil {
		for key, value := range data {
			mapData[key] = value

		}
		outData, err := yaml.Marshal(mapData)
		if err != nil {
			return fmt.Errorf("Could not marshal yamlfile")
		}
		err = ioutil.WriteFile(common.NexusDMPropertiesFile, outData, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Could not write to yamlfile")
		}
	}
	if cmd.Flags().Lookup("get").Changed {
		fmt.Printf("Properties\n ")
		for key, value := range mapData {
			fmt.Printf("\t%s:%s\n", key, value)
		}
		return nil
	}
	return nil
}

func init() {
	ConfigureCmd.PersistentFlags().StringToStringVar(&data, "set", nil, "set configuration")
	ConfigureCmd.PersistentFlags().BoolVarP(&get, "get", "", false, "get config object")
}
