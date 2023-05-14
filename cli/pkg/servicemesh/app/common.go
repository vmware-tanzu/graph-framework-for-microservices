package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type DmType string

type NexusDmProps struct {
	Location       string `yaml:"location"`
	IsDefault      bool   `yaml:"isDefault"`
	BuildDirectory string `yaml:"buildDirectory,omitempty"`
}

const nexusDmsFile = "nexus-dms.yaml"

func WriteToNexusDms(DmName string, DmProps NexusDmProps) error {
	_, err := os.Stat(nexusDmsFile)
	if err != nil {
		fmt.Printf("Creating %s\n", nexusDmsFile)
		_, err = os.Create(nexusDmsFile)
		if err != nil {
			return fmt.Errorf("Couldn't create file %s\n", nexusDmsFile)
		}
	}
	data, err := ioutil.ReadFile(nexusDmsFile)
	if err != nil {
		return fmt.Errorf("Could not read %s\n", nexusDmsFile)
	}

	var nexusDmMap map[string]NexusDmProps
	err = yaml.Unmarshal(data, &nexusDmMap)
	if err != nil {
		return fmt.Errorf("could not unmarshal %s\n", nexusDmsFile)
	}
	if nexusDmMap == nil {
		nexusDmMap = make(map[string]NexusDmProps)
		// set default=true if there's only 1 DM
		DmProps.IsDefault = true
	}

	if DmProps.IsDefault {
		for k, v := range nexusDmMap {
			v.IsDefault = false
			nexusDmMap[k] = v
		}
	}
	nexusDmMap[DmName] = DmProps

	data, err = yaml.Marshal(&nexusDmMap)
	if err != nil {
		return fmt.Errorf("Error while Marshaling nexus-dms. %v", err)
	}

	err = ioutil.WriteFile(nexusDmsFile, data, 0644)
	if err != nil {
		return fmt.Errorf("Could not write to %s: %v\n", nexusDmsFile, err)
	}

	return nil
}

func GetDefaultDm() (NexusDmProps, error) {
	_, err := os.Stat(nexusDmsFile)
	if err != nil {
		return NexusDmProps{}, err
	}

	data, err := ioutil.ReadFile(nexusDmsFile)
	if err != nil {
		return NexusDmProps{}, fmt.Errorf("Could not read %s\n", nexusDmsFile)
	}

	var nexusDmMap map[string]NexusDmProps
	err = yaml.Unmarshal(data, &nexusDmMap)
	if err != nil {
		return NexusDmProps{}, fmt.Errorf("could not unmarshal %s\n", nexusDmsFile)
	}

	for _, v := range nexusDmMap {
		if v.IsDefault {
			return v, nil
		}
	}
	return NexusDmProps{}, fmt.Errorf("did not find a default DM")
}

func SetDefaultDm(datamodelName string) error {
	_, err := os.Stat(nexusDmsFile)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(nexusDmsFile)
	if err != nil {
		return fmt.Errorf("Could not read %s\n", nexusDmsFile)
	}

	var nexusDmMap map[string]NexusDmProps
	err = yaml.Unmarshal(data, &nexusDmMap)
	if err != nil {
		return fmt.Errorf("could not unmarshal %s\n", nexusDmsFile)
	}

	_, contains := nexusDmMap[datamodelName]
	if !contains {
		return fmt.Errorf("Datamodel %s not found\n", datamodelName)
	}

	for k, v := range nexusDmMap {
		if k == datamodelName {
			v.IsDefault = true
		} else {
			v.IsDefault = false
		}
		nexusDmMap[datamodelName] = v
	}

	data, err = yaml.Marshal(&nexusDmMap)
	if err != nil {
		return fmt.Errorf("Error while Marshaling nexus-dms. %v", err)
	}

	err = ioutil.WriteFile(nexusDmsFile, data, 0644)
	if err != nil {
		return fmt.Errorf("Could not write to %s: %v\n", nexusDmsFile, err)
	}
	return nil
}
