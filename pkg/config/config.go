package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"connector/pkg/utils"
)

type Config struct {
	RemoteEndpoint           utils.NexusEndpoint `yaml:"remoteEndpoint"`
	Dispatcher               Dispatcher          `yaml:"dispatcher"`
	RemoteEndpointHost       string              `yaml:"-"`
	RemoteEndpointPort       string              `yaml:"-"`
	RemoteEndpointPath       string              `yaml:"-"`
	RemoteEndpointCert       string              `yaml:"-"`
	IgnoredNamespaces        IgnoredNamespaces   `yaml:"ignoredNamespaces"`
	StatusReplicationEnabled bool                `yaml:"-"`
}

type Dispatcher struct {
	RawWorkerTTL            string        `yaml:"workerTTL"`
	WorkerTTL               time.Duration `yaml:"-"`
	MaxWorkerCount          uint          `yaml:"maxWorkerCount"`
	CloseRequestsQueueSize  uint          `yaml:"closeRequestsQueueSize"`
	EventProcessedQueueSize uint          `yaml:"eventProcessedQueueSize"`
}

type IgnoredNamespaces struct {
	MatchNames []string `yaml:"matchNames"`
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

	config.Dispatcher.WorkerTTL, err = time.ParseDuration(config.Dispatcher.RawWorkerTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid worker ttl: %v", err)
	}
	config.RemoteEndpointHost = os.Getenv(utils.RemoteEndpointHost)
	config.RemoteEndpointPort = os.Getenv(utils.RemoteEndpointPort)
	config.RemoteEndpointPath = os.Getenv(utils.RemoteEndpointPath)
	config.RemoteEndpointCert = os.Getenv(utils.RemoteEndpointCert)
	config.StatusReplicationEnabled = os.Getenv(utils.StatusReplication) == utils.StatusEnabled

	return config, nil
}
