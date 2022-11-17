package rest

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type RestData struct {
	Spec    []SpecData          `yaml:"spec" json:"spec"`
	FuncMap map[string]SpecData `json:"func_map"`
}

type SpecData struct {
	Name   string `yaml:"name" json:"name"`
	Method string `yaml:"method" json:"method"`
	Path   string `yaml:"path" json:"path"`
	Url    string
	Data   string `yaml:"data" json:"data"`
}

func (r *RestData) GetRestData(filePath string) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, r)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func (r *RestData) ReadRestData(data io.Reader) {
	decoder := json.NewDecoder(data)
	err := decoder.Decode(r)
	if err != nil {
		log.Printf("Read err   #%v ", err)
	}
}

func (r *RestData) ProcessRestCalls(apiGateway string) {
	r.FuncMap = make(map[string]SpecData)
	for i := 0; i < len(r.Spec); i++ {
		spec := &r.Spec[i]
		url := apiGateway + spec.Path
		spec.Url = url
		spec.Data = strings.TrimSpace(spec.Data)
		r.FuncMap[spec.Name] = *spec
	}
}
