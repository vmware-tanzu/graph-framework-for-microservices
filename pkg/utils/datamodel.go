package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	common "gitlab.eng.vmware.com/nexus/cli/pkg/common"
	"gopkg.in/yaml.v2"
)

func GoToNexusDirectory() error {
	if _, err := os.Stat(common.NEXUS_DIR); os.IsNotExist(err) {
		return fmt.Errorf("%s directory not found", common.NEXUS_DIR)
	} else if err != nil {
		return fmt.Errorf("error %v trying to find directory %s", err, common.NEXUS_DIR)
	}

	if err := os.Chdir(common.NEXUS_DIR); err != nil {
		return fmt.Errorf("error %v trying to cd to directory %s", err, common.NEXUS_DIR)
	}
	return nil

}

func CheckDatamodelDirExists(datamodelName string) error {
	dmDir := datamodelName
	if _, err := os.Stat(dmDir); os.IsNotExist(err) {
		return fmt.Errorf("datamodel directory %s not found", dmDir)
	} else if err != nil {
		return fmt.Errorf("error %v trying to find datamodel directory %s", err, dmDir)
	}
	return nil
}

func StoreCurrentDatamodel(datamodelName string) error {
	_, err := os.Stat(common.NexusConfFile)
	if err != nil {
		_, err = os.Create(common.NexusConfFile)
		if err != nil {
			return err
		}
	}

	conf := common.NexusConfig{
		Name: datamodelName,
	}
	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	err = ioutil.WriteFile(common.NexusConfFile, yamlData, 0644)
	if err != nil {
		fmt.Println("Could not store current datamodel name")
		return err
	}
	return nil
}

func GetCurrentDatamodel() (string, error) {
	_, err := os.Stat(common.NexusConfFile)
	if err != nil {
		return "", fmt.Errorf("Could not get datamodelname : %s does not exists\n", common.NexusConfFile)
	}
	data, err := ioutil.ReadFile(common.NexusConfFile)
	if err != nil {
		return "", fmt.Errorf("Could not read yamlfile ")
	}
	var yamlConfig common.NexusConfig
	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		return "", fmt.Errorf("Could not unmarshal yamlfile")
	}
	return yamlConfig.Name, nil
}
