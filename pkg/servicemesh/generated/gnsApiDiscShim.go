package generated

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	model "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/nexus"
)

type GeneratedApiDiscovery struct{}

type APIDiscoveryProperties struct {
	Selectors          Selectors                               `json:"selectors,omitempty" mapstructure:"selectors"`
	Application        model.ApiDiscoveryApplicationConfigList `json:"application,omitempty" mapstructure:"application"`
	CompressionEnabled bool                                    `json:"compression_enabled,omitempty" mapstructure:"compression_enabled"`
}

type Selectors struct {
	Services          []string `mapstructure:"services" json:"services"`
	AllServices       bool     `mapstructure:"all_services" json:"all_services"`
	AllPublicServices bool     `mapstructure:"all_public_services" json:"all_public_services"`
}

func (g GeneratedApiDiscovery) GetUrl(input map[string]string) (string, error) {
	if val, ok := input["project"]; ok {
		base := fmt.Sprintf("/v1alpha2/projects/%s/global-namespaces", val)
		if val, ok := input["global-namespace"]; ok {
			base = fmt.Sprintf("%s/%s", base, val)
			if val, ok := input["api-discovery"]; ok {
				base = fmt.Sprintf("%s/api-discovery/%s", base, val)
			}
		}
		return base, nil
	} else {
		return "", errors.New("project not found in input args")
	}
}

func (g GeneratedApiDiscovery) GetBody(input interface{}) ([]byte, error) {
	adProp := APIDiscoveryProperties{}
	mapstructure.Decode(input, &adProp)
	fmt.Printf("API Discovery Prop: %+v", adProp)
	return json.Marshal(adProp)
}

func init() {
	GeneratedMap["project:global-namespace:api-discovery"] = GeneratedApiDiscovery{}
}
