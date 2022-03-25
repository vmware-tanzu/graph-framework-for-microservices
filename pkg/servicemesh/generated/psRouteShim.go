package generated

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type GeneratedPublicServiceRoute struct{}

type PublicServiceRoute struct {
	Paths      []string `mapstructure:"paths" json:"paths"`
	Target     string   `mapstructure:"target" json:"target"`
	TargetPort int      `mapstructure:"target_port" json:"target_port"`
}

func (g GeneratedPublicServiceRoute) GetUrl(input map[string]string) (string, error) {
	if val, ok := input["project"]; ok {
		base := fmt.Sprintf("/v1alpha2/projects/%s/global-namespaces", val)
		if val, ok := input["global-namespace"]; ok {
			base = fmt.Sprintf("%s/%s", base, val)
			if val, ok := input["public-service"]; ok {
				base = fmt.Sprintf("%s/public-service/%s", base, val)
				if val, ok := input["route"]; ok {
					base = fmt.Sprintf("%s/route/%s", base, val)
				}
			}
		}
		return base, nil
	} else {
		return "", errors.New("project not found in input args")
	}
}

func (g GeneratedPublicServiceRoute) GetBody(input interface{}) ([]byte, error) {
	fmt.Printf("PublicServiceRoute Prop Input: %+v \n", input)
	psRoute := PublicServiceRoute{}
	mapstructure.Decode(input, &psRoute)
	fmt.Printf("PublicServiceRoute Output: %+v \n", psRoute)
	return json.Marshal(psRoute)
}

func init() {
	GeneratedMap["project:global-namespace:public-service:route"] = GeneratedPublicServiceRoute{}
}
