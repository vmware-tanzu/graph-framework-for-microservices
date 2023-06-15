package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

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
