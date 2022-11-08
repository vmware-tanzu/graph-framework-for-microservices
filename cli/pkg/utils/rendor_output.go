package utils

import (
	"bytes"
	"encoding/json"
	"strings"

	"sigs.k8s.io/yaml"
)

var indentation = "    "

const (
	renderOutputTypeJSON = "json"
	renderOutputTypeYAML = "yaml"
)

// PrettyJSON ...Converts JSON in bytes to pretty JSON
func PrettyJSON(in []byte) ([]byte, error) {
	var out bytes.Buffer

	err := json.Indent(&out, in, "", indentation)
	if err != nil {
		return in, err
	}
	return out.Bytes(), err
}

// PrettyYAML ...Converts JSON in bytes to pretty YAML
func PrettyYAML(in []byte) ([]byte, error) {
	json, _ := PrettyJSON(in)
	yaml, err := yaml.JSONToYAML(json)
	if err != nil {
		return in, err
	}
	return yaml, err
}

// RenderOutput ... returns json or yaml based on render type
func RenderOutput(in []byte, renderType string) []byte {

	switch strings.ToLower(renderType) {
	case renderOutputTypeJSON:
		json, _ := PrettyJSON(in)
		return json
	case renderOutputTypeYAML:
		yaml, _ := PrettyYAML(in)
		return yaml
	}
	// default
	json, _ := PrettyJSON(in)
	return json
}
