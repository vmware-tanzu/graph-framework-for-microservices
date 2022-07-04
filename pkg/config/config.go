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
	RemoteEndpoint    utils.NexusEndpoint `yaml:"remoteEndpoint"`
	Dispatcher        Dispatcher          `yaml:"dispatcher"`
	IgnoredNamespaces IgnoredNamespaces   `yaml:"ignoredNamespaces"`
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
	return config, nil
}
