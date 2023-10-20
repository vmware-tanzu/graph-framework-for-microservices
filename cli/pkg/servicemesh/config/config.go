package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type NexusConfig struct {
	UpgradePromptDisable bool `yaml:"upgradePromptDisable"`
	DebugAlways          bool `yaml:"debugAlways"`
	SkipUpgradeCheck     bool `yaml:"skipUpgradeCheck"`
}

func getDefaultNexusConfig() NexusConfig {
	return NexusConfig{
		UpgradePromptDisable: false,
		DebugAlways:          false,
		SkipUpgradeCheck:     true,
	}
}

const nexusDir = ".nexus"
const nexusConfigFile = "config"

func initNexusConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	nexusDir := fmt.Sprintf("%s/%s", home, nexusDir)
	_, err = os.Stat(nexusDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(nexusDir, 0755)
		if err != nil {
			return err
		}
	}

	return writeNexusConfig(getDefaultNexusConfig())
}

// writeNexusConfig writes the provided nexus config to $HOME/.nexus/config
func writeNexusConfig(nexusConfig NexusConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFilePath := fmt.Sprintf("%s/%s/%s", home, nexusDir, nexusConfigFile)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		_, err := os.Create(configFilePath)
		if err != nil {
			return err
		}
	}
	data, err := yaml.Marshal(&nexusConfig)
	if err == nil {
		err = ioutil.WriteFile(configFilePath, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadNexusConfig returns the current nexus config (i.e., contents of $HOME/.nexus/config)
// if the nexus config file doesn't exist, initialize it to default values and return it
func LoadNexusConfig() (NexusConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error fetching user's home directory. Using default Nexus config...")
		return getDefaultNexusConfig(), nil
	}
	configFilePath := fmt.Sprintf("%s/%s/%s", home, nexusDir, nexusConfigFile)
	data, err := ioutil.ReadFile(configFilePath)
	if os.IsNotExist(err) {
		err = initNexusConfig()
		if err != nil {
			return NexusConfig{}, err
		}
	}
	var nexusConfig = getDefaultNexusConfig()
	err = yaml.Unmarshal(data, &nexusConfig)
	if err != nil {
		fmt.Printf("Failed to read contents of %s due to error: %s\n", configFilePath, err)
		return getDefaultNexusConfig(), nil
	}
	return nexusConfig, nil
}
