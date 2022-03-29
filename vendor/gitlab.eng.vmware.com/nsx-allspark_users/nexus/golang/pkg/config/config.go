package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"

	"github.com/ghodss/yaml"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// Load loads configuration from config file
func Load(fileName string, configEnv string) ([]byte, error) {
	f := configFile(fileName, configEnv)
	if _, err := os.Stat(f); os.IsNotExist(err) {
		log.Errorf("Failed to find file %s: %s", fileName, err.Error())
		return nil, err
	}

	file, err := os.Open(f)
	if err != nil {
		log.Errorf("Failed to open file %s: %s", fileName, err.Error())
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Failed to Read file %s: %s", fileName, err.Error())
		return nil, err
	}
	return b, nil
}

// Write creates and write to a file
func Write(fileName string, configEnv string, content []byte) error {
	f := configFile(fileName, configEnv)
	err := ioutil.WriteFile(f, content, 0600)
	if err != nil {
		log.Errorf("Failed write to file %s: %s", fileName, err.Error())
		return err
	}
	return nil
}

// ToJSON marshals a proto to canonical JSON
func ToJSON(msg proto.Message) (string, error) {
	return ToJSONWithIndent(msg, "")
}

// ToJSONWithIndent marshals a proto to canonical JSON with pretty printed string
func ToJSONWithIndent(msg proto.Message, indent string) (string, error) {
	if msg == nil {
		return "", errors.New("unexpected nil message")
	}

	// Marshal from proto to json bytes
	m := jsonpb.Marshaler{Indent: indent}
	return m.MarshalToString(msg)
}

// ToYAML marshals a proto to canonical YAML
func ToYAML(msg proto.Message) (string, error) {
	js, err := ToJSON(msg)
	if err != nil {
		log.Errorf("Failed to Marshal msg to Json: %s", err)
		return "", err
	}
	yml, err := yaml.JSONToYAML([]byte(js))
	return string(yml), err
}

// ToJSONMap converts a proto message to a generic map using canonical JSON encoding
// JSON encoding is specified here: https://developers.google.com/protocol-buffers/docs/proto3#json
func ToJSONMap(msg proto.Message) (map[string]interface{}, error) {
	js, err := ToJSON(msg)
	if err != nil {
		return nil, err
	}

	// Unmarshal from json bytes to go map
	var data map[string]interface{}
	err = json.Unmarshal([]byte(js), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ApplyJSON unmarshals a JSON string into a proto message. Unknown fields will produce an
//// error unless strict is set to false.
func ApplyJSON(js string, pb proto.Message) error {
	reader := strings.NewReader(js)
	m := jsonpb.Unmarshaler{}
	// Setting this to true will allow us to ignore unknown fields in the config
	// - a scenario thats likely to happen as the config evolves but the
	// code in components are still old.
	// The decoding will only fail when the types for existing fields are changed.
	m.AllowUnknownFields = true
	if err := m.Unmarshal(reader, pb); err != nil {
		log.Errorf("Failed to decode proto: %q\n", err)
		return err
	}
	return nil
}

// ApplyYAML unmarshals a YAML string into a proto message. Unknown fields will be ignored
func ApplyYAML(yml []byte, pb proto.Message) error {
	js, err := yaml.YAMLToJSON(yml)
	if err != nil {
		log.Errorf("Failed to convert YAML to Json: %s", err.Error())
		return err
	}
	return ApplyJSON(string(js), pb)
}

func configFile(fileName string, configEnv string) string {
	return filepath.Join(configDir(configEnv), fileName)
}

func configDir(configEnv string) string {
	if configEnv == "" {
		configEnv = "HOME"
	}
	return os.Getenv(configEnv)
}

// LoadNexusConf loads configuration as proto from config file
func LoadNexusConf(fileName string, configEnv string) (*NexusConfig, error) {
	f := configFile(fileName, configEnv)
	if _, err := os.Stat(f); os.IsNotExist(err) {
		log.Errorf("Failed to find NexusConfig file %s config %s: %s", fileName, configEnv, err.Error())
		return nil, err
	}

	file, err := os.Open(f)
	if err != nil {
		log.Errorf("Failed to open file %s: %s", fileName, err.Error())
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Failed read file %s: %s", fileName, err.Error())
		return nil, err
	}

	pb := &NexusConfig{}
	if err = ApplyYAML(b, pb); err != nil {
		log.Errorf("%v\n", err)
	}

	return pb, err
}

// ReadFile returns the content of a file.
func ReadFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Errorf("Failed to find file, path: %s, err: %v", path, err)
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		log.Errorf("Failed to open file, path: %s, err: %v", path, err)
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Failed to close file, path: %s, err: %v", path, err)
		}
	}()

	return ioutil.ReadAll(file)
}

// LoadConfig reads the file and unmarshals it into the provided proto.
func LoadConfig(configFile string, pb proto.Message) error {
	bytes, err := ReadFile(configFile)
	if err != nil {
		log.Errorf("Failed to read config file, file: %s, err: %v", configFile, err)
		return err
	}

	if err := ApplyYAML(bytes, pb); err != nil {
		log.Errorf("Failed to load config proto, file: %s, err: %v", configFile, err)
		return err
	}
	return nil
}
