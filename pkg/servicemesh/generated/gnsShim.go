package generated

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	model "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/nexus"
)

type GeneratedGlobalNamespace struct{}

func (g GeneratedGlobalNamespace) GetUrl(input map[string]string) (string, error) {
	if _, ok := input["project"]; ok {
		base := fmt.Sprintf("/v1alpha1/global-namespaces")
		if val, ok := input["global-namespace"]; ok {
			base = fmt.Sprintf("%s/%s", base, val)
		}
		return base, nil
	} else {
		return "", errors.New("project not found in input args")
	}
}

func (g GeneratedGlobalNamespace) GetBody(input interface{}) ([]byte, error) {
	gnsProp := model.GNSProperties{}
	mapstructure.Decode(input, &gnsProp)
	fmt.Printf("GNS Prop: %+v", gnsProp)
	return json.Marshal(gnsProp)
}

func init() {
	GeneratedMap["project:global-namespace"] = GeneratedGlobalNamespace{}
}
