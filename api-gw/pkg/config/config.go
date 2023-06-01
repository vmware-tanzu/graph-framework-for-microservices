package config

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server             ServerConfig `json:"server" yaml:"server"`
	EnableNexusRuntime bool         `json:"enable_nexus_runtime" yaml:"enable_nexus_runtime,omitempty"`
	BackendService     string       `json:"backend_service" yaml:"backend_service,omitempty"`
	TenantApiGwDomain  string       `json:"tenant_api_gw_domain" yaml:"tenant_api_gw_domain,omitempty"`
	CustomNotFoundPage string       `json:"custom_not_found_page" yaml:"custom_not_found_page,omitempty"`
}

type ServerConfig struct {
	HttpPort string `json:"httpPort" yaml:"httpPort"`
	Address  string `json:"address" yaml:"address"`
	CertPath string `json:"certPath" yaml:"certPath"`
	KeyPath  string `json:"keyPath" yaml:"keyPath"`
}

var Cfg *Config

type GlobalStaticRoutes struct {
	Prefix []string `json:"Prefix" yaml:"Prefix"`
	Suffix []string `json:"Suffix" yaml:"Suffix"`
}

var GlobalStaticRouteConfig *GlobalStaticRoutes

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

	log.Infof("read configmap values: %+v", config)

	if config.Server.Address == "" {
		return nil, fmt.Errorf("config doesn't contain Server.Address")
	}

	if config.Server.CertPath == "" {
		return nil, fmt.Errorf("config doesn't contain Server.CertPath")
	}

	if config.Server.KeyPath == "" {
		return nil, fmt.Errorf("config doesn't contain Server.KeyPath")
	}

	return config, nil
}

type SKUMap struct {
	SKU map[string][]string `json:"sku"`
}

var SKUConfig *SKUMap

func LoadSKUConfig(configFile string) (*SKUMap, error) {
	var config *SKUMap
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
	return config, nil
}

func LoadStaticUrlsConfig(configFile string) (*GlobalStaticRoutes, error) {
	var gsRoutes *GlobalStaticRoutes
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %s", err)
	}
	configStr, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}
	err = yaml.Unmarshal(configStr, &gsRoutes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}
	return gsRoutes, nil
}
