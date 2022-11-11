package rest

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type RestData struct {
	Spec    []SpecData `yaml:"spec"`
	FuncMap map[string]SpecData
}

type SpecData struct {
	Name   string `yaml:"name"`
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
	Data   string `yaml:"data"`
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

func (r *RestData) ProcessRestCalls(apiGateway string) {
	r.FuncMap = make(map[string]SpecData)
	for i := 0; i < len(r.Spec); i++ {
		spec := &r.Spec[i]
		url := apiGateway + spec.Path
		spec.Path = url
		spec.Data = strings.TrimSpace(spec.Data)
		r.FuncMap[spec.Name] = *spec
	}
}
