package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

type Config struct {
	GroupName     string `yaml:"groupName"`
	CrdModulePath string `yaml:"crdModulePath"`
}

func LoadConfig(configFile string) (*Config, error) {
	var config *Config
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %s", err)
	}
	configStr, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}
	err = yaml.Unmarshal(configStr, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}
	if config.CrdModulePath == "" {
		config.CrdModulePath = "gitlab.eng.vmware.com/nsx-allspark_users/m7/policymodel.git/pkg/apis/"
	}
	return config, nil
}
