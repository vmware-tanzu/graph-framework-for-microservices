package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	yamlv1 "github.com/ghodss/yaml"
	"gopkg.in/yaml.v2"
)

// Config ...Config will consume yaml and validate
type Config struct {
}

// GetConfig ... Returns a config object for resource validation
func GetConfig() *Config {
	return &Config{}
}

// ValidateYamlTags ... validates the yaml by checking if required tags are present
func (c *Config) ValidateYamlTags(yamlObject interface{}, tags []string) (bool, error) {
	switch v := yamlObject.(type) {
	case map[string]interface{}:
		// sanity checks if required tags are present
		for _, tag := range tags {
			if _, ok := v[tag]; !ok {
				msg := fmt.Sprintf("error: error validating data: \"%s\" is not set", tag)
				return false, errors.New(msg)
			}
		}
	default:
		return false, errors.New("error: bad yaml object")
	}
	return true, nil
}

// UnmarshalToYaml ... Reads the file and returns yaml object
func UnmarshalToYaml(filename string) (interface{}, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		msg := fmt.Sprintf("error: cannot read file \"%s\"", filename)
		return nil, errors.New(msg)
	}

	result := yaml.MapSlice{}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

// YAMLToJSON ... Converts YAML to JSON
func (c *Config) YAMLToJSON(in []byte) ([]byte, error) {
	out, err := yamlv1.YAMLToJSON(in)
	if err != nil {
		return out, err
	}
	return out, nil
}

// ValidateAndGetJSON ... Validates files and returns a JSON
func ValidateAndGetJSON(file string, tags []string) ([]byte, error) {
	var yamlObj interface{}

	cfg := GetConfig()
	// check if file present
	err := IsFileExist(file)
	if err != nil {
		return nil, err
	}

	// Unmarshall the file to yaml
	yamlObj, err = UnmarshalToYaml(file)
	if err != nil {
		return nil, err
	}

	// Validate if the tags are present
	valid, err := cfg.ValidateYamlTags(yamlObj, tags)
	if !valid {
		return nil, errors.New(err.Error())
	}

	out, _ := yaml.Marshal(yamlObj)
	// Convert YAML to JSON
	json, err := cfg.YAMLToJSON(out)

	return json, err
}

// IsFileExist ... Checks to see if the file exists
func IsFileExist(filename string) error {

	_, err := os.Stat(filename)
	if err != nil {
		msg := fmt.Sprintf("error: the path \"%s\" does not exist", filename)
		return errors.New(msg)
	}
	return nil
}
