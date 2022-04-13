package generated

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type GeneratedHealthCheck struct{}

type HealthCheck struct {
	Name            string `mapstructure:"name" json:"name"`
	Protocol        string `mapstructure:"protocol" json:"protocol"`
	Domain          string `mapstructure:"domain" json:"domain"`
	Port            int    `mapstructure:"port" json:"port"`
	Path            string `mapstructure:"path" json:"path"`
	HealthThreshold int    `mapstructure:"healthThreshold" json:"healthThreshold"`
	CertificateId   string `mapstructure:"certificate_id" json:"certificate_id"`
	ExternalPort    int    `mapstructure:"external_port" json:"external_port"`
	Interval        int    `mapstructure:"interval" json:"interval"`
}

func (g GeneratedHealthCheck) GetUrl(input map[string]string) (string, error) {
	if val, ok := input["project"]; ok {
		base := fmt.Sprintf("/v1alpha2/projects/%s/templates/health-checks", val)
		return base, nil
	} else {
		return "", errors.New("project not found in input args")
	}
}

func (g GeneratedHealthCheck) GetBody(input interface{}) ([]byte, error) {
	fmt.Printf("HealthCheck Prop Input: %+v \n", input)
	healthCheck := HealthCheck{}
	mapstructure.Decode(input, &healthCheck)
	fmt.Printf("HealthCheckProp Output: %+v \n", healthCheck)
	return json.Marshal(healthCheck)
}

func init() {
	GeneratedMap["project"] = GeneratedHealthCheck{}
}
