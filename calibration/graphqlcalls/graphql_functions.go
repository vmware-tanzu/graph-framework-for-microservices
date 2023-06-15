package graphqlcalls

import (
	"encoding/json"
	"io"
	"log"
)

type GQLData struct {
	Spec       []QueryData `yaml:"spec" json:"spec"`
	GQLFuncMap map[string]QueryData
}

type QueryData struct {
	Name  string `yaml:"name" json:"name"`
	Query string `yaml:"method" json:"method"`
}

func (r *GQLData) ReadQueryData(data io.Reader) {
	decoder := json.NewDecoder(data)
	err := decoder.Decode(r)
	if err != nil {
		log.Printf("Read err   #%v ", err)
	}
}

func (r *GQLData) ProcessGQLCalls() {
	r.GQLFuncMap = make(map[string]QueryData)
	for i := 0; i < len(r.Spec); i++ {
		spec := &r.Spec[i]
		r.GQLFuncMap[spec.Name] = *spec
	}
}
