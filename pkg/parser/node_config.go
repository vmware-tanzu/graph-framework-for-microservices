package parser

import (
	"errors"
	"go/doc"

	"gopkg.in/yaml.v3"
)

type NexusNodeConfig struct {
	NexusRestAPIGen            NexusRestAPIGen              `yaml:"nexus-rest-api-gen"`
	NexusAPIValidationEndpoint []NexusAPIValidationEndpoint `yaml:"nexus-api-validation-endpoint"`
	NexusVersion               string                       `yaml:"nexus-version"`
}

type NexusRestAPIGen struct {
	URI      string   `yaml:"uri"`
	Methods  []string `yaml:"methods,flow"`
	Response Response `yaml:"response"`
}

type Response struct {
	Num200 Num200 `yaml:"200"`
	Num400 Num400 `yaml:"400"`
	Num401 Num401 `yaml:"401"`
}

type Num200 struct {
	Message string `yaml:"message"`
}

type Num400 struct {
	Message string `yaml:"message"`
}

type Num401 struct {
	Message string `yaml:"message"`
}

type NexusAPIValidationEndpoint struct {
	Service  string `yaml:"service"`
	Endpoint string `yaml:"endpoint"`
}

// GetNexusNodeConfig will return NexusNodeConfig for given Package and nexus node name
func GetNexusNodeConfig(pkg Package, name string) (*NexusNodeConfig, error) {
	d := doc.New(&pkg.Pkg, pkg.Name, 0)

	for _, t := range d.Types {
		if t.Name == name {
			config := &NexusNodeConfig{}
			err := yaml.Unmarshal([]byte(t.Doc), config)
			if err != nil {
				return nil, err
			}
			return config, nil
		}
	}
	return nil, errors.New("could not find valid config")
}
